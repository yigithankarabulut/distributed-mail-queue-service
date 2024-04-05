package constant

import "time"

const (
	Test_To      = "dobrainmusic@gmail.com"
	Test_Subject = "Test Mail"
	Test_Body    = "Hello, this is a test email!"
)

const (
	RedisMailQueueChannel = "mail_queue"
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
	ShutdownTimeout      = 5 * time.Second
	ServerReadTimeout    = 5 * time.Second
	ServerWriteTimeout   = 5 * time.Second
	ServerIdleTimeout    = 5 * time.Second
)
