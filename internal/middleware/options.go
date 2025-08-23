package middleware

import (
	"context"
)

type CountryCodeConvert interface {
	AreaCode(context.Context, string) (string, error)
	TwoAreaCode(context.Context, string) (string, error)
}

type TextPlainReply interface {
	StringReply() string
}

type Options struct {
	convert CountryCodeConvert
}

type Option func(options Options)

func WithCountryCodeConvert(convert CountryCodeConvert) Option {
	return func(o Options) {
		o.convert = convert
	}
}
