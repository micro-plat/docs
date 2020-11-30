
服务运行
===================================
[toc]


支持开发模式与生产模式运行。

##### 开发模式
使用`run`命令运行，无须指定注册中心(默认`lm://.`），无须安装配置，日志清晰调试方便。

##### 生产模式
使用服务方式运行，使用分布式注册中心，配置预安装，以高可用方式管理服务。

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

#### 二、开发模式(run)

开发模式可直接使用`run`命令启动应用，应用将以前台方式运行，通过终端输出日志，通过CTL+C可退出当前应用。

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
以上参数可通过代码设置：
```go
package main

import (
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/hydra/servers/http"
)
  func main() {
	app := hydra.NewApp(
            hydra.WithPlatName("oms"), //必须
            hydra.WithServerTypes(http.API),//必须
            hydra.WithRegistry("lm://."), //可空，默认值为:lm://.		
            hydra.WithSystemName("apiserver"), //可空，默认值为当前应用名称
            hydra.WithClusterName("prod"), //可空，默认值为:prod
            hydra.WithVersion("1.2.0"), //可空，默认值为:1.0.0
            hydra.WithUsage("api接口服务"), //可空
	)
	app.API("/request", request)
	app.Start()
}

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

#### 二、生产模式(service)

以后台方式运行服务，异常关闭或重启OS后自动启动服务。

##### 1. 安装服务

```sh
$ ./apiserver install --help
NAME:
   apiserver install - 安装服务，以服务方式安装到本地系统

USAGE:
   apiserver install [command options] [arguments...]

OPTIONS:
   --registry value, -r value      -注册中心地址。格式：proto://host。如：zk://ip1,ip2  或 fs://../ [$registry]
   --name value, -n value          -服务全名，格式：/平台名称/系统名称/服务器类型/集群名称 [$name]
   --plat value, -p value          -平台名称
   --system value, -s value        -系统名称,默认为当前应用程序名称
   --server-types value, -S value  -服务类型，有api,web,rpc,cron,mqc,ws
   --cluster value, -c value       -集群名称，默认值为：prod
   --cover, -v                     -覆盖安装，本地已安装服务

```
安装参数与`run`参数基本一致，因为服务启动过程实际上是通过后台的方式执行`run`命令。

```sh
$  ./apiserver install
Install apiserver                                      [  OK  ]
```


##### 2. 启动服务

```sh
$ ./apiserver start --help
NAME:
   apiserver start - 启动服务，以后台方式运行服务
USAGE:
   apiserver start [arguments...]
```

无需任何参数即可启动服务：

```sh
$   ./apiserver start
Start apiserver                                        [  OK  ]
```
服务通过后台方式启动成功！
> 由于本示例并未执行配置安装(将服务器配置安装到注册中心)，所以注册中心(registry)只能是"lm://."(本地内存)才能启动成功。

>使用其它注册中心(zookeeper、etcd、redis、fs等)需先执行`apiserver conf install ....`命令安装配置或手工创建配置，方能启动成功。

##### 3. 查看状态
```sh
$   ./apiserver status
Status apiserver Running                                       [  OK  ]
```
服务正在运行！

##### 4. 停止服务

```sh
$  ./apiserver stop
Stop apiserver                                         [  OK  ]
```
再次查看服务已关闭:
$  ./apiserver status
Status apiserver Stopped                                       [  OK  ]

##### 5. 卸载服务
```sh
$  ./apiserver remove
Remove apiserver                                       [  OK  ]
```