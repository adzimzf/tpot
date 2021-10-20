package tsh

import (
	"fmt"
	"io"
	"os/exec"
)

// Forward run the tsh forwarding
func (t *TSH) Forward(userLogin, host, forwardAddress string, in io.Reader) error {
	args, err := t.getProxyFlags()
	if err != nil {
		return err
	}

	args = append(args, t.authFlags()...)
	args = append(args, fmt.Sprintf("%s@%s", userLogin, host))
	args = append([]string{"ssh", "-L", forwardAddress}, args...)
	cmd := exec.Command(t.tshBinary(), args...)
	cmd.Stdin = in
	return cmd.Run()
}
