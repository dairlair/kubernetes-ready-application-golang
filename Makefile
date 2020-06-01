BOILERPLATE_PROJECT=github.com/stepsisters/kgb
BOILERPLATE_ROOT?=vendor/$BOILERPLATE_PROJECT

RELEASE?=0.0.1
COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

GOOS?=linux
GOARCH?=amd64

# Docker settings
DOCKER_REGISTRY?=docker.io
DOCKER_IMAGE?=${AUTHOR}/${APP}
DOCKER_REGISTRY_IMAGE=${DOCKER_REGISTRY}/${DOCKER_IMAGE}

# Application runtime variables
PORT?=80
PROBES_PORT?=81

# Help variables
HELM_CHART_PATH=./helm

# This entry point provides functionality to check that required variable is set.
guard-%:
	@ if [ "${${*}}" = "" ]; then \
		echo "Environment variable $* not set"; \
		exit 1; \
	fi

.PHONY: clean
clean: guard-APP guard-APP_NAME
	rm -f ${APP}

.PHONY: build
build: clean
	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -ldflags '-s -w -X "${BOILERPLATE_PROJECT}/pkg/version.ApplicationName=${APP_NAME}" -X "${BOILERPLATE_PROJECT}/pkg/version.Release=${RELEASE}" -X "${BOILERPLATE_PROJECT}/pkg/version.Commit=${COMMIT}" -X "${BOILERPLATE_PROJECT}/pkg/version.BuildTime=${BUILD_TIME}"' -o ${APP}

# Just a target-specific variable, we always build a Linux binary to create the docker image.
# You can run `make image` on macOS or Windows without GOOS override.
.PHONY: image
image: GOOS=linux
image: build
	cp $(BOILERPLATE_ROOT)/Dockerfile ./Dockerfile
	docker build --build-arg APP=${APP} -t $(DOCKER_IMAGE):$(RELEASE) .
	rm -f ./Dockerfile

.PHONY: publish
publish: image
	docker push $(DOCKER_REGISTRY_IMAGE):$(RELEASE)

.PHONY: run
run: image
	docker stop $(DOCKER_IMAGE):$(RELEASE) || true && docker rm $(DOCKER_IMAGE):$(RELEASE) || true
	docker run --name ${APP} -p ${PORT}:${PORT} -p ${PROBES_PORT}:${PROBES_PORT} --rm -e "PORT=${PORT}" -e "PROBES_PORT=${PROBES_PORT}" $(DOCKER_IMAGE):$(RELEASE)

.PHONY: test
test:
	go test -v ./...

.PHONY: mocks
mocks:
	@echo " > Generate mocks..."
	mockery -all -keeptree -dir pkg -output ./mocks/pkg

.PHONY: helm
helm: guard-HELM_CHART_PATH publish
	@echo " > Run helm install --dry-run --debug..."
	helm install --debug --dry-run ${APP} ${HELM_CHART_PATH} --set Image="$(DOCKER_REGISTRY_IMAGE):$(RELEASE)"

.PHONY: deploy
deploy: guard-HELM_CHART_PATH publish
	@echo " > Run helm install"
	helm upgrade --install ${APP} ${HELM_CHART_PATH} --set Image="$(DOCKER_REGISTRY_IMAGE):$(RELEASE)"