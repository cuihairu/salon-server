package utils

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type JWT struct {
	SigningKey []byte
}

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

func NewJWT(key string) *JWT {
	return &JWT{
		SigningKey: []byte(key),
	}
}

func (j *JWT) GenerateToken(userID uint, expire time.Duration) (string, error) {
	nowTime := time.Now()
	if expire <= 5*time.Second {
		expire = 7 * 24 * time.Hour
	}
	expireTime := nowTime.Add(expire)

	claims := Claims{
		UserID: userID,
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

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
