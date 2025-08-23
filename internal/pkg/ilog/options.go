package ilog

type options struct {
	console   bool
	accesslog bool
}

type Option func(o *options)

func WithConsole() Option {
	return func(o *options) {
		o.console = true
	}
}

func WithAccessLog() Option {
	return func(o *options) {
		o.accesslog = true
	}
}
