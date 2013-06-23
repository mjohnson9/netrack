package app

import (
	"appengine"
)

type Future interface {
	GetResult() (interface{}, error)
	Wait()
}

func ResolveFutureError(errChan chan error) error {
	errors := make(appengine.MultiError, 0)

	for err := range errChan {
		if multiError, ok := err.(appengine.MultiError); ok {
			for _, thisErr := range multiError {
				if thisErr != nil {
					errors = append(errors, thisErr)
				}
			}
		} else {
			if err != nil {
				errors = append(errors, err)
			}
		}
	}

	if len(errors) == 1 {
		return errors[0]
	} else if len(errors) > 1 {
		return errors
	}

	return nil
}
