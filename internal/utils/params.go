package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

func ParseUintParam[T uint | uint32 | uint64](c *gin.Context, paramName string) (T, error) {
	param := c.Param(paramName)
	var err error
	if param == "" {
		err = fmt.Errorf("missing %s parameter", paramName)
		return 0, err
	}
	val, err := strconv.ParseUint(param, 10, 64)
	if err != nil {
		err = fmt.Errorf("parse %s parameter fail %w", paramName, err)
		return 0, err
	}
	return T(val), nil
}
