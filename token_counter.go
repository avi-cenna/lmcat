package internal

import (
	"github.com/rs/zerolog/log"
	"github.com/tiktoken-go/tokenizer"
	"sync"
)

var (
	GlobalTokenCounter TokenCounter
	once               sync.Once
)

type TokenCounter interface {
	CountTokens(bytes []byte) (int, error)
}
type ApproxTokenCounter struct {
}

type GptTokenCounter struct {
	codec tokenizer.Codec
}

func (a ApproxTokenCounter) CountTokens(bytes []byte) (int, error) {
	return len(string(bytes)) / 4, nil
}

func (t GptTokenCounter) CountTokens(bytes []byte) (int, error) {
	ids, _, err := t.codec.Encode(string(bytes))
	if err != nil {
		return 0, err
	}
	return len(ids), nil
}

func InitTokenizer(useApproximate bool) {
	once.Do(func() {
		if useApproximate {
			GlobalTokenCounter = ApproxTokenCounter{}
		} else {
			codec, err := tokenizer.Get(tokenizer.Cl100kBase)
			if err != nil {
				log.Fatal().Err(err).Msg("Error getting tokenizer")
			}
			GlobalTokenCounter = GptTokenCounter{codec}
		}
	})
}
