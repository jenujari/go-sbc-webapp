package main

import (
	"jenujari/go-sbc-webapp/config"
	"jenujari/go-sbc-webapp/server"

	rtc "github.com/jenujari/runtime-context"
)

var pc *rtc.ProcessContext

func init() {
	rtc.InitProcessContext(config.GetLogger())
}

func main() {
	pc = rtc.GetMainProcess()
	srv := server.GetServer()
	pc.Run(server.RunServer)
	config.GetLogger().Println("Server is running at ", srv.Addr)

	pc.WaitForFinish()
}
