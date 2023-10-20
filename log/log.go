package log

import (
	"context"

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

func NewLogger(opts ...Option) *Log {
	var o options
	for _, opt := range opts {
		opt(&o)
	}
	return &Log{o: o}
}

func (l *Log) Debug(args ...any) {
	ctx := logger.NewExtraKeyContext(context.Background(), l.o.Key)
	ctx = logger.NewGroupContext(ctx, l.o.Group)
	ctx = logger.NewModuleContext(ctx, l.o.Module)
	ctx = logger.NewProjectContext(ctx, l.o.Project)
	logger.WithContext(ctx).Debugln(args...)
}

func (l *Log) Info(args ...any) {
	ctx := logger.NewExtraKeyContext(context.Background(), l.o.Key)
	ctx = logger.NewGroupContext(ctx, l.o.Group)
	ctx = logger.NewModuleContext(ctx, l.o.Module)
	ctx = logger.NewProjectContext(ctx, l.o.Project)
	logger.WithContext(ctx).Infoln(args...)
}

func (l *Log) Warn(args ...any) {
	ctx := logger.NewExtraKeyContext(context.Background(), l.o.Key)
	ctx = logger.NewGroupContext(ctx, l.o.Group)
	ctx = logger.NewModuleContext(ctx, l.o.Module)
	ctx = logger.NewProjectContext(ctx, l.o.Project)
	logger.WithContext(ctx).Warnln(args...)
}

func (l *Log) Error(args ...any) {
	ctx := logger.NewExtraKeyContext(context.Background(), l.o.Key)
	ctx = logger.NewGroupContext(ctx, l.o.Group)
	ctx = logger.NewModuleContext(ctx, l.o.Module)
	ctx = logger.NewProjectContext(ctx, l.o.Project)
	logger.WithContext(ctx).Errorln(args...)
}
