package main

import (
	"bytes"
	"net"
	"testing"
)

func TestUIConfigStateRoundTripAndAliases(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	st := loadUIConfigState()
	if st.Host != defaultUIHost || st.Port != defaultUIPort || st.AllowRemote {
		t.Fatalf("unexpected default ui state: %+v", st)
	}

	custom := uiConfigState{
		Host:        "0.0.0.0",
		Port:        7899,
		AllowRemote: true,
	}
	saveUIConfigState(custom)

	loaded := loadUIConfigState()
	if loaded != custom {
		t.Fatalf("loaded = %+v, want %+v", loaded, custom)
	}

	if got := uiDisplayURL(loaded); got != "http://127.0.0.1:7899" {
		t.Fatalf("uiDisplayURL(0.0.0.0:7899) = %q, want http://127.0.0.1:7899", got)
	}

	for _, alias := range []string{"ui", "web", "webui", "web-ui", "UI"} {
		if !isUIServiceTarget(alias) {
			t.Errorf("expected %q to be recognized as UI service target", alias)
		}
	}
	if isUIServiceTarget("github") {
		t.Error("did not expect github to be recognized as UI service target")
	}

	var buf bytes.Buffer
	if err := printUIStatus(&buf); err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(buf.Bytes(), []byte("stopped")) {
		t.Fatalf("expected stopped status, got %q", buf.String())
	}

	// Verify isUIAddrListening detects a bound listener
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	port := ln.Addr().(*net.TCPAddr).Port
	if !isUIAddrListening(uiConfigState{Host: "127.0.0.1", Port: port}) {
		t.Fatalf("expected listener on port %d to be detected", port)
	}
}
