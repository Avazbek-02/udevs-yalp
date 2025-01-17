package handler

import (
	"fmt"
	"strconv"

	"github.com/Avazbek-02/udevslab-lesson6/config"
	"github.com/Avazbek-02/udevslab-lesson6/internal/entity"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreatetReview godoc
// @Router /review [post]
// @Summary Create a review
// @Description Create a review
// @Security BearerAuth
// @Tags review
// @Accept  json
// @Produce  json
// @Param review body entity.Review true "Review object"
// @Success 200 {object} entity.Review
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) CreateReview(ctx *gin.Context) {
	var (
		req entity.Review
	)

	err := ctx.BindJSON(&req)
	if h.HandleDbError(ctx, err, "Error getting review") {
		return
	}
	UserId := ctx.GetHeader("sub")
	req.UserID = UserId
	fmt.Println(UserId)
	res, err := h.UseCase.ReviewRepo.Create(ctx, req)
	if h.HandleDbError(ctx, err, "Error creating review") {
		return
	}
	ctx.JSON(200, res)
}

// GetReview godoc
// @Router /review/{id} [get]
// @Summary Get a review by ID
// @Description Get a review by ID
// @Security BearerAuth
// @Tags review
// @Accept  json
// @Produce  json
// @Param id path string true "Review ID"
// @Success 200 {object} entity.Review
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) GetReview(ctx *gin.Context) {
	var (
		req entity.Id
	)

	req.ID = ctx.Param("id")

	review, err := h.UseCase.ReviewRepo.GetSingle(ctx, req)
	if h.HandleDbError(ctx, err, "Error getting review") {
		return
	}

	ctx.JSON(200, review)
}

// GetReviews godoc
// @Router /review/list [get]
// @Summary Get a list of reviews
// @Description Get a list of reviews
// @Security BearerAuth
// @Tags review
// @Accept  json
// @Produce  json
// @Param page query number true "page"
// @Param limit query number true "limit"
// @Param business_id query string false "business_id"
// @Success 200 {object} entity.ReviewList
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) GetReviews(ctx *gin.Context) {
	var (
		req entity.GetListFilter
	)

	page := ctx.DefaultQuery("page", "1")
	limit := ctx.DefaultQuery("limit", "10")
	businessID := ctx.DefaultQuery("business_id", "")

	req.Page, _ = strconv.Atoi(page)
	req.Limit, _ = strconv.Atoi(limit)
	req.Filters = append(req.Filters,
		entity.Filter{
			Column: "business_id",
			Type:   "eq",
			Value:  businessID,
		},
	)

	req.OrderBy = append(req.OrderBy, entity.OrderBy{
		Column: "created_at",
		Order:  "desc",
	})
	if _, err := uuid.Parse(businessID); err != nil && businessID != "" {
		ctx.JSON(404,gin.H{"Error:":"Wrong format type please write UUID"})
		return
	}

	reviews, err := h.UseCase.ReviewRepo.GetList(ctx, req)
	if h.HandleDbError(ctx, err, "Error getting reviews") {
		return
	}

	ctx.JSON(200, reviews)
}

// UpdateReview godoc
// @Router /review [put]
// @Summary Update a review
// @Description Update a review
// @Security BearerAuth
// @Tags review
// @Accept  json
// @Produce  json
// @Param review body entity.Review true "Review object"
// @Success 200 {object} entity.Review
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) UpdateReview(ctx *gin.Context) {
	var (
		body entity.Review
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	review, err := h.UseCase.ReviewRepo.Update(ctx, body)
	if h.HandleDbError(ctx, err, "Error updating review") {
		return
	}

	ctx.JSON(200, review)
}

// DeleteReview godoc
// @Router /review/{id} [delete]
// @Summary Delete a review
// @Description Delete a review
// @Security BearerAuth
// @Tags review
// @Accept  json
// @Produce  json
// @Param id path string true "Review ID"
// @Success 200 {object} entity.SuccessResponse
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) DeleteReview(ctx *gin.Context) {
	var (
		req entity.Id
	)

	req.ID = ctx.Param("id")

	err := h.UseCase.ReviewRepo.Delete(ctx, req)
	if h.HandleDbError(ctx, err, "Error deleting review") {
		return
	}

	ctx.JSON(200, entity.SuccessResponse{
		Message: "Review deleted successfully",
	})
}
