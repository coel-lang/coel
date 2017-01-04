package main

import (
	f "../../lib/functions"
	_ "../../lib/parse"
	"github.com/docopt/docopt-go"
)

func main() {
	getArgs()

	f.Nil().Eval()
}

func getArgs() map[string]interface{} {
	usage := `Risp interpreter

Usage:
  risp <filename>

Options:
  -h --help     Show this help.`

	args, _ := docopt.Parse(usage, nil, true, "Risp 0.0.0", false)

	return args
}
