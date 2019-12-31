package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/kyf/systray"
	"github.com/lxn/walk"
)

var (
	StartMenu       *systray.MenuItem
	StopMenu        *systray.MenuItem
	startBt, exitBt *walk.PushButton

	dnsList = map[string]string{
		"OpenDNS_1":    "208.67.222.222:5353",
		"OpenDNS_2":    "208.67.220.220:443",
		"OpenDNS_2-fs": "208.67.220.123:443",
	}

	logger *log.Logger

	domainName = "steamcommunity.com"
	ui_index   = "http://127.0.0.1:" + UI_PORT
)

const (
	LOG_PREFIX    = "[游戏便当]"
	title         = "游戏便当Steam社区加速"
	width, height = 270, 330
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

func handleUI() {
	cmd := exec.Command("explorer.exe", ui_index)
	cmd.Start()
}

func trayReady() {
	systray.SetIcon(TrayIcon)
	systray.SetTitle(title)
	systray.SetTooltip(title)
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
	window := systray.GetWindow()
	if window != nil {
		window.SetMinMaxSizePixels(walk.Size{width, height}, walk.Size{width, height})
		winIco, err := walk.Resources.Image("youxibd.ico")
		if err == nil {
			window.SetIcon(winIco)
		}
	}
	notifyIcon := systray.GetNotifyIcon()
	if notifyIcon != nil {
		notifyIcon.MouseUp().Attach(func(x, y int, button walk.MouseButton) {
			if button == walk.LeftButton {
				//systray.ShowAppWindow(ui_index)
				handleUI()
			}
		})
	}
	//systray.ShowAppWindow(ui_index)
	checkFirst(window)
	go startUI()
	handleUI()
}

func trayExit() {

}

func checkFirst(window *walk.MainWindow) {
	_, err := ioutil.ReadFile("./deploy")
	if err != nil {
		if os.IsNotExist(err) {
			walk.MsgBox(window, "提示", "请确认添加授权证书，以保证访问正常", walk.MsgBoxIconInformation)
			output, err := exec.Command("./certs/certmgr.exe", "/c", "/add", "./certs/steamcommunityCA.pem", "/s", "root").Output()
			if err != nil {
				logger.Printf("add ca error:%s, output is %s", err.Error(), string(output))
				return
			}
			log.Print(string(output))
			ioutil.WriteFile("./deploy", []byte(""), os.ModePerm)
		} else {
			logger.Printf("read deploy error %s", err.Error())
		}
		return
	}
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
	systray.RunWithAppWindow(title, width, height, trayReady, trayExit)
}
