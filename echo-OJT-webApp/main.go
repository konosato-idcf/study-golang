package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/konosato-idcf/study-golang/echo-OJT-webApp/user"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.Validator = user.NewCustomValidator()

	// Open handle to database like normal
	db, err := user.ConnectDatabase(user.Config{
		Username: "admin",
		Password: "himitu",
		Host:     "localhost",
		Database: "sample_app",
		Port:     13306,
	})
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
