package releaseinfo

const (
	version       = "/" + Version
	MailTaskQueue = version + "/task/mail"
	User          = version + "/user"
)

const (
	RegisterUserApiPath = version + "/register"
	LoginUserApiPath    = version + "/login"
	GetUserApiPath      = User + "/:id"
	UpdateUserApiPath   = User + "/update"
)

const (
	EnqueueMailApiPath            = MailTaskQueue + "/enqueue"
	GetAllQueuedMailTasksApiPath  = MailTaskQueue + "/queue"
	GetAllFailedQueuedMailApiPath = MailTaskQueue + "/queue/fail"
)

// SmsQueue service apis
const (
// EnqueueSmsApiPath            = smsQueue + "/Enqueue"
// TriggerWorkerApiPath         = smsQueue + "/TriggerWorker"
// ReadAllSmsQueueApiPath       = smsQueue + "/ReadAll"
// ReadAllSmsQueueFailedApiPath = smsQueue + "/ReadAll/Fail"
// Ping                         = smsQueue + "/Ping"
// Pong                         = "Pong"
//
// Metric  = smsQueue + "/Metrics"
// Swagger = smsQueue + "/Swagger"
)
