package main

import (
	"bufio"
	"encoding/json"
	"flag"
	hoverfly "github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"log/slog"
	"os"
)

func processJSONResponse(view hoverfly.RequestResponsePairViewV1) {
	body := view.GetResponse().GetBody()
	slog.Info(body, "response Body")
}

var isJSONOutput bool

func init() {
	flag.BoolVar(&isJSONOutput, "json", false, "Output JSON")
}

func main() {

	flag.Parse()
	if isJSONOutput {
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	}

	s := bufio.NewScanner(os.Stdin)

	for s.Scan() {

		var payload hoverfly.RequestResponsePairViewV1

		err := json.Unmarshal(s.Bytes(), &payload)
		if err != nil {
			slog.Error("Failed to unmarshal payload from hoverfly")
		}

		processJSONResponse(payload)

		bts, err := json.Marshal(payload)
		if err != nil {
			slog.Error("Failed to marshal new payload")
		}

		os.Stdout.Write(bts)
	}
}
