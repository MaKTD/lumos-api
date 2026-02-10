package emails

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/sergeyandreenko/unisender"
	"github.com/sergeyandreenko/unisender/lists"
)

type UniSenderSrvCfg struct {
	AfterTrialExpiredListTitle         string
	AfterReccurrentPaymentListTitle    string
	AfterAutopaymentCancelledListTitle string
}

type UniSenderSrv struct {
	client *unisender.UniSender
	cfg    UniSenderSrvCfg
}

func NewUniSenderSrv(
	apiKey string,
	cfg UniSenderSrvCfg,
) *UniSenderSrv {
	client := unisender.New(apiKey)
	return &UniSenderSrv{client: client, cfg: cfg}
}

func (r *UniSenderSrv) ScheduleAfterTrialExpired(ctx context.Context, email string) error {
	return r.subsribeToList(ctx, email, r.cfg.AfterTrialExpiredListTitle)
}

func (r *UniSenderSrv) CancelTrialExpired(ctx context.Context, email string) error {
	return r.excludeFromList(ctx, email, r.cfg.AfterTrialExpiredListTitle)
}

func (r *UniSenderSrv) ScheduleAfterReccurrentPayment(ctx context.Context, email string) error {
	return r.subsribeToList(ctx, email, r.cfg.AfterReccurrentPaymentListTitle)
}

func (r *UniSenderSrv) ScheduleAfterAutopaymentCancelled(ctx context.Context, email string) error {
	return r.subsribeToList(ctx, email, r.cfg.AfterAutopaymentCancelledListTitle)
}

func (r *UniSenderSrv) subsribeToList(ctx context.Context, email string, listTitle string) error {
	lists, err := r.client.GetLists().Execute()
	if err != nil {
		return err
	}
	list := r.findListWithTitle(listTitle, lists)
	if list == nil {
		return fmt.Errorf("failed to find list with title %s", listTitle)
	}

	_, err = r.client.Subscribe(list.ID).
		Email(email).
		DoNotOverwrite().
		DoubleOptinConfirmed().
		Execute()

	if err != nil {
		return err
	}

	return nil
}

func (r *UniSenderSrv) excludeFromList(ctx context.Context, email string, listTitle string) error {
	lists, err := r.client.GetLists().Execute()
	if err != nil {
		return err
	}
	list := r.findListWithTitle(listTitle, lists)
	if list == nil {
		return fmt.Errorf("failed to find list with title %s", listTitle)
	}

	return r.client.Exclude(email).ContactTypeEmail().ListIDs(list.ID).Execute()
}

func (r *UniSenderSrv) findListWithTitle(targetTitle string, all []lists.GetListsResult) *lists.GetListsResult {
	norm := strings.ToLower(strings.TrimSpace(targetTitle))
	i := slices.IndexFunc(all, func(l lists.GetListsResult) bool {
		return strings.ToLower(strings.TrimSpace(l.Title)) == norm
	})

	if i == -1 {
		return nil
	}
	return &all[i]
}
