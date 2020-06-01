package cmd

import (
	"bytes"
	"github.com/codilime/floodgate/version"
	"io/ioutil"
	"testing"
)

func TestVersion(t *testing.T) {
	b := bytes.NewBufferString("")

	cmd := NewRootCmd(b)
	cmd.SetOut(b)
	cmd.SetArgs([]string{"version"})
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("cmd.Version() Execute err %v", err)
	}

	out, err := ioutil.ReadAll(b)
	if err != nil {
		t.Fatalf("cmd.Version() Read output err %v", err)
	}

	if string(out) != version.BuildInfo() {
		t.Fatalf("cmd.Version() got `%v` want `%v`", string(out), version.BuildInfo())
	}
}
