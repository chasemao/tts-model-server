package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/chasemao/tts-model-server/log"
	"github.com/chasemao/tts-model-server/model"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

func (s *Processer) getFields(c *gin.Context) {
	var modelOptions []*model.Option
	for _, m := range s.models {
		modelOptions = append(modelOptions, &model.Option{
			Value:         m.Name(),
			RelatedFields: m.Fields(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"fields": []*model.Field{
			{
				Name:         "model",
				DefaultValue: modelOptions[0].Value,
				Options:      modelOptions,
			},
		},
	})
}

func (p *Processer) getSubScribe(c *gin.Context) {
	// Unpack args, and put text field in it, and extract url
	args := make(map[string]string)
	for key, values := range c.Request.URL.Query() {
		if len(values) > 0 {
			args[key] = values[0]
		}
	}
	ttsHost := args["host"]
	delete(args, "host")
	args["text"] = "{{speakText}}"
	buf, err := json.Marshal(args)
	if err != nil {
		log.Logger.Error(err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Pack URL for Legado
	deatil := gin.H{
		"method": "POST",
		"body":   string(buf),
	}
	detailBuf, err := json.Marshal(deatil)
	if err != nil {
		log.Logger.Error(err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	invokeURL := ttsHost + "/tts/api/invoke," + string(detailBuf)

	c.JSON(200, gin.H{
		"concurrentRate":   "1000",
		"contentType":      p.formatContentType(args["format"]),
		"enabledCookieJar": false,
		"header":           "{\"Content-Type\":\"appliction/json\"}",
		"loginCheckJs":     "",
		"loginUi":          "",
		"loginUrl":         "",
		"name":             "Config from TTS Model Server",
		"url":              invokeURL,
	})
}

func (p *Processer) invokeTTS(c *gin.Context) {
	if err := p.invokeTTSCore(c); err != nil {
		log.Logger.Error(err)
		c.Header("error", err.Error())
	}
}

func (p *Processer) invokeTTSCore(c *gin.Context) error {
	// Unpack args
	args := make(map[string]string)
	if err := c.BindJSON(&args); err != nil {
		return errors.New("invalid JSON")
	}
	log.Logger.Info("args: ", args)

	// Check token
	if p.Token != "" && args["token"] != p.Token {
		return fmt.Errorf("invalid token: %s", args["token"])
	}

	// Find correct model
	m, ok := lo.Find(p.models, func(item model.Model) bool {
		return item.Name() == args["model"]
	})
	if !ok {
		return errors.New("invalid model")
	}

	// Prepare timeout context
	ctx, cancel := context.WithTimeout(c.Request.Context(), 120*time.Second)
	defer cancel()

	// Call model to get response
	buf, format, err := m.TTS(ctx, args["text"], args)
	if err != nil {
		log.Logger.Error(err)
		return err
	}

	// Pack response
	c.Header("Content-Length", strconv.FormatInt(int64(len(buf)), 10))
	c.Header("Connection", "keep-alive")
	c.Header("Keep-Alive", "timeout=5")
	c.Data(200, p.formatContentType(format), buf)

	return nil
}

func (p *Processer) formatContentType(format string) string {
	t := strings.Split(format, "-")[0]
	switch t {
	case "audio":
		return "audio/mpeg"
	case "webm":
		return "audio/webm; codec=opus"
	case "ogg":
		return "audio/ogg; codecs=opus; rate=16000"
	case "riff":
		return "audio/x-wav"
	case "raw":
		if strings.HasSuffix(format, "truesilk") {
			return "audio/SILK"
		} else {
			return "audio/basic"
		}
	}
	return ""
}
