package main

import (
	"embed"
	"flag"
	"io/fs"

	"github.com/chasemao/tts-model-server/server"
)

//go:embed webui/out/*
var webFiles embed.FS

var ip = flag.String("ip", "0.0.0.0", "自定义IP地址")
var port = flag.Int64("port", 1233, "自定义监听端口")
var token = flag.String("token", "", "使用token验证")

func main() {
	flag.Parse()
	webFilesFs, _ := fs.Sub(webFiles, "webui/out")
	srv := &server.Processer{
		IP:       *ip,
		Port:     *port,
		Token:    *token,
		WebFiles: webFilesFs,
	}
	srv.Serve()
}
