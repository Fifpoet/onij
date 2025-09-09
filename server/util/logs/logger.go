package logs

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config 日志配置
type Config struct {
	Dir           string // 日志存储目录
	Prefix        string // 日志文件前缀
	MaxSize       int    // 单个日志文件最大尺寸(MB)
	MaxAge        int    // 日志文件最大保存天数
	MaxBackups    int    // 最大备份文件数
	Compress      bool   // 是否压缩备份文件
	Level         string // 日志级别: debug, info, warn, error
	Development   bool   // 是否开发模式(开发模式下日志更详细)
	OutputConsole bool   // 是否同时输出到控制台
}

// Logger 日志实例
type Logger struct {
	*zap.Logger
	config Config
}

// 全局默认日志实例
var defaultLogger *Logger

// 初始化默认配置
func init() {
	var err error
	defaultLogger, err = New(Config{
		Dir:           "logs",
		Prefix:        "app",
		MaxSize:       100,  // 100MB
		MaxAge:        7,    // 保存7天
		MaxBackups:    30,   // 最多30个备份
		Compress:      true, // 压缩备份
		Level:         "info",
		Development:   false,
		OutputConsole: true,
	})
	if err != nil {
		panic("初始化默认日志失败: " + err.Error())
	}
}

// New 创建新的日志实例
func New(config Config) (*Logger, error) {
	// 确保日志目录存在
	if err := os.MkdirAll(config.Dir, 0755); err != nil {
		return nil, err
	}

	// 设置日志级别
	level := zap.NewAtomicLevel()
	switch config.Level {
	case "debug":
		level.SetLevel(zapcore.DebugLevel)
	case "warn":
		level.SetLevel(zapcore.WarnLevel)
	case "error":
		level.SetLevel(zapcore.ErrorLevel)
	default:
		level.SetLevel(zapcore.InfoLevel)
	}

	// 日志文件名格式: {prefix}-YYYYMMDD.log
	logFileName := filepath.Join(config.Dir, config.Prefix) + "-%Y%m%d.log"

	// 配置Lumberjack进行日志轮转和清理
	writeSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logFileName, // 注意: Lumberjack会自动处理时间格式化
		MaxSize:    config.MaxSize,
		MaxAge:     config.MaxAge,
		MaxBackups: config.MaxBackups,
		Compress:   config.Compress,
		LocalTime:  true, // 使用本地时间
	})

	// 配置Zap核心
	encoderConfig := zap.NewProductionEncoderConfig()
	// 自定义时间格式
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}
	// 自定义日志级别显示
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	var encoder zapcore.Encoder
	if config.Development {
		// 开发模式使用更易读的控制台格式
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		// 生产模式使用JSON格式
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	// 配置输出目标
	var cores []zapcore.Core
	cores = append(cores, zapcore.NewCore(encoder, writeSyncer, level))

	// 如果需要同时输出到控制台
	if config.OutputConsole {
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		cores = append(cores, zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), level))
	}

	// 创建Zap日志实例
	zapLogger := zap.New(zapcore.NewTee(cores...),
		zap.AddCaller(),                       // 显示调用者信息
		zap.AddCallerSkip(1),                  // 跳过当前层(因为我们封装了一层)
		zap.AddStacktrace(zapcore.ErrorLevel), // 错误级别以上添加堆栈跟踪
	)

	return &Logger{
		Logger: zapLogger,
		config: config,
	}, nil
}

// 以下是便捷的日志方法封装

// Debug 输出debug级别日志
func Debug(msg string, args ...any) {
	defaultLogger.Sugar().Debugf(msg, args...)
}

// Info 输出info级别日志
func Info(msg string, args ...any) {
	defaultLogger.Sugar().Infof(msg, args...)
}

// Warn 输出warn级别日志
func Warn(msg string, args ...any) {
	defaultLogger.Sugar().Warnf(msg, args...)
}

// Error 输出error级别日志
func Error(msg string, args ...any) {
	defaultLogger.Sugar().Errorf(msg, args...)
}

// Sync 刷新日志缓冲区
func Sync() error {
	return defaultLogger.Sync()
}

// GetDefault 获取默认日志实例
func GetDefault() *Logger {
	return defaultLogger
}
