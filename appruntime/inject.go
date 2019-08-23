package appruntime

import (
	"log"
	"os"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	// Env variables
	Env env
	// Logger log entity
	Logger *zap.Logger
)

func init() {
	// Load env variables
	{
		if err := envconfig.Process("", &Env); err != nil {
			log.Fatal(err.Error())
		}
	}

	// Init logger
	{
		hook := lumberjack.Logger{
			Filename:   Env.LogPath + "/" + Env.AppName + ".log",
			MaxSize:    128,
			MaxBackups: 30,
			MaxAge:     7,
			Compress:   true,
		}

		encoderConfig := zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "line",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeName:     zapcore.FullNameEncoder,
		}

		// Setting Log Level
		atomicLevel := zap.NewAtomicLevel()
		switch strings.ToUpper(Env.LogLevel) {
		case "DEBUG":
			atomicLevel.SetLevel(zap.DebugLevel)
		case "INFO":
			atomicLevel.SetLevel(zap.InfoLevel)
		case "WARN":
			atomicLevel.SetLevel(zap.WarnLevel)
		case "ERROR":
			atomicLevel.SetLevel(zap.ErrorLevel)
		case "DPANIC":
			atomicLevel.SetLevel(zap.DPanicLevel)
		case "PANIC":
			atomicLevel.SetLevel(zap.PanicLevel)
		case "FATAL":
			atomicLevel.SetLevel(zap.FatalLevel)
		default:
			atomicLevel.SetLevel(zap.InfoLevel)
		}

		core := zapcore.NewCore(
			//zapcore.NewJSONEncoder(encoderConfig), // Json
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // console
			atomicLevel,
		)

		// Trace code
		caller := zap.AddCaller()
		// open code line
		development := zap.Development()
		// construct
		Logger = zap.New(core, development, caller)
	}
}
