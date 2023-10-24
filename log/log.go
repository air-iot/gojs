package log

import (
	"context"
	"github.com/dop251/goja"

	"github.com/air-iot/logger"
)

const Key = "logger"

type Log struct {
	o options
}

type options struct {
	Project string
	Module  string
	Group   string
	Key     string
}

// Option 定义配置项
type Option func(*options)

func SetProject(s string) Option {
	return func(o *options) {
		o.Project = s
	}
}

func SetModule(s string) Option {
	return func(o *options) {
		o.Module = s
	}
}

func SetGroup(group string) Option {
	return func(o *options) {
		o.Group = group
	}
}

func SetKey(key string) Option {
	return func(o *options) {
		o.Key = key
	}
}

func SetVM(vm *goja.Runtime, opts ...Option) error {
	return vm.Set(Key, NewLogger(opts...))
}

func NewLogger(opts ...Option) *Log {
	var o options
	for _, opt := range opts {
		opt(&o)
	}
	return &Log{o: o}
}

func (l *Log) Debug(args ...any) {
	logger.WithContext(l.getCtx()).Debugln(args...)
}

func (l *Log) Info(args ...any) {
	logger.WithContext(l.getCtx()).Infoln(args...)
}

func (l *Log) Warn(args ...any) {
	logger.WithContext(l.getCtx()).Warnln(args...)
}

func (l *Log) Error(args ...any) {
	logger.WithContext(l.getCtx()).Errorln(args...)
}

func (l *Log) getCtx() context.Context {
	ctx := context.Background()
	if l.o.Key != "" {
		ctx = logger.NewExtraKeyContext(ctx, l.o.Key)
	}
	if l.o.Group != "" {
		ctx = logger.NewGroupContext(ctx, l.o.Group)
	}
	if l.o.Module != "" {
		ctx = logger.NewModuleContext(ctx, l.o.Module)
	}
	if l.o.Project != "" {
		ctx = logger.NewProjectContext(ctx, l.o.Project)
	}
	return ctx
}
