package main

import (
	"github.com/rs/zerolog/log"
	"github.com/tiktoken-go/tokenizer"
)

//var (
//	CountTokens TokenFunc
//	once        sync.Once
//)

//func InitTokenizer(useApproximate bool) {
//	once.Do(func() {
//		if useApproximate {
//			CountTokens = func(bytes []byte) (int, error) {
//				return len(string(bytes)) / 4, nil
//			}
//		} else {
//			codec, err := tokenizer.Get(tokenizer.Cl100kBase)
//			if err != nil {
//				log.Fatal().Err(err).Msg("Error getting tokenizer")
//			}
//			CountTokens = func(bytes []byte) (int, error) {
//				ids, _, err := codec.Encode(string(bytes))
//				if err != nil {
//					return 0, err
//				}
//				return len(ids), nil
//			}
//		}
//	})
//}

func GetTokenFunc(useApproximate bool) TokenFunc {
	var countTokens TokenFunc
	if useApproximate {
		countTokens = func(bytes []byte) int {
			// Using 4 bytes per token as an approximation
			// See https://platform.openai.com/docs/concepts/tokens for more info
			return len(bytes) / 4
		}
	} else {
		// This is the tokenizer used by gpt-4o
		codec, err := tokenizer.Get(tokenizer.O200kBase)
		if err != nil {
			log.Fatal().Err(err).Msg("Error getting tokenizer")
		}
		countTokens = func(bytes []byte) int {
			ids, _, err := codec.Encode(string(bytes))
			if err != nil {
				log.Err(err).Msg("Error encoding tokens")
				return 0
			}
			return len(ids)
		}
	}
	return countTokens
}

//var (
//	CountTokens TokenCounter
//	once               sync.Once
//)

// TokenFunc is a function that takes a byte slice and returns the number of tokens and an error
type TokenFunc func(bytes []byte) int

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

//func InitTokenizer(useApproximate bool) {
//	once.Do(func() {
//		if useApproximate {
//			CountTokens = ApproxTokenCounter{}
//		} else {
//			codec, err := tokenizer.Get(tokenizer.Cl100kBase)
//			if err != nil {
//				log.Fatal().Err(err).Msg("Error getting tokenizer")
//			}
//			CountTokens = GptTokenCounter{codec}
//		}
//	})
//}
