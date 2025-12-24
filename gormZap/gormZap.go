package gormZap

import (
	"context"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
	ZapLogger     *zap.Logger
	LogLevel      gormlogger.LogLevel
	SlowThreshold time.Duration
}

func NewGormZap(newLogger *zap.Logger, logLevel gormlogger.LogLevel, slowThreshold time.Duration) GormZap {
	return GormZap{
		ZapLogger:     newLogger,
		LogLevel:      logLevel,
		SlowThreshold: slowThreshold,
	}
}

// 和原来的gormlooger一样使用，设置log的级别后返回新的log对象
func (receiver GormZap) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return GormZap{
		ZapLogger:     receiver.ZapLogger,
		LogLevel:      level,
		SlowThreshold: receiver.SlowThreshold,
	}
}

func (receiver GormZap) Info(ctx context.Context, str string, args ...interface{}) {
	//gorm的logger等级过滤和zap相反
	//  info是最大的，error是最小的
	//  但是我们想要的是如果LogLevel设置的是info，那么debug无法打印
	//  如果设置的是warn，那么info和debug无法打印
	//  如果是设置的是error，那么warn,info,debug无法打印
	//这里如果receiver.LogLevel < gormlogger.Info，那么可能是debug
	//  那我们直接返回即可，否则就是打印info，因为这是info方法，是gorm调用的
	if receiver.LogLevel < gormlogger.Info {
		return
	}
	//这里可以用info也可以用debug，无所谓的
	receiver.ZapLogger.Sugar().Info(str, args)
}
func (receiver GormZap) Warn(ctx context.Context, str string, args ...interface{}) {
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
	//其他情况我们一般打印info
	//  但是我们的等级必须大于等于info才可以，如果是warn，那么就不应该打印info了
	//  同时你的慢sql时间要大于实际时间，这样才能打印info，不然就是慢日志了，那样就要打印warn了
	case receiver.LogLevel >= gormlogger.Info && receiver.SlowThreshold > latency:
		sql, rowsAffected := fc()
		fields := []zapcore.Field{
			zap.String("sql", sql),
			zap.Int64("rowsAffected", rowsAffected),
			zap.Duration("latency", latency),
		}
		receiver.ZapLogger.Info("", fields...)
	//然后是打印慢sql，要用warn，但是我们的等级必须要大于warn可行
	case receiver.LogLevel >= gormlogger.Warn && receiver.SlowThreshold <= latency:
		sql, rowsAffected := fc()
		fields := []zapcore.Field{
			zap.String("sql", sql),
			zap.Int64("rowsAffected", rowsAffected),
			zap.Duration("latency", latency),
		}
		receiver.ZapLogger.Warn("slowSql", fields...)
	}

}
