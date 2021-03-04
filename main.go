package main

import (
	"embed"

	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/conf/server/header"
	"github.com/micro-plat/hydra/conf/server/static"
	"github.com/micro-plat/hydra/hydra/servers/http"
)

var app = hydra.NewApp(
	hydra.WithPlatName("hydra", "hydra"),
	hydra.WithSystemName("docs", "文档中心"),
	hydra.WithServerTypes(http.Web),
)

func main() {
	app.Start()
}

//go:embed web
var mgrweb embed.FS

func init() {
	//设置配置参数
	hydra.Conf.Web("80").Header(header.WithCrossDomain()).
		Static(static.WithAutoRewrite(), static.WithEmbed("web", mgrweb))
}
