package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

type Profile struct {
	Did            string `json:"did"`
	Handle         string `json:"handle"`
	DisplayName    string `json:"displayName"`
	Description    string `json:"description"`
	FollowingCount int    `json:"followsCount"`
	FollowerCount  int    `json:"followersCount"`
	PostCount      int    `json:"postsCount"`
}

func (s *Service) GetProfile(handle string) (*Profile, error) {
	handle = strings.TrimPrefix(handle, "@")
	url := fmt.Sprintf("%s/app.bsky.actor.getProfile?actor=%s", baseurl, handle)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create get profile request")
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.userAuth.AccessJwt))

	res, err := s.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get profile")
	}

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response")
	}

	var p Profile
	err = json.Unmarshal(resBody, &p)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get unmarshal profile")
	}

	return &p, nil
}
