package kuiperbelt

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProxySendHandlerFunc__BulkSend(t *testing.T) {
	s1 := &TestSession{new(bytes.Buffer), "hogehoge", false, false}
	s2 := &TestSession{new(bytes.Buffer), "fugafuga", false, false}

	AddSession(s1)
	AddSession(s2)

	tc := TestConfig
	p := Proxy{tc}
	ts := httptest.NewServer(http.HandlerFunc(p.SendHandlerFunc))
	defer ts.Close()

	req, err := http.NewRequest("POST", ts.URL, bytes.NewBufferString("test message"))
	if err != nil {
		t.Fatal("proxy handler new request unexpected error:", err)
	}
	req.Header.Add(tc.SessionHeader, "hogehoge")
	req.Header.Add(tc.SessionHeader, "fugafuga")

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("proxy handler request unexpected error:", err)
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	result := struct {
		Result string `json:"result"`
	}{}
	err = dec.Decode(&result)
	if err != nil {
		t.Fatal("proxy handler response unexpected error:", err)
	}
	if result.Result != "OK" {
		t.Fatalf("proxy handler response unexpected response: %+v", result)
	}

	if s1.String() != "test message" {
		t.Fatalf("proxy handler s1 not receive message: %s", s1.String())
	}
	if s2.String() != "test message" {
		t.Fatalf("proxy handler s2 not receive message: %s", s2.String())
	}
}

func TestProxyCloseHandlerFunc__BulkClose(t *testing.T) {
	s1 := &TestSession{new(bytes.Buffer), "hogehoge", false, false}
	s2 := &TestSession{new(bytes.Buffer), "fugafuga", false, false}

	AddSession(s1)
	AddSession(s2)

	tc := TestConfig
	p := Proxy{tc}
	ts := httptest.NewServer(http.HandlerFunc(p.CloseHandlerFunc))
	defer ts.Close()

	req, err := http.NewRequest("POST", ts.URL, bytes.NewBufferString("test message"))
	if err != nil {
		t.Fatal("proxy handler new request unexpected error:", err)
	}
	req.Header.Add(tc.SessionHeader, "hogehoge")
	req.Header.Add(tc.SessionHeader, "fugafuga")

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("proxy handler request unexpected error:", err)
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	result := struct {
		Result string `json:"result"`
	}{}
	err = dec.Decode(&result)
	if err != nil {
		t.Fatal("proxy handler response unexpected error:", err)
	}
	if result.Result != "OK" {
		t.Fatalf("proxy handler response unexpected response: %+v", result)
	}

	if s1.String() != "test message" {
		t.Fatalf("proxy handler s1 is not receive message: %s", s1.String())
	}
	if s2.String() != "test message" {
		t.Fatalf("proxy handler s2 is not receive message: %s", s2.String())
	}

	if !s1.isClosed {
		t.Fatalf("proxy handler s1 is not closed")
	}
	if !s2.isClosed {
		t.Fatalf("proxy handler s1 is not closed")
	}
}
