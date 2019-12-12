// handle request from remote clients
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// function to log http request properties
func logRequest(req *http.Request) {
	log.Printf("http method <%s> received, url=<%s>, header=<%s>, ip=<%s>", req.Method, req.URL, req.Header, req.RemoteAddr)
}

// this function is called whenever the server receives a HTTP GET request
func CommandHandler(res http.ResponseWriter, req *http.Request) {
	// log request data
	logRequest(req)

	// get command from URL
	vars := mux.Vars(req)
	command := vars["command"]
	log.Printf("received command <%s>", command)

	// first, check whether the command is declared in the TOML configuration file
	// and if query parameters match command parameters
	runnableCommand, found := config.getCommandDetails(command) 
	if !found {
		log.Printf("command <%s> not found !", command)
		res.WriteHeader(http.StatusNotFound)
		return
	}

	// now substitue variables if any
	runnableCommand.substVarsInCommand(req.URL, config.General.Var_Prefix)

	// and execute command
	log.Printf("runnable command details: <%+v>", runnableCommand)
	output_data, exit_code := runnableCommand.exec()

	// marshall data to JSON
	responseData := CommandResponse{ ExitCode: exit_code, OutputData: output_data }

	// and return result to caller
	res.WriteHeader(http.StatusOK)
	fmt.Fprintf(res, responseData.toJSON())
}