package validator

import (
	"context"
	"regexp"
	"strings"
	"unicode/utf8"
)

type Validator interface {
	Valid(context.Context) Evaluator
}

var EmailRegex = regexp.MustCompile("[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Evaluator map[string]string

func (e *Evaluator) CheckFieldError(ok bool, key, message string) {
	if !ok {
		e.AddFieldError(key, message)
	}

}

func (e *Evaluator) AddFieldError(key, message string) {
	if *e == nil {
		*e = make(map[string]string)
	}

	if _, exists := (*e)[key]; !exists {
		(*e)[key] = message
	}

}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MaxChar(value string, max int) bool {
	return utf8.RuneCountInString(value) <= max
}

func MinChar(value string, min int) bool {
	return utf8.RuneCountInString(value) >= min
}

func ValidateEmail(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
