package handler

import (
	"strconv"

	"github.com/Avazbek-02/udevslab-lesson6/config"
	"github.com/Avazbek-02/udevslab-lesson6/internal/entity"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateReport godoc
// @Router /report [post]
// @Summary Create a report
// @Description Create a report
// @Security BearerAuth
// @Tags report
// @Accept  json
// @Produce  json
// @Param report body entity.Report true "Report object"
// @Success 200 {object} entity.Report
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) CreateReport(ctx *gin.Context) {
	var (
		req entity.Report
	)

	err := ctx.BindJSON(&req)
	if h.HandleDbError(ctx, err, "Error getting report") {
		return
	}
	req.UserID = ctx.GetHeader("sub")
	res, err := h.UseCase.ReportRepo.Create(ctx, req)
	if h.HandleDbError(ctx, err, "Error creating report") {
		return
	}
	ctx.JSON(200, res)
}

// GetReport godoc
// @Router /report/{id} [get]
// @Summary Get a report by ID
// @Description Get a report by ID
// @Security BearerAuth
// @Tags report
// @Accept  json
// @Produce  json
// @Param id path string true "Report ID"
// @Success 200 {object} entity.Report
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) GetReport(ctx *gin.Context) {
	var (
		req entity.Id
	)

	req.ID = ctx.Param("id")

	report, err := h.UseCase.ReportRepo.GetSingle(ctx, req)
	if h.HandleDbError(ctx, err, "Error getting report") {
		return
	}

	ctx.JSON(200, report)
}

// GetReports godoc
// @Router /report/list [get]
// @Summary Get a list of reports
// @Description Get a list of reports
// @Security BearerAuth
// @Tags report
// @Accept  json
// @Produce  json
// @Param page query number true "page"
// @Param limit query number true "limit"
// @Param business_id query string false "business_id"
// @Success 200 {object} entity.ReportList
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) GetReports(ctx *gin.Context) {
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
		ctx.JSON(404, gin.H{"Error:": "Wrong format type please write UUID"})
		return
	}

	reports, err := h.UseCase.ReportRepo.GetList(ctx, req)
	if h.HandleDbError(ctx, err, "Error getting reports") {
		return
	}

	ctx.JSON(200, reports)
}

// UpdateReport godoc
// @Router /report [put]
// @Summary Update a report
// @Description Update a report
// @Security BearerAuth
// @Tags report
// @Accept  json
// @Produce  json
// @Param report body entity.Report true "Report object"
// @Success 200 {object} entity.Report
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) UpdateReport(ctx *gin.Context) {
	var (
		body entity.Report
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	report, err := h.UseCase.ReportRepo.Update(ctx, body)
	if h.HandleDbError(ctx, err, "Error updating report") {
		return
	}

	ctx.JSON(200, report)
}

// DeleteReport godoc
// @Router /report/{id} [delete]
// @Summary Delete a report
// @Description Delete a report
// @Security BearerAuth
// @Tags report
// @Accept  json
// @Produce  json
// @Param id path string true "Report ID"
// @Success 200 {object} entity.SuccessResponse
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) DeleteReport(ctx *gin.Context) {
	var (
		req entity.Id
	)

	req.ID = ctx.Param("id")

	err := h.UseCase.ReportRepo.Delete(ctx, req)
	if h.HandleDbError(ctx, err, "Error deleting report") {
		return
	}

	ctx.JSON(200, entity.SuccessResponse{
		Message: "Report deleted successfully",
	})
}
