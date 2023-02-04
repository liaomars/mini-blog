package token

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v4"
	"sync"
	"time"
)

type Config struct {
	key         string
	identityKey string
}

// ErrMissingHeader 表示 `Authorization` 请求头为空.
var ErrMissingHeader = errors.New("the length of the `Authorization` header is zero")

var config = Config{key: "Rtg8BPKNEf2mB4mgvKONGPZZQSaJWNLijxR42qRgq0iBb5", identityKey: "identityKey"}
var once sync.Once

// Init 设置包级别的配置 config, config 会用于本包后面的 token 签发和解析.
func Init(key string, identityKey string) {
	once.Do(func() {
		if key != "" {
			config.key = key
		}
		if identityKey != "" {
			config.identityKey = identityKey
		}
	})
}

// Parse 使用指定的密钥 key 解析 token，解析成功返回 token 上下文，否则报错.
func Parse(tokenString string, key string) (string, error) {
	// 解析 token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(key), nil
	})

	if err != nil {
		return "", err
	}

	var identityKey string

	// 如果解析成功，从 token 中取出 token 的一些信息
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		identityKey = claims[config.identityKey].(string)
	}
	return identityKey, nil
}

// ParseRequest 解析请求带过来的token
func ParseRequest(c *gin.Context) (string, error) {
	// 从请求头部获取携带token字段
	header := c.Request.Header.Get("Authorization")

	if len(header) == 0 {
		return "", ErrMissingHeader
	}

	var t string

	fmt.Sscanf(header, "Bearer %s", &t)
	return Parse(t, config.key)
}

// Sign 用key生成一个token
func Sign(identityKey string) (tokenString string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		config.identityKey: identityKey,
		"nbf":              time.Now().Unix(),
		"iat":              time.Now().Unix(),
		"exp":              time.Now().Add(100000 * time.Hour).Unix(),
	})
	tokenString, err = token.SignedString([]byte(config.key))
	return
}
