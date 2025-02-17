package coquiaitts

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"sync"

	"github.com/chasemao/tts-model-server/log"
	"github.com/chasemao/tts-model-server/model"
	"github.com/google/uuid"
)

func NewModel() model.Model {
	return &impl{
		mu: &sync.Mutex{},
	}
}

type impl struct {
	mu *sync.Mutex
}

func (i *impl) Name() string {
	return "github.com/coqui-ai/TTS"
}

var submodels = []string{
	"tts_models/en/ljspeech/tacotron2-DDC",
	"tts_models/en/ljspeech/glow-tts",
	"tts_models/en/ljspeech/speedy-speech",
	"tts_models/en/ljspeech/vits",
	"tts_models/en/ljspeech/vits--neon",
	"tts_models/en/ljspeech/fast_pitch",
	"tts_models/en/ljspeech/overflow",
	"tts_models/en/ljspeech/neural_hmm",
}

func (i *impl) Fields() []*model.Field {
	var options []*model.Option
	for _, submodel := range submodels {
		options = append(options, &model.Option{
			Value: submodel,
		})
	}
	return []*model.Field{
		{
			Name:         "submodel",
			DefaultValue: "",
			Options:      options,
		},
	}
}

func (i *impl) TTS(ctx context.Context, text string, args map[string]string) ([]byte, string, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	id := uuid.New().String()
	fileName := fmt.Sprintf("/tmp/ttsmodelserver/%s.wav", id)
	commands := []string{
		"conda", "run", "-n", "base",
		"tts",
		"--text", text,
		"--model_name", args["submodel"],
		"--out_path", fileName,
	}
	cmd := exec.CommandContext(ctx, commands[0], commands[1:]...)
	log.Logger.Info(cmd.Args)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Logger.Error("output:", string(output), " ,err:", err)
		return nil, "", err
	}
	log.Logger.Info(string(output))
	buf, err := os.ReadFile(fileName)
	if err != nil {
		log.Logger.Error(err)
		return nil, "", err
	}
	if err := os.Remove(fileName); err != nil {
		log.Logger.Error(err)
	}
	return buf, "wav", nil
}
