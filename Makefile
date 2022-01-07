export PATH := $(PATH):$(shell go env GOPATH)/bin

localtoast:
	go install google.golang.org/protobuf/cmd/protoc-gen-go
	protoc -I=. --go_out=. scannerlib/proto/*.proto
	mv github.com/google/localtoast/scannerlib/proto/* scannerlib/proto/
	rm -r github.com
	go build localtoast.go

clean:
	rm -rf scannerlib/proto/*_go_proto
	rm localtoast
