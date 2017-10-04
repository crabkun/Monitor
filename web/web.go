package web

import (
	"github.com/crabkun/Monitor/plugin"
	"github.com/crabkun/crab"
)

var pluginMap map[string]*plugin.Plugin

func index(c *crab.Context){
	c.TplData["pluginArr"]=pluginMap
	c.Tpl("index.html")
}

func loadRoute(cfg *crab.Config){
	cfg.AddStaticPath("/static","static")
	cfg.Router.Add("GET","/",index)
	for _,v:=range pluginMap{

		cfg.AddStaticPath("/plugins/"+v.PluginInfo.Name+"/static",
			"plugins/"+v.PluginInfo.Name+"/static")

		cfg.Router.Add("GET","/plugins/"+v.PluginInfo.Name+"/",v.PluginIndex)

		for _,r:=range *v.PluginRoute{
			p,err:=v.Plugin.Lookup(r.FuncName)
			if err!=nil{
				continue
			}
			cfg.Router.Add(r.Method,"/plugins/"+v.PluginInfo.Name+"/"+r.Path,p.(func(context *crab.Context)))
		}
	}

}

func StartWebServer(m map[string]*plugin.Plugin){
	pluginMap=m
	cfg:=&crab.Config{}
	cfg.Address=":8080"
	cfg.UseMidware("session","2h")
	cfg.UseMidware("log","")
	loadRoute(cfg)
	crab.Listen(cfg)
}