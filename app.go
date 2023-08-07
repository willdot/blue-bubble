package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/willdot/bluebubble/input"
	"github.com/willdot/bluebubble/service"
)

type app struct {
	service *service.Service
}

func (a *app) newPost() error {
	postText, quit := input.PromptUserForText("What would you like to say?")
	if quit {
		fmt.Println("You quit. Goodbye")
		return nil
	}

	var mentionedHandle string
	if strings.Contains(postText, "@") {
		foundHandle := getMentionedHandleInMessage(postText)
		mention, quit := input.PromptUserForSingleChoice([]string{"yes", "no"}, fmt.Sprintf("detected potential handle '%s'. Would you like to mention?", foundHandle))
		if quit {
			fmt.Println("You quit. Goodbye")
			return nil
		}
		if mention == "yes" {
			mentionedHandle = foundHandle
		}
	}

	err := a.service.Post(postText, mentionedHandle)
	if err != nil {
		return errors.Wrap(err, "failed to make new post")
	}

	return nil
}

func (a *app) getProfile() error {
	user, quit := input.PromptUserForText("which user handle would you like to view?")
	if quit {
		fmt.Println("You quit. Goodbye")
		return nil
	}

	profile, err := a.service.GetProfile(user)
	if err != nil {
		return errors.Wrap(err, "failed to get handle DID")
	}

	fmt.Printf("DID: %s\n", profile.Did)
	fmt.Printf("Handle: %s\n", profile.Handle)
	fmt.Printf("Display Name: %s\n", profile.DisplayName)
	fmt.Println("----")
	fmt.Printf("Description: %s\n", profile.Description)
	fmt.Println("----")
	fmt.Printf("Following Count: %d\n", profile.FollowingCount)
	fmt.Printf("Follower count: %d\n", profile.FollowerCount)
	fmt.Printf("Post Count: %d\n", profile.PostCount)

	return nil
}

func (a *app) getUserFeed() error {
	user, quit := input.PromptUserForText("which user feed would you like to view?")
	if quit {
		fmt.Println("You quit. Goodbye")
		return nil
	}

	feed, err := a.service.GetUserFeed(user)
	if err != nil {
		return errors.Wrap(err, "failed to get user feed")
	}

	for _, feedItem := range feed.Feed {
		fmt.Println("-------------------")
		fmt.Printf("Post text: %s\n", feedItem.Post.Record.Text)
		fmt.Printf("Posted at: %s\n", feedItem.Post.Record.CreatedAt)
		fmt.Printf("Post likes: %d\n", feedItem.Post.LikeCount)
		fmt.Printf("Post reposts: %d\n", feedItem.Post.RepostCount)
		fmt.Printf("Post Replies: %d\n", feedItem.Post.ReplyCount)
	}

	return nil
}

func printRawData(data interface{}) error {
	rawJson, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return errors.Wrap(err, "failed to marshal indent data")
	}

	fmt.Println(string(rawJson))
	return nil
}

func getMentionedHandleInMessage(message string) string {
	idx := strings.Index(message, "@")
	words := strings.Fields(message[idx:])
	handle := words[0]

	return handle
}
