package httpclient

import (
	"testing"
)

func TestGET(t *testing.T) {
	resp, err := GET("http://google.com", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = resp.String()
	if err != nil {
		t.Error(err)
	}
}

func TestGETJSON(t *testing.T) {
	resp, err := GETJSON("https://shortesttrack.com/oauth/token_key", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	var m struct {
		Alg   string `json:"alg"`
		Value string `json:"value"`
	}
	err = resp.JSON(&m)
	if err != nil {
		t.Fatal(err)
	}
}
