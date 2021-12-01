package main

import (
	"StairUnlocker-Go/config"
	"StairUnlocker-Go/utils"
	"flag"
	"fmt"
	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/log"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"runtime"
	"time"
)

var (
	proxiesList []C.Proxy
	netflixCfg  config.RawConfig
	version     bool
	help        bool
	su          *config.SuConfig
)

func init() {
	su = config.Init()
	flag.BoolVar(&help, "h", false, "this help")
	flag.BoolVar(&version, "v", false, "show current version of StairUnlock")
	flag.StringVar(&su.SubURL, "u", su.SubURL, "Load config from subscription url")
	flag.StringVar(&su.Token, "t", su.Token, "The github token")
	flag.StringVar(&su.GistUrl, "g", su.GistUrl, "The gist api URL")
	flag.Parse()
	fmt.Printf("StairUnlock-Go %s %s %s with %s\n", utils.Version, runtime.GOOS, runtime.GOARCH, runtime.Version())
	if su.LocalFile {
		fmt.Println("Local file mode: on")
	} else {
		fmt.Println("Gist mode: on")
	}
}

func main() {
	//command-line
	if version {
		fmt.Printf("StairUnlock %s %s %s with %s\n", utils.Version, runtime.GOOS, runtime.GOARCH, runtime.Version())
		return
	}
	if help {
		fmt.Printf("StairUnlock %s %s %s with %s\n", utils.Version, runtime.GOOS, runtime.GOARCH, runtime.Version())
		flag.PrintDefaults()
		return
	}

	proxies, cfg, _ := config.GenerateProxies(su)
	for _, v := range proxies {
		proxiesList = append(proxiesList, v)
	}

	//同时连接数
	connNum := su.MaxConn
	if i := len(proxiesList); i < connNum {
		connNum = i
	}
	start := time.Now()
	netflixList := utils.BatchCheck(proxiesList, connNum)
	log.Warnln("Completed! Elapsed time: %s", time.Now().Sub(start).String())

	netflixCfg = config.NETFLIXFilter(netflixList, cfg)
	marshal, _ := yaml.Marshal(netflixCfg)
	if su.LocalFile {
		_ = ioutil.WriteFile("netflix.yaml", marshal, 0644)
		log.Infoln("Written to netflix.yaml.")
	} else {
		err := utils.Gist(marshal, su)
		if err != nil {
			return
		}
	}
}
