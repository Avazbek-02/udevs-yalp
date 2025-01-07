package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golanguzb70/system-design-openbudget/config"
	"github.com/golanguzb70/system-design-openbudget/internal/entity"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

func (h Handler) HandleDbError(c *gin.Context, err error, message string) bool {
	if err == nil {
		return false
	}

	h.Logger.Error(err, message)
	var errorResponse entity.ErrorResponse
	statusCode := http.StatusInternalServerError

	if err == pgx.ErrNoRows {
		errorResponse = entity.ErrorResponse{
			Message: "The requested resource was not found.",
			Code:    config.ErrorNotFound,
		}
		c.JSON(http.StatusNotFound, errorResponse)
		return true
	}

	switch e := err.(type) {
	case *pgconn.PgError:
		// Handle PostgreSQL-specific errors
		switch e.Code {
		case "23505":
			// Unique constraint violation
			errorResponse = entity.ErrorResponse{
				Message: "Duplicate key error (unique constraint violation).",
				Code:    config.ErrorDuplicateKey,
			}
			statusCode = http.StatusBadRequest
		case "23503":
			// Foreign key violation
			errorResponse = entity.ErrorResponse{
				Message: "The record could not be deleted because it is used in other records.",
				Code:    config.ErrorConflict,
			}
			statusCode = http.StatusBadRequest
		case "22001":
			// Value too long for column
			errorResponse = entity.ErrorResponse{
				Message: "Value too long for column.",
				Code:    config.ErrorInvalidRequest,
			}
			statusCode = http.StatusBadRequest

		default:
			// General PostgreSQL error
			errorResponse = entity.ErrorResponse{
				Message: "Ooops! Something went wrong.",
				Code:    config.ErrorInternalServer,
			}
		}
	default:
		// General PostgreSQL error
		errorResponse = entity.ErrorResponse{
			Message: "Ooops! Something went wrong.",
			Code:    config.ErrorInternalServer,
		}
	}

	c.JSON(statusCode, errorResponse)
	return true
}

func (h Handler) ReturnError(c *gin.Context, code string, message string, statusCode int) {
	h.Logger.Error(errors.New(message), code)
	errorResponse := entity.ErrorResponse{
		Message: message,
		Code:    code,
	}
	c.JSON(statusCode, errorResponse)
}
