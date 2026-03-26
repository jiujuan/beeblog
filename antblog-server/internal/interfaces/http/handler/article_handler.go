package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	appArticle "antblog/internal/application/article"
	"antblog/internal/interfaces/http/middleware"
	"antblog/pkg/response"
	"antblog/pkg/utils"
)

// ArticleHandler 前台文章 HTTP 处理器
type ArticleHandler struct {
	useCase appArticle.IArticleUseCase
}

// NewArticleHandler 创建文章前台处理器
func NewArticleHandler(uc appArticle.IArticleUseCase) *ArticleHandler {
	return &ArticleHandler{useCase: uc}
}

// ListArticles godoc
// @Summary  前台文章列表（分页）
// @Tags     article
// @Produce  json
// @Param    page        query int    false "页码"     default(1)
// @Param    page_size   query int    false "每页条数"  default(10)
// @Param    category_id query int    false "分类ID过滤"
// @Param    tag_id      query int    false "标签ID过滤"
// @Param    keyword     query string false "标题关键字"
// @Success  200 {object} response.Response{data=response.PageData[appArticle.ArticleListItemResp]}
// @Router   /api/v1/articles [get]
func (h *ArticleHandler) ListArticles(c *gin.Context) {
	var req appArticle.ListArticleReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误："+err.Error())
		return
	}
	req.Page = utils.NormalizePage(req.Page)
	req.PageSize = utils.NormalizePageSize(req.PageSize)

	// 尝试获取当前用户 ID（OptionalAuth 中间件注入，未登录时为 nil）
	userID := middleware.GetOptionalUserID(c)
	list, total, err := h.useCase.ListArticles(c.Request.Context(), &req, userID)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OK(c, response.NewPageData(list, total, req.Page, req.PageSize))
}

// GetArticleBySlug godoc
// @Summary  前台按 Slug 获取文章详情
// @Tags     article
// @Produce  json
// @Param    slug path string true "文章 Slug"
// @Success  200 {object} response.Response{data=appArticle.ArticleResp}
// @Router   /api/v1/articles/{slug} [get]
func (h *ArticleHandler) GetArticleBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		response.BadRequest(c, "slug 不能为空")
		return
	}

	// 尝试获取当前用户 ID（OptionalAuth 中间件注入，未登录时为 nil）
	userID := middleware.GetOptionalUserID(c)
	art, err := h.useCase.GetArticleBySlug(c.Request.Context(), slug, userID)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OK(c, art)
}

// GetFeaturedArticles godoc
// @Summary  精选文章列表
// @Tags     article
// @Produce  json
// @Param    limit query int false "数量限制" default(6)
// @Success  200 {object} response.Response{data=[]appArticle.ArticleListItemResp}
// @Router   /api/v1/articles/featured [get]
func (h *ArticleHandler) GetFeaturedArticles(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "6")
	limit, _ := strconv.Atoi(limitStr)
	userID := middleware.GetOptionalUserID(c)
	list, err := h.useCase.GetFeaturedArticles(c.Request.Context(), limit, userID)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OK(c, list)
}

// GetArchive godoc
// @Summary  文章归档时间线
// @Tags     article
// @Produce  json
// @Success  200 {object} response.Response{data=[]appArticle.ArchiveItemResp}
// @Router   /api/v1/articles/archive [get]
func (h *ArticleHandler) GetArchive(c *gin.Context) {
	items, err := h.useCase.GetArchive(c.Request.Context())
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OK(c, items)
}

// GetArchiveDetail godoc
// @Summary  归档详情（某年月文章列表）
// @Tags     article
// @Produce  json
// @Param    year      query int true  "年份"
// @Param    month     query int true  "月份（1-12）"
// @Param    page      query int false "页码"
// @Param    page_size query int false "每页条数"
// @Success  200 {object} response.Response{data=response.PageData[appArticle.ArticleListItemResp]}
// @Router   /api/v1/articles/archive/detail [get]
func (h *ArticleHandler) GetArchiveDetail(c *gin.Context) {
	var req appArticle.ArchiveDetailReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误："+err.Error())
		return
	}
	req.Page = utils.NormalizePage(req.Page)
	req.PageSize = utils.NormalizePageSize(req.PageSize)

	userID := middleware.GetOptionalUserID(c)
	list, total, err := h.useCase.GetArchiveDetail(c.Request.Context(), &req, userID)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OK(c, response.NewPageData(list, total, req.Page, req.PageSize))
}

// ─── 互动接口（需登录） ───────────────────────────────────────────────────────

// LikeArticle godoc
// @Summary  点赞文章
// @Tags     article
// @Security BearerAuth
// @Param    id path int true "文章 ID"
// @Success  200 {object} response.Response
// @Router   /api/v1/articles/{id}/like [post]
func (h *ArticleHandler) LikeArticle(c *gin.Context) {
	articleID, err := parseUintParam(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的文章 ID")
		return
	}
	userID := middleware.MustGetCurrentUserID(c)
	if err := h.useCase.LikeArticle(c.Request.Context(), articleID, userID); err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OKWithMsg(c, "点赞成功", nil)
}

// UnlikeArticle godoc
// @Summary  取消点赞
// @Tags     article
// @Security BearerAuth
// @Param    id path int true "文章 ID"
// @Success  200 {object} response.Response
// @Router   /api/v1/articles/{id}/like [delete]
func (h *ArticleHandler) UnlikeArticle(c *gin.Context) {
	articleID, err := parseUintParam(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的文章 ID")
		return
	}
	userID := middleware.MustGetCurrentUserID(c)
	if err := h.useCase.UnlikeArticle(c.Request.Context(), articleID, userID); err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OKWithMsg(c, "已取消点赞", nil)
}

// BookmarkArticle godoc
// @Summary  收藏文章
// @Tags     article
// @Security BearerAuth
// @Param    id path int true "文章 ID"
// @Success  200 {object} response.Response
// @Router   /api/v1/articles/{id}/bookmark [post]
func (h *ArticleHandler) BookmarkArticle(c *gin.Context) {
	articleID, err := parseUintParam(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的文章 ID")
		return
	}
	userID := middleware.MustGetCurrentUserID(c)
	if err := h.useCase.BookmarkArticle(c.Request.Context(), articleID, userID); err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OKWithMsg(c, "收藏成功", nil)
}

// UnbookmarkArticle godoc
// @Summary  取消收藏
// @Tags     article
// @Security BearerAuth
// @Param    id path int true "文章 ID"
// @Success  200 {object} response.Response
// @Router   /api/v1/articles/{id}/bookmark [delete]
func (h *ArticleHandler) UnbookmarkArticle(c *gin.Context) {
	articleID, err := parseUintParam(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的文章 ID")
		return
	}
	userID := middleware.MustGetCurrentUserID(c)
	if err := h.useCase.UnbookmarkArticle(c.Request.Context(), articleID, userID); err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OKWithMsg(c, "已取消收藏", nil)
}

// GetUserBookmarks godoc
// @Summary  获取当前用户的收藏列表
// @Tags     article
// @Security BearerAuth
// @Param    page      query int false "页码"
// @Param    page_size query int false "每页条数"
// @Success  200 {object} response.Response{data=response.PageData[appArticle.ArticleListItemResp]}
// @Router   /api/v1/user/bookmarks [get]
func (h *ArticleHandler) GetUserBookmarks(c *gin.Context) {
	userID := middleware.MustGetCurrentUserID(c)
	page := utils.NormalizePage(parseIntQuery(c, "page", 1))
	pageSize := utils.NormalizePageSize(parseIntQuery(c, "page_size", 10))

	list, total, err := h.useCase.GetUserBookmarks(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OK(c, response.NewPageData(list, total, page, pageSize))
}

// ─── 辅助函数 ────────────────────────────────────────────────────────────────

func parseUintParam(c *gin.Context, key string) (uint64, error) {
	return strconv.ParseUint(c.Param(key), 10, 64)
}

func parseIntQuery(c *gin.Context, key string, defaultVal int) int {
	s := c.Query(key)
	if s == "" {
		return defaultVal
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return v
}
