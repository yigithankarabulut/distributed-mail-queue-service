package releaseinfo

const (
	prefix        = "/api/" + Version
	MailTaskQueue = prefix + "/task"
	User          = prefix + "/user"
)

const (
	RegisterUserApiPath = prefix + "/register"
	LoginUserApiPath    = prefix + "/login"
	GetUserApiPath      = User + "/:id"
)

const (
	EnqueueMailApiPath            = MailTaskQueue + "/enqueue"
	GetAllQueuedMailTasksApiPath  = MailTaskQueue + "/queue"
	GetAllFailedQueuedMailApiPath = MailTaskQueue + "/queue/fail"
)
