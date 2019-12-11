package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/konosato-idcf/study-golang/echo-OJT-webApp/app/user"
	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	// Open handle to database like normal
	db, err := sql.Open("mysql", "admin:himitu@tcp(localhost:13306)/sample_app")
	if err != nil {
		e.Logger.Fatal(err.Error())
	}

	// ユーザーテーブルにアクセスするモデル
	u := user.NewUser(db)

	// ハンドラーの生成
	userHandler := user.NewUsersHandler(u)

	// ルーティング
	e.GET("/users", userHandler.Index)
	e.POST("/users", userHandler.Create)
	e.PUT("/users/:id", userHandler.Update)
	e.DELETE("/users/:id", userHandler.Delete)

	e.Logger.Fatal(e.Start(":1323"))
}
