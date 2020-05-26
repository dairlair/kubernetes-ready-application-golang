BOILERPLATE_PROJECT=github.com/dairlair/kubernetes-ready-application-golang
BOILERPLATE_ROOT?=vendor/$BOILERPLATE_PROJECT

RELEASE?=0.0.1
COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

GOOS?=linux
GOARCH?=amd64

# Application runtime variables
PORT?=80

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
image: GOOS=linux
image: build
	cp $(BOILERPLATE_ROOT)/Dockerfile ./Dockerfile
	docker build --build-arg APP=${APP} -t $(APP):$(RELEASE) .
	rm -f ./Dockerfile

run: image
	docker stop $(APP):$(RELEASE) || true && docker rm $(APP):$(RELEASE) || true
	docker run --name ${APP} -p ${PORT}:${PORT} --rm -e "PORT=${PORT}" $(APP):$(RELEASE)

.PHONY: test
test:
	go test -v ./...

.PHONY: mocks
mocks:
	@echo " > Generate mocks..."
	mockery -all -keeptree -dir pkg -output ./mocks/pkg