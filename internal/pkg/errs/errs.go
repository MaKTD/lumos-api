package errs

import (
	"fmt"
)

type ErrCode int

const (
	ErrCodeUnknown         ErrCode = 0
	ErrCodeNotFound        ErrCode = 404
	ErrCodeParsingFailed   ErrCode = 400
	ErrCodeInvalidArgument ErrCode = 401
	ErrCodeForbidden       ErrCode = 403
	ErrCodeInternal        ErrCode = 500
	ErrCodeTimeout         ErrCode = 408
)

func CodeToString(code ErrCode) string {
	switch code {
	case ErrCodeUnknown:
		return "unknown"
	case ErrCodeNotFound:
		return "not_found"
	case ErrCodeParsingFailed:
		return "parsing_failed"
	case ErrCodeInvalidArgument:
		return "invalid_argument"
	case ErrCodeForbidden:
		return "forbidden"
	case ErrCodeInternal:
		return "internal"
	case ErrCodeTimeout:
		return "timeout"
	default:
		return "unspecified"
	}
}

type CodeError struct {
	orig error
	code ErrCode
	msg  string
}

func (e *CodeError) Unwrap() error {
	return e.orig
}

func (e *CodeError) Error() string {
	if e.orig != nil {
		return fmt.Sprintf("%s, %v", e.msg, e.orig)
	}
	return e.msg
}

func (e *CodeError) Code() ErrCode {
	return e.code
}

func WrapErrorf(origin error, code ErrCode, format string, a ...any) error {
	return &CodeError{
		orig: origin,
		code: code,
		msg:  fmt.Sprintf(format, a...),
	}
}

func NewErrorf(code ErrCode, format string, a ...any) error {
	return WrapErrorf(nil, code, format, a...)
}

func Aggregate(errs ...error) error {
	if len(errs) == 0 {
		return nil
	}

	final := errs[0]
	for i := 1; i < len(errs); i++ {
		if errs[i] == nil {
			continue
		}
		final = fmt.Errorf("%w: %w", final, errs[i])
	}

	return final
}

func First(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}
