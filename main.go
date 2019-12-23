package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/getlantern/systray"
)

var (
	StartMenu *systray.MenuItem
	StopMenu  *systray.MenuItem

	dnsList = map[string]string{
		"OpenDNS_1":    "208.67.222.222:5353",
		"OpenDNS_2":    "208.67.220.220:443",
		"OpenDNS_2-fs": "208.67.220.123:443",
	}

	logger *log.Logger

	domainName = "steamcommunity.com"
)

const (
	LOG_PREFIX = "[游戏便当]"
)

func handleExit(exit chan struct{}) {
	for _ = range exit {
		logger.Print("exit ...")
		systray.Quit()
	}
}

func handleStart(start chan struct{}) {
	logger.Printf("start proxy ...")
	for _ = range start {
		var ipAddr string
		StartMenu.Disable()
		address, err := DnsLookUp(domainName, dnsList, logger)
		if err == nil && address != nil {
			ipAddr = address.String()
		} else {
			address, err = HttpLookup(domainName)
			if err == nil && address != nil {
				ipAddr = address.String()
			}
		}
		if ipAddr == "" {
			errMsg := ""
			if err != nil {
				errMsg = err.Error()
			}
			logger.Printf("get ip failed! error is %s", errMsg)
			StartMenu.Enable()
			continue
		}

		logger.Printf("get ip is %s", ipAddr)
		addHosts()
		go startService(ipAddr, logger)
		StopMenu.Enable()
	}
}

func addHosts() {
	if err := AddHosts("127.0.0.1", domainName); err != nil {
		logger.Fatal(err)
	}
}

func removeHosts() {
	if err := RemoveHosts("127.0.0.1", domainName); err != nil {
		logger.Fatal(err)
	}
}

func handleStop(stop chan struct{}) {
	for _ = range stop {
		StopMenu.Disable()
		go stopService(logger)
		removeHosts()
		StartMenu.Enable()
	}
}

func handleAboutus(about chan struct{}) {
	for _ = range about {
		cmd := exec.Command("explorer.exe", "https://www.youxibd.com")
		cmd.Start()
	}
}

func trayReady() {
	systray.SetIcon(TrayIcon)
	systray.SetTitle("开始")
	systray.SetTooltip("游戏补丁")
	abus := systray.AddMenuItem("关于我们", "")
	systray.AddSeparator()
	StopMenu = systray.AddMenuItem("停止代理", "")
	StartMenu = systray.AddMenuItem("启动代理", "")
	StartMenu.Disable()
	exit := systray.AddMenuItem("退出程序", "")

	go handleAboutus(abus.ClickedCh)
	go handleStart(StartMenu.ClickedCh)
	go handleStop(StopMenu.ClickedCh)
	go handleExit(exit.ClickedCh)
	StartMenu.ClickedCh <- struct{}{}
}

func trayExit() {

}

func main() {
	os.Mkdir("./log", 0666)
	today := time.Now().Format("2006-01-02")
	logName := fmt.Sprintf("./log/%s.log", today)
	writer, err := os.OpenFile(logName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
		return
	}
	defer func() {
		writer.Close()
		removeHosts()
	}()

	logger = log.New(writer, LOG_PREFIX, log.LstdFlags)
	if err != nil {
		panic(err)
		return
	}
	systray.Run(trayReady, trayExit)
}
