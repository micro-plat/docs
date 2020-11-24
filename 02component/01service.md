
服务运行
===================================
支持开发环境与生产环境两种方式运行。

##### 开发环境
默认使用本地内存作为注册中心，通过`run`前台运行服务，方便调试。

##### 生产环境
将应用安装到本地服务，以服务方式运行，避免意外关闭和服务器重启。

#### 一、命令概况

以编译后的二进制文件`apiserver`为例，输入 `apiserver --help`可查看支持的所有命令：


|----apiserver
|----------- main.go

```go
package main

import (
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/hydra/servers/http"
)

func main() {
	app := hydra.NewApp(
		hydra.WithServerTypes(http.API),
	)
	app.API("/request", request)
	app.Start()
}
func request(ctx hydra.IContext) interface{} {
	return "success"
}
```

```sh
$ ./apiserver --help
NAME:
   apiserver - apiserver(A new hydra application)

USAGE:
   apiserver [global options] command [command options] [arguments...]

VERSION:
   1.0.0

COMMANDS:
   conf     配置管理, 查看、安装配置信息
   install  安装服务，以服务方式安装到本地系统
   remove   删除服务，从本地服务器移除服务
   run      运行服务,以前台方式运行服务。通过终端输出日志，终端关闭后服务自动退出。
   update   更新应用，将服务发布到远程服务器
   start    启动服务，以后台方式运行服务
   status   查询状态，查询服务器运行、停止状态
   stop     停止服务，停止服务器运行
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     查看帮助信息
   --version, -v  查看版本信息
```

#### 二、开发环境(run)

开发环境可直接使用`run`命令启动应用，应用将以前台方式运行，通过终端输出日志，通过CTL+C可退出当前应用。

```sh
$ ./apiserver run --help
NAME:
   apiserver run - 运行服务,以前台方式运行服务。通过终端输出日志，终端关闭后服务自动退出。

USAGE:
   apiserver run [command options] [arguments...]

OPTIONS:
   --registry value, -r value      -注册中心地址。格式：proto://host。如：zk://ip1,ip2  或 fs://../ 
   --name value, -n value          -服务全名，格式：/平台名称/系统名称/服务器类型/集群名称 
   --plat value, -p value          -平台名称
   --system value, -s value        -系统名称,默认为当前应用程序名称
   --server-types value, -S value  -服务类型，有api,web,rpc,cron,mqc,ws
   --cluster value, -c value       -集群名称，默认值为：prod
   --trace value, -t value         -性能分析。支持:cpu,mem,block,mutex,web
   --tport value, --tp value       -性能分析服务端口号。用于trace为web模式时的端口号。默认：19999
```
以上参数可以通过代码方式设置：
```go
package main

import (
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/hydra/servers/http"
)
  func main() {
	app := hydra.NewApp(
		hydra.WithServerTypes(http.API),
		hydra.WithRegistry("lm://."),
		hydra.WithPlatName("oms"),
		hydra.WithSystemName("apiserver"),
		hydra.WithClusterName("prod"),
		hydra.WithVersion("1.2.0"),
		hydra.WithUsage("api接口服务"),
	)
	app.API("/request", request)
	app.Start()
}

```
为便于快速定位问题，日志根据级别显示不同颜色：

```go
func request(ctx hydra.IContext) interface{} {
	ctx.Log().Debug("调试信息")
	ctx.Log().Info("一般信息")
	ctx.Log().Warn("警告信息")
	ctx.Log().Error("错误信息")
	return "success"
}
```
运行服务：
```sh
$ go build
$ ./apiserver run
```

```sh
$ curl http://localhost:8080/request
```
![日志](./imgs/log.png)