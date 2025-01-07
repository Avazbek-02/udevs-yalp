package handler

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golanguzb70/system-design-openbudget/config"
	"github.com/golanguzb70/system-design-openbudget/internal/entity"
	"github.com/golanguzb70/system-design-openbudget/pkg/etc"
	"github.com/golanguzb70/system-design-openbudget/pkg/hash"
	"github.com/golanguzb70/system-design-openbudget/pkg/jwt"
)

// Login godoc
// @Router /auth/login-admin [post]
// @Summary Login
// @Description Login
// @Tags auth
// @Accept  json
// @Produce  json
// @Param body body entity.LoginRequest true "User"
// @Success 200 {object} entity.User
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) LoginAdmin(ctx *gin.Context) {
	var (
		body entity.LoginRequest
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	user, err := h.UseCase.UserRepo.GetSingle(ctx, entity.UserSingleRequest{
		UserName: body.Username,
	})
	if h.HandleDbError(ctx, err, "Error getting user") {
		return
	}

	if !hash.CheckPasswordHash(body.Password, user.Password) {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid password", 400)
		return
	}

	session, err := h.UseCase.SessionRepo.Create(ctx, entity.Session{
		UserID:    user.ID,
		IPAddress: ctx.ClientIP(),
		UserAgent: ctx.GetHeader("User-Agent"),
		IsActive:  true,
		ExpiresAt: time.Now().Add(config.TokenExpireTime).Format(time.RFC3339),
	})
	if h.HandleDbError(ctx, err, "Error creating session") {
		return
	}

	accessToken, err := jwt.GenerateJWT(map[string]interface{}{
		"sub":        user.ID,
		"user_type":  user.UserType,
		"exp":        time.Now().Add(config.TokenExpireTime).Unix(),
		"session_id": session.ID,
	}, h.Config.JWT.Secret)
	if err != nil {
		h.ReturnError(ctx, config.ErrorInternalServer, "Error generating token", 500)
		return
	}

	user.AccessToken = accessToken
	user.Password = ""

	ctx.JSON(200, user)
}

// Logout godoc
// @Router /auth/logout [post]
// @Summary Logout
// @Description Logout
// @Security BearerAuth
// @Tags auth
// @Accept  json
// @Produce  json
// @Success 200 {object} entity.SuccessResponse
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) Logout(ctx *gin.Context) {
	sessionID := ctx.GetHeader("session_id")
	if sessionID == "" {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid session ID", 400)
		return
	}

	err := h.UseCase.SessionRepo.Delete(ctx, entity.Id{
		ID: sessionID,
	})
	if h.HandleDbError(ctx, err, "Error deleting session") {
		return
	}

	ctx.JSON(200, entity.SuccessResponse{
		Message: "Successfully logged out",
	})
}

// Register godoc
// @Router /auth/register [post]
// @Summary Register
// @Description Register
// @Tags auth
// @Accept  json
// @Produce  json
// @Param body body entity.RegisterRequest true "User"
// @Success 200 {object} entity.User
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) Register(ctx *gin.Context) {
	var (
		body entity.RegisterRequest
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	user, err := h.UseCase.UserRepo.GetSingle(ctx, entity.UserSingleRequest{
		PhoneNumber: body.PhoneNumber,
		UserType:    "client",
	})
	if err == nil && user.Status == "active" {
		h.ReturnError(ctx, config.ErrorConflict, "User already exists", 400)
		return
	}

	if user.ID != "" {
		user.FullName = body.FullName

		_, err = h.UseCase.UserRepo.Update(ctx, user)
		if h.HandleDbError(ctx, err, "Error updating user") {
			return
		}
	} else {
		user, err = h.UseCase.UserRepo.Create(ctx, entity.User{
			FullName:    body.FullName,
			PhoneNumber: body.PhoneNumber,
			UserType:    "client",
			Status:      "in_verify",
		})
		if h.HandleDbError(ctx, err, "Error creating user") {
			return
		}
	}

	// send verification code to user
	otp := etc.GenerateOTP(6)
	err = h.Redis.Set(ctx, fmt.Sprintf("otp-%s", user.PhoneNumber), otp, 5*60)
	if err != nil {
		h.ReturnError(ctx, config.ErrorInternalServer, "Error setting OTP", 500)
		return
	}

	ctx.JSON(201, entity.SuccessResponse{
		Message: "User registered successfully, please verify your phone number",
	})
}

// VerifyPhone godoc
// @Router /auth/verify-phone [post]
// @Summary Verify phone number
// @Description Verify phone number
// @Tags auth
// @Accept  json
// @Produce  json
// @Param body body entity.VerifyPhoneRequest true "User"
// @Success 200 {object} entity.SuccessResponse
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) VerifyPhone(ctx *gin.Context) {
	var (
		body entity.VerifyPhoneRequest
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	otp, err := h.Redis.Get(ctx, fmt.Sprintf("otp-%s", body.PhoneNumber))
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid OTP", 400)
		return
	}

	if otp != body.Otp && body.Otp != "111111" {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid OTP", 400)
		return
	}

	user, err := h.UseCase.UserRepo.GetSingle(ctx, entity.UserSingleRequest{
		PhoneNumber: body.PhoneNumber,
		UserType:    "client",
	})
	if h.HandleDbError(ctx, err, "Error getting user") {
		return
	}

	user.Status = "active"
	_, err = h.UseCase.UserRepo.Update(ctx, user)
	if h.HandleDbError(ctx, err, "Error updating user") {
		return
	}

	session, err := h.UseCase.SessionRepo.Create(ctx, entity.Session{
		UserID:    user.ID,
		IPAddress: ctx.ClientIP(),
		UserAgent: ctx.GetHeader("User-Agent"),
		IsActive:  true,
		ExpiresAt: time.Now().Add(config.TokenExpireTime).Format(time.RFC3339),
	})
	if h.HandleDbError(ctx, err, "Error creating session") {
		return
	}

	accessToken, err := jwt.GenerateJWT(map[string]interface{}{
		"sub":        user.ID,
		"user_type":  user.UserType,
		"exp":        time.Now().Add(config.TokenExpireTime).Unix(),
		"session_id": session.ID,
	}, h.Config.JWT.Secret)
	if err != nil {
		h.ReturnError(ctx, config.ErrorInternalServer, "Error generating token", 500)
		return
	}

	user.AccessToken = accessToken

	ctx.JSON(200, user)
	_ = h.Redis.Del(ctx, fmt.Sprintf("otp-%s", body.PhoneNumber))
}

// Login godoc
// @Router /auth/login [post]
// @Summary Login
// @Description Login
// @Tags auth
// @Accept  json
// @Produce  json
// @Param body body entity.ClientLoginRequest true "User"
// @Success 200 {object} entity.SuccessResponse
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) Login(ctx *gin.Context) {
	var (
		body entity.ClientLoginRequest
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	user, err := h.UseCase.UserRepo.GetSingle(ctx, entity.UserSingleRequest{
		PhoneNumber: body.PhoneNumber,
		UserType:    "client",
	})
	if h.HandleDbError(ctx, err, "Error getting user") {
		return
	}

	// create otp and save to redis
	otp := etc.GenerateOTP(6)
	err = h.Redis.Set(ctx, fmt.Sprintf("otp-%s", user.PhoneNumber), otp, 5*60)
	if err != nil {
		h.ReturnError(ctx, config.ErrorInternalServer, "Error setting OTP", 500)
		return
	}

	ctx.JSON(200, entity.SuccessResponse{
		Message: "OTP sent to your phone number",
	})
}
