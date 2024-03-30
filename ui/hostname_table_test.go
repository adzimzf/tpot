package ui

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_stringResult_colorized(t *testing.T) {
	type fields struct {
		MatchPositions [][]int
		Keyword        string
		Data           string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "cartwheel",
			fields: fields{
				MatchPositions: func() [][]int {
					_, p := findMatchPositions("whl", "cartwheel")
					return p
				}(),
				Keyword: "whl",
				Data:    "cartwheel",
			},
			want: "cart\x1b[37;7mw\x1b[0m\x1b[37;7mh\x1b[0mee\x1b[37;7ml\x1b[0m",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := stringResult{
				MatchPositions: tt.fields.MatchPositions,
				Keyword:        tt.fields.Keyword,
				Data:           tt.fields.Data,
			}
			assert.Equalf(t, tt.want, s.colorizeMatchChars(), "colorizeMatchChars()")
		})
	}
}

func Test_findMatchPositions(t *testing.T) {
	type args struct {
		keyword string
		str     string
	}
	tests := []struct {
		name          string
		args          args
		wantMatches   []string
		wantPositions [][]int
	}{
		{
			name: "myapp-api",
			args: args{
				keyword: "api",
				str:     "myapp-api",
			},
			wantMatches: []string{"a", "p", "p", "a", "p", "i"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMatches, _ := findMatchPositions(tt.args.keyword, tt.args.str)
			assert.Equalf(t, tt.wantMatches, gotMatches, "findMatchPositions(%v, %v)", tt.args.keyword, tt.args.str)
		})
	}
}

func Test_lookup(t *testing.T) {
	type args struct {
		keyword string
		datum   []string
	}
	tests := []struct {
		name string
		args args
		want []stringResult
	}{
		{
			name: "empty keyword",
			args: args{
				keyword: "",
				datum: []string{
					"data-personalization",
					"donation",
					"fulfillment",
					"umrah",
				},
			},
			want: []stringResult{
				{
					Data: "data-personalization",
				},
				{
					Data: "donation",
				},
				{
					Data: "fulfillment",
				},
				{
					Data: "umrah",
				},
			},
		},
		{
			name: "jump",
			args: args{
				keyword: "donat",
				datum: []string{
					"data-personalization",
					"donation",
				},
			},
			want: []stringResult{
				{
					Data: "donation",
				},
				{
					Data: "data-personalization",
				},
			},
		},
		{
			name: "scoring number",
			args: args{
				keyword: "um",
				datum: []string{
					"fulfillment",
					"umrah",
				},
			},
			want: []stringResult{
				{
					Data: "umrah",
				},
				{
					Data: "fulfillment",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := lookup(tt.args.keyword, tt.args.datum)
			assert.Equal(t, len(tt.want), len(res))
			for i, result := range tt.want {
				assert.Equalf(t, result.Data, res[i].Data, fmt.Sprintf("index %d is not match", i))
			}
		})
	}
}
