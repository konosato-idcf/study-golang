package user

import (
	"context"
	"database/sql"
	"github.com/konosato-idcf/study-golang/echo-OJT-webApp/user/infra/models"
	"github.com/volatiletech/sqlboiler/boil"
)

type UsersInterface interface {
	All(context.Context) (models.UserSlice, error)
	Create(context.Context, *User) (*User, error)
}

func NewUser(db *sql.DB) *Users {
	return &Users{db : db}
}

type Users struct {
	db *sql.DB `json:"-"`
}

type User struct {
	ID    int    `boil:"id" json:"id" toml:"id" yaml:"id"`
	Name  string `boil:"name" json:"name" toml:"name" yaml:"name" validate:"required,gt=0,lt=45"`
	Email string `boil:"email" json:"email" toml:"email" yaml:"email" validate:"required,email"`
}

func (u Users) All(ctx context.Context) (models.UserSlice, error) {
	return models.Users().All(ctx, u.db)
}

func (u Users) Create(ctx context.Context, user *User) (*User, error) {
	var v models.User
	v.Name = user.Name
	v.Email = user.Email
	err := v.Insert(ctx, u.db, boil.Whitelist("name", "email"))
	if err != nil {
		return nil, err
	}
	user.ID = v.ID
	return user, nil
}
