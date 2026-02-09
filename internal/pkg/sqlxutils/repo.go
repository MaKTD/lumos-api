package sqlxutils

import (
	"context"
)

type SqlxBaseRepo interface {
	DB() SqlxQuerying
}

func ChooseQuerierX(ctx context.Context, repo SqlxBaseRepo) SqlxQuerying {
	if tx := ExtractTxx(ctx); tx != nil {
		return tx
	}
	return repo.DB()
}
