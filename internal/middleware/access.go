package middleware

import (
	"fmt"
	"github.com/cuihairu/salon/internal/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

const (
	ContextRole = "AccessRole"
)

type Access struct {
	logger         *zap.Logger
	getRoleHandler HandlerFunc
}

var (
	subordinate = make(map[string]*Role)
	Admin       = newAdmin()
	User        = NewRole("user")
	Anonymous   = NewRole("anonymous") // public url
	access      = Access{}
)

type HandlerFunc func(*gin.Context) (*Role, error)

func SetGetRoleHandler(handler HandlerFunc) {
	access.getRoleHandler = handler
}

func SetLogger(logger *zap.Logger) {
	access.logger = logger
}

var defaultHandler = func(c *gin.Context) (*Role, error) {
	claims, ok := utils.GetClaimsFormContext(c)
	var err = fmt.Errorf("unauthorized")
	if !ok {
		access.logger.Error("unauthorized missing claims", zap.String("path", c.Request.URL.Path))
		return nil, err
	}
	session := sessions.Default(c)
	oldToken := session.Get("token")
	if oldToken == nil {
		access.logger.Error("unauthorized missing old token", zap.String("path", c.Request.URL.Path))
		return nil, err
	}
	oldTokenStr, ok := oldToken.(string)
	if !ok {
		access.logger.Error("unauthorized old token is not string", zap.String("path", c.Request.URL.Path))
		return nil, err
	}
	token := utils.GetHeaderToken(c)
	if oldTokenStr != token {
		access.logger.Error("unauthorized old token is not equal new token", zap.String("path", c.Request.URL.Path))
		return nil, err
	}
	return getRoleName(claims.Role), nil
}

func GetRole(c *gin.Context) *Role {
	roleStr, ok := c.Get(ContextRole)
	if !ok {
		return Anonymous
	}
	role, ok := roleStr.(*Role)
	if !ok {
		return Anonymous
	}
	return role
}

func SetRole(c *gin.Context, role *Role) {
	c.Set(ContextRole, role)
}

func getRoleName(roleName string) *Role {
	switch roleName {
	case "admin":
		return Admin
	case "user":
		return User
	case "anonymous":
		return Anonymous
	default:
		role, ok := subordinate[roleName]
		if !ok {
			return Anonymous
		} else {
			return role
		}
	}
}
func SetRoleByName(c *gin.Context, roleName string) {
	c.Set(ContextRole, getRoleName(roleName))
}

type Role struct {
	name        string
	subordinate map[string]*Role
	leader      *Role
}

func (r *Role) Name() string {
	return r.name
}

func (r *Role) AddSubordinate(role *Role) {
	// check loop
	if role == nil {
		panic("role is nil")
	}
	if len(role.name) == 0 {
		panic("role name is empty")
	}
	if role.leader != nil {
		panic("role has leader already :" + role.leader.name)
	}
	for leader := r; leader != nil; leader = leader.leader {
		if leader.name == role.name {
			panic("role has loop :" + role.name)
		}
	}
	role.leader = r
	subordinate[role.name] = role
	r.subordinate[role.Name()] = role
	for leader := r.leader; leader != nil; leader = leader.leader {
		leader.subordinate[role.Name()] = role
	}
}

func (r *Role) HasPermission(role *Role) bool {
	if role == nil {
		return false
	}
	if r.IsSomeRole(role) {
		return true
	}
	_, ok := r.subordinate[role.name]
	return ok
}

func (r *Role) HasPermissionByName(roleName string) bool {
	role := getRoleName(roleName)
	return r.HasPermission(role)
}

func (r *Role) IsSomeRole(role *Role) bool {
	if role == nil {
		return false
	}
	if role == r || role.name == r.name {
		return true
	}
	return false
}

func NewRole(name string) *Role {
	if len(name) == 0 {
		panic("role name is empty")
	}
	_, ok := subordinate[name]
	if ok {
		panic("role has already existed :" + name)
	}
	r := &Role{
		name:        name,
		subordinate: make(map[string]*Role),
	}
	subordinate[name] = r
	return r
}

func newAdmin() *Role {
	admin := &Role{
		name:        "admin",
		subordinate: subordinate,
		leader:      nil,
	}
	return admin
}

func RequiredRole(requiredRole *Role) gin.HandlerFunc {
	if requiredRole == nil {
		panic("required role is nil")
	}
	if access.logger == nil {
		access.logger = zap.NewNop()
	}
	if access.getRoleHandler == nil {
		access.getRoleHandler = defaultHandler
	}
	// public url
	if requiredRole.IsSomeRole(Anonymous) {
		return func(c *gin.Context) {
			c.Next()
			return
		}
	}
	return func(c *gin.Context) {
		// check permission
		curRole, err := access.getRoleHandler(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"errorMessage": "unauthorized", "errorCode": http.StatusUnauthorized})
			c.Abort()
			return
		}
		if permission := curRole.HasPermission(requiredRole); !permission {
			access.logger.Error("unauthorized forbidden", zap.String("path", c.Request.URL.Path), zap.String("current role", curRole.name), zap.String("required role", requiredRole.name))
			c.JSON(http.StatusForbidden, gin.H{"errorMessage": "forbidden", "errorCode": http.StatusForbidden})
			c.Abort()
			return
		}
		SetRole(c, curRole)
		c.Next()
	}
}
