package emails

import (
	"context"
	"fmt"

	"github.com/sergeyandreenko/unisender"
	_ "github.com/sergeyandreenko/unisender"
)

type UniSenderSrv struct {
	client *unisender.UniSender
}

func NewUniSenderSrv(
	apiKey string,
) *UniSenderSrv {
	client := unisender.New(apiKey)
	return &UniSenderSrv{client: client}
}

func (r *UniSenderSrv) ScheduleAfterTrialExpired(ctx context.Context, email string) error {
	lists, err := r.client.GetLists().Execute()
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", lists)

	// s.client.List

	// s.client.Subscribe()

	return nil
}
