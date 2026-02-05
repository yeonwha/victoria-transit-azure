package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"google.golang.org/protobuf/proto"
)

type InvokeRequest struct {
	Data     map[string]interface{}
	Metadata map[string]interface{}
}

func loadStopNames(filename string) (map[string]string, error) {
	// Parse stop ids to stop name by reading the csv file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// create a map to hold stop IDs and names
	stopMap := make(map[string]string)
	// create a new CSV reader
	reader := csv.NewReader(file)

	_, _ = reader.Read() // skip header row

	// read each record from the CSV
	for {
		record, err := reader.Read()
		// if we reach the end of the file, break the loop
		if err == io.EOF {
			break
		}
		// handle any other errors
		if err != nil {
			return nil, err
		}
		// Map the stop ID to the stop name
		stopID := record[0]
		stopName := record[2]
		stopMap[stopID] = stopName
	}
	// return the map of stop IDs to names
	return stopMap, nil
}

func main() {
	// Load stop names map from the CSV file
	stopNames, err := loadStopNames("Victoria_Regional_Transit_System_stops.csv")
	if err != nil {
		log.Fatal("Failed to load the stops file:", err)
	}

	// Read the Protobuf file to get real-time transit data
	data, err := os.ReadFile("tripupdates.pb")
	if err != nil {
		log.Fatal("Failed to read the .pb file. Please check if the file exists in the folder:", err)
	}
	// Unmarshal the Protobuf data into a FeedMessage struct
	feed := &gtfs.FeedMessage{}
	if err := proto.Unmarshal(data, feed); err != nil {
		log.Fatal("Failed to unmarshal Protobuf data:", err)
	}
	// Iterate through each entity in the feed and print stop updates with names and delays
	fmt.Println("=== Victoria Transit Real-time Data ===")
	for _, entity := range feed.Entity {
		// Check if the entity has a TripUpdate
		if entity.TripUpdate != nil {
			// Update stop id to stop name mapping
			for _, stopUpdate := range entity.TripUpdate.StopTimeUpdate {
				stopID := *stopUpdate.StopId
				name, ok := stopNames[stopID]
				if !ok {
					name = "Unknown Stop"
				}

				// Determine the delay for the stop update
				var delay int32 = 0
				// Check if the departure is set, otherwise check arrival and update delay accordingly
				if stopUpdate.Departure != nil && stopUpdate.Departure.Delay != nil {
					delay = *stopUpdate.Departure.Delay
				} else if stopUpdate.Arrival != nil && stopUpdate.Arrival.Delay != nil {
					delay = *stopUpdate.Arrival.Delay
				}
				// Print the stop name and delay if delay is non-zero
				if delay != 0 {
					fmt.Printf("Stop: %s (%s) | Delay: %d Sec(s)\n", name, stopID, delay)
				}
			}
		}
	}
}
