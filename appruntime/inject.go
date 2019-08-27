package appruntime

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/Invisibi-nd/slack-bot/tool/db"

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
	// DB storage entity
	DB db.Storage
)

func init() {
	// Load env variables
	{
		log.Println("[Init] load env variables")
		if err := envconfig.Process("", &Env); err != nil {
			log.Fatalf("[Init] load env error = [%v]", err)
		}
	}

	// Init logger
	{
		log.Println("[Init] setting logger driver")
		hook := lumberjack.Logger{
			Filename:   Env.LogPath + "/" + Env.AppName + ".log",
			MaxSize:    128,
			MaxBackups: 30,
			MaxAge:     7,
			Compress:   true,
		}

		encoderConfig := zapcore.EncoderConfig{
			TimeKey:       "time",
			LevelKey:      "level",
			NameKey:       "logger",
			CallerKey:     "line",
			MessageKey:    "msg",
			StacktraceKey: "stacktrace",
			LineEnding:    zapcore.DefaultLineEnding,
			EncodeLevel:   zapcore.CapitalColorLevelEncoder,
			EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(t.Format("2006/01/02 15:04:05"))
			},
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

	// DB
	{
		log.Println("[Init] connect db driver")
		db, err := db.New(db.Driver(Env.DBDriver), &db.Connection{
			DBName: Env.DBName,
			Host:   Env.DBHost,
			Port:   Env.DBPort,
			User:   Env.DBPort,
			Pass:   Env.DBPass,
		})
		if err != nil {
			log.Fatalf("[Init] connect db error = [%v]", err)
		}
		DB = db
	}
}
