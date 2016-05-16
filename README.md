# nlog
Node Logger, fast and thread-safe logger for Golang

![Colored](http://i.imgur.com/4V3pR7B.png?1)

Node Logger is a KISS (Keep It Simple, Stupid) logger that  
suits most people. It is much faster than Logrus and Log15, and
asymtotically as fast as Logxi.

The motivation for writing this logger is that none of the other
loggers seems to be thread-safe and even they are buggy. 
I sacrificed fancy features in favour of performance and safety. 

Please feel free to suggest
features and pull requests. 

#### Install
	go get github.com/omidnikta/nlog

#### Example

```go
package main

import "github.com/omidnikta/nlog"

var log *nlog.Logger

func main() {
	formatter := nlog.NewTextFormatter(true, true)
	log = nlog.NewLogger(nlog.DebugLevel, formatter)
	log.Debugf("Logging without a node is possible")
	
	sNode := log.New("Server", nlog.Data{"Protocol": "tcp", "Port": 12542})
	sNode.Debugf("Server is serving")

	cNode := sNode.NewNode("Client", nlog.Data{"RemoteAddr": "[::1]9183", "error": nil})
	cNode.Infof("Client is serving")
}
```

#### Formatters

The built-in logging formatters are:

* `nlog.NewTextFormatter(show_caller, enable_color bool)`.
* `nlog.NewJsonFormatter(show_caller bool)`. Logs fields as JSON.


### Use NLog in your library

The best way to use nlog in your library is to define

```go
package customPackage

import "github.com/omidnikta/nlog"

var log *nlog.Logger

func SetLogger(logger *nlog.Logger) {
	log = logger
}
```

and then use this log in your library. Just remember to 
initialize a Logger in your application and call

```go
package main

import "customPackage"

var logger *nlog.Logger

func main() {
	formatter := nlog.NewTextFormatter(true, true)
	log = nlog.NewLogger(nlog.DebugLevel, formatter)
	
	customPackage.SetLogger(log)	
}
```
