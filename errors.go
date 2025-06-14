package yadi

import "github.com/pkg/errors"

var ErrNoBeanProvider = errors.New("no bean provider found")
var ErrNoValueFound = errors.New("no value found")
var ErrInjectNotSupported = errors.New("inject not supported for bean type")
var ErrNonBeanType = errors.New("not a bean type")
var ErrParseTag = errors.New("parse tag error")

func errNoContextProvided(err error) bool {
	return errors.Is(err, ErrNoBeanProvider) || errors.Is(err, ErrNoValueFound)
}
