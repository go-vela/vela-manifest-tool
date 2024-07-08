// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestVersion(t *testing.T) {
	cmd := versionCmd()
	cases := []struct {
		arg, expected string
	}{
		{cmd.Args[0], "manifest-tool"},
		{cmd.Args[1], "--version"},
	}

	for _, tc := range cases {
		if !strings.Contains(tc.arg, tc.expected) {
			t.Errorf(`Expected %v to contain %q`, tc.arg, tc.expected)
		}
	}
}

// Feels like execCmd should be written/tested in shared lib.
func TestExecution(t *testing.T) {
	cases := []struct {
		args           []string
		expout, experr string
	}{
		{[]string{"echo", "-n", "foo"}, "foo", ""},
	}

	oldStdout := stdout
	defer func() { stdout = oldStdout }()

	oldStderr := stderr
	defer func() { stderr = oldStderr }()

	for _, tc := range cases {
		var outbuf, errbuf bytes.Buffer

		stdout, stderr = &outbuf, &errbuf
		cmd := exec.Command(tc.args[0], tc.args[1:]...) //nolint:gosec // we control the test data

		err := execCmd(cmd)
		if err != nil {
			t.Errorf("Expected no error when creating command: %v", err)
		}

		if tc.expout != outbuf.String() {
			t.Errorf("Expected %q to be equal to %q", outbuf.String(), tc.expout)
		}

		if tc.experr != errbuf.String() {
			t.Errorf("Expected %q to be equal to %q", errbuf.String(), tc.experr)
		}
	}
}
