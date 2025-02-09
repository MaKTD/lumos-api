//go:build unit
// +build unit

package errs

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAggregate(t *testing.T) {
	err := Aggregate(nil, nil, nil, nil)
	require.Equal(t, nil, err)

	inputErr := NewErrorf(ErrCodeForbidden, "forbidden")
	err = Aggregate(inputErr)
	require.Equal(t, inputErr, err)

	err = Aggregate(
		inputErr,
		errors.New("failed, can not do"),
		errors.New("operation aborted xxx lib"),
	)
	require.Equal(t, "forbidden: failed, can not do: operation aborted xxx lib", err.Error())

	var targetErr *CodeError
	ok := errors.As(err, &targetErr)
	require.Equal(t, true, ok)
	require.Equal(t, targetErr, inputErr)
}
