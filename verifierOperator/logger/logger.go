package logger

import (
	"io"
	"os"
	"sync"

	"github.com/Cedro-Finance/CedroAVS/verifierOperator/config"
	sdklogging "github.com/Layr-Labs/eigensdk-go/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *zap.Logger
var once sync.Once

type ZapLogger struct {
	logger *zap.Logger
}

var _ sdklogging.Logger = (*ZapLogger)(nil)

func NewZapLogger(env sdklogging.LogLevel) (sdklogging.Logger, error) {
	config := zap.NewProductionConfig()

	if env == sdklogging.Development {
		config = zap.NewDevelopmentConfig()
	}

	config.DisableStacktrace = true
	config.Encoding = "console"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	if env == sdklogging.Development {
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	logger, err := config.Build(zap.AddCallerSkip(1))

	if err != nil {
		panic(err)
	}

	return &ZapLogger{
		logger: logger,
	}, nil

}

func init() {
	var err error
	logger, err = zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
}

func GetLogger() *zap.Logger {
	return logger
}

func InitLogging(mode string, cfg *config.LoggerConfig) {
	once.Do(func() {
		initLog(mode, cfg)
	})
}

func initLog(mode string, cfg *config.LoggerConfig) {
	logCfg := getLogConfig(mode)
	lumber := &lumberjack.Logger{
		Filename:   cfg.OutputPath,
		MaxSize:    logCfg.maxSize,
		MaxBackups: logCfg.maxBackup,
		MaxAge:     logCfg.maxAge,
		Compress:   logCfg.compress,
	}
	err := lumber.Rotate()
	if err != nil {
		panic(err)
	}
	var mW io.Writer // multi-writer
	if mode == config.Development {
		mW = io.MultiWriter(lumber, os.Stdout)
	} else {
		mW = io.MultiWriter(lumber)
	}

	w := zapcore.AddSync(mW)

	var (
		level  = zap.NewAtomicLevel()
		zapCfg zapcore.EncoderConfig
	)
	if mode == config.Production {
		level.SetLevel(zap.DebugLevel)
		zapCfg = zap.NewProductionEncoderConfig()
	} else {
		level.SetLevel(zap.DebugLevel)
		zapCfg = zap.NewDevelopmentEncoderConfig()
	}

	zapCfg.TimeKey = "ts"
	zapCfg.LevelKey = "level"
	zapCfg.CallerKey = "caller"
	zapCfg.StacktraceKey = "stacktrace"
	zapCfg.MessageKey = "msg"

	enc := zapcore.NewConsoleEncoder(zapCfg)
	if cfg.Encoding == "json" {
		enc = zapcore.NewJSONEncoder(zapCfg)
	}

	core := zapcore.NewCore(
		enc,
		w,
		level,
	)

	if mode == config.Production {
		logger = zap.New(core)
	} else {
		logger = zap.New(core, zap.Development(), zap.AddCaller())
	}
}

func getLogConfig(mode string) struct {
	maxSize   int
	maxBackup int
	maxAge    int
	compress  bool
} {

	if mode == config.Production {
		return struct {
			maxSize   int
			maxBackup int
			maxAge    int
			compress  bool
		}{
			maxSize:   100,
			maxBackup: 100,
			maxAge:    180,
			compress:  true,
		}
	}

	return struct {
		maxSize   int
		maxBackup int
		maxAge    int
		compress  bool
	}{
		maxSize:   10,
		maxBackup: 2,
		maxAge:    2,
		compress:  false,
	}
}

func (z *ZapLogger) Debug(msg string, tags ...any) {
	z.logger.Sugar().Debugw(msg, tags...)
}

func (z *ZapLogger) Info(msg string, tags ...any) {
	z.logger.Sugar().Infow(msg, tags...)
}

func (z *ZapLogger) Warn(msg string, tags ...any) {
	z.logger.Sugar().Warnw(msg, tags...)
}

func (z *ZapLogger) Error(msg string, tags ...any) {
	z.logger.Sugar().Errorw(msg, tags...)
}

func (z *ZapLogger) Fatal(msg string, tags ...any) {
	z.logger.Sugar().Fatalw(msg, tags...)
}

func (z *ZapLogger) Debugf(template string, args ...interface{}) {
	z.logger.Sugar().Debugf(template, args...)
}

func (z *ZapLogger) Infof(template string, args ...interface{}) {
	z.logger.Sugar().Infof(template, args...)
}

func (z *ZapLogger) Warnf(template string, args ...interface{}) {
	z.logger.Sugar().Warnf(template, args...)
}

func (z *ZapLogger) Errorf(template string, args ...interface{}) {
	z.logger.Sugar().Errorf(template, args...)
}

func (z *ZapLogger) Fatalf(template string, args ...interface{}) {
	z.logger.Sugar().Fatalf(template, args...)
}

func (z *ZapLogger) With(tags ...any) sdklogging.Logger {
	return &ZapLogger{
		logger: z.logger.Sugar().With(tags...).Desugar(),
	}
}
