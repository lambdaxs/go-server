package confu

import (
	"fmt"
	"testing"
)

type Config struct {
	Server struct {
		Ver string `toml:"ver"`
	} `toml:"server"`
}

func TestInitWithFilePath(t *testing.T) {
	data := Config{}
	if err := InitWithFilePath("./sample.toml", &data); err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(data.Server.Ver)
}

func TestInitWithRemotePath(t *testing.T) {
	data := Config{}
	if err := InitWithRemotePath("sample-remote.toml", &data, "127.0.0.1:8500"); err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(data.Server.Ver)
}
