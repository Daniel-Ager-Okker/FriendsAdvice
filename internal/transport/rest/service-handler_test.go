package rest

import (
	"FriendsAdvice/internal/services"
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetInvalidKey(t *testing.T) {
	r := CreateRouter(&services.Controller{})
	ts := httptest.NewServer(r.MuxRouter)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/objects/abc")
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("Status code for /GET/ is wrong. Have: %d, want: %d.", res.StatusCode, http.StatusBadRequest)
	}
}

func TestGet(t *testing.T) {
	r := CreateRouter(&services.Controller{})
	ts := httptest.NewServer(r.MuxRouter)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/objects/1")
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("Status code for /GET/1 is wrong. Have: %d, want: %d.", res.StatusCode, http.StatusOK)
	}
}

func TestPutInvalidKey(t *testing.T) {
	r := CreateRouter(&services.Controller{})
	ts := httptest.NewServer(r.MuxRouter)
	defer ts.Close()

	// Create request
	obj := "Oh!"
	reqBody := bytes.NewBufferString(obj)

	req, err := http.NewRequest("PUT", ts.URL+"/objects/AnyKey", reqBody)
	if err != nil {
		t.Fatalf("Error creating PUT request: %v", err)
	}

	// Send
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		t.Fatalf("Error sending PUT request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("Status code for /PUT/objects/AnyKey is wrong. Have: %d, want: %d.", res.StatusCode, http.StatusBadRequest)
	}
}
