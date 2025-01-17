package usecase

import (
	"github.com/Avazbek-02/udevslab-lesson6/config"
	"github.com/Avazbek-02/udevslab-lesson6/internal/usecase/repo"
	"github.com/Avazbek-02/udevslab-lesson6/pkg/logger"
	"github.com/Avazbek-02/udevslab-lesson6/pkg/postgres"
)

// UseCase -.
type UseCase struct {
	UserRepo    UserRepoI
	SessionRepo SessionRepoI
	BusinessRepo BusinessRepoI
	NotificationRepo NotificationRepoI
	ReviewRepo ReviewRepoI
	ReportRepo ReportRepoI
	EventRepo EventRepoI
}

// New -.
func New(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *UseCase {
	return &UseCase{
		UserRepo:    repo.NewUserRepo(pg, config, logger),
		SessionRepo: repo.NewSessionRepo(pg, config, logger),
		BusinessRepo: repo.NewBusinessRepo(pg, config, logger),
		NotificationRepo: repo.NewNotificationRepo(pg, config, logger),
		ReviewRepo: repo.NewReviewRepo(pg, config, logger),
		ReportRepo: repo.NewReportRepo(pg, config, logger),
		EventRepo: repo.NewEventRepo(pg, config, logger),
	}
}
