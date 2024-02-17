package rest

import (
	"FriendsAdvice/internal/database/postgresql"
	"FriendsAdvice/internal/services"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLiveness(t *testing.T) {
	r := CreateRouter(&services.Controller{})
	ts := httptest.NewServer(r.MuxRouter)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/probes/liveness")
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("Status code for /probes/liveness is wrong. Have: %d, want: %d.", res.StatusCode, http.StatusOK)
	}
}

func TestReadiness(t *testing.T) {
	storageManager, _ := postgresql.InitManager(createConnectionDTO())
	controller := services.InitController(storageManager)

	r := CreateRouter(controller)
	ts := httptest.NewServer(r.MuxRouter)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/probes/readiness")
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("Status code for /probes/readiness is wrong. Have: %d, want: %d.", res.StatusCode, http.StatusOK)
	}
}

func createConnectionDTO() *postgresql.ConnectionDTO {
	return &postgresql.ConnectionDTO{
		HostName: "172.25.32.1",
		User:     "postgres",
		Password: "postgres",
		Port:     "5432",
		DBName:   "postgres"}
}
