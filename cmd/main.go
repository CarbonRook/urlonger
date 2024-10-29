package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/carbonrook/urlonger/pkg/urlonger"
)

var OUTPUT_FORMATS = []string{"text", "json"}
var HEADER_FILTER = []string{"server", "via"}

func main() {

	outputFlag := flag.String("o", "json", fmt.Sprintf("output format; %s", strings.Join(OUTPUT_FORMATS, ", ")))
	urlFlag := flag.String("url", "", "url to resolve")
	verboseFlag := flag.Bool("v", false, "enable verbose logging")
	flag.Parse()

	if !*verboseFlag {
		log.SetOutput(io.Discard)
	}

	redirs, err := urlonger.Resolve(*urlFlag, HEADER_FILTER)
	if err != nil {
		log.Fatal(err)
	}

	var output string
	switch *outputFlag {
	case "json":
		redirsJson, err := json.MarshalIndent(redirs, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		output = string(redirsJson)
	case "text":
		strBuilder := []string{}
		for _, redir := range redirs {
			strBuilder = append(strBuilder, redir.String())
		}
		output = strings.Join(strBuilder, "\n")
	}

	fmt.Println(output)

}
