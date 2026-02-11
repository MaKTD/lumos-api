package payments

import (
	"context"
	"doctormakarhina/lumos/internal/core/domain"
	"time"
)

type CloudPayments interface {
	UpdateSubscription(ctx context.Context, ID string, startDate time.Time, inteval string, period int) (string, error)
	CancelSubscription(ctx context.Context, ID string) error
}

type UserRepo interface {
	FindByEmailOrCreate(ctx context.Context, user domain.User) (*domain.User, error)
	ByEmail(ctx context.Context, email string) (*domain.User, error)
	Create(ctx context.Context, user domain.User) (*domain.User, error)
	UpdateSub(ctx context.Context, user domain.User) (*domain.User, error)
	UpdateSubStatusBySubID(ctx context.Context, subscriptionID string, status string) error
}

type EmailsSrv interface {
	ScheduleAfterTrialExpired(ctx context.Context, email string) error
	CancelTrialExpired(ctx context.Context, email string) error
	ScheduleAfterReccurrentPayment(ctx context.Context, email string) error
	ScheduleAfterAutopaymentCancelled(ctx context.Context, email string) error
}

type Service interface {
	IsAccessAlowed(
		ctx context.Context,
		email string,
		projectID string,
	) (bool, error)
	User(
		ctx context.Context,
		email string,
		projectID string,
	) (*domain.User, error)
	RegisterFromTrial(
		ctx context.Context,
		email string,
		name string,
		trialDuration time.Duration,
	) error
	RegisterFromProdamus(
		ctx context.Context,
		subName string,
		email string,
		name string,
		price float32,
	) error
	RegisterFromCloudPayments(
		ctx context.Context,
		subName string,
		email string,
		name string,
		price float32,
		subscriptionID string,
	) error
	RegisterCloudPaymentReccurent(
		ctx context.Context,
		subscriptionID string,
		email string,
		status string,
	) error
}
