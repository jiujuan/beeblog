// Package admin 后台文章管理处理器（需要 Admin 角色）。
package admin

import (
	"strconv"

	"github.com/gin-gonic/gin"

	appArticle "antblog/internal/application/article"
	"antblog/internal/interfaces/http/middleware"
	"antblog/pkg/response"
	"antblog/pkg/utils"
	"antblog/pkg/validator"
)

// ArticleHandler 后台文章管理处理器
type ArticleHandler struct {
	useCase appArticle.IArticleUseCase
}

// NewAdminArticleHandler 创建后台文章处理器
func NewAdminArticleHandler(uc appArticle.IArticleUseCase) *ArticleHandler {
	return &ArticleHandler{useCase: uc}
}

// ListArticles godoc
// @Summary  后台文章列表（含草稿/归档）
// @Tags     admin-article
// @Security BearerAuth
// @Param    page        query int    false "页码"
// @Param    page_size   query int    false "每页条数"
// @Param    status      query int    false "状态过滤 1=草稿 2=已发布 3=归档"
// @Param    category_id query int    false "分类ID"
// @Param    keyword     query string false "标题关键字"
// @Success  200 {object} response.Response{data=response.PageData[appArticle.ArticleListItemResp]}
// @Router   /api/admin/articles [get]
func (h *ArticleHandler) ListArticles(c *gin.Context) {
	var req appArticle.AdminListArticleReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误："+err.Error())
		return
	}
	req.Page = utils.NormalizePage(req.Page)
	req.PageSize = utils.NormalizePageSize(req.PageSize)

	list, total, err := h.useCase.AdminListArticles(c.Request.Context(), &req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OK(c, response.NewPageData(list, total, req.Page, req.PageSize))
}

// GetArticle godoc
// @Summary  后台按 ID 获取文章（含草稿）
// @Tags     admin-article
// @Security BearerAuth
// @Param    id path int true "文章 ID"
// @Success  200 {object} response.Response{data=appArticle.ArticleResp}
// @Router   /api/admin/articles/{id} [get]
func (h *ArticleHandler) GetArticle(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的文章 ID")
		return
	}
	art, err := h.useCase.AdminGetArticle(c.Request.Context(), id)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OK(c, art)
}

// CreateArticle godoc
// @Summary  创建文章（后台 Markdown 编辑器提交）
// @Tags     admin-article
// @Security BearerAuth
// @Accept   json
// @Produce  json
// @Param    body body appArticle.CreateArticleReq true "文章内容"
// @Success  201  {object} response.Response{data=appArticle.ArticleResp}
// @Router   /api/admin/articles [post]
func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	var req appArticle.CreateArticleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求体解析失败："+err.Error())
		return
	}
	if err := validator.Validate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	authorID := middleware.MustGetCurrentUserID(c)
	art, err := h.useCase.CreateArticle(c.Request.Context(), authorID, &req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.Created(c, art)
}

// UpdateArticle godoc
// @Summary  更新文章内容
// @Tags     admin-article
// @Security BearerAuth
// @Accept   json
// @Produce  json
// @Param    id   path int                      true "文章 ID"
// @Param    body body appArticle.UpdateArticleReq true "文章内容"
// @Success  200  {object} response.Response{data=appArticle.ArticleResp}
// @Router   /api/admin/articles/{id} [put]
func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的文章 ID")
		return
	}

	var req appArticle.UpdateArticleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求体解析失败："+err.Error())
		return
	}
	if err := validator.Validate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	art, err := h.useCase.UpdateArticle(c.Request.Context(), id, &req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OK(c, art)
}

// UpdateArticleStatus godoc
// @Summary  变更文章状态（发布/归档/撤稿）
// @Tags     admin-article
// @Security BearerAuth
// @Accept   json
// @Produce  json
// @Param    id   path int                         true "文章 ID"
// @Param    body body appArticle.UpdateStatusReq  true "目标状态"
// @Success  200  {object} response.Response{data=appArticle.ArticleResp}
// @Router   /api/admin/articles/{id}/status [patch]
func (h *ArticleHandler) UpdateArticleStatus(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的文章 ID")
		return
	}

	var req appArticle.UpdateStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求体解析失败："+err.Error())
		return
	}
	if err := validator.Validate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	art, err := h.useCase.UpdateArticleStatus(c.Request.Context(), id, &req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OK(c, art)
}

// DeleteArticle godoc
// @Summary  删除文章（软删除）
// @Tags     admin-article
// @Security BearerAuth
// @Param    id path int true "文章 ID"
// @Success  200 {object} response.Response
// @Router   /api/admin/articles/{id} [delete]
func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的文章 ID")
		return
	}
	if err := h.useCase.DeleteArticle(c.Request.Context(), id); err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OKWithMsg(c, "删除成功", nil)
}

// ─── 辅助函数 ────────────────────────────────────────────────────────────────

func parseUintParam(c *gin.Context, key string) (uint64, error) {
	return strconv.ParseUint(c.Param(key), 10, 64)
}
