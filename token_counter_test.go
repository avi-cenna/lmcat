package main

import (
	"github.com/tiktoken-go/tokenizer"
	"reflect"
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
		// TODO: Add test cases.
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
		//TODO: fix these test cases, do two with approx=true and two with approx=false
		{
			name: "Returns approximate token function",
			args: args{
				useApproximate: true,
			},
			want: ApproximateTokenFunction, // Replace with actual function reference
		},
		{
			name: "Returns precise token function",
			args: args{
				useApproximate: false,
			},
			want: PreciseTokenFunction, // Replace with actual function reference
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetTokenFunc(tt.args.useApproximate); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTokenFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestGptTokenCounter_CountTokens(t1 *testing.T) {
	type fields struct {
		codec tokenizer.Codec
	}
	type args struct {
		bytes []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := GptTokenCounter{
				codec: tt.fields.codec,
			}
			got, err := t.CountTokens(tt.args.bytes)
			if (err != nil) != tt.wantErr {
				t1.Errorf("CountTokens() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t1.Errorf("CountTokens() got = %v, want %v", got, tt.want)
			}
		})
	}
}
