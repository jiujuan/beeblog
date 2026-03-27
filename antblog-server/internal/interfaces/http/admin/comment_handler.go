// Package admin 后台评论管理处理器（需要 Admin 角色）。
package admin

import (
	"github.com/gin-gonic/gin"

	appcomment "antblog/internal/application/comment"
	"antblog/pkg/response"
	"antblog/pkg/utils"
	"antblog/pkg/validator"
)

// CommentHandler 后台评论管理处理器
type CommentHandler struct {
	useCase appcomment.ICommentUseCase
}

// NewAdminCommentHandler 创建后台评论处理器
func NewAdminCommentHandler(uc appcomment.ICommentUseCase) *CommentHandler {
	return &CommentHandler{useCase: uc}
}

// ListComments godoc
// @Summary  后台评论列表（全状态，支持多条件过滤）
// @Tags     admin-comment
// @Security BearerAuth
// @Produce  json
// @Param    page       query int    false "页码"
// @Param    page_size  query int    false "每页条数"
// @Param    article_id query int    false "文章 ID 过滤"
// @Param    status     query int    false "状态过滤 1=待审核 2=已通过 3=已拒绝 4=垃圾"
// @Param    keyword    query string false "内容关键字"
// @Success  200 {object} response.Response{data=response.PageData[appcomment.CommentResp]}
// @Router   /api/admin/comments [get]
func (h *CommentHandler) ListComments(c *gin.Context) {
	var req appcomment.AdminListCommentReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误："+err.Error())
		return
	}
	req.Page = utils.NormalizePage(req.Page)
	req.PageSize = utils.NormalizePageSize(req.PageSize)

	list, total, err := h.useCase.AdminListComments(c.Request.Context(), &req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OK(c, response.NewPageData(list, total, req.Page, req.PageSize))
}

// GetComment godoc
// @Summary  后台按 ID 获取评论详情
// @Tags     admin-comment
// @Security BearerAuth
// @Produce  json
// @Param    id path int true "评论 ID"
// @Success  200 {object} response.Response{data=appcomment.CommentResp}
// @Router   /api/admin/comments/{id} [get]
func (h *CommentHandler) GetComment(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的评论 ID")
		return
	}
	comment, err := h.useCase.AdminGetComment(c.Request.Context(), id)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OK(c, comment)
}

// UpdateStatus godoc
// @Summary  审核评论（通过/拒绝/标记垃圾/重置待审核）
// @Tags     admin-comment
// @Security BearerAuth
// @Accept   json
// @Produce  json
// @Param    id   path int                            true "评论 ID"
// @Param    body body appcomment.AdminUpdateStatusReq true "目标状态"
// @Success  200  {object} response.Response{data=appcomment.CommentResp}
// @Router   /api/admin/comments/{id}/status [patch]
func (h *CommentHandler) UpdateStatus(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的评论 ID")
		return
	}

	var req appcomment.AdminUpdateStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求体解析失败："+err.Error())
		return
	}
	if err := validator.Validate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	comment, err := h.useCase.AdminUpdateStatus(c.Request.Context(), id, &req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OK(c, comment)
}

// DeleteComment godoc
// @Summary  删除评论（软删除）
// @Tags     admin-comment
// @Security BearerAuth
// @Produce  json
// @Param    id path int true "评论 ID"
// @Success  200 {object} response.Response
// @Router   /api/admin/comments/{id} [delete]
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的评论 ID")
		return
	}
	if err := h.useCase.AdminDeleteComment(c.Request.Context(), id); err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OKWithMsg(c, "删除成功", nil)
}
