package handler

import (
	"strconv"

	"github.com/Avazbek-02/udevslab-lesson6/config"
	"github.com/Avazbek-02/udevslab-lesson6/internal/entity"
	"github.com/Avazbek-02/udevslab-lesson6/pkg/etc"
	"github.com/gin-gonic/gin"
)

// CreateNotification godoc
// @Router /notification [post]
// @Summary Create a new notification
// @Description Create a notification and send an email to the user
// @Security BearerAuth
// @Tags notification
// @Accept  json
// @Produce  json
// @Param body body entity.Notification true "Notification data"
// @Success 201 {object} entity.Notification
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) CreateNotification(ctx *gin.Context) {
	var (
		body entity.Notification
	)
	if err := ctx.ShouldBindJSON(&body); err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid input data", 400)
		return
	}

	body.OwnerId = ctx.GetHeader("sub")
	body.Status = "unread"
	resUser,err := h.UseCase.UserRepo.GetSingle(ctx,entity.UserSingleRequest{ID: body.UserID})
	if h.HandleDbError(ctx, err, "Error getting notification") {
		return
	}

	createdNotification, err := h.UseCase.NotificationRepo.Create(ctx, body)
	if h.HandleDbError(ctx, err, "Error creating notification") {
		return
	}

	body.Email = resUser.Email
	emailBody, err := etc.GenerateNotificationEmailBody(body.Message)
	if err != nil {
		h.ReturnError(ctx, config.ErrorInternalServer, "Error generating email body", 500)
		return
	}

	err = etc.SendEmail(h.Config.Gmail.Host, h.Config.Gmail.Port, h.Config.Gmail.Email, h.Config.Gmail.EmailPass, body.Email, emailBody)
	if err != nil {
		h.ReturnError(ctx, config.ErrorInternalServer, "Error sending email", 500)
		return
	}

	
	ctx.JSON(201, createdNotification)
}

// GetNotification godoc
// @Router /notification/{id} [get]
// @Summary Get a notification by ID
// @Description Get a notification by ID
// @Security BearerAuth
// @Tags notification
// @Accept  json
// @Produce  json
// @Param id path string true "Notification ID"
// @Success 200 {object} entity.Notification
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) GetNotification(ctx *gin.Context) {
	var (
		req entity.Id
	)

	req.ID = ctx.Param("id")

	notification, err := h.UseCase.NotificationRepo.GetSingle(ctx, req)
	if h.HandleDbError(ctx, err, "Error getting notification") {
		return
	}

	ctx.JSON(200, notification)
}

// GetNotification godoc
// @Router /notification/update-status [put]
// @Summary Update status a notification by ID
// @Description Update status notification by ID
// @Security BearerAuth
// @Tags notification
// @Accept  json
// @Produce  json
// @Param notification body entity.Notification true "Notification object"
// @Success 200 {object} entity.Notification
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) UpdateStatusNotification(ctx *gin.Context) {
	var (
		req entity.Notification
	)

	err := ctx.ShouldBindJSON(&req)
	if h.HandleDbError(ctx, err, "Error getting notification") {
		return
	}

	notification, err := h.UseCase.NotificationRepo.UpdateStatus(ctx, req)
	if h.HandleDbError(ctx, err, "Error update status notification") {
		return
	}

	ctx.JSON(200, notification)
}


// GetNotifications godoc
// @Router /notification/list [get]
// @Summary Get a list of notifications
// @Description Get a list of notifications
// @Security BearerAuth
// @Tags notification
// @Accept  json
// @Produce  json
// @Param page query number true "page"
// @Param limit query number true "limit"
// @Param user_id query string false "user_id"
// @Success 200 {object} entity.NotificationList
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) GetNotifications(ctx *gin.Context) {
	var req entity.GetListFilter

	page := ctx.DefaultQuery("page", "1")
	limit := ctx.DefaultQuery("limit", "10")
	userId := ctx.DefaultQuery("user_id", "")

	if ctx.GetHeader("user_type") == "user" {
		userId = ctx.GetHeader("sub")
	}

	req.Page, _ = strconv.Atoi(page)
	req.Limit, _ = strconv.Atoi(limit)
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	if userId != "" {
		req.Filters = append(req.Filters, entity.Filter{
			Column: "user_id",
			Type:   "eq",
			Value:  userId,
		})
	}

	req.OrderBy = append(req.OrderBy, entity.OrderBy{
		Column: "created_at",
		Order:  "desc",
	})

	notifications, err := h.UseCase.NotificationRepo.GetList(ctx, req)
	if err != nil {
		// Handle database error properly, now passing the correct error type
		h.HandleDbError(ctx, err, "Error getting notifications")
		return
	}

	// Return response
	ctx.JSON(200, notifications)
}

// UpdateNotification godoc
// @Router /notification [put]
// @Summary Update a notification
// @Description Update a notification
// @Security BearerAuth
// @Tags notification
// @Accept  json
// @Produce  json
// @Param notification body entity.Notification true "Notification object"
// @Success 200 {object} entity.Notification
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) UpdateNotification(ctx *gin.Context) {
	var (
		body entity.Notification
	)
	
	body.OwnerRole = ctx.GetHeader("user_role")
	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	notification, err := h.UseCase.NotificationRepo.Update(ctx, body)
	if h.HandleDbError(ctx, err, "Error updating notification") {
		return
	}

	ctx.JSON(200, notification)
}

// DeleteNotification godoc
// @Router /notification/{id} [delete]
// @Summary Delete a notification
// @Description Delete a notification
// @Security BearerAuth
// @Tags notification
// @Accept  json
// @Produce  json
// @Param id path string true "Notification ID"
// @Success 200 {object} entity.SuccessResponse
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) DeleteNotification(ctx *gin.Context) {
	var (
		req entity.Id
	)

	req.ID = ctx.Param("id")

	err := h.UseCase.NotificationRepo.Delete(ctx, req)
	if h.HandleDbError(ctx, err, "Error deleting notification") {
		return
	}

	ctx.JSON(200, entity.SuccessResponse{
		Message: "Notification deleted successfully",
	})
}

// UpdateNotificationStatus godoc
// @Router /notification/status [put]
// @Summary Update a notification status
// @Description Update a notification's status
// @Security BearerAuth
// @Tags notification
// @Accept  json
// @Produce  json
// @Param notification body entity.Notification true "Notification object"
// @Success 200 {object} entity.Notification
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) UpdateNotificationStatus(ctx *gin.Context) {
	var (
		body entity.Notification
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	notification, err := h.UseCase.NotificationRepo.UpdateStatus(ctx, body)
	if h.HandleDbError(ctx, err, "Error updating notification status") {
		return
	}

	ctx.JSON(200, notification)
}
