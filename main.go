package main

import (
	"fmt"

	"github.com/willdot/bluebubble/input"
	"github.com/willdot/bluebubble/service"
)

func main() {
	options := []string{
		"create post",
		"users feed",
		"view profile",
	}

	selected, quit := input.PromptUserForSingleChoice(options, "what would you like to do?")
	if quit {
		fmt.Println("You quit. Goodbye")
		return
	}

	err := run(selected)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}
}

func run(selectedOption string) error {
	service, err := service.New()
	if err != nil {
		return err
	}

	app := app{
		service: service,
	}

	switch selectedOption {
	case "create post":
		err = app.newPost()
	case "users feed":
		err = app.getUserFeed()
	case "view profile":
		err = app.getProfile()
	default:
		fmt.Println("Invalid option selected")
		return nil
	}
	return err
}
