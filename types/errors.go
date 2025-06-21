package types

import "github.com/pkg/errors"

var ErrNoBeanProvider = errors.New("no bean provider found")
var ErrNoValueFound = errors.New("no value found")
var ErrInjectNotSupported = errors.New("inject not supported for bean type")
var ErrNonBeanType = errors.New("not a bean type")
var ErrParseTag = errors.New("parse tag error")
var ErrNilContext = errors.New("nil context")
var ErrContextAlreadyExists = errors.New("context already exists")
var ErrCycleDependencies = errors.New("detected cycle dependency")

func ErrNoInjectableProvided(err error) bool {
	return errors.Is(err, ErrNoBeanProvider) || errors.Is(err, ErrNoValueFound)
}
