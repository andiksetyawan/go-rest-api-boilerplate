package repository_test

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go-rest-api-boilerplate/internal/domain"
	"go-rest-api-boilerplate/internal/usecase/repository"
)

func newUserDBTest(t *testing.T) (db *sql.DB, mock sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	return
}

func TestNewUserRepository(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, _ := newUserDBTest(t)
		defer db.Close()

		repo := repository.NewUserRepository(db)
		assert.NotNil(t, repo)
	})
}

func TestUserRepository_FindAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock := newUserDBTest(t)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "updated_at", "created_at"}).
			AddRow(1, "john", "due", "john@mail.com", time.Now(), time.Now()).
			AddRow(2, "first", "name", "example@mail.com", time.Now(), time.Now())

		mock.ExpectQuery("SELECT * FROM users").WillReturnRows(rows)
		repo := repository.NewUserRepository(db)

		users, err := repo.FindAll(context.TODO())
		assert.NoError(t, err)
		assert.NotNil(t, users)

		assert.Equal(t, "john", (*users)[0].FirstName)
		assert.Equal(t, "first", (*users)[1].FirstName)
	})

	t.Run("error", func(t *testing.T) {
		db, mock := newUserDBTest(t)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "updated_at"}).
			AddRow(1, "john", "due", "john@mail.com", time.Now())

		mock.ExpectQuery("SELECT * FROM users").WillReturnRows(rows)
		repo := repository.NewUserRepository(db)

		users, err := repo.FindAll(context.TODO())
		assert.Error(t, err)
		assert.Nil(t, users)
	})
}

func TestUserRepository_FindByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock := newUserDBTest(t)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "updated_at", "created_at"}).
			AddRow(1, "john", "due", "john@mail.com", time.Now(), time.Now())

		mock.ExpectQuery("SELECT * FROM users WHERE id = $1").WithArgs(1).WillReturnRows(rows)

		repo := repository.NewUserRepository(db)
		user, err := repo.FindByID(context.TODO(), 1)
		assert.NoError(t, err)
		assert.NotNil(t, user)

		assert.Equal(t, user.ID, int64(1))
		assert.Equal(t, user.FirstName, "john")
	})

	t.Run("error", func(t *testing.T) {
		t.Run("errNoRows", func(t *testing.T) {
			db, mock := newUserDBTest(t)
			defer db.Close()

			mock.ExpectQuery("SELECT * FROM users WHERE id = $1").WithArgs(1).WillReturnError(sql.ErrNoRows)

			repo := repository.NewUserRepository(db)
			user, err := repo.FindByID(context.TODO(), 1)
			assert.Error(t, err)
			assert.Nil(t, user)

			assert.ErrorIs(t, err, sql.ErrNoRows)
		})
	})
}

func TestUserRepository_Save(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock := newUserDBTest(t)
		defer db.Close()

		user := domain.User{
			FirstName: "john",
			LastName:  "due",
			Email:     "john@email.test",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		expectSQL := "INSERT INTO users (first_name, last_name, email, updated_at, created_at) VALUES ($1, $2, $3, $4, $5)"
		mock.ExpectExec(expectSQL).WithArgs(user.FirstName, user.LastName, user.Email, user.UpdatedAt, user.CreatedAt).
			WillReturnResult(sqlmock.NewResult(0, 1))

		repo := repository.NewUserRepository(db)
		err := repo.Save(context.TODO(), &user)
		assert.NoError(t, err)
	})
}

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func TestUserRepository_UpdateByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock := newUserDBTest(t)
		defer db.Close()

		user := domain.User{
			FirstName: "john",
			LastName:  "due",
			Email:     "john@email.test",
		}

		expectSQL := "UPDATE users SET first_name = $1, last_name = $2, email = $3, updated_at = $4 where id = $5"
		mock.ExpectExec(expectSQL).WithArgs(user.FirstName, user.LastName, user.Email, AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		repo := repository.NewUserRepository(db)
		err := repo.UpdateByID(context.TODO(), 1, &user)
		assert.NoError(t, err)
	})
}

func TestUserRepository_DeleteByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock := newUserDBTest(t)
		defer db.Close()

		mock.ExpectExec("DELETE FROM users WHERE id = $1").WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		repo := repository.NewUserRepository(db)
		err := repo.DeleteByID(context.TODO(), 1)
		assert.NoError(t, err)
	})
}
