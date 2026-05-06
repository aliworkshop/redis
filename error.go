package redis

import "github.com/aliworkshop/errors"

var AlreadyLockedErr = errors.Validation().WithMessage("key has been already locked")