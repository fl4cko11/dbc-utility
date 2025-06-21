package main

import (
	"os"

	ce "github.com/fl4cko11/dbc-utility/internal/command_execution"
	"github.com/fl4cko11/dbc-utility/logs"
)

func main() {
	logger := logs.InitLogger(os.Stderr)

	ce.CommandExecution(logger)
}
