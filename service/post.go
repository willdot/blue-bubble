package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type index struct {
	Start int `json:"byteStart"`
	End   int `json:"byteEnd"`
}
type feature struct {
	Type string `json:"$type"`
	Did  string `json:"did"`
}
type facet struct {
	Type     string    `json:"$type"`
	Index    index     `json:"index"`
	Features []feature `json:"features"`
}

func (s *Service) Post(message, mentionedHandle string) error {
	record := map[string]interface{}{
		"text":      message,
		"createdAt": time.Now().Format(time.RFC3339),
	}

	if mentionedHandle != "" {
		mention, err := s.createFeatureFacet(mentionedHandle, message)
		if err != nil {
			return err
		}

		record["facets"] = []*facet{mention}
	}

	reqData := map[string]interface{}{
		"collection": "app.bsky.feed.post",
		"repo":       s.userAuth.Did,
		"record":     record,
	}

	data, err := json.Marshal(reqData)
	if err != nil {
		return errors.Wrap(err, "failed to marshal request")
	}

	r := bytes.NewReader(data)

	url := fmt.Sprintf("%s/com.atproto.repo.createRecord", baseurl)
	req, err := http.NewRequest("POST", url, r)
	if err != nil {
		return errors.Wrap(err, "failed to create new post request")
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.userAuth.AccessJwt))

	res, err := s.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to make create post request")
	}

	defer res.Body.Close()
	if res.StatusCode == 200 {
		return nil
	}

	return fmt.Errorf("failed to create post: %v", res.StatusCode)
}

func (s *Service) createFeatureFacet(handle, message string) (*facet, error) {
	profile, err := s.GetProfile(handle)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get DID of mentioned handle")
	}

	f := feature{
		Type: "app.bsky.richtext.facet#mention",
		Did:  profile.Did,
	}
	start, end := getStartEndOfHandle(message, handle)
	idx := index{
		Start: start,
		End:   end,
	}

	return &facet{
		Type:     "app.bsky.richtext.facet",
		Index:    idx,
		Features: []feature{f},
	}, nil
}

func getStartEndOfHandle(message, handle string) (int, int) {
	start := strings.Index(message, handle)
	end := start + len(handle)

	return start, end
}
