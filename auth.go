package pocket

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

type oauthCodeRequest struct {
	ConsumerKey string `json:"consumer_key"`
	RedirectURL string `json:"redirect_uri"`
}

type oauthCodeResponse struct {
	State int    `json:"state"`
	Code  string `json:"code"`
}

type oauthTokenRequest struct {
	ConsumerKey string `json:"consumer_key"`
	Code        string `json:"code"`
}

type oauthTokenResponse struct {
	AccessToken string `json:"access_token"`
	Username    string `json:"username"`
}

func (p *PocketClient) Auth(status chan string) error {
	// State 4
	serverChan := make(chan string)
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			serverChan <- "Got this far"
			w.WriteHeader(200)
			w.Write([]byte("Go back to the client"))
		},
	))

	// Stage 0
	if p.AccessToken != "" {
		return errors.New("Access token is already set")
	}
	defer ts.Close()

	// Stage 2
	request, err := json.Marshal(&oauthCodeRequest{
		ConsumerKey: p.ConsumerKey,
		RedirectURL: ts.URL,
	})
	if err != nil {
		return errors.New("Stage 2: Marshal: " + err.Error())
	}

	body := bytes.NewReader(request)
	req, err := http.NewRequest("POST", p.URL+"/v3/oauth/request", body)
	if err != nil {
		return errors.New("Stage 2: New Request: " + err.Error())
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("X-Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.New("Stage 2: Client Do: " + err.Error())
	}

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Stage 2: ReadAll: " + err.Error())
	}
	resp.Body.Close()

	var code oauthCodeResponse
	err = json.Unmarshal(response, &code)
	if err != nil {
		return errors.New("Stage 2: Unmarshal: " + err.Error())
	}

	// Stage 3
	visitURL := fmt.Sprintf("%s/auth/authorize?request_token=%s&redirect_uri=%s", p.URL, code.Code, ts.URL)
	status <- visitURL
	fmt.Printf("Please visit %s\n", visitURL)
	// TODO this needs to timeout
	<-serverChan

	// Stage 5
	request, err = json.Marshal(&oauthTokenRequest{
		ConsumerKey: p.ConsumerKey,
		Code:        code.Code,
	})
	if err != nil {
		return errors.New("Stage 5: Marshal: " + err.Error())
	}
	body = bytes.NewReader(request)
	req, err = http.NewRequest("POST", p.URL+"/v3/oauth/authorize", body)
	if err != nil {
		return errors.New("Stage 5: New Request: " + err.Error())
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("X-Accept", "application/json")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return errors.New("Stage 5: Client Do: " + err.Error())
	}

	response, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Stage 5: ReadAll: " + err.Error())
	}
	resp.Body.Close()

	var token oauthTokenResponse
	err = json.Unmarshal(response, &token)
	if err != nil {
		return errors.New("Stage 5: Unmarshal: " + err.Error())
	}

	p.AccessToken = token.AccessToken
	p.Username = token.Username
	return nil
}
