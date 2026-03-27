package handler

import (
	"github.com/gin-gonic/gin"

	appcategory "antblog/internal/application/category"
	"antblog/pkg/response"
)

// CategoryHandler 分类前台 HTTP 处理器
type CategoryHandler struct {
	useCase appcategory.ICategoryUseCase
}

// NewCategoryHandler 创建分类处理器
func NewCategoryHandler(uc appcategory.ICategoryUseCase) *CategoryHandler {
	return &CategoryHandler{useCase: uc}
}

// ListCategories godoc
// @Summary  获取分类列表
// @Tags     category
// @Produce  json
// @Success  200 {object} response.Response{data=[]appcategory.CategoryResp}
// @Router   /api/v1/categories [get]
func (h *CategoryHandler) ListCategories(c *gin.Context) {
	list, err := h.useCase.ListCategories(c.Request.Context())
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OK(c, list)
}

// GetCategoryBySlug godoc
// @Summary  按 Slug 获取分类详情
// @Tags     category
// @Produce  json
// @Param    slug path string true "分类 Slug"
// @Success  200 {object} response.Response{data=appcategory.CategoryResp}
// @Router   /api/v1/categories/{slug} [get]
func (h *CategoryHandler) GetCategoryBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		response.BadRequest(c, "slug 不能为空")
		return
	}

	cat, err := h.useCase.GetCategoryBySlug(c.Request.Context(), slug)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OK(c, cat)
}
