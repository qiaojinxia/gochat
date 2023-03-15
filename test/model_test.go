package test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"goChat/models"
	"strconv"
	"testing"
)

func Test_model(t *testing.T) {
	for i := 0; i < 10; i++ {
		err := models.SetUserChatContext("caomaoboy", strconv.Itoa(i))
		if err != nil {
			a := assert.New(t)
			a.Error(err,"test error")
		}
	}
	allContents, err := models.GetUserChatContext("caomaoboy")
	if err != nil {
		a := assert.New(t)
		a.Error(err,"test error")
	}
	fmt.Sprintf("%v\n", allContents)
}
