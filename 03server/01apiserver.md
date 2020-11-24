API服务器
----------------------
对外提供http api 服务。

[TOC]

### 一、简单的服务器
```go
package main

import (
    "github.com/micro-plat/hydra"
    "github.com/micro-plat/hydra/hydra/servers/http"
)

func main() {
  
    app := hydra.NewApp(hydra.WithServerTypes(http.API))

    //注册服务
    app.API("/api", api)
    
    //启动服务
    app.Start()
}

func api(ctx hydra.IContext) interface{} {
    ctx.Log().Info("--------api--------------")
    return map[string]interface{}{
        "name":"colin",
    }
}
```
> 响应内容: *{"name":"colin"}*  
> Content-Type: application/json; charset=utf-8



### 二、服务器配置

代码设置:
```go
hydra.Conf.API(":8081", api.WithTrace(),api.WithTimeout(5,5))
```


* 任何配置都可以操作注册中心添加、修改、删除。
* 固定不变的配置可以通过代码设置。
* 开发环境的配置可通过代码配置，但建议通过预编译命令对不同环境进行编译隔离。
