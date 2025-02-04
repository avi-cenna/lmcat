package main

import (
	_ "github.com/tiktoken-go/tokenizer"
	_ "reflect"
	"testing"
)

func TestApproxTokenCounter_CountTokens(t *testing.T) {
	type args struct {
		bytes []byte
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name:    "empty string",
			args:    args{bytes: []byte("")},
			want:    0,
			wantErr: false,
		},
		{
			name:    "4 character string",
			args:    args{bytes: []byte("test")},
			want:    1,
			wantErr: false,
		},
		{
			name:    "8 character string",
			args:    args{bytes: []byte("testtest")},
			want:    2,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := ApproxTokenCounter{}
			got, err := a.CountTokens(tt.args.bytes)
			if (err != nil) != tt.wantErr {
				t.Errorf("CountTokens() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CountTokens() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTokenFunc(t *testing.T) {
	tests := []struct {
		name        string
		approx      bool
		inputString string
		want        int
	}{
		{
			name:        "approximate short text",
			approx:      true,
			inputString: "Hello world",
			want:        2, // 11 chars / 4 = 2 tokens approx
		},
		{
			name:        "approximate longer text",
			approx:      true,
			inputString: "This is a longer piece of text for testing",
			want:        10, // 40 chars / 4 = 10 tokens approx
		},
		{
			name:        "precise short text",
			approx:      false,
			inputString: "Hello world",
			want:        2, // actual GPT tokens
		},
		{
			name:        "precise longer text",
			approx:      false,
			inputString: "This is a longer piece of text for testing",
			want:        9, // actual GPT tokens
		},
		{
			name:        "empty string",
			approx:      false,
			inputString: "",
			want:        0, // actual GPT tokens
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			countFunc := GetTokenFunc(tt.approx)
			got := countFunc([]byte(tt.inputString))
			if got != tt.want {
				t.Errorf("GetTokenFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

//func TestGptTokenCounter_CountTokens(t1 *testing.T) {
//	type fields struct {
//		codec tokenizer.Codec
//	}
//	type args struct {
//		bytes []byte
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		want    int
//		wantErr bool
//	}{
//		{
//			name: "empty string",
//			fields: fields{
//				codec: tokenizer.NewCLIP(),
//			},
//			args: args{
//				bytes: []byte(""),
//			},
//			want:    0,
//			wantErr: false,
//		},
//		{
//			name: "simple text",
//			fields: fields{
//				codec: tokenizer.NewCLIP(),
//			},
//			args: args{
//				bytes: []byte("Hello world"),
//			},
//			want:    2,
//			wantErr: false,
//		},
//		{
//			name: "longer text",
//			fields: fields{
//				codec: tokenizer.NewCLIP(),
//			},
//			args: args{
//				bytes: []byte("This is a longer piece of text for testing"),
//			},
//			want:    9,
//			wantErr: false,
//		},
//	}
//	for _, tt := range tests {
//		t1.Run(tt.name, func(t1 *testing.T) {
//			t := GptTokenCounter{
//				codec: tt.fields.codec,
//			}
//			got, err := t.CountTokens(tt.args.bytes)
//			if (err != nil) != tt.wantErr {
//				t1.Errorf("CountTokens() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if got != tt.want {
//				t1.Errorf("CountTokens() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
