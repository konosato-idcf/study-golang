package user

import (
	"context"
	"database/sql"
	"github.com/konosato-idcf/study-golang/echo-OJT-webApp/user/infra/models"
)

type UserInterface interface {
	All(context.Context) (models.UserSlice, error)
}

func NewUser(db *sql.DB) *User {
	return &User{db : db}
}

type User struct {
	ID    int    `boil:"id" json:"id" toml:"id" yaml:"id"`
	Name  string `boil:"name" json:"name" toml:"name" yaml:"name" validate:"required,gt=0,lt=45"`
	Email string `boil:"email" json:"email" toml:"email" yaml:"email" validate:"required,email"`
	db *sql.DB `json:"-"`
}

func (u User) All(ctx context.Context) (models.UserSlice, error) {
	return models.Users().All(ctx, u.db)
}