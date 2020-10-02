package config

import (
	"os"
	"reflect"
	"testing"
)

func TestProxy_AppendNode(t *testing.T) {
	type fields struct {
		Address  string
		UserName string
		Env      string
		TwoFA    bool
		Node     Node
	}
	type args struct {
		n Node
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Node
		wantErr bool
	}{
		{
			name: "",
			fields: fields{
				Env: "prod",
			},
			args: args{
				n: Node{
					Items: []Item{
						{
							Hostname: "proxy-172.20.1.3",
							Address:  "172.20.1.3:3022",
						},
					},
				},
			},
			want: Node{
				Items: []Item{
					{
						Hostname: "proxy-172.20.1.1",
						Address:  "172.20.1.1:3022",
					},
					{
						Hostname: "proxy-172.20.1.2",
						Address:  "172.20.1.2:3022",
					},
					{
						Hostname: "proxy-172.20.1.3",
						Address:  "172.20.1.3:3022",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Proxy{
				Address:  tt.fields.Address,
				UserName: tt.fields.UserName,
				Env:      tt.fields.Env,
				TwoFA:    tt.fields.TwoFA,
				Node:     tt.fields.Node,
			}

			// this not proper solution
			// just override the configuration
			wd, err := os.Getwd()
			if err != nil {
				t.Fatal(err)
			}

			configDir = wd + "/test/"

			got, err := p.AppendNode(tt.args.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("AppendNode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppendNode() got = %v, want %v", got, tt.want)
			}
		})
	}
}
