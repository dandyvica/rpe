// this is where all TOML configuration structures, mapping the TOML topics, are defined
package main

import (
	"strings"
	"log"
	"os/exec"
	"net/url"

	"github.com/BurntSushi/toml"
)

// Global configuration
type tomlConfig struct {
	General generalConfig
	Network networkConfig
	Command []commandDetails
}

// General section 
type generalConfig struct {
	Logfile	string
	Var_Prefix string
}

// Topic/Item configuration
type networkConfig struct {
	Port	int16
	Address	string
}

// Commands as defined in the TOML file
type commandDetails struct {
	Name	string
	Path	string
	Args	[]string
}

//-------------------------------------------------------------------------------------------------
// Methods & func
//-------------------------------------------------------------------------------------------------

// loads configuration from TOML file
func (config *tomlConfig) loadConfig(tomlFileName string) error {
	_, err := toml.DecodeFile(tomlFileName, &config)
	return err
}

// fetch command details from command name. Return false if not found
func (config *tomlConfig) getCommandDetails(cmdName string) (commandDetails, bool) {
	for _, command := range config.Command {
		if command.Name == cmdName {
			// need to deep copy slice
			copiedCommand := command
			copiedCommand.Args = make([]string, len(command.Args))
			copy(copiedCommand.Args, command.Args)

			return copiedCommand, true
		}
	}

	// no command found
	var emptyCommand commandDetails
	return emptyCommand, false
}

// replace variables with their values. Values comes from query string
func (command *commandDetails) substVarsInCommand(url *url.URL, prefix string) {

	// get query values. If no values, no change
	query := url.Query()
	if len(query) == 0 {
		return
	}	

	// if query var starts with the prefix found in the configuration file, replace by its value
	for i, arg := range command.Args {

		// variables should start with a prefix
		if strings.HasPrefix(arg, prefix) {

			// first, delete prefix from argument to be able to compare it
			trimmedArg := strings.TrimLeft(arg, config.General.Var_Prefix)

			// now check this trimmed string into query and replace by its value if any
			if value, inQuery := query[trimmedArg]; inQuery {
				log.Printf("arg %s is in query %s", arg, query)

				// just replace the argument
				command.Args[i] = value[0]
			}
		}
	}	
}

// runs command
func (commandToExecute *commandDetails) exec() (string, int) {
	// execute external command and just print out to stdout
	cmd := exec.Command(commandToExecute.Path)

	// as the first argument will be poped by exec, use this trick by adding a dummy argument as element 0
	cmd.Args = append([]string{" "}, commandToExecute.Args...)

	cmdOut, err := cmd.Output()
	if err != nil {
		log.Printf("error <%v> trying to execute command <%s>", err, commandToExecute)
	}
	log.Printf("cmd=<%+v>", cmd)

	// just return output data from commands
	return string(cmdOut), cmd.ProcessState.ExitCode()

}