package user

import (
	"catalog-digital-product/internal/custom"
	"context"
	"database/sql"
)

type UserRepository interface {
	FindById(ctx context.Context, tx *sql.Tx, id int) (User, error)
	FindByUsername(ctx context.Context, tx *sql.Tx, username string) (User, error)
	Update(ctx context.Context, tx *sql.Tx, user User) error
}

type UserRepositoryImpl struct {
}

func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{}
}

func (u *UserRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id int) (User, error) {
	sql := "select id, name, username, password, created_at, updated_at from users where id = ?"
	row, err := tx.QueryContext(ctx, sql, id)
	if err != nil {
		return User{}, err
	}
	defer row.Close()

	user := User{}
	if row.Next() {
		err := row.Scan(&user.Id, &user.Name, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return User{}, err
		}

		return user, nil
	}

	return user, custom.ErrNotFound
}

func (u *UserRepositoryImpl) FindByUsername(ctx context.Context, tx *sql.Tx, username string) (User, error) {
	sql := "select id, name, username, password, created_at, updated_at from users where username = ?"
	row, err := tx.QueryContext(ctx, sql, username)
	if err != nil {
		return User{}, err
	}
	defer row.Close()

	user := User{}
	if row.Next() {
		err := row.Scan(&user.Id, &user.Name, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return User{}, err
		}
	}

	return user, custom.ErrNotFound
}

func (u *UserRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, user User) error {
	sql := "update users set name = ?, username = ?, password = ? where id = ?"
	_, err := tx.ExecContext(ctx, sql, user.Name, user.Username, user.Password, user.Id)
	if err != nil {
		return err
	}

	return nil
}
