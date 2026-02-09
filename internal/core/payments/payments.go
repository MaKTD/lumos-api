package payments

import (
	"doctormakarhina/lumos/internal/core/domain"
	"time"
)

type UserRepo interface {
	FindByEmailOrCreate(user domain.User) (*domain.User, error)
	ByEmail(email string) (*domain.User, error)
	Create(user domain.User) (*domain.User, error)
}

type EmailsSrv interface {
	ScheduleAfterTrialExpired(email string) error
}

type NotificationsSrv interface {
	ForAdmin(msg string)
}

type Payments interface {
	RegisterFromTrial(
		email string,
		name string,
		phone string,
		trialDuration time.Duration,
	) error
	// RegisterFromProdamus()
	// RegisterFromCP()
	// AcceptCPReccurentNotification()
}
