# 定义变量
ifndef GOPATH
	GOPATH := $(shell go env GOPATH)
endif

GOBIN=$(GOPATH)/bin
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) mod tidy

GOCTL=$(GOBIN)/goctl ## goctl

# 安装goctl代码生成工具
$(shell if [ ! -d $(GOCTL) ]; then \
	$(GOCMD) install github.com/zeromicro/go-zero/tools/goctl@latest; \
fi; \
)


gen:	## 生成所有模块代码
	# 合并rpc代码 & 生成sys-rpc代码
	$(GOCMD) run rpc/sys/proto/main.go
	$(GOCTL) rpc protoc rpc/sys/sys.proto --go_out=./rpc/sys/ --go-grpc_out=./rpc/sys/ --zrpc_out=./rpc/sys/ -m

model: ## 生成model代码
	$(GOCMD) run rpc/sys/db/mysql/generator.go