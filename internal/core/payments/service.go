package payments

import (
	"context"
	"doctormakarhina/lumos/internal/core/domain"
	"doctormakarhina/lumos/internal/core/notify"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrUserAlreadyRegistered = errors.New("user already registered")
	oneMonthSubNames         = []string{
		strings.ToLower(strings.TrimSpace("Доступ на месяц")),
		strings.ToLower(strings.TrimSpace("Продление 1 месяц")),
	}
	threeMonthsSubNames = []string{
		strings.ToLower(strings.TrimSpace("Доступ на 3 месяца")),
		strings.ToLower(strings.TrimSpace("Продление 3 месяца")),
	}
	sixMonthsSubNames = []string{
		strings.ToLower(strings.TrimSpace("Доступ на 6 месяцев")),
		strings.ToLower(strings.TrimSpace("Продление 6 месяцев")),
	}
)

type service struct {
	repo          UserRepo
	emails        EmailsSrv
	notif         notify.Service
	cloudPayments CloudPayments
}

func NewPaymentsService(
	repo UserRepo,
	emails EmailsSrv,
	notif notify.Service,
	cloudPayments CloudPayments,
) Service {
	return &service{
		repo:          repo,
		emails:        emails,
		notif:         notif,
		cloudPayments: cloudPayments,
	}
}

func (s *service) RegisterFromTrial(
	ctx context.Context,
	email string,
	name string,
	trialDuration time.Duration,
) error {
	emailNorm := s.normalizeStr(email)

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

func (s *service) RegisterFromProdamus(
	ctx context.Context,
	subName string,
	email string,
	name string,
	price float32,
) error {
	emailNorm := s.normalizeStr(email)
	subNorm := s.normalizeStr(subName)

	var tariffName string
	if slices.Contains(oneMonthSubNames, subNorm) {
		tariffName = domain.UserTariff1Month
	} else if slices.Contains(threeMonthsSubNames, subNorm) {
		tariffName = domain.UserTariff3Months
	} else if slices.Contains(sixMonthsSubNames, subNorm) {
		tariffName = domain.UserTariff6Months
	}

	if tariffName == "" {
		s.notif.ForAdmin(fmt.Sprintf("[RegisterFromProdamus]: unknown tariff (%s), can not update user %s subscription", subNorm, email))
		return nil
	}

	candidate := domain.User{
		ID:                 uuid.New().String(),
		Email:              emailNorm,
		Name:               name,
		Tariff:             tariffName,
		ExpiresAt:          time.Now().Add(-time.Minute),
		SubscriptionID:     "",
		SubscriptionStatus: "",
		LastSubPrice:       price,
	}

	user, err := s.repo.FindByEmailOrCreate(ctx, candidate)
	if err != nil {
		s.notif.ForAdmin(fmt.Sprintf("[RegisterFromProdamus]: error finding or creating user (%s): %v", emailNorm, err))
		return err
	}
	if user == nil {
		s.notif.ForAdmin(fmt.Sprintf("[RegisterFromProdamus]: failed to find or create user (%s)", emailNorm))
		return fmt.Errorf("failed to find or create user")
	}

	err = s.emails.CancelTrialExpired(ctx, user.Email)
	if err != nil {
		s.notif.ForAdmin(fmt.Sprintf("[RegisterFromProdamus]: error canceling trial for email (%s): %v", emailNorm, err))
	}

	newExpiresAt, err := user.NewSubEndedAt(time.Now(), tariffName)
	if err != nil {
		s.notif.ForAdmin(fmt.Sprintf("[RegisterFromProdamus]: failed to calculate new expiration date for user (%s): %v", emailNorm, err))
		return nil
	}

	user.Tariff = tariffName
	user.ExpiresAt = newExpiresAt
	user.SubscriptionID = ""
	user.SubscriptionStatus = ""
	user.LastSubPrice = price

	_, err = s.repo.UpdateSub(ctx, *user)
	if err != nil {
		s.notif.ForAdmin(fmt.Sprintf("[RegisterFromProdamus]: error updating user (%s): %v", emailNorm, err))
		return err
	}

	s.notif.ForAdmin(fmt.Sprintf(
		"[RegisterFromProdamus]: user payment is accepted (%s), tariff: %s, expires at: %s",
		emailNorm,
		tariffName,
		newExpiresAt.Format(time.RFC3339),
	))

	return nil
}

func (s *service) RegisterFromCloudPayments(
	ctx context.Context,
	subName string,
	email string,
	name string,
	price float32,
	subscriptionID string,
) error {
	subscriptionID = strings.TrimSpace(subscriptionID)
	emailNorm := s.normalizeStr(email)
	subNorm := s.normalizeStr(subName)

	var tariffName string
	var subInterval int
	if slices.Contains(oneMonthSubNames, subNorm) {
		tariffName = domain.UserTariff1Month
		subInterval = 1
	} else if slices.Contains(threeMonthsSubNames, subNorm) {
		tariffName = domain.UserTariff3Months
		subInterval = 3
	} else if slices.Contains(sixMonthsSubNames, subNorm) {
		tariffName = domain.UserTariff6Months
		subInterval = 6
	}

	if tariffName == "" || subInterval == 0 {
		s.notif.ForAdmin(fmt.Sprintf("[RegisterFromCloudPayments]: unknown tariff (%s), can not update user %s subscription", subNorm, email))
		return nil
	}

	candidate := domain.User{
		ID:                 uuid.New().String(),
		Email:              emailNorm,
		Name:               name,
		Tariff:             tariffName,
		ExpiresAt:          time.Now().Add(-time.Minute),
		SubscriptionID:     "",
		SubscriptionStatus: "",
		LastSubPrice:       price,
	}
	user, err := s.repo.FindByEmailOrCreate(ctx, candidate)
	if err != nil {
		s.notif.ForAdmin(fmt.Sprintf("[RegisterFromCloudPayments]: error finding or creating user (%s): %v", emailNorm, err))
		return err
	}
	if user == nil {
		s.notif.ForAdmin(fmt.Sprintf("[RegisterFromCloudPayments]: failed to find or create user (%s)", emailNorm))
		return fmt.Errorf("failed to find or create user")
	}

	err = s.emails.CancelTrialExpired(ctx, user.Email)
	if err != nil {
		s.notif.ForAdmin(fmt.Sprintf("[RegisterFromCloudPayments]: error canceling trial for email (%s): %v", emailNorm, err))
	}

	newExpiresAt, err := user.NewSubEndedAt(time.Now(), tariffName)
	if err != nil {
		s.notif.ForAdmin(fmt.Sprintf("[RegisterFromCloudPayments]: failed to calculate new expiration date for user (%s): %v", emailNorm, err))
		return nil
	}

	oldSubId := strings.TrimSpace(user.SubscriptionID)

	subStatus, err := s.cloudPayments.UpdateSubscription(ctx, subscriptionID, newExpiresAt, "Month", subInterval)
	if err != nil {
		s.notif.ForAdmin(fmt.Sprintf("[RegisterFromCloudPayments]: error updating subscription for user (%s): %v", emailNorm, err))
		return err
	}

	user.Tariff = tariffName
	user.ExpiresAt = newExpiresAt
	user.SubscriptionID = subscriptionID
	user.SubscriptionStatus = subStatus
	user.LastSubPrice = price

	_, err = s.repo.UpdateSub(ctx, *user)
	if err != nil {
		s.notif.ForAdmin(fmt.Sprintf("[RegisterFromCloudPayments]: error updating user (%s): %v", emailNorm, err))
		return err
	}

	if oldSubId != "" && oldSubId != strings.TrimSpace(subscriptionID) {
		err := s.cloudPayments.CancelSubscription(ctx, oldSubId)
		if err != nil {
			s.notif.ForAdmin(fmt.Sprintf("[RegisterFromCloudPayments]: error canceling old subscription for user (%s): %v", emailNorm, err))
		}
	}

	if oldSubId != "" && oldSubId == strings.TrimSpace(subscriptionID) {
		err := s.emails.ScheduleAfterReccurrentPayment(ctx, user.Email)
		if err != nil {
			s.notif.ForAdmin(fmt.Sprintf("[RegisterFromCloudPayments]: error scheduling after recurrent email for user (%s): %v", emailNorm, err))
		}
	}

	s.notif.ForAdmin(fmt.Sprintf(
		"[RegisterFromCloudPayments]: user payment is accepted (%s), tariff: %s, expires at: %s",
		emailNorm,
		tariffName,
		newExpiresAt.Format(time.RFC3339),
	))

	return nil
}

func (s *service) RegisterCloudPaymentReccurent(
	ctx context.Context,
	subscriptionID string,
	email string,
	status string,
) error {
	status = strings.TrimSpace(status)
	subscriptionID = strings.TrimSpace(subscriptionID)
	emailNorm := s.normalizeStr(email)

	if subscriptionID == "" {
		s.notif.ForAdmin(fmt.Sprintf("[RegisterCloudPaymentReccurent]: cloud payment recurrent notification for user (%s): subscription ID is empty", emailNorm))

		return nil
	}

	if status == "Cancelled" {
		err := s.emails.ScheduleAfterAutopaymentCancelled(ctx, emailNorm)
		if err != nil {
			s.notif.ForAdmin(fmt.Sprintf("[RegisterCloudPaymentReccurent]: error scheduling after autopayment cancelled email for user (%s): %v", emailNorm, err))
		}
	}

	err := s.repo.UpdateSubStatusBySubID(ctx, subscriptionID, status)
	if err != nil {
		s.notif.ForAdmin(fmt.Sprintf("[RegisterCloudPaymentReccurent]: error updating subscription status for user (%s): %v", emailNorm, err))
	}

	s.notif.ForAdmin(fmt.Sprintf(
		"[RegisterCloudPaymentReccurent]: recurrent notification accepted for user (%s), status = %s",
		emailNorm,
		status,
	))

	return nil
}

func (s *service) normalizeStr(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
