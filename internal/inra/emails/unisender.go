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
	AfterTrialExpiredListTitle string
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
	return &UniSenderSrv{client: client}
}

func (r *UniSenderSrv) ScheduleAfterTrialExpired(ctx context.Context, email string) error {
	return r.subsribeToList(ctx, email, r.cfg.AfterTrialExpiredListTitle)
}

func (r *UniSenderSrv) subsribeToList(ctx context.Context, email string, listTitle string) error {
	lists, err := r.client.GetLists().Execute()
	if err != nil {
		return err
	}
	list := r.findListWithTitle(listTitle, lists)
	if list == nil {
		return fmt.Errorf("failed to find list with title %s", r.cfg.AfterTrialExpiredListTitle)
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
