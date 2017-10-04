package plugin

import (
	"io/ioutil"
	"bufio"
	"log"
	"bytes"
	"plugin"
	"github.com/crabkun/MonitorKits"
	"github.com/crabkun/crab"
	"errors"
)

type Plugin struct {
	Plugin *plugin.Plugin
	PluginInfo *MonitorKits.PluginInfo
	PluginRoute *MonitorKits.PluginRoute
	LoadPlugin func()error
	UnloadPlugin func()error
	PluginIndex func(*crab.Context)
}
func checkLoadErr(err error){
	if err!=nil{
		panic(err)
	}
}

func loadPlugin(filename string,pluginMap map[string]*Plugin){
	defer func(){
		if err:=recover();err!=nil{
			log.Printf("插件%s加载错误,原因%v\n",filename,err)
		}
	}()
	var err error
	var i interface{}
	t:=Plugin{}
	p,err:=plugin.Open("plugins/"+filename+"/"+filename+".so")
	checkLoadErr(err)

	i,err=p.Lookup("GetPluginInfo")
	checkLoadErr(err)
	t.PluginInfo=i.(func()(*MonitorKits.PluginInfo))()
	if _,ok:=pluginMap[t.PluginInfo.Name];ok{
		checkLoadErr(errors.New("此插件已被加载过"))
	}

	i,err=p.Lookup("GetPluginRoute")
	checkLoadErr(err)
	t.PluginRoute=i.(func()(*MonitorKits.PluginRoute))()

	i,err=p.Lookup("LoadPlugin")
	checkLoadErr(err)
	t.LoadPlugin=i.(func()error)
	t.LoadPlugin()

	i,err=p.Lookup("UnloadPlugin")
	checkLoadErr(err)
	t.UnloadPlugin=i.(func()error)

	i,err=p.Lookup("PluginIndex")
	checkLoadErr(err)
	t.PluginIndex=i.(func(*crab.Context))

	log.Printf("[插件]%s %s加载成功！\n",t.PluginInfo.Name,t.PluginInfo.Version)

	t.Plugin=p
	pluginMap[filename]=&t
}

func LoadAllPlugin()(map[string]*Plugin){
	var err error
	pluginMap:=make(map[string]*Plugin)
	buf,err:=ioutil.ReadFile("plugin.lst")
	if err!=nil{
		panic("插件清单文件(plugin.lst)加载错误！"+err.Error())
	}
	reader:=bufio.NewReader(bytes.NewBuffer(buf))
	for{
		l,_,err:=reader.ReadLine()
		if err!=nil{
			break
		}
		loadPlugin(string(l),pluginMap)
	}
	return pluginMap
}