package user

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)


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

//
//// POST("/users", userHandler.Create) Form
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
////$ curl -v -F "name=Yamada Hanako" -F "email=hyamada@labstack.com" http://localhost:1323/users
//
//
// POST("/users", userHandler.Create)
func (h UsersHandler) Create (c echo.Context) error {
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
//// $ curl http://localhost:1323/users -v -X POST -H "Content-Type: application/json" -d '{"name":"aaa", "email":"aaa@idcf.jp"}'
//// $ curl http://localhost:1323/users -v -X POST -H "Content-Type: application/json" -d '{"name":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "email":"kotaro@idcf.jp"}'
//// -> "Key: 'User.Name' Error:Field validation for 'Name' failed on the 'lt' tag"
//
//
//
//// PUT("/users/:id", userHandler.Update)
//func (u UsersHandler) Update (c echo.Context) error {
//	// Validate
//	u := new(User)
//	if err := c.Bind(u); err != nil {
//		return err
//	}
//	if err = c.Validate(u); err != nil {
//		return c.JSON(http.StatusBadRequest, err.Error())
//	}
//
//	// Retrieve data
//	var id int
//	id, _ = strconv.Atoi(c.Param("id"))
//	ctx := context.Background()
//	findUser, err := models.FindUser(ctx, db, id)
//	if err != nil {
//		message := &Message{
//			Message: fmt.Sprintf("The user id:%d does not exist.", id),
//		}
//		return c.JSON(http.StatusNotFound, message)
//	}
//
//	// Update
//	findUser.Name = u.Name
//	findUser.Email = u.Email
//	rowsAff, err := findUser.Update(ctx, db, boil.Infer())
//	if rowsAff == 0 {
//		return c.String(http.StatusTeapot, "Differences are not found.")
//	}
//
//	return c.JSON(http.StatusOK, findUser)
//}
//// $ curl http://localhost:1323/users/1 -v -X PUT -H "Content-Type: application/json" -d '{"name":"konosato", "email":"konosato@idcf.jp"}'
//
//// DELETE("/users/:id", userHandler.Delete)
//func (u UsersHandler) Delete (c echo.Context) error {
//	u := new(User)
//	if err := c.Bind(u); err != nil {
//		return err
//	}
//
//	// Retrieve data
//	var id int
//	id, _ = strconv.Atoi(c.Param("id"))
//	ctx := context.Background()
//	findUser, err := models.FindUser(ctx, db, id)
//	if err != nil {
//		message := &Message{
//			Message: fmt.Sprintf("The user id:%d does not exist.", id),
//		}
//		return c.JSON(http.StatusNotFound, message)
//	}
//
//	// Delete
//	_, err = findUser.Delete(ctx, db)
//	if err != nil {
//		return err
//	}
//	return c.String(http.StatusNoContent, "")
//})
//// $ curl http://localhost:1323/users/20 -v -X DELETE
