package list

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestVolumeListCommand_noTabs(t *testing.T) {
	t.Parallel()
	if strings.ContainsRune(New(nil).Help(), '\t') {
		t.Fatal("help has tabs")
	}
}

func TestVolumeListCommand_Validation(t *testing.T) {
	t.Parallel()
	ui := cli.NewMockUi()
	c := New(ui)

	cases := map[string]struct {
		args   []string
		output string
	}{
		"missing parameter": {
			[]string{"volume", "list"},
			"Error on parsing parameters: Project need to be specified",
		},
		"unknown parameter": {
			[]string{"volume", "list", "-project", "toto", "-unknown-parameter", "dummy-value"},
			"flag provided but not defined: -unknown-parameter",
		},
	}

	for name, tc := range cases {
		// Ensure our buffer is always clear
		if ui.ErrorWriter != nil {
			ui.ErrorWriter.Reset()
		}
		if ui.OutputWriter != nil {
			ui.OutputWriter.Reset()
		}

		code := c.Run(tc.args)
		if code == 0 {
			t.Errorf("%s: expected non-zero exit", name)
		}

		output := ui.ErrorWriter.String()
		if !strings.Contains(output, tc.output) {
			t.Errorf("`%s`: expected:\n`%q`\nto contain:\n`%q`", name, output, tc.output)
		}
	}
}

// func TestVolumeListCommand(t *testing.T) {
// 	t.Parallel()
// 	a := agent.NewTestAgent(t, ``)
// 	defer a.Shutdown()
// 	client := a.Client()
//
// 	ui := cli.NewMockUi()
// 	c := New(ui)
//
// 	pair := &api.KVPair{
// 		Key:   "foo",
// 		Value: []byte("bar"),
// 	}
// 	_, err := client.KV().Put(pair, nil)
// 	if err != nil {
// 		t.Fatalf("err: %#v", err)
// 	}
//
// 	args := []string{
// 		"-http-addr=" + a.HTTPAddr(),
// 		"foo",
// 	}
//
// 	code := c.Run(args)
// 	if code != 0 {
// 		t.Fatalf("bad: %d. %#v", code, ui.ErrorWriter.String())
// 	}
//
// 	output := ui.OutputWriter.String()
// 	if !strings.Contains(output, "bar") {
// 		t.Errorf("bad: %#v", output)
// 	}
// }
//
// func TestVolumeListCommand_Base64(t *testing.T) {
// 	t.Parallel()
// 	a := agent.NewTestAgent(t, ``)
// 	defer a.Shutdown()
// 	client := a.Client()
//
// 	ui := cli.NewMockUi()
// 	c := New(ui)
//
// 	pair := &api.KVPair{
// 		Key:   "foo",
// 		Value: []byte("bar"),
// 	}
// 	_, err := client.KV().Put(pair, nil)
// 	if err != nil {
// 		t.Fatalf("err: %#v", err)
// 	}
//
// 	args := []string{
// 		"-http-addr=" + a.HTTPAddr(),
// 		"-base64",
// 		"foo",
// 	}
//
// 	code := c.Run(args)
// 	if code != 0 {
// 		t.Fatalf("bad: %d. %#v", code, ui.ErrorWriter.String())
// 	}
//
// 	output := ui.OutputWriter.String()
// 	if !strings.Contains(output, base64.StdEncoding.EncodeToString(pair.Value)) {
// 		t.Errorf("bad: %#v", output)
// 	}
// }
//
// func TestVolumeListCommand_Missing(t *testing.T) {
// 	t.Parallel()
// 	a := agent.NewTestAgent(t, ``)
// 	defer a.Shutdown()
//
// 	ui := cli.NewMockUi()
// 	c := New(ui)
//
// 	args := []string{
// 		"-http-addr=" + a.HTTPAddr(),
// 		"not-a-real-key",
// 	}
//
// 	code := c.Run(args)
// 	if code == 0 {
// 		t.Fatalf("expected bad code")
// 	}
// }
//
// func TestVolumeListCommand_Empty(t *testing.T) {
// 	t.Parallel()
// 	a := agent.NewTestAgent(t, ``)
// 	defer a.Shutdown()
// 	client := a.Client()
//
// 	ui := cli.NewMockUi()
// 	c := New(ui)
//
// 	pair := &api.KVPair{
// 		Key:   "empty",
// 		Value: []byte(""),
// 	}
// 	_, err := client.KV().Put(pair, nil)
// 	if err != nil {
// 		t.Fatalf("err: %#v", err)
// 	}
//
// 	args := []string{
// 		"-http-addr=" + a.HTTPAddr(),
// 		"empty",
// 	}
//
// 	code := c.Run(args)
// 	if code != 0 {
// 		t.Fatalf("bad: %d. %#v", code, ui.ErrorWriter.String())
// 	}
// }
//
// func TestVolumeListCommand_Detailed(t *testing.T) {
// 	t.Parallel()
// 	a := agent.NewTestAgent(t, ``)
// 	defer a.Shutdown()
// 	client := a.Client()
//
// 	ui := cli.NewMockUi()
// 	c := New(ui)
//
// 	pair := &api.KVPair{
// 		Key:   "foo",
// 		Value: []byte("bar"),
// 	}
// 	_, err := client.KV().Put(pair, nil)
// 	if err != nil {
// 		t.Fatalf("err: %#v", err)
// 	}
//
// 	args := []string{
// 		"-http-addr=" + a.HTTPAddr(),
// 		"-detailed",
// 		"foo",
// 	}
//
// 	code := c.Run(args)
// 	if code != 0 {
// 		t.Fatalf("bad: %d. %#v", code, ui.ErrorWriter.String())
// 	}
//
// 	output := ui.OutputWriter.String()
// 	for _, key := range []string{
// 		"CreateIndex",
// 		"LockIndex",
// 		"ModifyIndex",
// 		"Flags",
// 		"Session",
// 		"Value",
// 	} {
// 		if !strings.Contains(output, key) {
// 			t.Fatalf("bad %#v, missing %q", output, key)
// 		}
// 	}
// }
//
// func TestVolumeListCommand_Keys(t *testing.T) {
// 	t.Parallel()
// 	a := agent.NewTestAgent(t, ``)
// 	defer a.Shutdown()
// 	client := a.Client()
//
// 	ui := cli.NewMockUi()
// 	c := New(ui)
//
// 	keys := []string{"foo/bar", "foo/baz", "foo/zip"}
// 	for _, key := range keys {
// 		if _, err := client.KV().Put(&api.KVPair{Key: key}, nil); err != nil {
// 			t.Fatalf("err: %#v", err)
// 		}
// 	}
//
// 	args := []string{
// 		"-http-addr=" + a.HTTPAddr(),
// 		"-keys",
// 		"foo/",
// 	}
//
// 	code := c.Run(args)
// 	if code != 0 {
// 		t.Fatalf("bad: %d. %#v", code, ui.ErrorWriter.String())
// 	}
//
// 	output := ui.OutputWriter.String()
// 	for _, key := range keys {
// 		if !strings.Contains(output, key) {
// 			t.Fatalf("bad %#v missing %q", output, key)
// 		}
// 	}
// }
//
// func TestVolumeListCommand_Recurse(t *testing.T) {
// 	t.Parallel()
// 	a := agent.NewTestAgent(t, ``)
// 	defer a.Shutdown()
// 	client := a.Client()
//
// 	ui := cli.NewMockUi()
// 	c := New(ui)
//
// 	keys := map[string]string{
// 		"foo/a": "a",
// 		"foo/b": "b",
// 		"foo/c": "c",
// 	}
// 	for k, v := range keys {
// 		pair := &api.KVPair{Key: k, Value: []byte(v)}
// 		if _, err := client.KV().Put(pair, nil); err != nil {
// 			t.Fatalf("err: %#v", err)
// 		}
// 	}
//
// 	args := []string{
// 		"-http-addr=" + a.HTTPAddr(),
// 		"-recurse",
// 		"foo",
// 	}
//
// 	code := c.Run(args)
// 	if code != 0 {
// 		t.Fatalf("bad: %d. %#v", code, ui.ErrorWriter.String())
// 	}
//
// 	output := ui.OutputWriter.String()
// 	for key, value := range keys {
// 		if !strings.Contains(output, key+":"+value) {
// 			t.Fatalf("bad %#v missing %q", output, key)
// 		}
// 	}
// }
//
// func TestVolumeListCommand_RecurseBase64(t *testing.T) {
// 	t.Parallel()
// 	a := agent.NewTestAgent(t, ``)
// 	defer a.Shutdown()
// 	client := a.Client()
//
// 	ui := cli.NewMockUi()
// 	c := New(ui)
//
// 	keys := map[string]string{
// 		"foo/a": "Hello World 1",
// 		"foo/b": "Hello World 2",
// 		"foo/c": "Hello World 3",
// 	}
// 	for k, v := range keys {
// 		pair := &api.KVPair{Key: k, Value: []byte(v)}
// 		if _, err := client.KV().Put(pair, nil); err != nil {
// 			t.Fatalf("err: %#v", err)
// 		}
// 	}
//
// 	args := []string{
// 		"-http-addr=" + a.HTTPAddr(),
// 		"-recurse",
// 		"-base64",
// 		"foo",
// 	}
//
// 	code := c.Run(args)
// 	if code != 0 {
// 		t.Fatalf("bad: %d. %#v", code, ui.ErrorWriter.String())
// 	}
//
// 	output := ui.OutputWriter.String()
// 	for key, value := range keys {
// 		if !strings.Contains(output, key+":"+base64.StdEncoding.EncodeToString([]byte(value))) {
// 			t.Fatalf("bad %#v missing %q", output, key)
// 		}
// 	}
// }
//
// func TestVolumeListCommand_DetailedBase64(t *testing.T) {
// 	t.Parallel()
// 	a := agent.NewTestAgent(t, ``)
// 	defer a.Shutdown()
// 	client := a.Client()
//
// 	ui := cli.NewMockUi()
// 	c := New(ui)
//
// 	pair := &api.KVPair{
// 		Key:   "foo",
// 		Value: []byte("bar"),
// 	}
// 	_, err := client.KV().Put(pair, nil)
// 	if err != nil {
// 		t.Fatalf("err: %#v", err)
// 	}
//
// 	args := []string{
// 		"-http-addr=" + a.HTTPAddr(),
// 		"-detailed",
// 		"-base64",
// 		"foo",
// 	}
//
// 	code := c.Run(args)
// 	if code != 0 {
// 		t.Fatalf("bad: %d. %#v", code, ui.ErrorWriter.String())
// 	}
//
// 	output := ui.OutputWriter.String()
// 	for _, key := range []string{
// 		"CreateIndex",
// 		"LockIndex",
// 		"ModifyIndex",
// 		"Flags",
// 		"Session",
// 		"Value",
// 	} {
// 		if !strings.Contains(output, key) {
// 			t.Fatalf("bad %#v, missing %q", output, key)
// 		}
// 	}
//
// 	if !strings.Contains(output, base64.StdEncoding.EncodeToString([]byte("bar"))) {
// 		t.Fatalf("bad %#v, value is not base64 encoded", output)
// 	}
// }
