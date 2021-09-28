package IPFilter

import "errors"

var (
	ErrInvalidIndex     = errors.New("the '.' index is out of range")
	ErrInvalidIPAddress = errors.New("the IPAddress is invalid")
	ErrChildNode        = errors.New("the child node could not be added")
)
