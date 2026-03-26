package logger

// Option logger 构建选项（函数选项模式）
type Option func(*options)

type options struct {
	level      string // debug | info | warn | error
	format     string // console | json
	output     string // stdout | 文件路径
	maxSize    int    // MB，文件轮转大小
	maxBackups int    // 最大保留备份数
	maxAge     int    // 天，最大保留天数
	compress   bool   // 是否压缩归档
	caller     bool   // 是否记录调用位置
	stacktrace bool   // 是否在 error 级别记录堆栈
}

func defaultOptions() *options {
	return &options{
		level:      "info",
		format:     "console",
		output:     "stdout",
		maxSize:    100,
		maxBackups: 10,
		maxAge:     30,
		compress:   true,
		caller:     true,
		stacktrace: true,
	}
}

// WithLevel 设置日志级别 (debug|info|warn|error)
func WithLevel(level string) Option {
	return func(o *options) { o.level = level }
}

// WithFormat 设置日志格式 (console|json)
func WithFormat(format string) Option {
	return func(o *options) { o.format = format }
}

// WithOutput 设置输出目标 ("stdout" 或文件路径)
func WithOutput(output string) Option {
	return func(o *options) { o.output = output }
}

// WithMaxSize 设置日志文件最大体积（MB）
func WithMaxSize(mb int) Option {
	return func(o *options) { o.maxSize = mb }
}

// WithMaxBackups 设置最大保留备份文件数
func WithMaxBackups(n int) Option {
	return func(o *options) { o.maxBackups = n }
}

// WithMaxAge 设置日志文件保留天数
func WithMaxAge(days int) Option {
	return func(o *options) { o.maxAge = days }
}

// WithCompress 是否压缩轮转的日志文件
func WithCompress(compress bool) Option {
	return func(o *options) { o.compress = compress }
}

// WithCaller 是否在日志中记录调用文件和行号
func WithCaller(caller bool) Option {
	return func(o *options) { o.caller = caller }
}

// WithStacktrace 是否在 error 级别自动记录堆栈
func WithStacktrace(enable bool) Option {
	return func(o *options) { o.stacktrace = enable }
}
