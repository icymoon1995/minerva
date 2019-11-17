# minerva
Welcome to your advice

## 使用方式:  
使用了go-micro的网关 
注册中心使用的consul(使用别的请替换)
 

1. 更改config/env.yaml的配置
 
2. 开启 go-micro 网关   
```MICRO_REGISTRY=consul micro api --handler=http --address={address} --namespace=personal.icymoon```

3. 运行服务     
 ```MICRO_REGISTRY=consul go run src/server/main.go```

## 目录

*   config - 配置文件  
*   docs - 文档  
*   log - 日志  
*   **src - 源码**    
    *   common - 公用的一些代码
    *   http
        *   controller - 控制器
        *   middleware - 中间间 
    *   logic - 处理逻辑 相当于services 
    *   model - model层
    *   routes - 路由文件
    *   server - 启动文件 
    *   service - 对接其他服务的文件

## 更新记录
*   v0.1   
        搭建项目    
*   v0.2    
        使用Echo做路由的转发器
*   v0.3    
        集成xorm + mysql
*   v0.4    
        中间件jwt登录及session存储
*   v0.5    
        redis对接,并完成test文件
*   v0.6    
        go-micro 服务注册完成,替代原始的echo启动
*   v0.7    
        基于rabbitmq做tcc事务
*   v.8     
       集成logrus 并将日志按照日期/类型输出到文件中
### Used
xorm    
echo    
go-micro    
rabbitmq    
consul
logrus
..




