A remote program execution in Go

This time, I'll be tackling the Go programming language. Even though its synntax is not the most inspiring and it's missing some modern features like generics (but this is well requested by the Go community and it'll be probably part of the _Go 2_ umbrella project), I was ready to give a try. Being backed by Google makes this a key advantage.

I was curious of how long I could write a decent WEB API project, starting from zero knowledge. Of course being proficient at C/C++ and intermediate in Rust and D, helps.

# The project

As I'm currently dabbling with Nagios those days, I discovered how the Nagios Remote Execution Program works. But also how terse the source code was. I decided to experiment a _Remote Execution Program_ in Go, by developping a REST simple WEB API mockup. The interesting fact in Go resides in is its speed and in the fact it's generating a standalone executable, which is a key element when distributing your app.

The general idea is this: execute a remote program or script in the remote server running the app, by sending a HTTP GET command with a path and optionnally a query string. It's easy to test using _wget_ or even _postman_. On the server side, the list of available commands is read from a _TOML_ file and result is delivered to the client as a JSON string.

Of course, this could be done using Python _Flask_ or Ruby _Sinatra_ but the point here was to experiment Go and see how far I can land.

# The server side

Basically, the server app will:

* read the TOML file and marshal its data into a Go struct
* verify commands defined in the TOML file are found locally
* define a route to handle GET request
* start the WEB server on port specified in the TOML file
* whenever a GET request is received, exec and return output data to the client, along with the command exit code

## Managing TOML data
First of all, the _TOML_ file could look like this:

```toml
# RPE configuration file sample
[general]
logfile = "/tmp/rpe.log"
var_prefix = "$"

# network specifics
[network]
address = "127.0.0.1"
port = 8080

[[command]]
name = "list_files"
path = "/bin/ls"
args = ["-l", "$DIR"]

[[command]]
name = "find_files"
path = "/usr/bin/find"
args = ["$ROOT", "-name", "$PATTERN"]
```

The _[[command]]_ section defines command path and filename, its arguments with potential variable arguments and a name to refer to. That name will be used as the GET route, and its query string as variable arguments. The Go _TOML_ package is used to marshal data to a Go struct, with field names mapping TOML data and sections:

```go
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
```

The marshalling TOML data is straighforward:

```go
// now marshall toml data
if err := config.loadConfig(tomlFileName); err != nil {
    fmt.Println(err)
    os.Exit(2)
}
```

## WEB framework

It seems _mux_ is the _de facto_ standard as a HTTP WEB framework. It works similarly to what's usual in _Flask_ or _Sinatra_, binding routes to method handlers:

```go
// http loop and receive commands
router := mux.NewRouter()

// only handle GET requests
router.HandleFunc("/{command}", CommandHandler).Methods("GET")

// build address:port and listen HTTP incoming requests
addr := fmt.Sprintf("%s:%d", config.Network.Address, config.Network.Port)
http.ListenAndServe(addr, router)
```

It merely binds the _CommandHandler_ Go function, only for _GET_ requests and for the variable URL path _{command}_. The handler is then responsible for getting command characteristics from the TOML configuration, and if found, it spawns the command and return the results as JSON. 

All the details here:

## Using the app

Just use _wget_ like this:

```console
$ wget -O - http://127.0.0.1:8080/find_files\?ROOT\=/home\&PATTERN\=\*.go 2>/dev/null 
```

which gives a list of _*.go_ files from the _/home_ directory.

# Possible refinements

It's easy to add other security features:

* use https
* use authentication
* use host filter to allow connexion only from specific hosts
* implement a timer to kill a lengthy command

# Conclusion

I was really surprised at how fast I learned Go. Being a C expert, it took only 2 days googling Go documentation and Stackoverflow posts to be productive. 

Of course, it's not as powerful as Rust, because it lacks lots of modern features like generics or sum types. And I found some constructs somewhat weird like initializing a map, and error management is rather old style.

Being a garbage collected language, it's not meant to be a true system's programming language, but rather targeted at **devops**. I think it has a great potential and a bright future when Go 2 will be rolled out with generics. Its standard library is large and covers lots of features.

You can get the whole source code at: 
