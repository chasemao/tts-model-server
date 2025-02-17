package server

import (
	"fmt"
	"io/fs"
	"net/http"

	"github.com/chasemao/tts-model-server/log"
	"github.com/chasemao/tts-model-server/model"
	"github.com/chasemao/tts-model-server/model/coquiaitts"
	"github.com/chasemao/tts-model-server/model/edge"
	"github.com/gin-gonic/gin"
)

type Processer struct {
	IP       string
	Port     int64
	Token    string
	WebFiles fs.FS

	models    []model.Model
	ginServer gin.Engine
}

func (s *Processer) Serve() {
	r := gin.Default()

	s.models = append(s.models,
		edge.NewModel(),
		coquiaitts.NewModel(),
	)

	r.StaticFS("/tts/webui", http.FS(s.WebFiles))

	r.GET("/tts/api/fields", s.getFields)
	r.GET("/tts/api/subscribe", s.getSubScribe)
	r.POST("/tts/api/invoke", s.invokeTTS)

	addr := fmt.Sprintf("%s:%d", s.IP, s.Port)
	log.Logger.Infof("listen on %s...", addr)
	r.Run(addr)
}
