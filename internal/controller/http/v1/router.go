// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"net/http"

	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// Swagger docs.
	"github.com/Avazbek-02/udevslab-lesson6/config"
	_ "github.com/Avazbek-02/udevslab-lesson6/docs"
	"github.com/Avazbek-02/udevslab-lesson6/internal/controller/http/v1/handler"
	"github.com/Avazbek-02/udevslab-lesson6/internal/usecase"
	minio "github.com/Avazbek-02/udevslab-lesson6/pkg/MinIO"
	"github.com/Avazbek-02/udevslab-lesson6/pkg/logger"
	rediscache "github.com/golanguzb70/redis-cache"
)

// NewRouter -.
// Swagger spec:
// @title       Go Clean Template API
// @description This is a sample server Go Clean Template server.
// @version     1.0
// @BasePath    /v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func NewRouter(engine *gin.Engine, l *logger.Logger, config *config.Config, useCase *usecase.UseCase, redis rediscache.RedisCache, minio *minio.MinIO) {
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	handlerV1 := handler.NewHandler(l, config, useCase, redis,minio)

	e := casbin.NewEnforcer("config/rbac.conf", "config/policy.csv")
	engine.Use(handlerV1.AuthMiddleware(e))

	url := ginSwagger.URL("swagger/doc.json")
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	engine.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	engine.GET("/metrics", gin.WrapH(promhttp.Handler()))

	v1 := engine.Group("/v1")

	user := v1.Group("/user")
	{
		user.POST("/", handlerV1.CreateUser)
		user.GET("/list", handlerV1.GetUsers)
		user.GET("/:id", handlerV1.GetUser)
		user.PUT("/", handlerV1.UpdateUser)
		user.POST("/avatar", handlerV1.SetUserAvatar)
		user.DELETE("/:id", handlerV1.DeleteUser)
	}

	session := v1.Group("/session")
	{
		session.GET("/list", handlerV1.GetSessions)
		session.GET("/:id", handlerV1.GetSession)
		session.PUT("/", handlerV1.UpdateSession)
		session.DELETE("/:id", handlerV1.DeleteSession)
	}

	auth := v1.Group("/auth")
	{
		auth.POST("/logout", handlerV1.Logout)
		auth.POST("/register", handlerV1.Register)
		auth.POST("/verify-email", handlerV1.VerifyEmail)
		auth.POST("/login", handlerV1.Login)
	}
	business := v1.Group("/business")
	{
		business.POST("/", handlerV1.CreateBuisiness)
		business.GET("/:id", handlerV1.GetBusiness)
		business.GET("/list", handlerV1.GetBusinesses)
		business.PUT("/", handlerV1.UpdateBusiness)
		business.DELETE("/:id", handlerV1.DeleteBusiness)
		business.POST("/:id/image", handlerV1.SetBusinessImage)
	}

	review := v1.Group("/review")
	{
		review.POST("/", handlerV1.CreateReview)
		review.GET("/:id", handlerV1.GetReview)
		review.GET("/list", handlerV1.GetReviews)
		review.PUT("/", handlerV1.UpdateReview)
		review.DELETE("/:id", handlerV1.DeleteReview)
		review.POST("/:id/image", handlerV1.SetReviewImage)
	}

	report := v1.Group("/report")
	{
		report.POST("/", handlerV1.CreateReport)
		report.GET("/:id", handlerV1.GetReport)
		report.GET("/list", handlerV1.GetReports)
		report.PUT("/", handlerV1.UpdateReport)
		report.DELETE("/:id", handlerV1.DeleteReport)
	}

	notification := v1.Group("/notification")
	{
		notification.POST("/", handlerV1.CreateNotification)
		notification.PUT("/", handlerV1.UpdateNotification)
		notification.PUT("/update-status", handlerV1.UpdateStatusNotification)
		notification.GET("/list", handlerV1.GetNotifications)
		notification.GET("/:id", handlerV1.GetNotification)
		notification.DELETE("/:id", handlerV1.DeleteNotification)
	}
	event := v1.Group("/event")
	{
		event.POST("/", handlerV1.CreateEvent)
		event.PUT("/", handlerV1.UpdateEvent)
		event.GET("/list", handlerV1.GetEvents)
		event.GET("/:id", handlerV1.GetEvent)
		event.DELETE("/:id", handlerV1.DeleteEvent)
		event.POST("/add-participant", handlerV1.AddParticipant)
		event.DELETE("/remove-participant", handlerV1.RemoveParticipant)
		event.GET("/:id/participants", handlerV1.GetParticipants)
	}
}