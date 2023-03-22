SERVICE_NAME := "backend"
CURRENT_DIR = $(shell pwd)
GOPATH = $(shell echo ${HOME})
ifdef BUILD_VERSION
	VERSION := "-$(BUILD_VERSION)"
else
	VERSION := ""
endif

.SILENT:

deps:
	cd ${CURRENT_DIR}/app && go mod download

clean:
	rm -rf ./.bin/${SERVICE_NAME}


build: clean deps
	cd ${CURRENT_DIR}/app && CGO_ENABLED=0 GOOS=linux go build -o ./.bin/${SERVICE_NAME}${VERSION} ./cmd/${SERVICE_NAME}/main.go
