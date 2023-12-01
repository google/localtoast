export PATH := $(PATH):$(shell go env GOPATH)/bin

comma:= ,
empty:=
space:= $(empty) $(empty)
comma-separate = $(subst ${space},${comma},$1)

# Files under configs/defs/
CONFIG_DEFS = $(call comma-separate,$(wildcard configs/defs/*.textproto))
# Files under configs/reduced/*/
REDUCED_CONFIGS = $(call comma-separate,$(wildcard configs/reduced/*/*.textproto))
# Full configs go under configs/full/*/
FULL_CONFIGS = $(subst /reduced/,/full/,${REDUCED_CONFIGS})

localtoast: protos
	go build localtoast.go

localtoast_sql: protos
	cd localtoast_sql && go build localtoast_sql.go

test: protos
	go test ./...

configs: protos
	mkdir -p configs/full && cp -rf configs/reduced/* configs/full
	go build configs/genfullconfig/gen_full_config.go
	./gen_full_config --in=$(REDUCED_CONFIGS),$(CONFIG_DEFS) --out=$(FULL_CONFIGS) --omit-descriptions

protos:
	./build_protos.sh

clean:
	rm -rf scannerlib/proto/*_go_proto
	rm -rf scannerlib/proto/v1
	rm -f localtoast
	rm -f localtoast_sql/localtoast_sql
	rm -rf configs/full
	rm -f gen_full_config
