package repository

import (
	"context"
	"database/sql"
	"time"

	log "github.com/sirupsen/logrus"
	"go-rest-api-boilerplate/internal/domain"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (u *userRepository) Save(ctx context.Context, user *domain.User) error {
	q := "INSERT INTO users (first_name, last_name, email, updated_at, created_at) VALUES ($1, $2, $3, $4, $5)"
	_, err := u.db.ExecContext(ctx, q, user.FirstName, user.LastName, user.Email, user.UpdatedAt, user.CreatedAt)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error save user repository")
		return err
	}

	return nil
}

func (u *userRepository) UpdateByID(ctx context.Context, id int64, user *domain.User) error {
	q := "UPDATE users SET first_name = $1, last_name = $2, email = $3, updated_at = $4 where id = $5"
	_, err := u.db.ExecContext(ctx, q, user.FirstName, user.LastName, user.Email, time.Now(), id)
	if err != nil {
		log.WithError(err).Error("error UpdateByID user repository")
		return err
	}

	return nil
}

func (u *userRepository) DeleteByID(ctx context.Context, id int64) error {
	q := "DELETE FROM users WHERE id = $1"
	_, err := u.db.ExecContext(ctx, q, id)
	if err != nil {
		log.WithError(err).Error("error DeleteByID user repository")
		return err
	}

	return nil
}

func (u *userRepository) FindAll(ctx context.Context) (*[]domain.User, error) {
	q := "SELECT * FROM users"
	rows, err := u.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]domain.User, 0)
	for rows.Next() {
		var user domain.User
		err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			log.WithError(err).Error("error while scan row")
			return nil, err
		}

		result = append(result, user)
	}

	if err = rows.Err(); err != nil {
		log.WithError(err).Error("error FindAll user repository")
		return nil, err
	}

	return &result, nil
}

func (u *userRepository) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	var user domain.User
	q := "SELECT * FROM users WHERE id = $1"
	err := u.db.QueryRowContext(ctx, q, id).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		log.WithError(err).Error("error FindByID user repository")
		return nil, err
	}

	return &user, nil
}
