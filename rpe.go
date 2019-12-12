// read TOML config file
package main

import (
	"fmt"
	//"bufio"
	"os"
	"log"
	"net/http"

	// imported packages
	"github.com/gorilla/mux"
)

// global configuration: this is where whole config is kept
var config tomlConfig;

func main() {
	// check number of arguments
	if len(os.Args) < 2 {
		fmt.Println("needs the TOML file as an argument")
		os.Exit(1)
	}

	// get first argument: TOML file name
	tomlFileName := os.Args[1]

	// now marshall toml data
	if err := config.loadConfig(tomlFileName); err != nil {
		fmt.Printf("error <%v> when opening toml file <%s>\n", err)
		os.Exit(2)
	}

	// open or create log file
	logfile, err := os.OpenFile(config.General.Logfile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
		os.Exit(3)
	}

	// close log at main's exit and set this file for log
	defer logfile.Close()
	log.SetOutput(logfile)

	// print out configuration file 
	log.Printf("config=%+v", config)

	// Check whether commands are found
	for _, cmd := range config.Command {
		if _, err := os.Stat(cmd.Path); err != nil {
			log.Printf("unable to find command <name=%s, path=%s>", cmd.Name, cmd.Path)
			os.Exit(4)
		}
	}

	// http loop and receive commands
	router := mux.NewRouter()

	// only handle GET requests
	router.HandleFunc("/{command}", CommandHandler).Methods("GET")

	// build address:port and listen HTTP incoming requests
	addr := fmt.Sprintf("%s:%d", config.Network.Address, config.Network.Port)
	http.ListenAndServe(addr, router)

}