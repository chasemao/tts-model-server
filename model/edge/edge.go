package edge

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/chasemao/tts-model-server/log"
	"github.com/chasemao/tts-model-server/model"
	"github.com/chasemao/tts-model-server/wspool"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/samber/lo"
)

const (
	edgeMetadataURL = "https://speech.platform.bing.com/consumer/speech/synthesize/readaloud/voices/list?trustedclienttoken=6A5AA1D4EAFF4E9FB37E23D68491D6F4"
	wssUrl          = `wss://speech.platform.bing.com/consumer/speech/synthesize/readaloud/edge/v1?TrustedClientToken=6A5AA1D4EAFF4E9FB37E23D68491D6F4&ConnectionId=`
)

func NewModel() model.Model {
	header := http.Header{}
	header.Set("Accept-Encoding", "gzip, deflate, br")
	header.Set("Origin", "chrome-extension://jdiccldimpdaibmpdkjnbmckianbfold")
	header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.5060.66 Safari/537.36 Edg/103.0.1264.44")
	return &impl{
		connPool: wspool.New(wssUrl, true, 2, header, 20*time.Second),
	}
}

type impl struct {
	connPool wspool.Client
}

type edgeMetaData struct {
	ShortName      string
	Locale         string
	SuggestedCodec string
}

func (i *impl) Name() string {
	return "Edge"
}

func (i *impl) Fields() []*model.Field {
	resp, err := http.Get(edgeMetadataURL)
	if err != nil {
		log.Logger.Errorf("get edge metadata failed err: %v", err)
		return nil
	}
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Logger.Errorf("read edge metadata failed err: %v", err)
		return nil
	}
	var metas []*edgeMetaData
	if err := json.Unmarshal(buf, &metas); err != nil {
		log.Logger.Errorf("unmarshal edge metadata failed err: %v", err)
		return nil
	}

	metaGroups := lo.PartitionBy(metas, func(item *edgeMetaData) string {
		return item.Locale
	})
	var languageOptions []*model.Option
	for _, metaGroup := range metaGroups {
		var voiceOptions []*model.Option
		for _, meta := range metaGroup {
			voiceOptions = append(voiceOptions, &model.Option{
				Value: meta.ShortName,
				RelatedFields: []*model.Field{
					{
						Name:         "format",
						DefaultValue: meta.SuggestedCodec,
						Options: []*model.Option{
							{
								Value: meta.SuggestedCodec,
							},
						},
					},
				},
			})
		}
		languageOption := &model.Option{
			Value: metaGroup[0].Locale,
			RelatedFields: []*model.Field{
				{
					Name:         "voice",
					DefaultValue: metaGroup[0].ShortName,
					Options:      voiceOptions,
				},
			},
		}
		languageOptions = append(languageOptions, languageOption)
	}
	return []*model.Field{
		{
			Name:         "language",
			DefaultValue: "en-US",
			Options:      languageOptions,
		},
	}
}

func (i *impl) TTS(ctx context.Context, text string, args map[string]string) ([]byte, string, error) {
	voice := args["voice"]
	ssml := i.createSSML(text, voice)
	conn, err := i.connPool.Get(ctx)
	defer i.connPool.Putback(conn)
	if err != nil {
		log.Logger.Error(err)
		return nil, "", err
	}
	if err := i.sendConfigMessage(conn, args["format"]); err != nil {
		log.Logger.Error(err)
		return nil, "", err
	}
	if err := i.sendSsmlMessage(conn, ssml); err != nil {
		log.Logger.Error(err)
		return nil, "", err
	}
	var res []byte
	for {
		msgType, buf, err := conn.ReadMessage()
		if err != nil {
			log.Logger.Error(err)
			return nil, "", err
		}
		if msgType == websocket.BinaryMessage {
			index := strings.Index(string(buf), "Path:audio")
			data := []byte(string(buf)[index+12:])
			res = append(res, data...)
		}
		if msgType == websocket.TextMessage && strings.Contains(string(buf), "Path:turn.end") {
			break
		}
	}
	return res, args["format"], nil
}

func (t *impl) sendConfigMessage(conn *websocket.Conn, format string) error {
	cfgMsg := "X-Timestamp:" + time.Now().Format("2006-01-02T15:04:05Z") + "\r\nContent-Type:application/json; charset=utf-8\r\n" + "Path:speech.config\r\n\r\n" +
		`{"context":{"synthesis":{"audio":{"metadataoptions":{"sentenceBoundaryEnabled":"false","wordBoundaryEnabled":"false"},"outputFormat":"` + format + `"}}}}`
	return conn.WriteMessage(websocket.TextMessage, []byte(cfgMsg))
}

func (i *impl) escapeXML(text string) string {
	sb := &strings.Builder{}
	xml.Escape(sb, []byte(text))
	return sb.String()
}

func (i *impl) createSSML(text string, voice string) string {
	return `<speak xmlns="http://www.w3.org/2001/10/synthesis" xmlns:mstts="http://www.w3.org/2001/mstts" xmlns:emo="http://www.w3.org/2009/10/emotionml" version="1.0" xml:lang="en-US">` +
		`<voice name="` + i.escapeXML(voice) + `">` +
		`<prosody rate="0%" pitch="0%">` + i.escapeXML(text) + `</prosody ></voice ></speak >`
}

func (i *impl) sendSsmlMessage(conn *websocket.Conn, ssml string) error {
	id := uuid.New().String()
	msg := "Path: ssml\r\nX-RequestId: " + id + "\r\nX-Timestamp: " + time.Now().Format("2006-01-02T15:04:05Z") + "\r\nContent-Type: application/ssml+xml\r\n\r\n" + ssml
	return conn.WriteMessage(websocket.TextMessage, []byte(msg))
}
