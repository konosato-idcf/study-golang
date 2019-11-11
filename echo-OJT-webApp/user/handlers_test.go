package user

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/mock/gomock"
	"github.com/konosato-idcf/study-golang/echo-OJT-webApp/user/infra/models"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var config = Config {
	Username: "root",
	Password: "himitu",
	Host:     "loaclchost",
	Database: "sample_app_test",
	Port:     13306,
}


func GetEchoContext(path string, requestMethod string, requestJson string) (echo.Context, *httptest.ResponseRecorder){
	e := echo.New()
	e.Validator = NewCustomValidator()
	req := httptest.NewRequest(requestMethod, path, strings.NewReader(requestJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath(path)

	return c, rec
}

// バリデーションチェックが全て通った場合、ユーザーが登録される。
func TestUsersHandler_Create(t *testing.T) {
	requestJson := `{"name":"Joe","email":"joe@idcf.jp"}`
	c, rec := GetEchoContext("/user", http.MethodPost, requestJson)

	db, err := ConnectDatabase(config)
	if err != nil {
		log.Fatal(err.Error())
	}

	//ctrl := gomock.NewController(t)
	//defer ctrl.Finish()

	//u := NewMockUsersInterface(ctrl)
	//u.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&User{
	//	ID: 0,
	//	Name: "Joe",
	//	Email: "joe@idcf.jp",
	//}, nil)
	u := NewUser(db)
	h := NewUsersHandler(u)

	// Assertion
	// want := `{"id":0,"name":"Joe","email":"joe@idcf.jp"}`
	if assert.NoError(t, h.Create(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		// assert.Equal(t, want, strings.TrimSpace(rec.Body.String()))
	}
}

func TestUsersHandler_Update(t *testing.T) {
	requestJson := `{"name":"Joe","email":"joe2@idcf.jp"}`
	c, rec := GetEchoContext("/user", http.MethodPut, requestJson)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := NewMockUsersInterface(ctrl)
	u.EXPECT().FindById(gomock.Any(), 0).Return(&User{
		ID: 0,
		Name: "Joe",
		Email: "joe@idcf.jp",
	}, nil)
	h := NewUsersHandler(u)

	// Assertion
	want := `{"id":0,"name":"Joe","email":"joe2@idcf.jp"}`
	if assert.NoError(t, h.Update(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, want, strings.TrimSpace(rec.Body.String()))
	}
}


func TestMain(m *testing.M) {
	db, err := ConnectDatabase(config)
	if err != nil {
		panic(err)
	}
	deleteUser(db)
	code := m.Run()

	// deleteUser()
	os.Exit(code)
}

func deleteUser(db *sql.DB) {
	models.Users().DeleteAll(context.Background(), db)
}
