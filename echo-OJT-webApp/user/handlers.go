package user

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type Message struct {
	Message string `json:"message"`
}

type UsersHandler struct {
	Users UsersInterface
}

func NewUsersHandler(user UsersInterface) *UsersHandler {
	return &UsersHandler{Users: user}
}

// GET("/users", userHandler.Index)
func (h *UsersHandler) Index(c echo.Context) error {
	ctx := context.Background()
	users, err := h.Users.All(ctx)
	if err != nil {
		c.Logger().Fatal(err.Error())
	}
	userList := make([]*User, 0)
	for _, v := range users {
		//fmt.Printf("%#v\n", v)
		u := &User{}
		//u := &models.User{}
		u.ID = v.ID
		u.Name = v.Name
		u.Email = v.Email
		fmt.Printf("%#v\n", u)
		userList = append(userList, u)
	}
	return c.JSON(http.StatusOK, userList)
}

// POST("/users", userHandler.Create) Catch parameter through form
//func (u UsersHandler) Create (c echo.Context) error {
//	ctx := context.Background()
//	var p models.User
//	p.Name = c.FormValue("name")
//	p.Email = c.FormValue("email")
//	err := p.Insert(ctx, db, boil.Whitelist("name", "email"))
//	if err != nil {
//		e.Logger.Fatal(err.Error())
//	}
//
//	return c.String(http.StatusCreated, "name:" + p.Name + ", email:" + p.Email)
//}
//$ curl http://localhost:1323/users -v -F "name=Yamada Hanako" -F "email=hyamada@labstack.com"

// POST("/users", userHandler.Create)
func (h UsersHandler) Create(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return err
	}
	if err := c.Validate(u); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := context.Background()
	u, err := h.Users.Create(ctx, u)
	if err != nil {
		c.Logger().Fatal(err.Error())
	}
	return c.JSON(http.StatusCreated, u)
}

// $ curl http://localhost:1323/users -v -X POST -H "Content-Type: application/json" -d '{"name":"aaa", "email":"aaa@idcf.jp"}'
// $ curl http://localhost:1323/users -v -X POST -H "Content-Type: application/json" -d '{"name":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "email":"kotaro@idcf.jp"}'
// -> "Key: 'User.Name' Error:Field validation for 'Name' failed on the 'lt' tag"

// PUT("/users/:id", userHandler.Update)
func (h UsersHandler) Update(c echo.Context) error {
	// Validate
	u := new(User)
	if err := c.Bind(u); err != nil {
		return err
	}
	if err := c.Validate(u); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// Retrieve data
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}
	ctx := context.Background()
	_, err = h.Users.FindById(ctx, id)
	if err != nil {
		message := &Message{
			Message: fmt.Sprintf("The user id:%d does not exist.", id),
		}
		return c.JSON(http.StatusNotFound, message)
	}

	// Update
	rowsAff, err := h.Users.Update(ctx, u)
	if err != nil {
		return err
	}
	if rowsAff == 0 {
		return c.String(http.StatusBadRequest, "Differences are not found.")
	}

	return c.JSON(http.StatusOK, u)
}

// $ curl http://localhost:1323/users/1 -v -X PUT -H "Content-Type: application/json" -d '{"name":"konosato", "email":"konosato@idcf.jp"}'

// DELETE("/users/:id", userHandler.Delete)
func (h UsersHandler) Delete(c echo.Context) error {
	// Retrieve data
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}
	ctx := context.Background()
	u, err := h.Users.FindById(ctx, id)
	if err != nil {
		message := &Message{
			Message: fmt.Sprintf("The user id:%d does not exist.", id),
		}
		return c.JSON(http.StatusNotFound, message)
	}

	// Delete
	rowsAff, err := h.Users.Delete(ctx, u)
	if err != nil {
		return err
	}
	if rowsAff == 0 {
		return c.String(http.StatusBadRequest, "Delete object is not defined.")
	}
	return c.String(http.StatusNoContent, "")
}

// $ curl http://localhost:1323/users/20 -v -X DELETE
