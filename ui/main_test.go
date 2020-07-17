package ui

import (
	"reflect"
	"testing"
)

func Test_matchWord(t *testing.T) {
	type args struct {
		data    []string
		keyword string
	}
	tests := []struct {
		name string
		args args
		want map[string]stringResult
	}{
		{
			name: "",
			args: args{
				data: []string{
					"userdata", "userogod", "userondol", "notmatch",
				},
				keyword: "user",
			},
			want: map[string]stringResult{
				"userdata": {
					FormatedData: "userdata",
				},
				"userogod": {
					FormatedData: "userogod",
				},
				"userondol": {
					FormatedData: "userondol",
				},
			},
		},
		{
			name: "",
			args: args{
				data:    sampleData,
				keyword: "user",
			},
			want: map[string]stringResult{
				"user@localhost": {
					FormatedData: "user@localhost",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filterData(tt.args.data, tt.args.keyword); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_toString(t *testing.T) {
	type args struct {
		d map[string]stringResult
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{
				d: map[string]stringResult{
					"user":  {FormatedData: "user@user"},
					"user2": {FormatedData: "user2@user"},
				},
			},
			want: "user@user\nuser2@user\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatResult(tt.args.d); got != tt.want {
				t.Errorf("formatResult() = %v, want %v", got, tt.want)
			}
		})
	}
}
