package main

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"math"
	"os"
	"time"
)

type Metadata struct {
	Min   string `json:"min"`
	Max   string `json:"max"`
	Count int64  `json:"count"`
	Total string `json:"total"`
	Avg   string `json:"avg"`
}

func main() {
	formatFlag := flag.String("format", "key-value", "Output format (key-value, json or xml)")
	flag.Parse()

	var (
		minDuration   = time.Duration(math.MaxInt64) //set max duration, so any entered value will be less
		maxDuration   time.Duration
		totalDuration time.Duration
		count         int64
	)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		duration, err := time.ParseDuration(line)
		if err != nil && line == "" { //in task description not described sequence stop process, so stop receiving durations once empty string is entered
			break
		} else if err != nil { //if duration cannot be parsed, and it is not empty, lets consider a human error, so continue receiving durations
			fmt.Printf("Error parsing duration %s: %v\n", line, err)
			continue
		}

		if duration < minDuration {
			minDuration = duration
		}

		if duration > maxDuration {
			maxDuration = duration
		}

		totalDuration += duration
		count++
	}

	if count == 0 {
		fmt.Println("No valid durations found.")
		return
	}

	avgDuration := totalDuration / time.Duration(count)

	metadata := Metadata{
		Min:   minDuration.String(),
		Max:   maxDuration.String(),
		Count: count,
		Total: totalDuration.String(),
		Avg:   avgDuration.String(),
	}

	switch *formatFlag {
	case "json":
		jsonData, err := json.Marshal(metadata)
		if err != nil {
			fmt.Printf("Error marshalling to JSON: %v\n", err)
			return
		}
		fmt.Println(string(jsonData))
	case "xml":
		xmlData, err := xml.Marshal(metadata)
		if err != nil {
			fmt.Printf("Error marshalling to XML: %v\n", err)
			return
		}
		fmt.Println(string(xmlData))
	default: // default is key-value
		fmt.Printf("min=%s\nmax=%s\ncount=%d\ntotal=%s\navg=%s\n", metadata.Min, metadata.Max, metadata.Count, metadata.Total, metadata.Avg)
	}
}
