IMAGE_NAME = gpsd/gpsd-user-mgmt-svc
TAG ?= latest  # If no tag is provided, default to 'latest'

build-image:
	docker build -f docker/Dockerfile -t $(IMAGE_NAME):$(TAG) .

push-image:
	docker push $(IMAGE_NAME):$(TAG)

run-image:
    docker run -d -p 5500:5500 --name test_container $(IMAGE_NAME):$(TAG)

clean-image:
    docker stop test_container || true
    docker rm test_container || true

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