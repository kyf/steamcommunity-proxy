package main

import (
	"crypto/tls"
	"io"
	"log"
	"net"
	"net/http"
	"os/exec"

	"github.com/kyf/martini"
	mylog "github.com/kyf/util/log"
)

const (
	cert = "./certs/steamcommunity.crt"
	key  = "./certs/steamcommunity.key"

	UI_PORT = "18060"
)

var (
	ls, ls80 net.Listener
	sig      = make(chan int, 1)
)

func startService(ip string, logger *log.Logger) {
	cfg := &tls.Config{}
	cfg.Certificates = make([]tls.Certificate, 1)
	var err error
	cfg.Certificates[0], err = tls.LoadX509KeyPair(cert, key)
	if err != nil {
		logger.Printf("load cert error is %s", err.Error())
		return
	}
	ls, err = tls.Listen("tcp", ":443", cfg)
	if err != nil {
		logger.Printf("listen local error %s", err.Error())
		return
	}

	go start80()

	for {
		conn, err := ls.Accept()
		if err != nil {
			logger.Printf("accept connection error:%s", err.Error())
			break
		}

		go handleConn(conn, ip, logger)
	}
}

type Server80 struct{}

func (this *Server80) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://"+domainName, http.StatusFound)
	return
}

func start80() {
	var err error
	ls80, err = net.Listen("tcp", ":80")
	if err != nil {
		logger.Printf("start 80service error %s", err.Error())
		return
	}
	svr := &Server80{}
	err = http.Serve(ls80, svr)
	if err != nil {
		logger.Printf("start 80service error %s", err.Error())
		return
	}
}

func handleServiceStatus(w http.ResponseWriter, r *http.Request, params martini.Params) {
	status := params["status"]
	if status == "true" {
		if !StartMenu.Disabled() {
			StartMenu.ClickedCh <- struct{}{}
		}
	} else {
		if StartMenu.Disabled() {
			StopMenu.ClickedCh <- struct{}{}
		}
	}

	jsonTo(w, true, "ok")
}

func handleGetStatus(w http.ResponseWriter, r *http.Request) {
	status := StartMenu.Disabled()
	jsonTo(w, true, "ok", status)
}

func handleAddCA(w http.ResponseWriter, r *http.Request) {
	output, err := exec.Command("./certs/certmgr.exe", "/c", "/add", "./certs/steamcommunityCA.pem", "/s", "root").Output()
	if err != nil {
		logger.Printf("add ca error:%s, output is %s", err.Error(), string(output))
		jsonTo(w, false, err.Error())
		return
	}

	jsonTo(w, true, "ok", string(output))
}

func startUI() {
	m := martini.Classic()
	m.Map(mylog.DefaultLogger)
	m.Use(martini.Static("./ui"))
	m.Get("/api/service/status/:status", handleServiceStatus)
	m.Get("/api/service/getstatus", handleGetStatus)
	m.Get("/api/service/addca", handleAddCA)
	m.RunOnAddr(":" + UI_PORT)
}

func handleConn(conn net.Conn, ip string, logger *log.Logger) {
	defer conn.Close()

	cfg := &tls.Config{
		ServerName:         "",
		InsecureSkipVerify: true,
	}

	remote, err := tls.Dial("tcp", ip+":443", cfg)
	if err != nil {
		logger.Printf("dial remote error %s", err.Error())
		return
	}
	defer remote.Close()

	go func() {
		for _ = range sig {
			conn.Close()
			remote.Close()
			logger.Printf("close connections ...")
		}
	}()
	go io.Copy(remote, conn)
	io.Copy(conn, remote)
}

func stopService(logger *log.Logger) {
	var err error
	if ls != nil {
		err = ls.Close()
		if err != nil {
			logger.Printf("stop service error is %s", err.Error())
			return
		}
	}
	if ls80 != nil {
		err = ls80.Close()
		if err != nil {
			logger.Printf("stop 80 service error is %s", err.Error())
			return
		}
	}
	sig <- 1

	logger.Printf("stop service ...")
}
