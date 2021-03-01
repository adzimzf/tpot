package tsh

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    Version
		wantErr bool
	}{
		{
			name: "v2.6.0-alpha.4",
			args: args{
				s: "Teleport v2.6.0-alpha.4",
			},
			want: Version{
				Major: 2,
				Minor: 6,
				Patch: 0,
			},
		},
		{
			name: "v2.6.0-beta.3",
			args: args{
				s: "Teleport v2.6.0-beta.3",
			},
			want: Version{
				Major: 2,
				Minor: 6,
				Patch: 0,
			},
		},
		{
			name: "v2.6.0-beta.3",
			args: args{
				s: "Teleport v2.6.0-beta.3",
			},
			want: Version{
				Major: 2,
				Minor: 6,
				Patch: 0,
			},
		},
		{
			name: "v2.6.0-rc.1",
			args: args{
				s: "Teleport v2.6.0-rc.1",
			},
			want: Version{
				Major: 2,
				Minor: 6,
				Patch: 0,
			},
		},
		{
			name: "v2.65.123",
			args: args{
				s: "Teleport v2.65.123",
			},
			want: Version{
				Major: 2,
				Minor: 65,
				Patch: 123,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewVersion(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestVersion_IsSupported(t *testing.T) {
	tests := []struct {
		v           string
		cv          string
		isSupported bool
	}{
		{
			v:           "v2.6.0-beta.3",
			cv:          "v2.6.0",
			isSupported: true,
		},
		{
			v:           "v5.6.0",
			cv:          "v2.6.0",
			isSupported: false,
		},
		{
			v:           "v5.6.0",
			cv:          "v5.4.0",
			isSupported: false,
		},
		{
			v:           "v5.6.0",
			cv:          "v4.6.0",
			isSupported: false,
		},
		{
			v:           "v5.6.3",
			cv:          "v4.6.0",
			isSupported: false,
		},
		{
			v:           "v5.6.3",
			cv:          "v5.6.7",
			isSupported: true,
		},
		{
			v:           "v2.6.1",
			cv:          "v4.1.11",
			isSupported: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.v+" x "+tt.cv, func(t *testing.T) {
			v, err := NewVersion("Teleport " + tt.v)
			assert.NoError(t, err)
			cv, err := NewVersion("Teleport " + tt.cv)
			assert.NoError(t, err)
			assert.Equal(t, tt.isSupported, v.IsSupported(cv))
		})
	}
}
