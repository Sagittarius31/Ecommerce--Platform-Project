package repository

import (
	"database/sql"
	"errors"
	"time"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/yourname/ecommerce/user-service/internal/domain"
)

type userRepo struct{ db *sqlx.DB }

func NewUserRepository(db *sqlx.DB) domain.UserRepository { return &userRepo{db: db} }

func (r *userRepo) Create(u *domain.User) error {
	_, err := r.db.NamedExec(`INSERT INTO users (id,email,password_hash,first_name,last_name,role,is_active,created_at,updated_at) VALUES (:id,:email,:password_hash,:first_name,:last_name,:role,:is_active,:created_at,:updated_at)`, u)
	return err
}
func (r *userRepo) FindByID(id uuid.UUID) (*domain.User, error) {
	var u domain.User
	err := r.db.Get(&u, `SELECT * FROM users WHERE id=$1 AND is_active=true`, id)
	if errors.Is(err, sql.ErrNoRows) { return nil, domain.ErrNotFound }
	return &u, err
}
func (r *userRepo) FindByEmail(email string) (*domain.User, error) {
	var u domain.User
	err := r.db.Get(&u, `SELECT * FROM users WHERE email=$1 AND is_active=true`, email)
	if errors.Is(err, sql.ErrNoRows) { return nil, domain.ErrNotFound }
	return &u, err
}
func (r *userRepo) Update(u *domain.User) error {
	_, err := r.db.NamedExec(`UPDATE users SET first_name=:first_name,last_name=:last_name,updated_at=:updated_at WHERE id=:id`, u)
	return err
}
func (r *userRepo) Delete(id uuid.UUID) error {
	_, err := r.db.Exec(`UPDATE users SET is_active=false,updated_at=$1 WHERE id=$2`, time.Now().UTC(), id)
	return err
}
func (r *userRepo) ExistsByEmail(email string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)`, email).Scan(&exists)
	return exists, err
}
