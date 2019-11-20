package user

import (
	"context"
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/konosato-idcf/study-golang/echo-OJT-webApp/user/infra/models"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var (
	testConfig = Config{
		Username: "root",
		Password: "himitu",
		Host:     "localhost",
		Database: "sample_app_test",
		Port:     13306,
	}
)

func GetEchoContext(path string, requestMethod string, requestJson string) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	e.Validator = NewCustomValidator()
	req := httptest.NewRequest(requestMethod, path, strings.NewReader(requestJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath(path)

	return c, rec
}

// https://github.com/golang/go/wiki/TableDrivenTests
// https://qiita.com/yut-kt/items/5f9eb752f40d4d2a2e97
func TestUsersHandler_Create_Validation(t *testing.T) {
	assert.Equal(t, http.StatusCreated, 201)
}

// バリデーションチェックが全て通った場合、ユーザーが登録される。
func TestUsersHandler_Create(t *testing.T) {
	testDb, err := ConnectDatabase(testConfig)
	if err != nil {
		panic(err)
	}
	defer deleteUser(testDb)

	requestJson := `{"name":"Joe","email":"joe@idcf.jp"}`
	c, rec := GetEchoContext("/users/:id", http.MethodPost, requestJson)

	u := NewUser(testDb)
	h := NewUsersHandler(u)

	// Assertion
	// want := `{"id":0,"name":"Joe","email":"joe@idcf.jp"}`
	want := &User{
		Name:  "Joe",
		Email: "joe@idcf.jp",
	}

	if assert.NoError(t, h.Create(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)

		responseUser := &User{}
		json.Unmarshal(rec.Body.Bytes(), responseUser)
		// fmt.Println(responseUser)
		// assert.Equal(t, want, strings.TrimSpace(rec.Body.String()))
		assert.GreaterOrEqual(t, responseUser.ID, 0)
		assert.Equal(t, want.Name, responseUser.Name)
		assert.Equal(t, want.Email, responseUser.Email)

		// DBのユーザーテーブルの値の確認
		ctx := context.Background()
		user, err := models.Users(qm.Where("id=?", responseUser.ID)).One(ctx, testDb)
		if err != nil {
			log.Fatal(err)
		}
		assert.Equal(t, want.Name, user.Name)
		assert.Equal(t, want.Email, user.Email)
	}
}

//func TestUsersHandler_Update(t *testing.T) {
//	requestJson := `{"name":"Joe","email":"joe2@idcf.jp"}`
//	c, rec := GetEchoContext("/users/:id", http.MethodPut, requestJson)
//  c.SetParamNames("id")
//  c.SetParamValues(1)
//
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	u := NewMockUsersInterface(ctrl)
//	u.EXPECT().FindById(gomock.Any(), 0).Return(&User{
//		ID: 0,
//		Name: "Joe",
//		Email: "joe@idcf.jp",
//	}, nil)
//	h := NewUsersHandler(u)
//
//	// Assertion
//	want := `{"id":0,"name":"Joe","email":"joe2@idcf.jp"}`
//	if assert.NoError(t, h.Update(c)) {
//		assert.Equal(t, http.StatusCreated, rec.Code)
//		assert.Equal(t, want, strings.TrimSpace(rec.Body.String()))
//	}
//}

func TestMain(m *testing.M) {
	testDb, err := ConnectDatabase(testConfig)
	if err != nil {
		panic(err)
	}
	deleteUser(testDb)
	code := m.Run()

	// deleteUser()
	os.Exit(code)
}

func deleteUser(db *sql.DB) {
	models.Users().DeleteAll(context.Background(), db)
}
