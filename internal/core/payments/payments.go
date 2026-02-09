package payments

import (
	"context"
	"doctormakarhina/lumos/internal/core/domain"
	"time"
)

type UserRepo interface {
	FindByEmailOrCreate(ctx context.Context, user domain.User) (*domain.User, error)
	ByEmail(ctx context.Context, email string) (*domain.User, error)
	Create(ctx context.Context, user domain.User) (*domain.User, error)
}

type EmailsSrv interface {
	ScheduleAfterTrialExpired(ctx context.Context, email string) error
}

type Service interface {
	RegisterFromTrial(
		ctx context.Context,
		email string,
		name string,
		phone string,
		trialDuration time.Duration,
	) error
	// RegisterFromProdamus()
	// RegisterFromCP()
	// AcceptCPReccurentNotification()
}
