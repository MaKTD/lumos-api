package payments

import (
	"context"
	"doctormakarhina/lumos/internal/core/domain"
	"doctormakarhina/lumos/internal/core/notify"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrUserAlreadyRegistered = errors.New("user already registered")
)

type service struct {
	repo   UserRepo
	emails EmailsSrv
	notif  notify.Service
}

func NewPaymentsService(
	repo UserRepo,
	emails EmailsSrv,
	notif notify.Service,
) Service {
	return &service{
		repo:   repo,
		emails: emails,
		notif:  notif,
	}
}

func (s *service) RegisterFromTrial(
	ctx context.Context,
	email string,
	name string,
	phone string,
	trialDuration time.Duration,
) error {
	emailNorm := s.normalizeEmail(email)

	user, err := s.repo.ByEmail(ctx, emailNorm)
	if err != nil {
		s.notif.ForAdmin(fmt.Sprintf("[RegisterFromTrial]: error fetching user by email (%s): %v", emailNorm, err))
		return err
	}
	if user != nil {
		s.notif.ForAdmin(fmt.Sprintf("[RegisterFromTrial]: user (%s) already registered, skipping trial registration", emailNorm))
		return ErrUserAlreadyRegistered
	}

	user = &domain.User{
		ID:                 uuid.New().String(),
		Email:              emailNorm,
		Name:               name,
		Tariff:             domain.UserTariffTrial,
		ExpiresAt:          time.Now().Add(trialDuration),
		SubscriptionID:     "",
		SubscriptionStatus: "",
		LastSubPrice:       0,
	}
	_, err = s.repo.Create(ctx, *user)
	if err != nil {
		s.notif.ForAdmin(fmt.Sprintf("[RegisterFromTrial]: error creating user (%s): %v", emailNorm, err))
		return err
	}

	err = s.emails.ScheduleAfterTrialExpired(ctx, emailNorm)
	if err != nil {
		s.notif.ForAdmin(fmt.Sprintf("[RegisterFromTrial]: error scheduling email after trial expired for email (%s): %v", emailNorm, err))
		return err
	}

	s.notif.ForAdmin(
		fmt.Sprintf("[RegisterFromTrial]: user (%s) registered for trial, duration = %d days",
			emailNorm,
			int(trialDuration.Hours()/24),
		))

	return nil
}

func (s *service) normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
