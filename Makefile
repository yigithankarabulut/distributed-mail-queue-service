POSTGRES_YAML := ./deployment/postgresql/deployment.yml
POSTGRES_SERVICE_YAML := ./deployment/postgresql/service.yml

REDIS_YAML := ./deployment/redis/deployment.yml
REDIS_SERVICE_YAML := ./deployment/redis/service.yml

GO_API_YAML := ./deployment/app/deployment.yml
GO_API_SERVICE_YAML := ./deployment/app/service.yml

HPA_YAML := ./deployment/hpa/hpa.yml


all: cluster apply

cluster:
	minikube start
	eval $(minikube docker-env)
	minikube addons enable metrics-server
	minikube addons enable dashboard

clean:
	minikube stop
	minikube delete

re: delete apply

apply:
	kubectl apply -f $(POSTGRES_YAML)
	kubectl apply -f $(POSTGRES_SERVICE_YAML)
	kubectl apply -f $(REDIS_YAML)
	kubectl apply -f $(REDIS_SERVICE_YAML)
	kubectl wait --for=condition=available deployment/postgres --timeout=300s
	kubectl wait --for=condition=available deployment/redis --timeout=300s
	kubectl apply -f $(GO_API_YAML)
	kubectl apply -f $(GO_API_SERVICE_YAML)
	kubectl wait --for=condition=available deployment/go-api --timeout=300s
	eval $(minikube docker-env)
	minikube addons enable metrics-server
	kubectl apply -f $(HPA_YAML)

delete:
	kubectl delete -f $(POSTGRES_YAML)
	kubectl delete -f $(POSTGRES_SERVICE_YAML)
	kubectl delete -f $(REDIS_YAML)
	kubectl delete -f $(REDIS_SERVICE_YAML)
	kubectl delete -f $(GO_API_YAML)
	kubectl delete -f $(GO_API_SERVICE_YAML)
	kubectl delete -f $(HPA_YAML)

PHONY: all cluster clean re apply delete