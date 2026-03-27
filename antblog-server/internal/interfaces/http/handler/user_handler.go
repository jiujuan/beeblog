// Package handler HTTP 接口处理层，负责请求解析、校验和响应格式化。
package handler

import (
	"github.com/gin-gonic/gin"

	appuser "antblog/internal/application/user"
	"antblog/internal/interfaces/http/middleware"
	"antblog/pkg/response"
	"antblog/pkg/validator"
)

// UserHandler 用户模块 HTTP 处理器
type UserHandler struct {
	useCase appuser.IUserUseCase
}

// NewUserHandler 创建用户处理器
func NewUserHandler(uc appuser.IUserUseCase) *UserHandler {
	return &UserHandler{useCase: uc}
}

// ─── Auth 接口 ────────────────────────────────────────────────────────────────

// Register godoc
// @Summary  用户注册
// @Tags     auth
// @Accept   json
// @Produce  json
// @Param    body body appuser.RegisterReq true "注册信息"
// @Success  201  {object} response.Response{data=appuser.UserResp}
// @Router   /api/v1/auth/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req appuser.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求体解析失败："+err.Error())
		return
	}
	if err := validator.Validate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	user, err := h.useCase.Register(c.Request.Context(), &req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.Created(c, user)
}

// Login godoc
// @Summary  用户登录
// @Tags     auth
// @Accept   json
// @Produce  json
// @Param    body body appuser.LoginReq true "登录信息"
// @Success  200  {object} response.Response{data=appuser.TokenResp}
// @Router   /api/v1/auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req appuser.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求体解析失败："+err.Error())
		return
	}

	if err := validator.Validate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	token, err := h.useCase.Login(
		c.Request.Context(),
		&req,
		c.Request.UserAgent(),
		c.ClientIP(),
	)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OK(c, token)
}

// Logout godoc
// @Summary  用户登出
// @Tags     auth
// @Security BearerAuth
// @Accept   json
// @Produce  json
// @Param    body body appuser.RefreshTokenReq true "Refresh Token"
// @Success  200  {object} response.Response
// @Router   /api/v1/auth/logout [post]
func (h *UserHandler) Logout(c *gin.Context) {
	var req appuser.RefreshTokenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求体解析失败："+err.Error())
		return
	}

	if err := h.useCase.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OKWithMsg(c, "已退出登录", nil)
}

// RefreshToken godoc
// @Summary  刷新 Token
// @Tags     auth
// @Accept   json
// @Produce  json
// @Param    body body appuser.RefreshTokenReq true "Refresh Token"
// @Success  200  {object} response.Response{data=appuser.TokenResp}
// @Router   /api/v1/auth/refresh [post]
func (h *UserHandler) RefreshToken(c *gin.Context) {
	var req appuser.RefreshTokenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求体解析失败："+err.Error())
		return
	}
	if err := validator.Validate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	token, err := h.useCase.RefreshToken(
		c.Request.Context(),
		&req,
		c.Request.UserAgent(),
		c.ClientIP(),
	)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OK(c, token)
}

// ─── 用户个人接口 ─────────────────────────────────────────────────────────────

// GetProfile godoc
// @Summary  获取个人资料
// @Tags     user
// @Security BearerAuth
// @Produce  json
// @Success  200 {object} response.Response{data=appuser.UserResp}
// @Router   /api/v1/user/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := middleware.MustGetCurrentUserID(c)

	user, err := h.useCase.GetProfile(c.Request.Context(), userID)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OK(c, user)
}

// UpdateProfile godoc
// @Summary  更新个人资料
// @Tags     user
// @Security BearerAuth
// @Accept   json
// @Produce  json
// @Param    body body appuser.UpdateProfileReq true "资料信息"
// @Success  200  {object} response.Response{data=appuser.UserResp}
// @Router   /api/v1/user/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := middleware.MustGetCurrentUserID(c)

	var req appuser.UpdateProfileReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求体解析失败："+err.Error())
		return
	}
	if err := validator.Validate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	user, err := h.useCase.UpdateProfile(c.Request.Context(), userID, &req)
	if err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OK(c, user)
}

// ChangePassword godoc
// @Summary  修改密码
// @Tags     user
// @Security BearerAuth
// @Accept   json
// @Produce  json
// @Param    body body appuser.ChangePasswordReq true "密码信息"
// @Success  200  {object} response.Response
// @Router   /api/v1/user/password [put]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID := middleware.MustGetCurrentUserID(c)

	var req appuser.ChangePasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求体解析失败："+err.Error())
		return
	}
	if err := validator.Validate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.useCase.ChangePassword(c.Request.Context(), userID, &req); err != nil {
		response.FailWithError(c, err)
		return
	}
	response.OKWithMsg(c, "密码修改成功，请重新登录", nil)
}
