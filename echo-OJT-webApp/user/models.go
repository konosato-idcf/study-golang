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
	FindById(context.Context, int) (*User, error)
	Update(context.Context, *User) (int64, error)
	Delete(context.Context, *User) (int64, error)
}

func NewUser(db *sql.DB) *Users {
	return &Users{db: db}
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

func (u Users) FindById(ctx context.Context, id int) (*User, error) {
	v, err := models.FindUser(ctx, u.db, id)
	if err != nil {
		return nil, err
	}
	//var user *User
	//user := &User{}
	user := new(User)
	user.ID = v.ID
	user.Name = v.Name
	user.Email = v.Email
	return user, nil
}

func (u Users) Update(ctx context.Context, user *User) (int64, error) {
	var v models.User
	v.ID = user.ID
	v.Name = user.Name
	v.Email = user.Email
	rowsAff, err := v.Update(ctx, u.db, boil.Infer())
	if err != nil {
		return rowsAff, err
	}
	return rowsAff, nil
}

func (u Users) Delete(ctx context.Context, user *User) (int64, error) {
	var v models.User
	v.ID = user.ID
	v.Name = user.Name
	v.Email = user.Email
	rowsAff, err := v.Delete(ctx, u.db)
	if err != nil {
		return rowsAff, err
	}
	return rowsAff, nil
}
