package util

import (
	"maintainman/config"
	"time"

	"github.com/iris-contrib/middleware/jwt"
)

var (
	key    []byte
	expire time.Duration
)

func init() {
	exp, err := time.ParseDuration(config.AppConfig.GetString("token.exp"))
	if err != nil {
		panic(err)
	}
	expire = exp
	key = []byte(config.AppConfig.GetString("token.key"))
}

func GetJwtString(uid uint, role string) (string, error) {
	token := jwt.NewTokenWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   uid,
		"user_role": role,

		"iss": config.AppConfig.GetString("app.name"),
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(expire).Unix(),
	})

	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
