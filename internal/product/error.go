package product

import "errors"

var (
	ErrProductNotFound = errors.New("product not found")
	ErrForbiddenAccess = errors.New("user has no access to the data")
)
