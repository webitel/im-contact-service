package main

import (
	"fmt"

	"github.com/webitel/im-contact-service/cmd"
)

func main() {
	if err := cmd.Run(); err != nil {
		fmt.Println(err.Error())
		return
	}
}
