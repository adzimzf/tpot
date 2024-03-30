package tsh

import (
	"bytes"
	"testing"

	"github.com/adzimzf/tpot/config"
	"github.com/stretchr/testify/assert"
)

func TestTSH_parseStringToStatus(t1 *testing.T) {
	t1.Skip("will update later")
	tests := []struct {
		name   string
		str    string
		status *config.ProxyStatus
	}{
		{
			name: "",
			str: `
> Profile URL:  https:/my.teleport.com
  Logged in as: ikhsan@my.com
  Cluster:      main
  Roles:        engineer-frontend-role, engineer-role*
  Logins:       ikhsan@my.com, root, readonly
  Valid until:  2021-02-14 22:56:30 +0700 WIB [valid for 12h0m0s]
  Extensions:   permit-X11-forwarding, permit-agent-forwarding, permit-port-forwarding, permit-pty


* RBAC is only available in Teleport Enterprise
  https://gravitational.com/teleport/docs/enterprise
`,
			status: &config.ProxyStatus{
				LoginAs:    "ikhsan@my.com",
				Roles:      []string{"engineer-frontend-role", "engineer-role*"},
				UserLogins: []string{"ikhsan@my.com", "root", "readonly"},
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &TSH{}
			got := t.parseStringToStatus(tt.str)
			assert.Equal(t1, tt.status, got)
		})
	}
}

type cmdMock struct {
	cmdResult
	err error
}

func (c *cmdMock) Run() (cmdResult, error) {
	return c.cmdResult, c.err
}

func TestTSH_isLogin(t1 *testing.T) {

	t1.Skip("will update later")
	type fields struct {
		proxy      *config.Proxy
		userLogin  string
		dstHost    string
		cmdExec    func(name string, arg ...string) CmdExecutor
		minVersion Version
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "got success login",
			fields: fields{
				proxy: &config.Proxy{
					Address: "https://staging.teleport.net:3080",
				},
				cmdExec: func(name string, arg ...string) CmdExecutor {
					stdOut := bytes.NewBufferString(`
> Profile URL:        https://staging.teleport.net:3080
  Logged in as:       youremail@domain.com
  Cluster:            main
  Roles:              engineer
  Logins:             root
  Kubernetes:         disabled
  Valid until:        2023-07-08 21:36:23 +0700 WIB [valid for 11h59m0s]
  Extensions:         permit-X11-forwarding, permit-agent-forwarding, permit-port-forwarding, permit-pty

  Profile URL:        https://teleport.net:3080
  Logged in as:       youremail@domain.com
  Cluster:            main
  Roles:              engineer
  Logins:             non-root, root
  Kubernetes:         disabled
  Valid until:        2023-07-08 00:46:44 +0700 WIB [EXPIRED]
  Extensions:         permit-X11-forwarding, permit-agent-forwarding, permit-port-forwarding, permit-pty
`)

					return &cmdMock{
						cmdResult: cmdResult{
							stdOut: stdOut,
							stdErr: &bytes.Buffer{},
						},
					}
				},
			},
			want: true,
		},
		{
			name: "got success non-login",
			fields: fields{
				proxy: &config.Proxy{
					Address: "https://teleport.net:3080",
				},
				cmdExec: func(name string, arg ...string) CmdExecutor {
					stdOut := bytes.NewBufferString(`
> Profile URL:        https://staging.teleport.net:3080
  Logged in as:       youremail@domain.com
  Cluster:            main
  Roles:              engineer
  Logins:             root
  Kubernetes:         disabled
  Valid until:        2023-07-08 21:36:23 +0700 WIB [valid for 11h59m0s]
  Extensions:         permit-X11-forwarding, permit-agent-forwarding, permit-port-forwarding, permit-pty

  Profile URL:        https://teleport.net:3080
  Logged in as:       youremail@domain.com
  Cluster:            main
  Roles:              engineer
  Logins:             non-root, root
  Kubernetes:         disabled
  Valid until:        2023-07-08 00:46:44 +0700 WIB [EXPIRED]
  Extensions:         permit-X11-forwarding, permit-agent-forwarding, permit-port-forwarding, permit-pty
`)

					return &cmdMock{
						cmdResult: cmdResult{
							stdOut: stdOut,
							stdErr: &bytes.Buffer{},
						},
					}
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &TSH{
				proxy:      tt.fields.proxy,
				userLogin:  tt.fields.userLogin,
				dstHost:    tt.fields.dstHost,
				cmdExec:    tt.fields.cmdExec,
				minVersion: tt.fields.minVersion,
			}
			assert.Equalf(t1, tt.want, t.isLogin(), "isLogin()")
		})
	}
}
