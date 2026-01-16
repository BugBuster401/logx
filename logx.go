package logx

import "time"

type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	Trace(msg string, fields ...Field)
	Close() error
	With(fields ...Field) Logger
	WithGroup(name string) Logger
}

type Field struct {
	Key   string
	Value interface{}
}

func Err(err error) Field {
	return Field{
		Key:   "error",
		Value: err.Error(),
	}
}

func String(key, value string) Field {
	return Field{
		Key:   key,
		Value: value,
	}
}

func Int(key string, value int) Field {
	return Field{
		Key:   key,
		Value: value,
	}
}

func Duration(key string, value time.Duration) Field {
	return Field{
		Key:   key,
		Value: value,
	}
}

func Int64(key string, value int64) Field {
	return Field{
		Key:   key,
		Value: value,
	}
}
