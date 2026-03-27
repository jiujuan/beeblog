package handler

import (
	"github.com/gin-gonic/gin"

	appmedia "antblog/internal/application/media"
	"antblog/internal/interfaces/http/middleware"
	"antblog/pkg/response"
	"antblog/pkg/utils"
)

// MediaHandler 用户侧媒体 HTTP 处理器（需登录，查看自己上传的文件）
type MediaHandler struct {
	useCase appmedia.IMediaUseCase
}

// NewMediaHandler 创建用户侧媒体处理器
func NewMediaHandler(uc appmedia.IMediaUseCase) *MediaHandler {
	return &MediaHandler{useCase: uc}
}

// ListMyMedia godoc
// @Summary  获取当前用户的媒体列表（分页）
// @Tags     media
// @Security BearerAuth
// @Produce  json
// @Param    page      query int false "页码"
// @Param    page_size query int false "每页条数"
// @Success  200 {object} response.Response{data=response.PageData[appmedia.MediaResp]}
// @Router   /api/v1/user/media [get]
func (h *MediaHandler) ListMyMedia(c *gin.Context) {
	userID := middleware.MustGetCurrentUserID(c)
	page := utils.NormalizePage(parseIntQuery(c, "page", 1))
	pageSize := utils.NormalizePageSize(parseIntQuery(c, "page_size", 20))

	list, total, err := h.useCase.ListMyMedia(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OK(c, response.NewPageData(list, total, page, pageSize))
}
