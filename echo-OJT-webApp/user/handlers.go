package user

import (
	"context"
	"fmt"
	"github.com/konosato-idcf/study-golang/echo-OJT-webApp/user/infra/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

type UsersHandler struct {
	User UserInterface
}

func NewUsersHandler(user UserInterface) *UsersHandler {
	return &UsersHandler{User: user}
}

func (u *UsersHandler) Index(c echo.Context) error {
	ctx := context.Background()
	users, err := u.User.All(ctx)
	if err != nil {
		c.Logger().Fatal(err.Error())
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
}
//// GET JSON response
//e.GET("/users", func(c echo.Context) error {
//	ctx := context.Background()
//	users, err := models.Users().All(ctx, db)
//	if err != nil {
//		e.Logger.Fatal(err.Error())
//	}
//	userList := make([]*models.User, 0)
//	for _, v := range users {
//		//fmt.Printf("%#v\n", v)
//		u := &models.User{}
//		u.ID = v.ID
//		u.Name = v.Name
//		u.Email = v.Email
//		fmt.Printf("%#v\n", u)
//		userList = append(userList, u)
//	}
//	return c.JSON(http.StatusOK, userList)
//})
//
//// POST/Form insert user info.
///*e.POST("/users", func(c echo.Context) error {
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
//})
////$ curl -v -F "name=Yamada Hanako" -F "email=hyamada@labstack.com" http://localhost:1323/users
//*/
//
//// POST /JSON insert user info.
//e.POST("/users", func(c echo.Context) (err error) {
//	u := new(User)
//	if err := c.Bind(u); err != nil {
//		return err
//	}
//	if err = c.Validate(u); err != nil {
//		return c.JSON(http.StatusBadRequest, err.Error())
//	}
//
//	ctx := context.Background()
//	var user models.User
//	user.Name = u.Name
//	user.Email = u.Email
//	err = user.Insert(ctx, db, boil.Whitelist("name", "email"))
//	if err != nil {
//		e.Logger.Fatal(err.Error())
//	}
//	return c.JSON(http.StatusCreated, user)
//})
//// $ curl http://localhost:1323/users -v -X POST -H "Content-Type: application/json" -d '{"name":"aaa", "email":"aaa@idcf.jp"}'
//// $ curl http://localhost:1323/users -v -X POST -H "Content-Type: application/json" -d '{"name":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "email":"kotaro@idcf.jp"}'
//// -> "Key: 'User.Name' Error:Field validation for 'Name' failed on the 'lt' tag"
//
//// PUT /users/:id Update user.
//e.PUT("/users/:id", func(c echo.Context) (err error) {
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
//})
//// $ curl http://localhost:1323/users/1 -v -X PUT -H "Content-Type: application/json" -d '{"name":"konosato", "email":"konosato@idcf.jp"}'
//
//// DELETE /users/:id Delete user.
//e.DELETE("/users/:id", func(c echo.Context) (err error) {
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
