配置中心
=================================
[toc]

通过配置中心，集中管理服务器配置，并具有如下特点:
1. 统一配置，本地零配置，解决配置零散分布在每台服务器，导致的不一致问题及服务器数量较多难以维护的问题。
2. 热更新，所有配置通过配置中心拉取，变更后实时通知到各服务器，服务器根据排号自动更新配置。
3. 可审计，远程配置通过多级审核，避免参数配置错误导致安全事故。

目前支持以下配置中心:
* 1. zookeeper,示例: zk://192.168.0.101,192.168.0.102
* 2. etcd,示例:etcd://192.168.0.100
* 3. redis,示例:redis://192.168.0.101,192.168.0.102
* 4. filesystem,示例: fs://../
* 5. local memory 示例: lm://.


注意：目前配置中心与配置中心共用同一服务器。

#### 1). 配置分类
服务器启动、运行期间的配置信息，分为三类：
* 1). 服务器启动参数，如：端口号、连接超时信息、集群模式(对等、主从、分片等)
* 2). 服务器组件参数，如：黑名单、限流、安全验证、跨域配置、metric、render，任务信息等
* 3). 应用配置参数，如：数据库配置、缓存配置、消息队列服务器配置等


#### 2). 代码设置
通过`hydra.Conf`进行设置。

除`MQC`消息消费服务启动时需要连接到消息队列服务器外，其它服务均有默认启动参数，无须指定任何配置即可启动。

配置设置方式:
```go
    hydra.Conf.API(":8071", api.WithTimeout(10, 10)).
         .BlackList(blacklist.WithIP("222.202.**"))
    hydra.Conf.Vars().Redis("redis", redis.New("redis", redis.WithAddrs("192.168.0.106")))
```

配置依赖于平台名称、系统名称、集群名称等，而这些参数是通过cli命令行指定的，则需要将配置初始化代码放到`hydra.OnReady`中延迟执行：
```go
hydra.OnReady(func() {    
    hydra.Conf.Vars()
        .Redis("redis", redis.New(hydra.G.GetPlatName(), redis.WithAddrs("192.168.0.106")))
})

```


#### 3). 配置安装

除`lm`无须安装外，其它配置中心需先安装才能启动。

注意：配置安装与服务启动是两个独立的过程。配置安装的目的是在配置中心生成正确可用的配置，服务启动是从配置中心拉取配置进行服务初始化。两个过程发生在不同时间，一前一后。 由于`lm`具有本地调试的便利性，故将两个过程合二为一，由系统自动完成。



* 1). 命令安装

以`apiserver`为例:
```sh
$ ./apiserver conf install --help
NAME:
   apiserver conf install - -安装配置，将配置信息安装到配置中心

USAGE:
   apiserver conf install [command options] [arguments...]

OPTIONS:
   --registry value, -r value      -配置中心地址。格式：proto://host。如：zk://ip1,ip2  或 fs://../ 
   --name value, -n value          -服务全名，格式：/平台名称/系统名称/服务器类型/集群名称
   --plat value, -p value          -平台名称
   --system value, -s value        -系统名称,默认为当前应用程序名称
   --server-types value, -S value  -服务类型，有api,web,rpc,cron,mqc,ws
   --cluster value, -c value       -集群名称，默认值为：prod
   --cover, -v                     -覆盖配置，覆盖配置中心和本地服务
```

根据参数要求即可将本地代码设置的配置信息安装到配置中心。

示例:
```go
package main

import (
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/hydra/servers/http"
)

func main() {
	app := hydra.NewApp(
		hydra.WithServerTypes(http.API),
		hydra.WithRegistry("zk://192.168.0.109"),
		hydra.WithPlatName("oms"),
		hydra.WithSystemName("apiserver"),
		hydra.WithClusterName("prod"),
		hydra.WithVersion("1.2.0"),
		hydra.WithUsage("api接口服务"),
	)
	app.API("/request", request)
	app.Start()
}
func request(ctx hydra.IContext) interface{} {
	return "success"
}
```

安装配置:

```sh
$ go build
$ ./apiserver conf install
安装到配置中心:                                 [OK]
```
查看配置：
```sh
$ ./apiserver conf show
└─oms
  └─apiserver
    └─api
      └─prod
        └─conf
          └─main[1]
          └─router[2]
请输入数字序号 > 
```

* 2). 手动配置
 指通过三方工具远程连接到配置中心配置：
 zookeeper推荐使用工具:

![zooInspector](/02component/imgs/zooInspector.png)


