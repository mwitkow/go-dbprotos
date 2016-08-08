# Copyright 2016 Michal Witkowski. All Rights Reserved.
# See LICENSE for licensing terms.

export PATH := ${GOPATH}/bin:${PATH}

install:
	@echo "Installing dbprotos to GOPATH"
	go install github.com/mwitkow/go-dbprotos/protoc-gen-dbprotos

regenerate_test:
	@echo "Regenerating test .proto files"
	(protoc  \
	--proto_path=${GOPATH}/src \
	--proto_path=${GOPATH}/src/github.com/gogo/protobuf/protobuf \
 	--proto_path=. \
	--go_out=. \
	--dbprotos_out=. \
	test/*.proto)


test: install regenerate_test
	@echo "Running tests"
	(go test -v ./...)

core_proto: dbprotos.proto
	@echo "Regenerating dbprotos.proto"
	(protoc \
	--proto_path=${GOPATH}/src \
	--proto_path=${GOPATH}/src/github.com/gogo/protobuf/protobuf \
	--proto_path=. \
	--gogo_out=Mgoogle/protobuf/descriptor.proto=github.com/gogo/protobuf/protoc-gen-gogo/descriptor:. \
	dbprotos.proto)
