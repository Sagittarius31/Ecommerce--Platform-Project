package service

import (
	"errors"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/yourname/ecommerce/user-service/internal/domain"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailExists       = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type userSvc struct {
	repo      domain.UserRepository
	jwtSecret string
	logger    *zap.Logger
}

func NewUserService(repo domain.UserRepository, secret string, logger *zap.Logger) domain.UserService {
	return &userSvc{repo: repo, jwtSecret: secret, logger: logger}
}

func (s *userSvc) Register(in domain.RegisterInput) (*domain.AuthResponse, error) {
	exists, err := s.repo.ExistsByEmail(in.Email)
	if err != nil { return nil, err }
	if exists { return nil, ErrEmailExists }
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), 12)
	if err != nil { return nil, err }
	u := &domain.User{
		ID: uuid.New(), Email: in.Email, PasswordHash: string(hash),
		FirstName: in.FirstName, LastName: in.LastName,
		Role: domain.RoleCustomer, IsActive: true,
		CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(),
	}
	if err := s.repo.Create(u); err != nil { return nil, err }
	s.logger.Info("user registered", zap.String("id", u.ID.String()))
	return s.auth(u)
}

func (s *userSvc) Login(in domain.LoginInput) (*domain.AuthResponse, error) {
	u, err := s.repo.FindByEmail(in.Email)
	if err != nil { return nil, ErrInvalidCredentials }
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(in.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}
	return s.auth(u)
}

func (s *userSvc) GetProfile(id uuid.UUID) (*domain.User, error) { return s.repo.FindByID(id) }

func (s *userSvc) UpdateProfile(id uuid.UUID, in domain.UpdateProfileInput) (*domain.User, error) {
	u, err := s.repo.FindByID(id)
	if err != nil { return nil, err }
	if in.FirstName != "" { u.FirstName = in.FirstName }
	if in.LastName != ""  { u.LastName = in.LastName }
	u.UpdatedAt = time.Now().UTC()
	if err := s.repo.Update(u); err != nil { return nil, err }
	return u, nil
}

func (s *userSvc) DeleteAccount(id uuid.UUID) error { return s.repo.Delete(id) }

type claims struct {
	UserID string          `json:"user_id"`
	Email  string          `json:"email"`
	Role   domain.UserRole `json:"role"`
	jwt.RegisteredClaims
}

func (s *userSvc) auth(u *domain.User) (*domain.AuthResponse, error) {
	at, err := s.tok(u, 15*time.Minute)
	if err != nil { return nil, err }
	rt, err := s.tok(u, 7*24*time.Hour)
	if err != nil { return nil, err }
	return &domain.AuthResponse{AccessToken: at, RefreshToken: rt, ExpiresIn: 900, User: u}, nil
}

func (s *userSvc) tok(u *domain.User, d time.Duration) (string, error) {
	c := claims{
		UserID: u.ID.String(), Email: u.Email, Role: u.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(d)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString([]byte(s.jwtSecret))
}
