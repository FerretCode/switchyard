package main

import (
	scheduler "github.com/ferretcode/switchyard/scheduler/pkg/types"
)

var ConfigRegistry = map[string]any{
	"scheduler-service": scheduler.Config{},
}

func main() {

}
