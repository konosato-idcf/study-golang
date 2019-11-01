package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/konosato-idcf/study-golang/echo-OJT-webApp/models"

	//"github.com/labstack/echo"

	// Import this so we don't have to use qm.Limit etc.
	//. "github.com/volatiletech/sqlboiler/queries/qm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"net/http"
)

func main() {
	e := echo.New()
	// Open handle to database like normal
	db, err := sql.Open("mysql", "admin:himitu@tcp(localhost:13306)/sample_app")
	if err != nil {
		e.Logger.Fatal(err.Error())
	}
	e.GET("/", func(c echo.Context) error {
		ctx := context.Background()
		users, err := models.Users().All(ctx, db)
		if err != nil {
			e.Logger.Fatal(err.Error())
		}
		for i, v := range users {
			fmt.Println(i, v)
			fmt.Printf("%#v", v)
		}
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Logger.Fatal(e.Start(":1323"))
}
