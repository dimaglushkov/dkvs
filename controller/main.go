package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/dimaglushkov/dkvs/controller/internal"
)

func run() error {
	readFileLines := func(path string) ([]string, error) {
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		var lines []string
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		return lines, scanner.Err()
	}

	rf, err := strconv.ParseInt(os.Getenv("REPLICA_FACTOR"), 10, 32)
	if err != nil {
		return fmt.Errorf("error while parsing REPLICA_FACTOR env var (%s): %s", os.Getenv("REPLICA_FACTOR"), err)
	}

	storageLocs, err := readFileLines(os.Getenv("STORAGES_FILE"))
	if err != nil {
		return fmt.Errorf("error while reading %s: %s", os.Getenv("STORAGES_FILE"), err)
	}

	handler := internal.NewHandler(storageLocs, int(rf))

	return http.ListenAndServe(":"+os.Getenv("APP_PORT"), handler)
}

func main() {
	requiredEnvVars := []string{"STORAGES_FILE", "APP_PORT", "REPLICA_FACTOR"}
	for _, i := range requiredEnvVars {
		if os.Getenv(i) == "" {
			log.Fatalf("required %s env variable is not set", i)
		}
	}

	if err := run(); err != nil {
		log.Fatal(err)
	}
}
