/**
log 包装了 logrus 来进行日志输出

默认情况下 log 出来的内容包含为以下的 {level logrus.Level, msg string, time int64}
level: 对应对了调用的 log 函数 trace/debug/info/warning/error/fatal/panic
msg: 对引用 log 出来的内容，如果是 obj 会被转为 json
time: 日志时间戳
如果希望增加其他字段，通过 log.WithField 来添加对应的 field，其中 log.WithError 可以添加 error 字段
可以通过添加特定的 field 来与后续的 hook 进行联动
*/
package log
