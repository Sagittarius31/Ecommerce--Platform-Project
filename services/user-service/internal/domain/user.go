package domain

import (
	"errors"
	"time"
	"github.com/google/uuid"
)

var ErrNotFound = errors.New("not found")

type UserRole string
const (
	RoleCustomer UserRole = "customer"
	RoleAdmin    UserRole = "admin"
)

type User struct {
	ID           uuid.UUID `db:"id"            json:"id"`
	Email        string    `db:"email"         json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"`
	FirstName    string    `db:"first_name"    json:"first_name"`
	LastName     string    `db:"last_name"     json:"last_name"`
	Role         UserRole  `db:"role"          json:"role"`
	IsActive     bool      `db:"is_active"     json:"is_active"`
	CreatedAt    time.Time `db:"created_at"    json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"    json:"updated_at"`
}

type RegisterInput struct {
	Email     string `json:"email"      validate:"required,email"`
	Password  string `json:"password"   validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name"  validate:"required"`
}

type LoginInput struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UpdateProfileInput struct {
	FirstName string `json:"first_name" validate:"omitempty"`
	LastName  string `json:"last_name"  validate:"omitempty"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	User         *User  `json:"user"`
}

type UserRepository interface {
	Create(u *User) error
	FindByID(id uuid.UUID) (*User, error)
	FindByEmail(email string) (*User, error)
	Update(u *User) error
	Delete(id uuid.UUID) error
	ExistsByEmail(email string) (bool, error)
}

type UserService interface {
	Register(in RegisterInput) (*AuthResponse, error)
	Login(in LoginInput) (*AuthResponse, error)
	GetProfile(id uuid.UUID) (*User, error)
	UpdateProfile(id uuid.UUID, in UpdateProfileInput) (*User, error)
	DeleteAccount(id uuid.UUID) error
}
