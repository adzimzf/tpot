package tsh

import (
	"testing"

	"github.com/adzimzf/tpot/config"
	"github.com/stretchr/testify/assert"
)

func TestTSH_parseStringToStatus(t1 *testing.T) {
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
