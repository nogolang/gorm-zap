

go get -u github.com/nogolang/gorm-zap



```
config := &gorm.Config{
   Logger: gormZap.NewGormZap(NewZapConfig(), gormLogger.Info, time.Second*3),
}

func NewZapConfig() *zap.Logger {
	config := zap.NewDevelopmentEncoderConfig()
	config.EncodeTime = func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(t.Format(time.DateTime))
	}
	return zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(config),
		zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(os.Stdout),
		),
		zap.NewAtomicLevelAt(zap.InfoLevel),
	))
}
```

