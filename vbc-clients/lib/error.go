package lib

import "github.com/pkg/errors"

func ErrorWrap(err error, s *string) error {
	if s == nil {
		return err
	}
	if err == nil {
		return errors.New("nil | " + *s)
	}
	return errors.New(err.Error() + " | " + *s)
}
