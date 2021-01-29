#!/usr/bin/make -f

# Copyright (c) 2020-present, The cloudquery authors
#
# This source code is licensed as defined by the LICENSE file found in the
# root directory of this source tree.
#
# SPDX-License-Identifier: (Apache-2.0 OR GPL-2.0-only)

INSTALL-DIR ?= /opt/cloudquery

all: deps test build

deps:
	@go mod download

test:
	@go test -v -race -cover ./...

build:
	@go build -ldflags="-s -w" -o . ./...

install:
	@cp cloudquery /usr/local/bin/cloudquery.ext ; \
	mkdir -p ${INSTALL-DIR}/config ; \
	cp extension/extension_config.json.sample ${INSTALL-DIR}/config/extension_config.json ; \
	for f in $$(find extension -name table_config.json); do \
		DIR=$$(echo $$f | cut -d / -f 2-3) ; \
		mkdir -p ${INSTALL-DIR}/$${DIR}   ; \
		cp $$f ${INSTALL-DIR}/$${DIR}/    ; \
	done

clean:
	@rm -f cloudquery

.PHONY: all
