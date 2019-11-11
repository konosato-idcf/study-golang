package user

import (
	"context"
	"github.com/konosato-idcf/study-golang/echo-OJT-webApp/user/infra/models"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUsers_Create(t *testing.T) {
	db, err := ConnectDatabase(config)
	if err != nil {
		log.Error(err)
	}

	u := NewUser(db)
	u.Create(User{

	})

	// 確認
	ctx := context.Background()
	users, err := models.Users().All(ctx, db)
	if err != nil {
		log.Error(err)
	}
	assert.Equal(t, 1, len(users))
	user := users[0]
	assert.Equal(t, "xxx", user.Name)
	assert.Equal(t, "xxx@xxx.com", user.Email)
}