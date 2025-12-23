package gormUtils

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)
import (
	gormlogger "gorm.io/gorm/logger"
)

/*
  @author me
  @email  xxx@qq.com
  @create 2025/10/15 03:43
  @param
  @description  gorm结合zap使用
*/

type GormZap struct {
	ZapLogger *zap.Logger
	LogLevel  gormlogger.LogLevel
}

func NewGormZap(newLogger *zap.Logger, logLevel gormlogger.LogLevel) GormZap {
	return GormZap{
		ZapLogger: newLogger,
		LogLevel:  logLevel,
	}
}

// 和原来的gormlooger一样使用，设置log的级别后返回新的log对象
func (receiver GormZap) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return GormZap{
		ZapLogger: receiver.ZapLogger,
		LogLevel:  level,
	}
}

func (receiver GormZap) Info(ctx context.Context, str string, args ...interface{}) {
	//判定现在指定的级别，当我们设置的级别是warn(3)，小于4，那么肯定打印不出来info信息
	if receiver.LogLevel < gormlogger.Info {
		return
	}
	//这里可以用info也可以用debug，无所谓的
	receiver.ZapLogger.Sugar().Info(str, args)
}
func (receiver GormZap) Warn(ctx context.Context, str string, args ...interface{}) {
	//当我们设置的级别是error(2)，小于3(warn)，那么肯定打印不出来warn的信息
	if receiver.LogLevel < gormlogger.Warn {
		return
	}
	receiver.ZapLogger.Sugar().Warn(str, args)
}

func (receiver GormZap) Error(ctx context.Context, str string, args ...interface{}) {
	if receiver.LogLevel < gormlogger.Error {
		return
	}
	receiver.ZapLogger.Sugar().Error(str, args)
}

// Trace Trace是gorm自动调用的，在执行完一条语句后会调用Trace
func (receiver GormZap) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	//花费时间
	latency := time.Since(begin)
	switch {
	//如果发现错误，则打印error
	case err != nil:
		sql, rowsAffected := fc()
		fields := []zapcore.Field{
			zap.String("sql", sql),
			zap.Int64("rowsAffected", rowsAffected),
			zap.Duration("latency", latency),
			zap.Error(err),
		}
		receiver.ZapLogger.Error("", fields...)
	case receiver.LogLevel >= gormlogger.Warn:
		sql, rowsAffected := fc()
		fields := []zapcore.Field{
			zap.String("sql", sql),
			zap.Int64("rowsAffected", rowsAffected),
			zap.Duration("latency", latency),
		}
		receiver.ZapLogger.Info("", fields...)
	case receiver.LogLevel >= gormlogger.Info:
		sql, rowsAffected := fc()
		fields := []zapcore.Field{
			zap.String("sql", sql),
			zap.Int64("rowsAffected", rowsAffected),
			zap.Duration("latency", latency),
		}
		receiver.ZapLogger.Info("", fields...)
	}
}
