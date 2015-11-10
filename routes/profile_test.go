package routes_test

import (
	"fmt"
	"github.com/maleck13/locals-api"
	"net/http"
	"net/http/httptest"
	"testing"

	"io/ioutil"
)

func TestGetProfile(t *testing.T) {
	server := httptest.NewServer(main.SetUpRoutes())
	defer server.Close()

	profileGet := fmt.Sprintf("%s/profile/17", server.URL)
	resp, err := http.Get(profileGet)
	if nil != err {
		t.Fatalf("failed to get " + profileGet)
	}

	data, err := ioutil.ReadAll(resp.Body)

	content := string(data)

	t.Logf("data %s ", content, "status code  %i", resp.StatusCode)

}
