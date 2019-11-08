package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/konosato-idcf/study-golang/echo-OJT-webApp/user"
	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
)

type Message struct {
	Message string `json:"message"`
}

type CustomValidator struct {
	validator *validator.Validate
}

func main() {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	// Open handle to database like normal
	db, err := sql.Open("mysql", "admin:himitu@tcp(localhost:13306)/sample_app")
	if err != nil {
		e.Logger.Fatal(err.Error())
	}

	// GET "/" Return "Hello, World!"
	//e.GET("/", helloWorld)

	// ユーザーテーブルにアクセスするモデル
	u := user.NewUser(db)

	// ハンドラーの生成
	userHandler := user.NewUsersHandler(u)

	e.GET("/users", userHandler.Index)
	//e.POST("/users", userHandler.Create)
	//e.PUT("/users/:id", userHandler.Update)
	//e.DELETE("/users/:id", userHandler.Delete)

	e.Logger.Fatal(e.Start(":1323"))
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

/*func helloWorld(c echo.Context, e echo.Echo, db *sql.DB, ctx context.Context) error {
	users, err := models.Users().All(ctx, db)
	if err != nil {
		e.Logger.Fatal(err.Error())
	}
	for i, v := range users {
		fmt.Println(i, v)
		fmt.Printf("%#v", v)
	}
	return c.String(http.StatusOK, "Hello, World!")
}*/
