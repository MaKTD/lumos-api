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
	UpdateSub(ctx context.Context, user domain.User) error
}

type EmailsSrv interface {
	ScheduleAfterTrialExpired(ctx context.Context, email string) error
	CancelTrialExpired(ctx context.Context, email string) error
}

type Service interface {
	RegisterFromTrial(
		ctx context.Context,
		email string,
		name string,
		trialDuration time.Duration,
	) error
	RegisterFromProdamus(
		ctx context.Context,
		tariff string,
		email string,
		name string,
		price float32,
	) error
	// RegisterFromCP()
	// AcceptCPReccurentNotification()
}
