package util

import (
	"time"

	"github.com/xaxys/maintainman/core/config"

	"github.com/iris-contrib/middleware/jwt"
)

var (
	key    []byte
	expire time.Duration
)

func init() {
	exp, err := time.ParseDuration(config.AppConfig.GetString("token.expire"))
	if err != nil {
		panic(err)
	}
	expire = exp
	key = []byte(config.AppConfig.GetString("token.key"))
}

func GetJwtString(id uint, name, role string) (string, error) {
	return GetJwtStringWithClaims(id, name, role, nil)
}

func GetJwtStringWithClaims(id uint, name, role string, custom map[string]any) (string, error) {
	claims := jwt.MapClaims{
		"user_id":   id,
		"user_name": name,
		"user_role": role,

		"iss": config.AppConfig.GetString("app.name"),
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(expire).Unix(),
	}
	for k, v := range custom {
		claims[k] = v
	}
	return GetRawJwtString(claims)
}

func GetRawJwtString(claims jwt.MapClaims) (string, error) {
	token := jwt.NewTokenWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
