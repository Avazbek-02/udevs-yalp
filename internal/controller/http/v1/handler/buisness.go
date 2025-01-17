package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/Avazbek-02/udevslab-lesson6/config"
	"github.com/Avazbek-02/udevslab-lesson6/internal/entity"
)

// CreateBusiness godoc
// @Router /business [post]
// @Summary Create a new business
// @Description Create a new business entity
// @Tags business
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param business body entity.Business true "Business data"
// @Success 200 {object} entity.Business
// @Failure 400 {object} entity.ErrorResponse
// @Failure 500 {object} entity.ErrorResponse
func (h *Handler) CreateBuisiness(c *gin.Context) {
	var (
		req entity.Business
	)
	err := c.BindJSON(&req) // Fixed pointer issue
	if h.HandleDbError(c, err, "Error getting business") {
		return
	}
	req.OwnerID = c.GetHeader("sub")
	res, err := h.UseCase.BusinessRepo.Create(c, req)
	if h.HandleDbError(c, err, "Error creating business") {
		return
	}
	var(
		newUpdateType entity.User
	)
	newUpdateType.UserType = "businessman"
	newUpdateType.ID = req.OwnerID
	_, err = h.UseCase.UserRepo.Update(c,newUpdateType)
	if h.HandleDbError(c, err, "Error update type user in business") {
		return
	}
	c.Request.Header.Set("user-type", "businessman")
	c.JSON(200, res)
}


// GetBusiness godoc
// @Router /business/{id} [get]
// @Summary Get a business by ID
// @Description Get a business by ID
// @Security BearerAuth
// @Tags business
// @Accept  json
// @Produce  json
// @Param id path string true "Business ID"
// @Success 200 {object} entity.Business
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) GetBusiness(ctx *gin.Context) {
	var (
		req entity.Id
	)

	req.ID = ctx.Param("id")

	business, err := h.UseCase.BusinessRepo.GetSingle(ctx, req)
	if h.HandleDbError(ctx, err, "Error getting business") {
		return
	}

	ctx.JSON(200, business)
}

// GetBusinesses godoc
// @Router /business/list [get]
// @Summary Get a list of businesses
// @Description Get a list of businesses
// @Security BearerAuth
// @Tags business
// @Accept  json
// @Produce  json
// @Param page query number true "page"
// @Param limit query number true "limit"
// @Param owner_id query string false "owner_id"
// @Success 200 {object} entity.BusinessList
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) GetBusinesses(ctx *gin.Context) {
	var (
		req entity.GetListFilter
	)

	page := ctx.DefaultQuery("page", "1")
	limit := ctx.DefaultQuery("limit", "10")
	ownerId := ctx.DefaultQuery("owner_id", "")

	if ctx.GetHeader("user_type") == "user" {
		ownerId = ctx.GetHeader("sub")
	}

	req.Page, _ = strconv.Atoi(page)
	req.Limit, _ = strconv.Atoi(limit)
	req.Filters = append(req.Filters,
		entity.Filter{
			Column: "owner_id",
			Type:   "eq",
			Value:  ownerId,
		},
	)

	req.OrderBy = append(req.OrderBy, entity.OrderBy{
		Column: "created_at",
		Order:  "desc",
	})

	businesses, err := h.UseCase.BusinessRepo.GetList(ctx, req)
	if h.HandleDbError(ctx, err, "Error getting businesses") {
		return
	}

	ctx.JSON(200, businesses)
}

// UpdateBusiness godoc
// @Router /business [put]
// @Summary Update a business
// @Description Update a business
// @Security BearerAuth
// @Tags business
// @Accept  json
// @Produce  json
// @Param business body entity.Business true "Business object"
// @Success 200 {object} entity.Business
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) UpdateBusiness(ctx *gin.Context) {
	var (
		body entity.Business
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	business, err := h.UseCase.BusinessRepo.Update(ctx, body)
	if h.HandleDbError(ctx, err, "Error updating business") {
		return
	}

	ctx.JSON(200, business)
}

// DeleteBusiness godoc
// @Router /business/{id} [delete]
// @Summary Delete a business
// @Description Delete a business
// @Security BearerAuth
// @Tags business
// @Accept  json
// @Produce  json
// @Param id path string true "Business ID"
// @Success 200 {object} entity.SuccessResponse
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) DeleteBusiness(ctx *gin.Context) {
	var (
		req entity.Id
	)

	req.ID = ctx.Param("id")

	err := h.UseCase.BusinessRepo.Delete(ctx, req)
	if h.HandleDbError(ctx, err, "Error deleting business") {
		return
	}

	ctx.JSON(200, entity.SuccessResponse{
		Message: "Business deleted successfully",
	})
}
