package middleware

import (
	"github.com/cuihairu/salon/internal/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	ContextRole = "AccessRole"
)

var (
	subordinate = make(map[string]*Role)
	Admin       = newAdmin()
	User        = NewRole("user")
	Anonymous   = NewRole("anonymous")
)

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
	if role == r || role.name == r.name {
		return true
	}
	_, ok := r.subordinate[role.name]
	return ok
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
	return func(c *gin.Context) {
		// check permission
		if requiredRole == Anonymous || requiredRole.name == Anonymous.name {
			c.Next()
			return
		}
		claims, ok := utils.GetClaimsFormContext(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"errorMessage": "unauthorized"})
			c.Abort()
			return
		}
		session := sessions.Default(c)
		oldToken := session.Get("token")
		if oldToken == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"errorMessage": "unauthorized"})
			c.Abort()
			return
		}
		oldTokenStr, ok := oldToken.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"errorMessage": "unauthorized"})
			c.Abort()
			return
		}
		token := utils.GetHeaderToken(c)
		if oldTokenStr != token {
			c.JSON(http.StatusUnauthorized, gin.H{"errorMessage": "unauthorized"})
			c.Abort()
			return
		}
		// check role
		curRole := getRoleName(claims.Role)
		if permission := curRole.HasPermission(requiredRole); !permission {
			c.JSON(http.StatusForbidden, gin.H{"errorMessage": "forbidden", "errorCode": http.StatusForbidden})
			c.Abort()
			return
		}
		SetRole(c, curRole)
		c.Next()
	}
}
