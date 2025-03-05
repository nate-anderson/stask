INSTALL_DIR ?= /usr/bin
BUILD_TARGET ?= ./stask

.PHONY: build
build: go build -o ${BUILD_TARGET} ./cmd/cli

.PHONY: install
install: cp ${BUILD_TARGET} ${INSTALL_DIR}