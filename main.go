package main

import (
	"github.com/crabkun/Monitor/web"
	"github.com/crabkun/Monitor/plugin"
)

func main(){
	web.StartWebServer(plugin.LoadAllPlugin())
}
