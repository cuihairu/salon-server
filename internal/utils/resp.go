package utils

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ErrorResponse struct {
	ErrorCode    int    `json:"errorCode"`    // 业务约定的错误码
	ErrorMessage string `json:"errorMessage"` // 业务约定的错误信息
	Success      bool   `json:"success"`      // 业务上是否请求成功
}

type RespList struct {
	Data    any  `json:"data"`    // 正常的数据
	Total   int  `json:"total"`   // 总数 只有让 data 为 List 生效
	Success bool `json:"success"` // 业务上是否请求成功
}

type Context struct {
	*gin.Context
	claims  *Claims
	session sessions.Session
	token   *string
}

func NewContext(c *gin.Context) *Context {
	return &Context{
		Context: c,
		session: sessions.Default(c),
	}
}

func (c *Context) Session() sessions.Session {
	return c.session
}

func (c *Context) ClearSession() {
	session := sessions.Default(c.Context)
	session.Clear()
}

func (c *Context) SetSession(key string, value any) error {
	c.Session().Set(key, value)
	return c.session.Save()
}

func (c *Context) GetSession(key string) any {
	return c.session.Get(key)
}

func (c *Context) SetToken(token string) error {
	c.token = &token
	err := c.SetSession("token", token)
	if err != nil {
		c.ServerError(err)
	}
	SetHeaderToken(c.Context, token)
	return nil
}

func (c *Context) GetToken() string {
	if c.token != nil {
		return *c.token
	}
	tokenObj := c.session.Get("token")
	if tokenObj == nil {
		return ""
	}
	token, ok := tokenObj.(string)
	if !ok {
		token = ""
	}
	c.token = &token
	return token
}

func (c *Context) loadClaims() {
	if c.claims != nil {
		return
	}
	claims, ok := MustGetClaimsFormContext(c.Context)
	if !ok {
		c.claims = &Claims{}
	} else {
		c.claims = claims
	}
}

func (c *Context) Claims() (*Claims, bool) {
	c.loadClaims()
	return c.claims, c.claims.UserID != 0
}

func (c *Context) Id() uint {
	c.loadClaims()
	return c.claims.UserID
}

func (c *Context) Role() string {
	c.loadClaims()
	return c.claims.Role
}

func (c *Context) Success(data any) {
	c.JSON(http.StatusOK, data)
}

func (c *Context) OK() {
	c.JSON(http.StatusOK, gin.H{})
}

func (c *Context) ReturnList(l any, total int) {
	c.JSON(http.StatusOK, RespList{Data: l, Total: total, Success: true})
}

func (c *Context) Fail(errorCode int, errorMessage string) {
	c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
		ErrorCode:    errorCode,
		ErrorMessage: errorMessage,
		Success:      true,
	})
}

func (c *Context) Error(err error) {
	c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
		ErrorCode:    http.StatusUnprocessableEntity,
		ErrorMessage: err.Error(),
		Success:      false,
	})
}

func (c *Context) BadRequest(err error) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		ErrorCode:    http.StatusBadRequest,
		ErrorMessage: err.Error(),
		Success:      false,
	})
}

func (c *Context) NotFound(err error) {
	c.JSON(http.StatusNotFound, ErrorResponse{
		ErrorCode:    http.StatusNotFound,
		ErrorMessage: err.Error(),
		Success:      false,
	})
}

func (c *Context) Unauthorized() {
	c.JSON(http.StatusUnauthorized, ErrorResponse{
		ErrorCode:    http.StatusUnauthorized,
		ErrorMessage: "unauthorized",
		Success:      false,
	})
}

func (c *Context) Forbidden() {
	c.JSON(http.StatusForbidden, ErrorResponse{
		ErrorCode:    http.StatusForbidden,
		ErrorMessage: "forbidden",
		Success:      false,
	})
}

func (c *Context) ServerError(err error) {
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		ErrorCode:    http.StatusInternalServerError,
		ErrorMessage: err.Error(),
		Success:      false,
	})
}
