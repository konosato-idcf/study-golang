package user

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/mock/gomock"
	"github.com/konosato-idcf/study-golang/echo-OJT-webApp/user/infra/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
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
    type user struct {
        Name  string `json:"name"`
        Email string `json:"email"`
    }

// https://github.com/golang/go/wiki/TableDrivenTests
// https://qiita.com/yut-kt/items/5f9eb752f40d4d2a2e97
func TestUsersHandler_Create_Validation(t *testing.T) {
	// 種別：正常
        casesOk := []struct {
            testName string
            user User
        }{
            {
                testName: "Nameが一文字の場合",
                user: User{ Name: "k", Email: "k@gmail.com"},
            },
//             {"aaaaaaaaaabbbbbbbbbbccccccccccddddddddddeeeee",
//                 "aaaaaaaaaabbbbbbbbbbccccccccccddddddddddeeeee@idcf.jp"},
//             {"longEmail",
//                 "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa@idcf.jp"},
        }
        for _, tt := range casesOk {
	        t.Run(tt.testName, func(t *testing.T) {
                bytes, err := json.Marshal(tt.user)
                if err != nil {
                    fmt.Println(err)
                    return
                }
                requestJson := string(bytes)
                c, _ := GetEchoContext("/users", http.MethodPost, requestJson)
                u := new(User)
                assert.NoError(t, c.Bind(u))
                assert.NoError(t, c.Validate(u))
	        })
        }

	// 種別：異常
	casesNg := []struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}{
		{"aaaaaaaaaabbbbbbbbbbccccccccccddddddddddeeeeee",
			"tooLongName@idcf.jp"},
		{"",
			"emptyName@idcf.jp"},
		{"notEmailFormat",
			"aaaidcf.jp"},
		{"tooLongEmail",
			"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa@idcf.jp"},
		{"emptyEmail",
			""},
	}
	for _, tt := range casesNg {
		bytes, err := json.Marshal(tt)
		if err != nil {
			fmt.Println(err)
			return
		}
		requestJson := string(bytes)
		c, _ := GetEchoContext("/users", http.MethodPost, requestJson)
		u := new(User)
		assert.NoError(t, c.Bind(u))
		assert.Error(t, c.Validate(u))
	}

	// 種別：異常、JSON形式でないリクエスト
	caseNotJson := `{name:Joe,email:joe@idcf.jp}`
	c, _ := GetEchoContext("/users", http.MethodPost, caseNotJson)
	u := new(User)
	assert.Error(t, c.Bind(u))
	assert.Error(t, c.Validate(u))

	// 種別：異常、リクエストパラーメータにIDを加えて送る
	//caseAddId :=  `{"id":10, "name":"Joe","email":"joe@idcf.jp", "hoge":"huga"}`
	//c, _ = GetEchoContext("/users", http.MethodPost, caseAddId)
	//u = new(User)
	//assert.NoError(t, c.Bind(u))
	//assert.NoError(t, c.Validate(u))
}

// バリデーションチェックが全て通った場合、ユーザーが登録される。
func TestUsersHandler_Create(t *testing.T) {
	testDb, err := ConnectDatabase(testConfig)
	if err != nil {
		panic(err)
	}
	defer deleteUser(testDb)

	// Case: Pass test
	requestJson := `{"name":"Joe","email":"joe@idcf.jp"}`
	c, rec := GetEchoContext("/users", http.MethodPost, requestJson)

	u := NewUser(testDb)
	h := NewUsersHandler(u)

	// Assertion
	// want := `{"id":1,"name":"Joe","email":"joe@idcf.jp"}`
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

	// Case: NG test
	//requestJsonNg := `{"name":"Joe","email":"joe@idcf.jp"}`
	//cNg, recNg := GetEchoContext("/users", http.MethodGet, requestJsonNg)
	//fmt.Println(cNg)
	//
	//u = NewUser(testDb)
	//h = NewUsersHandler(u)
	//
	//if assert.Error(t, h.Create(cNg)) {
	//	assert.Equal(t, http.StatusMethodNotAllowed, recNg.Code)
	//}
}

func ToJson(t *testing.T, i interface{}) string {
	t.Helper()
	bytes, err := json.Marshal(i)
	if err != nil {
		t.Fatal(err)
	}
	return string(bytes)
}

func TestUsersHandler_Update(t *testing.T) {
	ctx := context.Background()
	testDb, err := ConnectDatabase(testConfig)
	if err != nil {
		panic(err)
	}
	defer deleteUser(testDb)

	user := models.User{Name: "Joe", Email: "joe@idcf.jp"}
	err = user.Insert(ctx, testDb, boil.Whitelist("name", "email"))
	if err != nil {
		t.Fatal(err)
	}

	requestUser := User{
		ID: user.ID,
		Name: "Joe",
		Email: "joe2@idcf.jp",
	}
	requestString := ToJson(t, requestUser)

	// requestJson := `{"name":"Joe","email":"joe2@idcf.jp"}`
	c, rec := GetEchoContext("/users/:id", http.MethodPut, requestString)
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(user.ID))

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := NewUser(testDb)
// 	u := NewMockUsersInterface(ctrl)
// 	u.EXPECT().FindById(gomock.Any(), 2).Return(&User{
// 		ID:    2,
// 		Name:  "Joe",
// 		Email: "joe@idcf.jp",
// 	}, nil)
	h := NewUsersHandler(u)

	// Assertion
	// want := `{"id":2,"name":"Joe","email":"joe2@idcf.jp"}`
	if assert.NoError(t, h.Update(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, requestString, strings.TrimSpace(rec.Body.String()))

		updatedUser, err := models.Users(qm.Where("id=?", user.ID)).One(ctx, testDb)
		if err != nil {
			log.Fatal(err)
		}
		want := requestUser
		assert.Equal(t, want.Name, updatedUser.Name)
		assert.Equal(t, want.Email, updatedUser.Email)
	}
}

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
