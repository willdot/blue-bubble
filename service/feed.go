package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type Feed struct {
	Feed []FeedItem `json:"feed"`
}

type FeedItem struct {
	Post Post `json:"post"`
}

type Post struct {
	Record      Record `json:"record"`
	LikeCount   int    `json:"likeCount"`
	ReplyCount  int    `json:"replyCount"`
	RepostCount int    `json:"repostCount"`
}

type Record struct {
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"createdAt"`
}

func (s *Service) GetUserFeed(handle string) (*Feed, error) {
	handle = strings.TrimPrefix(handle, "@")
	url := fmt.Sprintf("%s/app.bsky.feed.getAuthorFeed?actor=%s", baseurl, handle)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create get user feed request")
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.userAuth.AccessJwt))

	res, err := s.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user feed")
	}

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response")
	}

	var feed Feed
	err = json.Unmarshal(resBody, &feed)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get unmarshal users feed")
	}

	return &feed, nil
}
