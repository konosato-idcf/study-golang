package user

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	models2 "github.com/konosato-idcf/study-golang/echo-OJT-webApp/app/user/infra/models"
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

func GetEchoContext(path string, requestMethod string, requestJson string) (*echo.Echo, echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	e.Validator = NewCustomValidator()
	req := httptest.NewRequest(requestMethod, path, strings.NewReader(requestJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath(path)
	return e, c, rec
}

// https://github.com/golang/go/wiki/TableDrivenTests
// https://qiita.com/yut-kt/items/5f9eb752f40d4d2a2e97
func TestUsersHandler_Validation(t *testing.T) {
	/*
		/ カテゴリー：バリデーションチェック
		/ サブカテゴリー：ユーザー名、メールアドレス
		/ 種別：正常
		/ 内容：-
	*/
	casesOk := []struct {
		testName string
		user     User
	}{
		{
			testName: "Nameが1文字の場合",
			user: User{
				Name:  "k",
				Email: "k@gmail.com",
			},
		},
		{
			testName: "Nameが45文字の場合",
			user: User{
				Name:  "aaaaaaaaaabbbbbbbbbbccccccccccddddddddddeeeee",
				Email: "aaaaaaaaaabbbbbbbbbbccccccccccddddddddddeeeee@idcf.jp",
			},
		},
		{
			testName: "Emailが255文字の場合",
			user: User{
				Name:  "longEmail",
				Email: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa@idcf.jp",
			},
		},
	}
	for _, tt := range casesOk {
		t.Run(tt.testName, func(t *testing.T) {
			bytes, err := json.Marshal(tt.user)
			if err != nil {
				fmt.Println(err)
				return
			}
			requestJson := string(bytes)
			_, c, _ := GetEchoContext("/users", http.MethodPost, requestJson)
			u := new(User)
			assert.NoError(t, c.Bind(u))
			assert.NoError(t, c.Validate(u))
		})
	}

	/*
		/ カテゴリー：バリデーションチェック
		/ サブカテゴリー：ユーザー名、メールアドレス
		/ 種別：異常
		/ 内容：-
	*/
	casesNg := []struct {
		testName string
		user     User
	}{
		{
			testName: "Nameが46文字の場合",
			user: User{
				Name:  "aaaaaaaaaabbbbbbbbbbccccccccccddddddddddeeeeee",
				Email: "tooLongName@idcf.jp",
			},
		},
		{
			testName: "Nameが空文字の場合",
			user: User{
				Name:  "",
				Email: "emptyName@idcf.jp",
			},
		},
		{
			testName: "Emailのフォーマットが正しくない場合",
			user: User{
				Name:  "aaa",
				Email: "aaa.jp",
			},
		},
		{
			testName: "Emailが256文字の場合",
			user: User{
				Name:  "tooLongEmail",
				Email: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa@idcf.jp",
			},
		},
		{
			testName: "Emailが空文字の場合",
			user: User{
				Name:  "emptyEmail",
				Email: "",
			},
		},
	}
	for _, tt := range casesNg {
		t.Run(tt.testName, func(t *testing.T) {
			bytes, err := json.Marshal(tt.user)
			if err != nil {
				fmt.Println(err)
				return
			}
			requestJson := string(bytes)
			_, c, _ := GetEchoContext("/users", http.MethodPost, requestJson)
			u := new(User)
			assert.NoError(t, c.Bind(u))
			assert.Error(t, c.Validate(u))
		})
	}

	/*
		/ カテゴリー：バリデーションチェック
		/ サブカテゴリー：メールアドレス
		/ 種別：異常
		/ 内容：JSON形式でないリクエスト
	*/
	caseNotJson := `{name:Joe,email:joe@idcf.jp}`
	_, c, _ := GetEchoContext("/users", http.MethodPost, caseNotJson)
	u := new(User)
	assert.Error(t, c.Bind(u))
	assert.Error(t, c.Validate(u))

	/*
		/ カテゴリー：バリデーションチェック
		/ サブカテゴリー：ID
		/ 種別：異常
		/ 内容：リクエストパラメータにIDを加えて送る
	*/
	//caseAddId :=  `{"id":10, "name":"Joe","email":"joe@idcf.jp", "hoge":"huga"}`
	//c, _ = GetEchoContext("/users", http.MethodPost, caseAddId)
	//u = new(User)
	//assert.NoError(t, c.Bind(u))
	//assert.NoError(t, c.Validate(u))
}

func TestUsersHandler_Read(t *testing.T) {
	testDb, err := ConnectDatabase(testConfig)
	if err != nil {
		panic(err)
	}
	defer deleteUser(testDb)

	/*
		/ カテゴリー：一覧取得処理
		/ サブカテゴリー：DB接続
		/ 種別：正常
		/ 内容：-
	*/
	ctx := context.Background()
	preUser := models2.User{Name: "Joe", Email: "joe@idcf.jp"}
	err = preUser.Insert(ctx, testDb, boil.Whitelist("name", "email"))
	if err != nil {
		t.Fatal(err)
	}
	_, c, rec := GetEchoContext("/users", http.MethodGet, "")

	u := NewUser(testDb)
	h := NewUsersHandler(u)

	// Assertion
	wantValue := User{
		Name:  "Joe",
		Email: "joe@idcf.jp",
	}
	wantArray := []User{wantValue}

	if assert.NoError(t, h.Index(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var responseUser []User
		err := json.Unmarshal(rec.Body.Bytes(), &responseUser)
		if err != nil {
			t.Fatal(err)
		}
		assert.GreaterOrEqual(t, responseUser[0].ID, 0)
		assert.Equal(t, wantArray[0].Name, responseUser[0].Name)
		assert.Equal(t, wantArray[0].Email, responseUser[0].Email)

		// DBのユーザーテーブルの値の確認
		ctx := context.Background()
		user, err := models2.Users(qm.Where("id=?", responseUser[0].ID)).One(ctx, testDb)
		if err != nil {
			log.Fatal(err)
		}
		assert.Equal(t, wantArray[0].Name, user.Name)
		assert.Equal(t, wantArray[0].Email, user.Email)
	}
}

// バリデーションチェックが全て通った場合、ユーザーが登録される。
func TestUsersHandler_Create(t *testing.T) {
	testDb, err := ConnectDatabase(testConfig)
	if err != nil {
		panic(err)
	}
	defer deleteUser(testDb)

	/*
		/ カテゴリー：登録処理
		/ サブカテゴリー：DB接続
		/ 種別：正常
		/ 内容：-
	*/
	requestJson := `{"name":"Joe","email":"joe@idcf.jp"}`
	_, c, rec := GetEchoContext("/users", http.MethodPost, requestJson)

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
		err = json.Unmarshal(rec.Body.Bytes(), responseUser)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(rec.Body.String())
		// assert.Equal(t, want, strings.TrimSpace(rec.Body.String()))
		assert.GreaterOrEqual(t, responseUser.ID, 0)
		assert.Equal(t, want.Name, responseUser.Name)
		assert.Equal(t, want.Email, responseUser.Email)

		// DBのユーザーテーブルの値の確認
		ctx := context.Background()
		user, err := models2.Users(qm.Where("id=?", responseUser.ID)).One(ctx, testDb)
		if err != nil {
			log.Fatal(err)
		}
		assert.Equal(t, want.Name, user.Name)
		assert.Equal(t, want.Email, user.Email)
	}

	/*
		/ カテゴリー：登録処理
		/ サブカテゴリー：メソッド
		/ 種別：異常
		/ 内容：不適切なメソッドにリクエストを送信
	*/
	//requestJsonNg := `{"name":"Joe","email":"joe@idcf.jp"}`
	////cNg, recNg := GetEchoContext("/users", http.MethodPut, requestJsonNg)
	//cNg, recNg := GetEchoContext("/users", http.MethodPut, requestJsonNg)
	//
	//u = NewUser(testDb)
	//h = NewUsersHandler(u)
	//
	////fmt.Println(recNg.Code)
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

	/*
		/ カテゴリー：更新処理
		/ サブカテゴリー：DB接続
		/ 種別：正常
		/ 内容：-
	*/
	preUser := models2.User{Name: "Joe", Email: "joe@idcf.jp"}
	err = preUser.Insert(ctx, testDb, boil.Whitelist("name", "email"))
	if err != nil {
		t.Fatal(err)
	}

	requestUser := User{
		ID:    preUser.ID,
		Name:  "Joe",
		Email: "joe2@idcf.jp",
	}
	requestJson := ToJson(t, requestUser)

	_, c, rec := GetEchoContext("/users/:id", http.MethodPut, requestJson)
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(preUser.ID))

	u := NewUser(testDb)
	h := NewUsersHandler(u)

	// Assertion
	if assert.NoError(t, h.Update(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, requestJson, strings.TrimSpace(rec.Body.String()))

		updatedUser, err := models2.Users(qm.Where("id=?", preUser.ID)).One(ctx, testDb)
		if err != nil {
			log.Fatal(err)
		}
		want := requestUser
		assert.Equal(t, want.Name, updatedUser.Name)
		assert.Equal(t, want.Email, updatedUser.Email)
	}

	/*
		/ カテゴリー：更新処理
		/ サブカテゴリー：DB登録
		/ 種別：異常
		/ 内容：存在しないIDを指定して更新処理を実行
	*/
	requestUser = User{
		ID:    9999,
		Name:  "Joe",
		Email: "joe2@idcf.jp",
	}
	requestJson = ToJson(t, requestUser)
	e, c, rec := GetEchoContext("/users/:id", http.MethodPut, requestJson)
	c.SetParamNames("id")
	c.SetParamValues("9999")

	u = NewUser(testDb)
	h = NewUsersHandler(u)

	//assert.Equal(t, http.StatusNotFound, rec.Code)
	err = h.Update(c)
	if assert.Error(t, err) {
		e.HTTPErrorHandler(err, c)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	}

	/*
		/ カテゴリー：更新処理
		/ サブカテゴリー：URL
		/ 種別：異常
		/ 内容：登録済みと情報とリクエストされた情報に差分がない
	*/
	requestUser = User{
		ID:    preUser.ID,
		Name:  "Joe",
		Email: "joe2@idcf.jp",
	}
	requestJson = ToJson(t, requestUser)

	e, c, rec = GetEchoContext("/users/:id", http.MethodPut, requestJson)
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(preUser.ID))

	u = NewUser(testDb)
	h = NewUsersHandler(u)

	// Assertion
	err = h.Update(c)
	if assert.Error(t, err) {
		e.HTTPErrorHandler(err, c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}

	/*
		/ カテゴリー：更新処理
		/ サブカテゴリー：URL
		/ 種別：異常
		/ 内容；idパラメーターが空
	*/
	//requestUser = User{
	//	ID:    preUser.ID,
	//	Name:  "Joe",
	//	Email: "joe@idcf.jp",
	//}
	//requestJson = ToJson(t, requestUser)
	//
	//e, c, rec = GetEchoContext("/users/:id", http.MethodPut, requestJson)
	//c.SetParamNames("id")
	//c.SetParamValues(strconv.Itoa(preUser.ID))
	//
	//u = NewUser(testDb)
	//h = NewUsersHandler(u)
	//
	//// Assertion
	//if assert.Error(t, h.Update(c)) {
	//	assert.Equal(t, http.StatusBadRequest, rec.Code)
	//}
}

func TestUsersHandler_Delete(t *testing.T) {
	ctx := context.Background()
	testDb, err := ConnectDatabase(testConfig)
	if err != nil {
		panic(err)
	}
	defer deleteUser(testDb)

	/*
		/ カテゴリー：削除処理
		/ サブカテゴリー：DB接続
		/ 種別：正常
		/ 内容；-
	*/
	user := models2.User{Name: "Joe", Email: "joe@idcf.jp"}
	err = user.Insert(ctx, testDb, boil.Whitelist("name", "email"))
	if err != nil {
		t.Fatal(err)
	}

	_, c, rec := GetEchoContext("/users/:id", http.MethodDelete, "")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(user.ID))

	u := NewUser(testDb)
	h := NewUsersHandler(u)

	// Assertion
	if assert.NoError(t, h.Delete(c)) {
		assert.Equal(t, http.StatusNoContent, rec.Code)

		deletedUser, err := models2.Users(qm.Where("id=?", user.ID)).One(ctx, testDb)
		if err != nil {
			if fmt.Sprint(err) != "sql: no rows in result set" {
				log.Fatal(err)
			}
		}
		assert.Empty(t, deletedUser)
	}

	/*
		/ カテゴリー：削除処理
		/ サブカテゴリー：DB接続
		/ 種別：削除
		/ 内容；DBに登録されていないIDを指定してリクエストを送る
	*/
	e, c, rec := GetEchoContext("/users/:id", http.MethodDelete, "")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(user.ID))

	u = NewUser(testDb)
	h = NewUsersHandler(u)

	// Assertion
	err = h.Delete(c)
	if assert.Error(t, err) {
		e.HTTPErrorHandler(err, c)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	}

	/*
		/ カテゴリー：削除処理
		/ サブカテゴリー：DB接続
		/ 種別：削除
		/ 内容；空のテーブルに対してDeleteのリクエストを送る
	*/
	deleteUser(testDb)

	e, c, rec = GetEchoContext("/users/:id", http.MethodDelete, "")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(user.ID))

	u = NewUser(testDb)
	h = NewUsersHandler(u)

	// Assertion
	err = h.Delete(c)
	if assert.Error(t, err) {
		e.HTTPErrorHandler(err, c)
		assert.Equal(t, http.StatusNotFound, rec.Code)
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
	_, err := models2.Users().DeleteAll(context.Background(), db)
	if err != nil {
		fmt.Println(err)
	}
}
