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
		//fmt.Printf("%#v\n", u)
		userList = append(userList, u)
	}
	return c.JSON(http.StatusOK, userList)
}

//$ curl http://localhost:1323/users -v

func (h UsersHandler) Create(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return err
	}
	if err := c.Validate(u); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
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

func (h UsersHandler) Update(c echo.Context) error {
	// Validate
	u := new(User)
	if err := c.Bind(u); err != nil {
		return err
	}
	if err := c.Validate(u); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Retrieve data
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		message := &Message{
			Message: fmt.Sprintf("strconv.ParseInt: parsing \"a\": invalid syntax"),
		}
		return echo.NewHTTPError(http.StatusBadRequest, message)
	}
	ctx := context.Background()
	_, err = h.Users.FindById(ctx, id)
	if err != nil {
		message := &Message{
			Message: fmt.Sprintf("The user id:%d does not exist.", id),
		}
		return echo.NewHTTPError(http.StatusNotFound, message)
	}

	// Update
	rowsAff, err := h.Users.Update(ctx, u)
	if err != nil {
		return err
	}
	if rowsAff == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Differences are not found.")
	}

	return c.JSON(http.StatusOK, u)
}

// $ curl http://localhost:1323/users/1 -v -X PUT -H "Content-Type: application/json" -d '{"name":"konosato", "email":"konosato@idcf.jp"}'

func (h UsersHandler) Delete(c echo.Context) error {
	// Retrieve data
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		message := &Message{
			Message: fmt.Sprintf("ID must be numeric."),
		}
		return echo.NewHTTPError(http.StatusBadRequest, message)
	}
	ctx := context.Background()
	u, err := h.Users.FindById(ctx, id)
	if err != nil {
		message := &Message{
			Message: fmt.Sprintf("The user id:%d does not exist.", id),
		}
		return echo.NewHTTPError(http.StatusNotFound, message)
	}

	// Delete
	rowsAff, err := h.Users.Delete(ctx, u)
	if err != nil {
		return err
	}
	if rowsAff == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Delete object is not defined.")
	}
	return c.String(http.StatusNoContent, "")
}

// $ curl http://localhost:1323/users/20 -v -X DELETE
