package test

import (
	"github.com/stretchr/testify/assert"
	"goChat/conf"
	"testing"
)

func Test_LoadConfig(t *testing.T) {
	conf.InitGlobalConfig("test_config.ini")
	a := assert.New(t)
	a.EqualValues("hello", conf.GetConfig("test1", "url"))
	a.EqualValues("chatExcel", conf.GetConfig("test2", "key"))
	a.EqualValues("", conf.GetConfig("test2", "default_prompt"))

}
