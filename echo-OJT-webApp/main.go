package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/konosato-idcf/study-golang/echo-OJT-webApp/models"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/boil"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"strconv"
)

type User struct {
	ID    int    `boil:"id" json:"id" toml:"id" yaml:"id"`
	Name  string `boil:"name" json:"name" toml:"name" yaml:"name" validate:"required,gt=0,lt=45"`
	Email string `boil:"email" json:"email" toml:"email" yaml:"email" validate:"required,email"`
}

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

	// GET JSON response
	e.GET("/users", func(c echo.Context) error {
		ctx := context.Background()
		users, err := models.Users().All(ctx, db)
		if err != nil {
			e.Logger.Fatal(err.Error())
		}
		userList := make([]*models.User, 0)
		for _, v := range users {
			//fmt.Printf("%#v\n", v)
			u := &models.User{}
			u.ID = v.ID
			u.Name = v.Name
			u.Email = v.Email
			fmt.Printf("%#v\n", u)
			userList = append(userList, u)
		}
		return c.JSON(http.StatusOK, userList)
	})

	// POST/Form insert user info.
	/*e.POST("/users", func(c echo.Context) error {
		ctx := context.Background()
		var p models.User
		p.Name = c.FormValue("name")
		p.Email = c.FormValue("email")
		err := p.Insert(ctx, db, boil.Whitelist("name", "email"))
		if err != nil {
			e.Logger.Fatal(err.Error())
		}

		return c.String(http.StatusCreated, "name:" + p.Name + ", email:" + p.Email)
	})
	//$ curl -v -F "name=Yamada Hanako" -F "email=hyamada@labstack.com" http://localhost:1323/users
	*/

	// POST /JSON insert user info.
	e.POST("/users", func(c echo.Context) (err error) {
		u := new(User)
		if err := c.Bind(u); err != nil {
			return err
		}
		if err = c.Validate(u); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		ctx := context.Background()
		var user models.User
		user.Name = u.Name
		user.Email = u.Email
		err = user.Insert(ctx, db, boil.Whitelist("name", "email"))
		if err != nil {
			e.Logger.Fatal(err.Error())
		}
		return c.JSON(http.StatusCreated, user)
	})
	// $ curl http://localhost:1323/users -v -X POST -H "Content-Type: application/json" -d '{"name":"aaa", "email":"aaa@idcf.jp"}'
	// $ curl http://localhost:1323/users -v -X POST -H "Content-Type: application/json" -d '{"name":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "email":"kotaro@idcf.jp"}'
	// -> "Key: 'User.Name' Error:Field validation for 'Name' failed on the 'lt' tag"

	// PUT /users/:id Update user.
	e.PUT("/users/:id", func(c echo.Context) (err error) {
		// Validate
		u := new(User)
		if err := c.Bind(u); err != nil {
			return err
		}
		if err = c.Validate(u); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		// Retrieve data
		var id int
		id, _ = strconv.Atoi(c.Param("id"))
		ctx := context.Background()
		findUser, err := models.FindUser(ctx, db, id)
		if err != nil {
			message := &Message{
				Message: fmt.Sprintf("The user id:%d does not exist.", id),
			}
			return c.JSON(http.StatusNotFound, message)
		}

		// Update
		findUser.Name = u.Name
		findUser.Email = u.Email
		rowsAff, err := findUser.Update(ctx, db, boil.Infer())
		if rowsAff == 0 {
			return c.String(http.StatusTeapot, "Differences are not found.")
		}

		return c.JSON(http.StatusOK, findUser)
	})
	// $ curl http://localhost:1323/users/1 -v -X PUT -H "Content-Type: application/json" -d '{"name":"konosato", "email":"konosato@idcf.jp"}'

	// DELETE /users/:id Delete user.
	e.DELETE("/users/:id", func(c echo.Context) (err error) {
		u := new(User)
		if err := c.Bind(u); err != nil {
			return err
		}

		// Retrieve data
		var id int
		id, _ = strconv.Atoi(c.Param("id"))
		ctx := context.Background()
		findUser, err := models.FindUser(ctx, db, id)
		if err != nil {
			message := &Message{
				Message: fmt.Sprintf("The user id:%d does not exist.", id),
			}
			return c.JSON(http.StatusNotFound, message)
		}

		// Delete
		_, err = findUser.Delete(ctx, db)
		if err != nil {
			return err
		}
		return c.String(http.StatusNoContent, "")
	})
	// $ curl http://localhost:1323/users/20 -v -X DELETE

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
