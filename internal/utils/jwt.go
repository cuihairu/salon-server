package utils

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"sync"
	"time"
)

type JWT struct {
	SigningKey []byte
	expire     time.Duration
}

type Claims struct {
	UserID uint   `json:"user_id"`
	Group  string `json:"group"`
	jwt.StandardClaims
}

func (c *Claims) SessionKey() string {
	return fmt.Sprintf("%s:%d", c.Group, c.UserID)
}

var one sync.Once

func NewJWT(key string, expire time.Duration) *JWT {
	var j *JWT
	one.Do(func() {
		j = &JWT{
			SigningKey: []byte(key),
			expire:     expire,
		}
	})
	return j
}

func (j *JWT) GenerateToken(userID uint) (string, error) {
	return j.GenerateTokenWithGroup(userID, "user")
}

func (j *JWT) GenerateTokenWithGroup(userID uint, group string) (string, error) {
	return j.GenerateTokenWithExpire(userID, group, j.expire)
}

func (j *JWT) GenerateTokenWithExpire(userID uint, group string, expire time.Duration) (string, error) {
	nowTime := time.Now()
	if expire <= 5*time.Second {
		expire = 7 * 24 * time.Hour
	}
	expireTime := nowTime.Add(expire)

	claims := Claims{
		UserID: userID,
		Group:  group,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			IssuedAt:  nowTime.Unix(),
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenClaims.SignedString(j.SigningKey)
}

func (j *JWT) ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return nil, err
	}
	if tokenClaims == nil || !tokenClaims.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return tokenClaims.Claims.(*Claims), nil
}
