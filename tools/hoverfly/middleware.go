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
	slog.Info("response Body", "body", body)
}

var (
	isJSONOutput bool
)

func init() {
	flag.BoolVar(&isJSONOutput, "json", false, "Output JSON")
}

func main() {

	flag.Parse()
	if isJSONOutput {
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)).WithGroup("middleware"))
	}

	slog.Info("Hoverfly middleware started")

	s := bufio.NewScanner(os.Stdin)

	processResponse := os.Getenv("PROCESS_RESPONSE") == "true"

	for s.Scan() {
		// If client don't want to process responses - infinite loop
		if !processResponse {
			os.Stdout.Write(s.Bytes())
			continue
		}

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
