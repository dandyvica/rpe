// JSON data returned to client
package main

import (
	"log"
	"encoding/json"
)

// structure marshalled to JSON for sending response to client
type CommandResponse struct {
	ExitCode	int  	`json:"exit_code"`
	OutputData	string	`json:"output_data"`
}

// marshall data to JSON
func (res CommandResponse) toJSON() string {

	jsonData, err := json.Marshal(res)
	if err != nil {
		log.Printf("error encoding to JSON")
		return ""
	}

	return string(jsonData)
}