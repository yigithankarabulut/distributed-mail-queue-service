package constant

import "time"

const (
	RedisMailQueueChannel = "mail_queue"
	QueueConsumerCount    = 10
	WorkerCount           = 10
	MaxTryCount           = 3
)

const (
	StatusQueued = iota
	StatusProcessing
	StatusSuccess
	StatusFailed
	StatusCancelled
	StatusScheduled
)

const (
	ContentType    = "Content-Type"
	Authorization  = "Authorization"
	AllowedOrigins = "*"
)

const (
	ContextCancelTimeout = 5 * time.Second
	ShutdownTimeout      = 2 * time.Second
	ServerReadTimeout    = 5 * time.Second
	ServerWriteTimeout   = 5 * time.Second
	ServerIdleTimeout    = 5 * time.Second
	TaskCancelTimeout    = 5 * time.Second
)
