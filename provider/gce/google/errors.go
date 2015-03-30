// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package google

import (
	"fmt"

	"github.com/juju/errors"
)

// InvalidConfigValue indicates that one of the config values failed validation.
type InvalidConfigValue struct {
	errors.Err
	cause error

	// Key is the OS env var corresponding to the field with the bad value.
	Key string

	// Value is the invalid value.
	Value interface{}

	// Reason is the underlying error.
	Reason error
}

// NewInvalidConfigValue returns a new InvalidConfigValue for the given
// info. If the provided reason is an error then Reason is set to that
// error. Otherwise a non-nil value is treated as a string and Reason is
// set to a non-nil value that wraps it.
func NewInvalidConfigValue(key string, value, reason interface{}) error {
	var underlying error
	switch reason := reason.(type) {
	case error:
		underlying = reason
	default:
		if reason != nil {
			underlying = errors.Errorf("%v", reason)
		}
	}
	err := &InvalidConfigValue{
		cause:  errors.NewNotValid(underlying, "GCE config value"),
		Key:    key,
		Value:  value,
		Reason: underlying,
	}
	msg := "config value for " + key
	if value != nil {
		if strValue, ok := value.(string); ok {
			if strValue != "" {
				msg = fmt.Sprintf("%s %q", msg, strValue)
			}
		} else {
			msg = fmt.Sprintf("%s (%v)", msg, value)
		}
	}
	err.Err = errors.NewErr("config value")
	err.Err.SetLocation(1)
	return err
}

// NewMissingConfigValue returns a new error for a missing config field.
func NewMissingConfigValue(key, field string) error {
	return NewInvalidConfigValue(key, "", "missing "+field)
}

// Cause implements errors.causer. This is necessary so that
// errors.IsNotValid works.
func (err *InvalidConfigValue) Cause() error {
	return err.cause
}

// Underlying implements errors.wrapper.
func (err InvalidConfigValue) Underlying() error {
	return err.cause
}

// Error implements error.
func (err InvalidConfigValue) Error() string {
	return fmt.Sprintf("invalid config value (%s) for %q: %v", err.Value, err.Key)
}
