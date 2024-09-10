package controller

import (
	"fmt"
	"github.com/cuihairu/salon/internal/biz"
	"github.com/cuihairu/salon/internal/config"
	"github.com/cuihairu/salon/internal/middleware"
	"github.com/cuihairu/salon/internal/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/url"
	"path"
)

type StaticAPI struct {
	config *config.Config
	logger *zap.Logger
}

func (s *StaticAPI) Initialize(config *config.Config, bizStore *biz.BizStore, logger *zap.Logger) error {
	s.logger = logger
	s.config = config
	return nil
}

func (s *StaticAPI) RegisterRoutes(router *gin.RouterGroup) {
	staticGroup := router.Group("/static")
	{
		staticGroup.GET("/list", middleware.RequiredRole(middleware.Admin), s.GetStaticList)
		staticGroup.GET("/files", middleware.RequiredRole(middleware.Admin), s.ListFilesWithPermissions)
		staticGroup.POST("/upload", middleware.RequiredRole(middleware.Admin), s.Upload)
		staticGroup.DELETE("/file/*filename", middleware.RequiredRole(middleware.Admin), s.DeleteFile)
	}
}

type FileWithPermissionsView struct {
	utils.FileWithPermissions
	Url string `json:"url" yaml:"url"`
}

func (s *StaticAPI) ListFilesWithPermissions(c *gin.Context) {
	ctx := utils.NewContext(c)
	paging := ctx.Paging()
	staticConfig, err := s.config.GetStaticConfig()
	if err != nil {
		ctx.ServerError(err)
		return
	}
	filesWithPermissions, total, err := utils.ListFilesWithPermissionsAndPaging(staticConfig.StaticPath, paging.Page, paging.PageSize)
	if err != nil {
		ctx.ServerError(err)
		return
	}
	filesWithPermissionsViews := make([]FileWithPermissionsView, len(filesWithPermissions))
	for i, file := range filesWithPermissions {
		filesWithPermissionsViews[i] = FileWithPermissionsView{
			FileWithPermissions: file,
			Url:                 fmt.Sprintf("%s/%s", staticConfig.Domain, path.Join(staticConfig.UrlPath, file.Filename)),
		}
	}
	ctx.Paginated(filesWithPermissionsViews, total)
}

func (s *StaticAPI) GetStaticList(c *gin.Context) {
	ctx := utils.NewContext(c)
	staticConfig, err := s.config.GetStaticConfig()
	if err != nil {
		ctx.ServerError(err)
		return
	}
	files, err := utils.ListFiles(staticConfig.StaticPath)
	if err != nil {
		ctx.ServerError(err)
		return
	}
	fileUrls := make([]string, len(files))
	for i, file := range files {
		fileUrls[i] = fmt.Sprintf("%s%s/%s", staticConfig.Domain, "/static", file)
	}
	ctx.Success(fileUrls)
}

func (s *StaticAPI) Upload(c *gin.Context) {
	ctx := utils.NewContext(c)
	uploadFile, uploadFileHeader, err := c.Request.FormFile("file")
	if err != nil {
		ctx.BadRequest(err)
		return
	}
	staticConfig, err := s.config.GetStaticConfig()
	if err != nil {
		ctx.ServerError(err)
		return
	}
	absolutePath, err := staticConfig.GetStaticAbsolutePath()
	if err != nil {
		ctx.ServerError(err)
		return
	}
	if err := utils.SaveFile(uploadFile, absolutePath, uploadFileHeader.Filename); err != nil {
		ctx.ServerError(err)
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
	})
}

func (s *StaticAPI) DeleteFile(c *gin.Context) {
	ctx := utils.NewContext(c)
	filename := c.Param("filename")
	if filename == "" {
		ctx.BadRequest(nil)
		return
	}
	decodedFilename, err := url.QueryUnescape(filename)
	if err != nil {
		ctx.BadRequest(err)
		return
	}
	staticConfig, err := s.config.GetStaticConfig()
	if err != nil {
		ctx.ServerError(err)
		return
	}
	absolutePath, err := staticConfig.GetStaticAbsolutePath()
	if err != nil {
		ctx.ServerError(err)
		return
	}
	if err := utils.DeleteFile(absolutePath, decodedFilename); err != nil {
		ctx.ServerError(err)
		return
	}
	ctx.Success(gin.H{"status": "ok"})
}
