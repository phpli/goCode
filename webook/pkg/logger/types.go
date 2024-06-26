package logger

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

func LoggerExample() {
	var l Logger
	phone := "12316*****12"
	l.Info("用户未注册，，手机号是%s", phone)
}

type LoggerV1 interface {
	Debug(msg string, args ...Field)
	Info(msg string, args ...Field)
	Warn(msg string, args ...Field)
	Error(msg string, args ...Field)
}

type Field struct {
	Key   string
	Value any
}

func LoggerV1Example() {
	var l LoggerV1
	phone := "12316*****12"
	l.Info("用户未注册", Field{
		Key:   "phone",
		Value: phone,
	})
}

type LoggerV2 interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

func LoggerV2Example() {
	var l LoggerV2
	phone := "12316*****12"
	l.Info("用户未注册", Field{
		Key:   "phone",
		Value: phone,
	})
}
