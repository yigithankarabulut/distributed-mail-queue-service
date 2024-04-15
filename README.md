# Distributed Mail Queue Service
This project represents an email service. Users can register, login and retrieve user information. When registering, the user provides email SMTP information (host, port, password, etc.). The user can create an email sending task by requesting the task/enqueue endpoint with the "Bearer" token. This task is saved in the database and also added to the Redis queue. Redis consumers receive these tasks and send them through the channel to other workers.

## Requirements

-  Docker, minikube, kubectl installed on your system.
- Make installed on your system.
- You need to change the values of the environment variables in the deployment files.

### Deployments
The application runs on Kubernetes. K8s ymls are located in the /deployment directory. Horizontal Pod Autoscaling feature of K8s is used. The pod is given a certain resource value to use and when a certain percentage of this value is reached, a new instance is created and the load is balanced.

### Make Commands
You can use Makefile to compile and run project files. Example commands:

* You should run make cluster -> make apply -> make hpa in order of.
```bash
make all      # Call cluster and apply command.
make cluster  # Create a minikube cluster. Adds metric server and dashboard.
make clean    # Stop and delete cluster.
make re       # All .yml's delete and apply via kubectl.
make apply    # Creates all Kubernetes objects. It waits until Postgresql and Redis deployments are ready. After that start the deployment and service of our Go application.
make hpa      # Start Horizontal Pod Autoscaling deployments.
make delete   # Delete created k8s objects (deployments, services) via kubectl.
```

## Usage

Users can register, login and retrieve user information. You can add a task to the queue, retrieve tasks in the queue and retrieve tasks that have failed due to an error.
### Endpoints
```http
GET     /healthz/live
GET     /healthz/ping

POST    /api/v1/register
POST    /api/v1/login
GET     /api/v1/user/:id

POST    /api/v1/task/enqueue
GET     /api/v1/task/queue
GET     /api/v1/task/queue/fail
```
The json body required to register is as follows.
```json
{
  "email": 		"example@example.com",
  "password": 		"password123",
  "smtp_host": 		"smtp.example.com",
  "smtp_port": 		587,
  "smtp-username": 	"smtp_username@smtphost.com",
  "smtp-password": 	"smtp_password"
}
```
#### The json body required for login in is as follows.
```json
{
  "email": 	"example@example.com",
  "password": 	"password123",
}
```
#### Once you have logged in, you must send a request to the task endpoints with the Bearer token in the login response body.
The json body required to add a task to the queue is as follows.
```json
{
  "recipient_email": 	"recipient@example.com",
  "subject": 		"Example Subject",
  "body": 		"Example Body Content",
  "scheduled_at": 	"2024-04-15T12:00:00"
}
```

## Operation
* There are two worker count values in pkg/constant when the system starts.
* These values are the values of the workers that will run concurrently and consume our queue and the workers that will handle the tasks taken from the queue.
```go
const (
	QueueConsumerCount    = N
	WorkerCount           = N
)
```
* For the communication of our Queue consumers and Workers, there is a channel of Task model type and the consumed tasks are given to the worker through this channel.
* When users submit a task, the first step is to create a record for the task in postgresql and publish the task to redis.
* Consumers receive the task from the queue. Then they unmarshal the task and send it to the channel.
* Our workers that receive the task from the channel process the task, that is, they send mail. 
* Status is updated in Postgres according to the result of the task.
* If the task has failed, first the value of the TryCount filed is compared with the MaxTryCount value in pkg/constant.
* If the value is not exceeded, the task is sent to the queue again.
* Since we change the status of the failed task and update it in postgres and then send it to the queue again, it goes through the same pipeline and when MaxTryCount is exceeded, it is not sent to the queue and its status is updated as Cancelled.

This pipeline uses cron service to process leaked tasks that need to be processed but are not. Cron service running a method called FindUnprocessedTasksAndEnqueue every 5 minutes. This method takes tasks that are StatusQueued in postgres and hasn't been processed for the last 5 minutes and sends them to the queue.

