// Package admin 后台媒体管理处理器（需要 Admin 角色）。
package admin

import (
	"github.com/gin-gonic/gin"

	appmedia "antblog/internal/application/media"
	"antblog/internal/interfaces/http/middleware"
	"antblog/pkg/response"
	"antblog/pkg/utils"
)

// MediaHandler 后台媒体管理处理器
type MediaHandler struct {
	useCase appmedia.IMediaUseCase
}

// NewAdminMediaHandler 创建后台媒体处理器
func NewAdminMediaHandler(uc appmedia.IMediaUseCase) *MediaHandler {
	return &MediaHandler{useCase: uc}
}

// Upload godoc
// @Summary  上传媒体文件（multipart/form-data，字段名 file）
// @Tags     admin-media
// @Security BearerAuth
// @Accept   multipart/form-data
// @Produce  json
// @Param    file       formData file  true  "上传文件"
// @Param    article_id formData int   false "关联文章 ID（可选）"
// @Success  201 {object} response.Response{data=appmedia.MediaResp}
// @Router   /api/admin/media/upload [post]
func (h *MediaHandler) Upload(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "请选择要上传的文件（字段名：file）")
		return
	}

	var req appmedia.UploadReq
	if err = c.ShouldBind(&req); err != nil {
		response.BadRequest(c, "参数错误："+err.Error())
		return
	}

	uploaderID := middleware.MustGetCurrentUserID(c)
	media, err := h.useCase.Upload(c.Request.Context(), uploaderID, fileHeader, req.ArticleID)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.Created(c, media)
}

// ListMedia godoc
// @Summary  后台媒体列表（分页，支持多条件过滤）
// @Tags     admin-media
// @Security BearerAuth
// @Produce  json
// @Param    page        query int    false "页码"
// @Param    page_size   query int    false "每页条数"
// @Param    uploader_id query int    false "上传者 ID"
// @Param    article_id  query int    false "文章 ID"
// @Param    mime_type   query string false "MIME 类型过滤，如 image/jpeg"
// @Success  200 {object} response.Response{data=response.PageData[appmedia.MediaResp]}
// @Router   /api/admin/media [get]
func (h *MediaHandler) ListMedia(c *gin.Context) {
	var req appmedia.AdminListMediaReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误："+err.Error())
		return
	}
	req.Page = utils.NormalizePage(req.Page)
	req.PageSize = utils.NormalizePageSize(req.PageSize)

	list, total, err := h.useCase.AdminListMedia(c.Request.Context(), &req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OK(c, response.NewPageData(list, total, req.Page, req.PageSize))
}

// GetMedia godoc
// @Summary  后台按 ID 获取媒体详情
// @Tags     admin-media
// @Security BearerAuth
// @Produce  json
// @Param    id path int true "媒体 ID"
// @Success  200 {object} response.Response{data=appmedia.MediaResp}
// @Router   /api/admin/media/{id} [get]
func (h *MediaHandler) GetMedia(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的媒体 ID")
		return
	}
	media, err := h.useCase.GetMedia(c.Request.Context(), id)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OK(c, media)
}

// BindArticle godoc
// @Summary  将媒体绑定/解绑文章（article_id 为 null 则解绑）
// @Tags     admin-media
// @Security BearerAuth
// @Accept   json
// @Produce  json
// @Param    id   path int                     true "媒体 ID"
// @Param    body body appmedia.BindArticleReq true "绑定信息"
// @Success  200  {object} response.Response{data=appmedia.MediaResp}
// @Router   /api/admin/media/{id}/bind [patch]
func (h *MediaHandler) BindArticle(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的媒体 ID")
		return
	}

	var req appmedia.BindArticleReq
	if err = c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求体解析失败："+err.Error())
		return
	}

	media, err := h.useCase.BindArticle(c.Request.Context(), id, &req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OK(c, media)
}

// DeleteMedia godoc
// @Summary  删除媒体资源（软删除记录 + 物理文件删除）
// @Tags     admin-media
// @Security BearerAuth
// @Produce  json
// @Param    id path int true "媒体 ID"
// @Success  200 {object} response.Response
// @Router   /api/admin/media/{id} [delete]
func (h *MediaHandler) DeleteMedia(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的媒体 ID")
		return
	}
	if err = h.useCase.DeleteMedia(c.Request.Context(), id); err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OKWithMsg(c, "删除成功", nil)
}
