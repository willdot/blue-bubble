package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMention(t *testing.T) {
	tt := map[string]struct {
		message string
	}{
		"handle at begining": {message: "@willdot hello"},
		"handle in middle":   {message: "hello @willdot how are you"},
		"handle at end":      {message: "hello @willdot"},
		// "comma after handle":  {message: "hello @willdot, how are you"},
		// "comma before handle": {message: "hello ,@willdot"},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			res := getMentionedHandleInMessage(tc.message)
			assert.Equal(t, "@willdot", res)
		})
	}
}
