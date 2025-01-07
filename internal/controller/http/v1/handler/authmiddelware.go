package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
	"github.com/golanguzb70/system-design-openbudget/internal/entity"
	"github.com/golanguzb70/system-design-openbudget/pkg/jwt"
)

func (h *Handler) AuthMiddleware(e *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			usertype string

			act = c.Request.Method
			obj = c.FullPath()
		)

		token := c.GetHeader("Authorization")
		if token == "" {
			usertype = "unauthorized"
		}

		if usertype == "" {
			token = strings.TrimPrefix(token, "Bearer ")

			claims, err := jwt.ParseJWT(token, h.Config.JWT.Secret)
			if err != nil {
				usertype = "unauthorized"
			}

			v, ok := claims["user_type"].(string)
			if !ok {
				usertype = "unauthorized"
			} else {
				usertype = v
			}

			for key, value := range claims {
				c.Request.Header.Set(key, fmt.Sprintf("%v", value))
			}
		}

		// TO DO: Check if session is valid

		if usertype != "unauthorized" {
			session, err := h.UseCase.SessionRepo.GetSingle(c, entity.Id{ID: c.GetHeader("session_id")})
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Session is invalid"})
				return
			}

			if !session.IsActive {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Session is invalid"})
				return
			}
		}

		ok, err := e.EnforceSafe(usertype, obj, act)
		if err != nil {
			h.Logger.Error(err, "Error enforcing policy")
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied"})
			return
		}

		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied"})
			return
		}

		c.Next()
	}
}
