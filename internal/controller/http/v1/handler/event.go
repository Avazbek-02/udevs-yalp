package handler

import (
	"strconv"
	"time"

	"github.com/Avazbek-02/udevslab-lesson6/config"
	"github.com/Avazbek-02/udevslab-lesson6/internal/entity"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateEvent godoc
// @Router /event [post]
// @Summary Create an event
// @Description Create an event
// @Security BearerAuth
// @Tags event
// @Accept  json
// @Produce  json
// @Param event body entity.Event true "Event object"
// @Success 200 {object} entity.Event
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) CreateEvent(ctx *gin.Context) {
	var req entity.Event

	if ctx.GetHeader("user_type") != "businessman" {
		ctx.JSON(404, gin.H{"error": "Only businessmen can create an event"})
		return
	}

	err := ctx.BindJSON(&req)
	if h.HandleDbError(ctx, err, "Error creating event") {
		return
	}

	res, err := h.UseCase.EventRepo.Create(ctx, req)
	if h.HandleDbError(ctx, err, "Error creating event") {
		return
	}

	ctx.JSON(200, res)
}

// GetEvent godoc
// @Router /event/{id} [get]
// @Summary Get an event by ID
// @Description Get an event by ID
// @Security BearerAuth
// @Tags event
// @Accept  json
// @Produce  json
// @Param id path string true "Event ID"
// @Success 200 {object} entity.Event
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) GetEvent(ctx *gin.Context) {
	id := ctx.Param("id")

	event, err := h.UseCase.EventRepo.GetSingle(ctx, entity.Id{ID: id})
	if h.HandleDbError(ctx, err, "Error fetching event") {
		return
	}

	ctx.JSON(200, event)
}

// GetEvents godoc
// @Router /event/list [get]
// @Summary Get a list of events
// @Description Get a list of events
// @Security BearerAuth
// @Tags event
// @Accept  json
// @Produce  json
// @Param page query number true "Page"
// @Param limit query number true "Limit"
// @Param business_id query string false "Business ID"
// @Success 200 {object} entity.EventList
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) GetEvents(ctx *gin.Context) {
	var req entity.GetListFilter

	page := ctx.DefaultQuery("page", "1")
	limit := ctx.DefaultQuery("limit", "10")
	businessID := ctx.DefaultQuery("business_id", "")

	req.Page, _ = strconv.Atoi(page)
	req.Limit, _ = strconv.Atoi(limit)

	if businessID != "" {
		if _, err := uuid.Parse(businessID); err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid id format"})
			return
		}
		req.Filters = append(req.Filters, entity.Filter{
			Column: "business_id",
			Type:   "eq",
			Value:  businessID,
		})
	}

	req.OrderBy = append(req.OrderBy, entity.OrderBy{
		Column: "date",
		Order:  "asc",
	})

	events, err := h.UseCase.EventRepo.GetList(ctx, req)
	if h.HandleDbError(ctx, err, "Error fetching events") {
		return
	}

	ctx.JSON(200, events)
}

// UpdateEvent godoc
// @Router /event [put]
// @Summary Update an event
// @Description Update an event
// @Security BearerAuth
// @Tags event
// @Accept  json
// @Produce  json
// @Param event body entity.Event true "Event object"
// @Success 200 {object} entity.Event
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) UpdateEvent(ctx *gin.Context) {
	var req entity.Event

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}
	if _, err := uuid.Parse(req.ID); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid UUID format for ID", "detail": req.ID})
		return
	}


	updatedEvent, err := h.UseCase.EventRepo.Update(ctx, req)
	if h.HandleDbError(ctx, err, "Error updating event") {
		return
	}

	ctx.JSON(200, updatedEvent)
}

// DeleteEvent godoc
// @Router /event/{id} [delete]
// @Summary Delete an event
// @Description Delete an event
// @Security BearerAuth
// @Tags event
// @Accept  json
// @Produce  json
// @Param id path string true "Event ID"
// @Success 200 {object} entity.SuccessResponse
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) DeleteEvent(ctx *gin.Context) {
	id := ctx.Param("id")

	err := h.UseCase.EventRepo.Delete(ctx, entity.Id{ID: id})
	if h.HandleDbError(ctx, err, "Error deleting event") {
		return
	}

	ctx.JSON(200, entity.SuccessResponse{
		Message: "Event deleted successfully",
	})
}

// AddParticipant godoc
// @Router /event/add-participant [post]
// @Summary Add a participant to an event
// @Description Add a participant to an event
// @Security BearerAuth
// @Tags event
// @Accept  json
// @Produce  json
// @Param participant body entity.EventParticipant true "Participant object"
// @Success 200 {object} entity.EventParticipant
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) AddParticipant(ctx *gin.Context) {
	var req entity.EventParticipant

	err := ctx.ShouldBindJSON(&req)
	if h.HandleDbError(ctx, err, "Error adding participant") {
		return
	}

	req.UserID = ctx.GetHeader("sub") 
	req.JoinedAt = time.Now().Format(time.RFC3339)
	participant, err := h.UseCase.EventRepo.AddParticipant(ctx, req)
	if h.HandleDbError(ctx, err, "Error adding participant") {
		return
	}

	ctx.JSON(200, participant)
}

// RemoveParticipant godoc
// @Router /event/remove-participant [delete]
// @Summary Remove a participant from an event
// @Description Remove a participant from an event
// @Security BearerAuth
// @Tags event
// @Accept  json
// @Produce  json
// @Param participant body entity.EventParticipant true "Participant object"
// @Success 200 {object} entity.SuccessResponse
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) RemoveParticipant(ctx *gin.Context) {
	var req entity.EventParticipant

	err := ctx.ShouldBindJSON(&req)
	if h.HandleDbError(ctx, err, "Error removing participant") {
		return
	}

	err = h.UseCase.EventRepo.RemoveParticipant(ctx, req)
	if h.HandleDbError(ctx, err, "Error removing participant") {
		return
	}

	ctx.JSON(200, entity.SuccessResponse{
		Message: "Participant removed successfully",
	})
}

// GetParticipants godoc
// @Router /event/:id/participants [get]
// @Summary Get participants of an event
// @Description Get participants of an event
// @Security BearerAuth
// @Tags event
// @Accept  json
// @Produce  json
// @Param page query number true "Page"
// @Param limit query number true "Limit"
// @Param event_id query string true "Event ID"
// @Success 200 {object} entity.EventParticipantList
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) GetParticipants(ctx *gin.Context) {
	var req entity.GetListFilter

	eventID := ctx.Query("event_id")
	if _, err := uuid.Parse(eventID); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid event_id format"})
		return
	}

	page := ctx.DefaultQuery("page", "1")
	limit := ctx.DefaultQuery("limit", "10")

	req.Page, _ = strconv.Atoi(page)
	req.Limit, _ = strconv.Atoi(limit)
	req.Filters = append(req.Filters, entity.Filter{
		Column: "event_id",
		Type:   "eq",
		Value:  eventID,
	})

	participants, err := h.UseCase.EventRepo.GetParticipants(ctx, req)
	if h.HandleDbError(ctx, err, "Error fetching participants") {
		return
	}

	ctx.JSON(200, participants)
}
