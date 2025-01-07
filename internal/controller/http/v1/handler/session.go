package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golanguzb70/system-design-openbudget/config"
	"github.com/golanguzb70/system-design-openbudget/internal/entity"
)

// GetSession godoc
// @Router /session/{id} [get]
// @Summary Get a session by ID
// @Description Get a session by ID
// @Security BearerAuth
// @Tags session
// @Accept  json
// @Produce  json
// @Param id path string true "Session ID"
// @Success 200 {object} entity.Session
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) GetSession(ctx *gin.Context) {
	var (
		req entity.Id
	)

	req.ID = ctx.Param("id")

	session, err := h.UseCase.SessionRepo.GetSingle(ctx, req)
	if h.HandleDbError(ctx, err, "Error getting session") {
		return
	}

	ctx.JSON(200, session)
}

// GetSessions godoc
// @Router /session/list [get]
// @Summary Get a list of users
// @Description Get a list of users
// @Security BearerAuth
// @Tags session
// @Accept  json
// @Produce  json
// @Param page query number true "page"
// @Param limit query number true "limit"
// @Param user_id query string false "user_id"
// @Success 200 {object} entity.SessionList
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) GetSessions(ctx *gin.Context) {
	var (
		req entity.GetListFilter
	)

	page := ctx.DefaultQuery("page", "1")
	limit := ctx.DefaultQuery("limit", "10")
	userId := ctx.DefaultQuery("user_id", "")

	if ctx.GetHeader("user_type") == "client" {
		userId = ctx.GetHeader("sub")
	}

	req.Page, _ = strconv.Atoi(page)
	req.Limit, _ = strconv.Atoi(limit)
	req.Filters = append(req.Filters,
		entity.Filter{
			Column: "user_id",
			Type:   "eq",
			Value:  userId,
		},
	)

	req.OrderBy = append(req.OrderBy, entity.OrderBy{
		Column: "created_at",
		Order:  "desc",
	})

	users, err := h.UseCase.SessionRepo.GetList(ctx, req)
	if h.HandleDbError(ctx, err, "Error getting users") {
		return
	}

	ctx.JSON(200, users)
}

// UpdateSession godoc
// @Router /session [put]
// @Summary Update a session
// @Description Update a session
// @Security BearerAuth
// @Tags session
// @Accept  json
// @Produce  json
// @Param session body entity.Session true "Session object"
// @Success 200 {object} entity.Session
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) UpdateSession(ctx *gin.Context) {
	var (
		body entity.Session
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	session, err := h.UseCase.SessionRepo.Update(ctx, body)
	if h.HandleDbError(ctx, err, "Error updating session") {
		return
	}

	ctx.JSON(200, session)
}

// DeleteSession godoc
// @Router /session/{id} [delete]
// @Summary Delete a session
// @Description Delete a session
// @Security BearerAuth
// @Tags session
// @Accept  json
// @Produce  json
// @Param id path string true "Session ID"
// @Success 200 {object} entity.SuccessResponse
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) DeleteSession(ctx *gin.Context) {
	var (
		req entity.Id
	)

	req.ID = ctx.Param("id")

	err := h.UseCase.SessionRepo.Delete(ctx, req)
	if h.HandleDbError(ctx, err, "Error deleting session") {
		return
	}

	ctx.JSON(200, entity.SuccessResponse{
		Message: "Session deleted successfully",
	})
}
