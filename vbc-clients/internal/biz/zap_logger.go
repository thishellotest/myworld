package biz

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"vbc/configs"
)

func NewZapLog() *zap.Logger {
	encoderConfig := zap.NewProductionEncoderConfig() //指定时间格式
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	fileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   configs.AppRuntimePath + "/logs/zap.log", //日志文件存放目录
		MaxSize:    1,                                        //文件大小限制,单位MB
		MaxBackups: 5,                                        //最大保留日志文件数量
		MaxAge:     30,                                       //日志文件保留天数
		Compress:   false,                                    //是否压缩处理
	})
	core := zapcore.NewCore(encoder, fileWriteSyncer, zapcore.DebugLevel) // //第三个及之后的参数为写入文件的日志级别,ErrorLevel模式只记录error级别的日志
	logger := zap.New(core, zap.AddCaller())                              //AddCaller()为显示文件名和行号
	return logger
}
