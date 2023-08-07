package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
)

const (
	baseurl = "https://bsky.social/xrpc"

	httpClientTimeoutDuration        = time.Second * 5
	transportIdleConnTimeoutDuration = time.Second * 90
)

type auth struct {
	AccessJwt string `json:"accessJwt"`
	Did       string `json:"did"`
}

type Service struct {
	userAuth *auth
	client   http.Client
}

func New() (*Service, error) {
	client := http.Client{
		Timeout: httpClientTimeoutDuration,
		Transport: &http.Transport{
			IdleConnTimeout: transportIdleConnTimeoutDuration,
		},
	}

	auth, err := login(client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to authenticate")
	}
	return &Service{
		userAuth: auth,
		client:   client,
	}, nil
}

func login(client http.Client) (*auth, error) {
	handle := os.Getenv("BSKY_HANDLE")
	appPass := os.Getenv("BSKY_PASS")

	url := fmt.Sprintf("%s/com.atproto.server.createsession", baseurl)

	requestData := map[string]interface{}{
		"identifier": handle,
		"password":   appPass,
	}

	data, err := json.Marshal(requestData)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal request")
	}

	r := bytes.NewReader(data)

	req, err := http.NewRequest("POST", url, r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make request")
	}

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response")
	}

	var loginResp auth
	err = json.Unmarshal(resBody, &loginResp)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal response")
	}
	return &loginResp, nil
}
