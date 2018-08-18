package pocket_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	pocket "github.com/brimstone/go-pocket"
)

func TestAuth(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			/* TODO actually check the body
			body, err := ioutil.ReadAll(r.Body)
			r.Body.Close()
			if err != nil {
				w.WriteHeader(400)
				return
			}
			*/
			fmt.Printf("Mock: %s\n", r.URL.Path)
			if r.URL.Path == "/v3/oauth/request" {
				w.WriteHeader(200)
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.Write([]byte(`{"state":null,"code":"dcba4321-dcba-4321-dcba-4321dc"}`))
			} else if r.URL.Path == "/v3/oauth/authorize" {
				w.WriteHeader(200)
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.Write([]byte(`{"access_token":"5678defg-5678-defg-5678-defg56",
				"username":"pocketuser"}`))
			}
		},
	))
	defer ts.Close()

	p := pocket.NewPocketClient(&pocket.PocketClientOptions{
		ConsumerKey: "55086-b24420f727e2d3a80014b34d",
		URL:         ts.URL,
	})

	status := make(chan string)
	go func() {
		url := <-status
		urlParts := strings.Split(url, "=")
		// TODO test for out of range
		_, err := http.Get(urlParts[2])
		if err != nil {
			fmt.Println(err)
			t.Error(err)
		}
	}()

	err := p.Auth(status)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("client: %#v\n", p)
}
