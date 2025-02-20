IMAGE_NAME = gpsd/gpsd-user-mgmt
TAG ?= latest  # If no tag is provided, default to 'latest'

build-image:
	docker build -f docker/Dockerfile -t $(IMAGE_NAME):$(TAG) .

push-image:
	docker push $(IMAGE_NAME):$(TAG)

run-image:
	docker run -d -p 5500:5500 --name test_container $(IMAGE_NAME):$(TAG)
