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

all: deps build ## 默认的构建目标


clean: ## 清理目标
	#$(GOCLEAN)
	rm -rf target


deps: ## 安装依赖目标
	@export GOPROXY=https://goproxy.cn,direct
	$(GOGET) -v

copy_config:
	mkdir -p target/sys-rpc && cp rpc/sys/etc/sys.yaml target/sys-rpc/sys-rpc.yaml
	mkdir -p target/admin-api && cp api/admin/etc/admin-api.yaml target/admin-api/admin-api.yaml

build: copy_config ## 构建目标
	$(GOBUILD) -o target/sys-rpc/sys-rpc -v ./rpc/sys/sys.go
	$(GOBUILD) -o target/admin-api/admin-api -v ./api/admin/admin.go


start: ## 运行目标
	nohup ./target/sys-rpc/sys-rpc -f ./target/sys-rpc/sys-rpc.yaml  > /dev/null 2>&1 &
	nohup ./target/admin-api/admin-api -f ./target/admin-api/admin-api.yaml > /dev/null 2>&1 &


stop: ## 停止目标
	-pkill -f admin-api
	-pkill -f sys-rpc
	@for i in 5 4 3 2 1; do\
      echo -n "stop $$i";\
      sleep 1; \
      echo " "; \
    done


restart: stop start ## 重启项目

.DEFAULT_GOAL := all ## 默认构建目标是

format: ## 格式化代码
	$(GOCTL) api format --dir api/admin/doc/api
	$(GOCTL) api format --dir api/front/doc/api
	$(GOCTL) api format --dir api/web/doc/api

gen:	## 生成所有模块代码
	$(GOCTL) api go -api ./api/admin/doc/api/admin.api -dir ./api/admin/ -home ./script/.goctl -style go_zero

	# 合并rpc代码 & 生成sys-rpc代码
	$(GOCMD) run rpc/sys/proto/main.go
	$(GOCTL) rpc protoc rpc/sys/sys.proto --go_out=./rpc/sys/ --go-grpc_out=./rpc/sys/ --zrpc_out=./rpc/sys/ -home ./script/.goctl -m

model: ## 生成model代码
	$(GOCMD) run rpc/sys/db/mysql/generator.go

help: ## show help message
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m\033[0m\n"} /^[$$()% 0-9a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
