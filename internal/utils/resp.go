package utils

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Resp struct {
	ErrorCode    int    `json:"errorCode"`    // 业务约定的错误码
	ErrorMessage string `json:"errorMessage"` // 业务约定的错误信息
	Success      bool   `json:"success"`      // 业务上是否请求成功
	Data         any    `json:"data"`         // 正常的数据
	Total        int    `json:"total"`        // 总数 只有让 data 为 List 生效
}

type Context struct {
	*gin.Context
	claims  *Claims
	session sessions.Session
}

func NewContext(c *gin.Context) *Context {
	return &Context{
		Context: c,
	}
}

func (c *Context) Session() sessions.Session {
	if c.session != nil {
		return c.session
	}
	session := sessions.Default(c.Context)
	c.session = session
	return c.session
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

func (c *Context) IsAdmin() bool {
	c.loadClaims()
	return c.claims.IsAdmin()
}

func (c *Context) Success(data any) {
	c.JSON(http.StatusOK, Resp{
		ErrorCode:    0,
		ErrorMessage: "",
		Success:      true,
		Data:         data,
		Total:        0,
	})
}

func (c *Context) OK() {
	c.JSON(http.StatusOK, Resp{
		ErrorCode:    0,
		ErrorMessage: "",
		Success:      true,
		Data:         gin.H{},
		Total:        0,
	})
}

func (c *Context) ReturnList(l any, total int) {
	c.JSON(http.StatusOK, Resp{
		ErrorCode:    0,
		ErrorMessage: "",
		Success:      true,
		Data:         l,
		Total:        total,
	})
}

func (c *Context) Fail(errorCode int, errorMessage string) {
	c.JSON(http.StatusUnprocessableEntity, Resp{
		ErrorCode:    errorCode,
		ErrorMessage: errorMessage,
		Success:      true,
		Data:         nil,
		Total:        0,
	})
}

func (c *Context) Error(err error) {
	c.JSON(http.StatusUnprocessableEntity, Resp{
		ErrorCode:    http.StatusUnprocessableEntity,
		ErrorMessage: err.Error(),
		Success:      false,
		Data:         nil,
		Total:        0,
	})
}

func (c *Context) BadRequest(err error) {
	c.JSON(http.StatusBadRequest, Resp{
		ErrorCode:    http.StatusBadRequest,
		ErrorMessage: err.Error(),
		Success:      false,
		Data:         nil,
		Total:        0,
	})
}

func (c *Context) NotFound(err error) {
	c.JSON(http.StatusNotFound, Resp{
		ErrorCode:    http.StatusNotFound,
		ErrorMessage: err.Error(),
		Success:      false,
		Data:         nil,
		Total:        0,
	})
}

func (c *Context) Unauthorized() {
	c.JSON(http.StatusUnauthorized, Resp{
		ErrorCode:    http.StatusUnauthorized,
		ErrorMessage: "unauthorized",
		Success:      false,
		Data:         nil,
		Total:        0,
	})
}

func (c *Context) Forbidden() {
	c.JSON(http.StatusForbidden, Resp{
		ErrorCode:    http.StatusForbidden,
		ErrorMessage: "forbidden",
		Success:      false,
		Data:         nil,
		Total:        0,
	})
}

func (c *Context) ServerError(err error) {
	c.JSON(http.StatusInternalServerError, Resp{
		ErrorCode:    http.StatusInternalServerError,
		ErrorMessage: err.Error(),
		Success:      false,
		Data:         nil,
		Total:        0,
	})
}
