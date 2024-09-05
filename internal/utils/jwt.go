package utils

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"time"
)

const (
	ClaimsKey = "claims"
)

type JWT struct {
	SigningKey []byte
	expire     time.Duration
}

type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

func GetClaimsFormContext(c *gin.Context) (*Claims, bool) {
	claims, ok := c.Get(ClaimsKey)
	if !ok {
		return nil, false
	}
	cl, ok := claims.(*Claims)
	if !ok {
		return nil, false
	}
	return cl, true
}

func MustGetClaimsFormContext(c *gin.Context) (*Claims, bool) {
	claims, ok := GetClaimsFormContext(c)
	if !ok || claims == nil || claims.UserID == 0 || claims.Role == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"errorMessage": "unauthorized", "errorCode": http.StatusUnauthorized})
		return nil, false
	}
	return claims, ok
}

func SetClaimsToContext(c *gin.Context, claims *Claims) {
	c.Set(ClaimsKey, claims)
}

func (c *Claims) SessionKey() string {
	return fmt.Sprintf("%s:%d", c.Role, c.UserID)
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

func (j *JWT) GenerateTokenWithGroup(userID uint, role string) (string, error) {
	return j.GenerateTokenWithExpire(userID, role, j.expire)
}

func (j *JWT) GenerateTokenWithExpire(userID uint, role string, expire time.Duration) (string, error) {
	nowTime := time.Now()
	if expire <= 5*time.Second {
		expire = 7 * 24 * time.Hour
	}
	expireTime := nowTime.Add(expire)

	claims := Claims{
		UserID: userID,
		Role:   role,
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
	claims, ok := tokenClaims.Claims.(*Claims)
	if claims == nil || !ok || claims.UserID == 0 || claims.Role == "" {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}
