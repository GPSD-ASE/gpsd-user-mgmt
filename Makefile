TAG ?= 0.1.1  # If no tag is provided, default to 'latest'
NAMESPACE = gpsd
DEPLOYMENT = gpsd-user-mgmt
SERVICE_NAME = $(DEPLOYMENT)
IMAGE_NAME = $(NAMESPACE)/$(DEPLOYMENT)
CHART_DIRECTORY = helm
LOCAL_CHART_NAME = $(DEPLOYMENT)-$(TAG).tgz
LOCAL_INDEX_FILE = index.yaml
REMOTE_CHART_REPOSITORY = gpsd-ase.github.io

# Use `make develop` for local testing
develop: helm-uninstall build-image push-image helm

docker: build-image push-image

build-image:
	docker build -f docker/Dockerfile -t $(IMAGE_NAME):$(TAG) --platform linux/amd64 .

push-image:
	docker push $(IMAGE_NAME):$(TAG)

run-image:
	docker run -p 5500:5500 $(DEPLOYMENT)

clean-image:
	docker rmi $(docker images --filter "dangling=true" -q) -f

helm:
	helm upgrade --install demo ./helm --set image.tag=$(TAG) --namespace $(NAMESPACE)

helm-uninstall:
	helm uninstall demo -n $(NAMESPACE) 

clean:
	kubectl delete all --all -n $(NAMESPACE)  || true
	kubectl delete namespace $(NAMESPACE)  || true
	sleep 2

gh-pages-publish:
	@echo "Publishing Helm chart for $(SERVICE_NAME) to GitHub Pages..."
	rm -rf /tmp/$(LOCAL_CHART_NAME) /tmp/$(LOCAL_INDEX_FILE)
	helm package ./$(CHART_DIRECTORY) -d /tmp
	helm repo index /tmp --url https://$(REMOTE_CHART_REPOSITORY)/$(SERVICE_NAME)/ --merge /tmp/index.yaml
	git checkout gh-pages
	cp  /tmp/$(LOCAL_CHART_NAME) /tmp/$(LOCAL_INDEX_FILE) .
	git add .
	git commit -m "fix: commit to update GitHub Pages"
	git push origin gh-pages -f
	watch -k https://$(REMOTE_CHART_REPOSITORY)/$(SERVICE_NAME)/index.yaml