### DDD 实践websocket 

目录解释: 

```shell
├── application                 # 应用层 
│   ├── converter         # transformer
│   └── websocket         # 应用入口
├── asset                       # 不相关的静态资源
│   ├── index.html
│   └── index.js
├── cmd                         # main 包
│   ├── main.go 
│   └── test
├── config                      # 配置
│   └── config.go
├── config.yaml
├── design                      # 架构设计资源
│   └── design.puml
├── domain                      # 领域层，复杂的业务逻辑都在此层实现、不关心具体存储，不关心具体事件
│   ├── service
│   └── websocket
├── endpoints                   # 入口、接口、事件层
│   ├── consumer
│   ├── http-server.go
│   └── network
├── go.mod
├── go.sum
├── infrastructure             # 持久层 
│   ├── gateway_impl
│   └── nsq
├── pkg
│   └── utils
├── readme.md
└── req.http              http 测试文件
```


依赖关系如下：

```shell

endpoints -> application -> domain -> gateway 
gateway_impl -> infrastructure

```

```shell
endpoints: 具体的http、事件封装、解析参数，转成dto，交给app处理

application: 应用层，调用不同领域实现业务

domain: 只认 model , 所有业务围绕model、对象具备生命周期 
  - service: 贫血模型、由application 调用 
  - factory: 聚合根工厂
  - convertor: model 转换器 
  - vo_xxx: 值对象、不可修改、内存级别. 
  - gateway: 接口，domain依赖接口，而不是依赖持久层，该接口由持久层实现. 

infrastructure: 持久层，实现gateway接口，由domain调用, 由DI注入 

pkg: 与业务无关的工具类实现
```

