package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID        int       `db:"id" fieldtag:"pk" json:"id"`
	Login     string    `db:"login" json:"login"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	Role      string    `db:"role" json:"role"`
	Password  string    `db:"password" json:"password,omitempty"`
	Active    bool      `db:"active" json:"active,omitempty"`
}

var UserStruct = sqlbuilder.NewStruct(new(User))

type IUserRepo interface {
	Create(ctx context.Context, user *User) *User
	BlockByID(ctx context.Context, userID string) bool
	GetByID(ctx context.Context, userID int) *User
	GetByLogin(ctx context.Context, userLogin string) (*User, error)
	GetAll(ctx context.Context) []User
}

type UserRepo struct {
	dbPool *pgxpool.Pool
}

func NewUserRepo(dbPool *pgxpool.Pool) *UserRepo {
	return &UserRepo{dbPool: dbPool}
}

func (u *UserRepo) Create(ctx context.Context, user *User) *User {
	ID, err := u.generateNextUserID(ctx)
	if err != nil {
		return nil
	}
	user.ID = ID
	user.Active = true
	user.CreatedAt = time.Now()

	sql, args := UserStruct.InsertInto("users", user).
		BuildWithFlavor(sqlbuilder.PostgreSQL) // todo

	row := u.dbPool.QueryRow(ctx, sql, args...)
	rowScanErr := row.Scan()
	if rowScanErr != nil && !errors.Is(rowScanErr, pgx.ErrNoRows) {
		return nil
	}

	return user
}

func (u *UserRepo) BlockByID(ctx context.Context, userID string) bool {

	ub := sqlbuilder.Update("users")
	sql, args := ub.Where(ub.Equal("id", userID)).
		Set(ub.Assign("active", false)).
		BuildWithFlavor(sqlbuilder.PostgreSQL)

	_, err := u.dbPool.Exec(ctx, sql, args...)
	if err != nil {
		return false
	}

	return true
}

func (u *UserRepo) GetByID(ctx context.Context, userID int) *User {
	sb := UserStruct.SelectFrom("users")
	sql, args := sb.Where(sb.Equal("id", userID)).
		BuildWithFlavor(sqlbuilder.PostgreSQL)
	row := u.dbPool.QueryRow(ctx, sql, args...)

	var user User
	rowScanErr := row.Scan(UserStruct.Addr(&user)...)
	if rowScanErr != nil && errors.Is(rowScanErr, pgx.ErrNoRows) {
		return nil
	}

	return &user
}

func (u *UserRepo) GetByLogin(ctx context.Context, userLogin string) (*User, error) {
	sb := UserStruct.SelectFrom("users")
	sql, args := sb.Where(sb.Equal("login", userLogin)).
		BuildWithFlavor(sqlbuilder.PostgreSQL)
	row := u.dbPool.QueryRow(ctx, sql, args...)

	var user User
	if err := row.Scan(UserStruct.Addr(&user)...); err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserRepo) GetAll(ctx context.Context) []User {
	sql, _ := UserStruct.SelectFrom("users").
		OrderBy("id").
		BuildWithFlavor(sqlbuilder.PostgreSQL)

	rows, err := u.dbPool.Query(ctx, sql)
	if err != nil {
		return nil
	}
	defer rows.Close()

	res := make([]User, 0)
	for rows.Next() {
		var user User
		rowScanErr := rows.Scan(UserStruct.Addr(&user)...)
		if rowScanErr != nil {
			return nil
		}
		res = append(res, user)
	}

	return res
}

func (u *UserRepo) generateNextUserID(ctx context.Context) (int, error) {
	rows, err := u.dbPool.Query(ctx, fmt.Sprintf("SELECT nextval('%s')", "users_sequence"))
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	if rows.Next() {
		var id int
		rowScanErr := rows.Scan(&id)
		if rowScanErr != nil {
			return 0, rowScanErr
		}
		return id, nil
	}
	return 0, fmt.Errorf("something was wrong. there is no next user id")
}
