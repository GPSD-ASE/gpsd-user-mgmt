IMAGE_NAME = gpsd/gpsd-user-mgmt-svc
TAG ?= latest  # If no tag is provided, default to 'latest'

build-image:
	docker build -f docker/Dockerfile -t $(IMAGE_NAME):$(TAG) .

push-image:
	docker push $(IMAGE_NAME):$(TAG)

run-image:
	docker run -p 5500:5500 -it $(IMAGE_NAME):$(TAG)

clean-image:
	docker rmi $(docker images --filter "dangling=true" -q) -f

build:
	kubectl create namespace gpsd || true

setup:
	kubectl apply -f deployments/user-mgmt-deployment.yaml
	kubectl apply -f services/user-mgmt-service.yaml

all: build-image push-image build setup

clean:
	kubectl delete all --all -n gpsd || true

	kubectl delete namespace gpsd || true
	sleep 2