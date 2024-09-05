package controller

import (
	"github.com/cuihairu/salon/internal/biz"
	"github.com/cuihairu/salon/internal/middleware"
	"github.com/cuihairu/salon/internal/model"
	"github.com/cuihairu/salon/internal/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type MemberAPI struct {
	memberBiz *biz.MemberBiz
	logger    *zap.Logger
}

func NewMemberAPI(memberBiz *biz.MemberBiz, logger *zap.Logger) *MemberAPI {
	return &MemberAPI{
		memberBiz: memberBiz,
		logger:    logger,
	}
}

func (api *MemberAPI) RegisterRoutes(router *gin.RouterGroup) {
	memberGroup := router.Group("/members")
	{
		memberGroup.GET("/", middleware.RequiredRole(middleware.User), api.GetAllMembers)
		memberGroup.GET("/:id", middleware.RequiredRole(middleware.User), api.GetMemberByID)
		memberGroup.PUT("/:id", middleware.RequiredRole(middleware.User), api.UpdateMember)
		memberGroup.DELETE("/:id", middleware.RequiredRole(middleware.User), api.DeleteMember)
	}
}

func (api *MemberAPI) GetAllMembers(c *gin.Context) {

	members, err := api.memberBiz.GetAllMembers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, members)
}

func (api *MemberAPI) GetMemberByID(c *gin.Context) {
	id, err := utils.ParseUintParam[uint](c, "id")
	if err != nil {
		return
	}
	member, err := api.memberBiz.GetMemberByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, member)
}

func (api *MemberAPI) UpdateMember(c *gin.Context) {
	var member model.Member
	if err := c.BindJSON(&member); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := api.memberBiz.UpdateMember(&member); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, member)
}

func (api *MemberAPI) DeleteMember(c *gin.Context) {
	id, err := utils.ParseUintParam[uint](c, "id")
	if err != nil {
		return
	}
	if err := api.memberBiz.DeleteMember(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
