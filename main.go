package main

import (
	_ "sheets-plugin/plugin"

	"github.com/Jeffail/benthos/v3/lib/service"
)

func main() {
	service.Run()
}
