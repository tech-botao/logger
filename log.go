// 创建一个Zap Log
package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strings"
	"sync"
)

type Logger struct {
	Debug *zap.Logger
	Err   *zap.Logger
}

var Log Logger
var once sync.Once

func init() {
	once.Do(func() {
		file := os.Getenv("LOG_FILE")
		Log = Logger{
			Debug: NewLogger(file),
			Err:   NewErrorLog(file),
		}
	})
}

func Any(v interface{}) zap.Field {
	return zap.Any("data", &v)
}

func createLogger(path string) lumberjack.Logger {
	return lumberjack.Logger{
		Filename:   path, // 日志文件路径
		MaxSize:    10,   // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 30,   // 日志文件最多保存多少个备份
		MaxAge:     7,    // 文件最多保存多少天
		Compress:   true, // 是否压缩
	}
}

func createEncoding() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		MessageKey:     "m",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
}

// NewLogger Debug, Info的Logger
func NewLogger(path string) *zap.Logger {
	var encoderConfig = createEncoding()
	var core zapcore.Core

	// stdoutに出力
	if len(path) == 0 {
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)),
			createAtomicLevel(zap.DebugLevel),
		)
	} else {
		hook := createLogger(path)
		core = zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(&hook)),
			createAtomicLevel(zap.InfoLevel),
		)
	}

	caller := zap.AddCaller()
	callerSkip := zap.AddCallerSkip(1)
	return zap.New(core, caller, callerSkip)
}

// NewErrorLog Warning, Error的日志句柄
func NewErrorLog(path string) *zap.Logger {
	var core zapcore.Core
	var encoderConfig = createEncoding()
	encoderConfig.StacktraceKey = "trace"

	if len(path) == 0 {
		core = zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)),
			createAtomicLevel(zap.WarnLevel),
		)
	} else {
		path = strings.Replace(path, ".log", "-err.log", 1)
		hook := createLogger(path)
		core = zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(&hook)),
			createAtomicLevel(zap.WarnLevel),
		)
	}

	return zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2), zap.AddStacktrace(zap.WarnLevel))
}

func createAtomicLevel(level zapcore.Level) zap.AtomicLevel {
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(level)
	return atomicLevel
}

func Debug(msg string, value interface{}) {
	Log.Debug.Debug(msg, Any(value))
}
func Info(msg string, value interface{}) {
	Log.Debug.Info(msg, Any(value))
}
func Warn(msg string, value interface{}) {
	Log.Err.Warn(msg, Any(value))
}
func Error(msg string, value interface{}) {
	Log.Err.Error(msg, Any(value))
}
func Panic(msg string, value interface{}) {
	Log.Err.Panic(msg, Any(value))
}
func Fatal(msg string, value interface{}) {
	Log.Err.Fatal(msg, Any(value))
}
