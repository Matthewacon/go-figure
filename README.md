# go-figure
An environment configuration bus with both synchronous and asynchronous implementations.

## Getting started
Add `github.com/Matthewacon/go-figure v0.0.3` to the require section in your `go.mod`.

## Setup
Somewhere in your application, define an environment type that fits your needs and implement `config.IEnvironment`:
```go
package environment

import (
 "github.com/Matthewacon/go-figure"
 "github.com/Matthewacon/go-figure/config"
)

//Define your environment type and implement config.IEnvironment
type Environment struct {
 name string
 config.IConfigBus
}
//config.IEnvironment
func (env *Environment) SetConfig(cfg config.IConfigBus) { env.IConfigBus = cfg }
func (env *Environment) GetConfig() config.IConfigBus    { return env.IConfigBus }
func (env *Environment) String() string                  { return env.name }

//Declare your environments
var (
 ALPHA,
 BETA,
 STAGING,
 PRODUCTION =
 Environment{ name: "ALPHA" },
 Environment{ name: "BETA" },
 Environment{ name: "STAGING" },
 Environment{ name: "PRODUCTION" }
)

//Optionally initialize all of your environments beforehand
func init() {
 for _, env := range []config.IEnvironment{ALPHA, BETA, STAGING, PRODUCTION} {
  go_figure.NewSynchronousConfig(env)
 }
}
```
Then at your application entry point, choose the environment 
```go
package main

import (
 "internal/environments"
 "github.com/Matthewacon/go-figure"
)

func main() {
 //Select your environment and initialize your config bus
 env := environments.ALPHA
 go_figure.NewSynchronousConfig(env)
 //You're now all set to start adding configuration parameters and
 //parameter access listeners!
}
```

### Adding and retrieving configuration parameters
Your configuration requirements may vary depending on the usage within your application. To make the most of your
configuration options, you may define multiple parameter key and value types.
```go
package "sql"

import "github.com/Matthewacon/go-figure/config"

type SqlConfigKey string
//config.IParameterKey
func (k *SqlConfigKey) Key() interface{} { return k }
func (k *SqlConfigKey) String() string { return string(*k) }

type SqlConfigValue int
//config.IParameterValue
func (v *SqlConfigValue) Value() interface{} { return v }
func (v *SqlConfigValue) String() string { return fmt.Sprintf("%d", *v) }

//Sql configuration options
var (
 CONNECTION_TIMEOUT,
 MAX_CONNECTIONS,
 MAX_IDLE_TIME SqlConfigKey =
 "some key identifier",
 "that matches the type",
 "of your custom key"
)

func NewSqlComponent(env config.IEnvironment) {
 cfg := env.GetConfig()
 var timeout SqlConfigValue
 var ok bool
 if timeout, ok = cfg.GetParameter(CONNECTION_TIMEOUT); !ok {
  timeout = SqlConfigValue(32)
  cfg.SetParameter(CONNECTION_TIMEOUT, timeout)
 }
}
```

### Adding configuration parameter listeners
You may want to listen for configuration changes to make live updates to your environment. There are two main
events that you can listen to: `PARAMETER_ACCESS_READ` and `PARAMETER_ACCESS_WRITE`. If you want to handle both
events, you may also use `PARAMETER_ACCESS_ANY`.
```go
package sql

import (
 "fmt"

 "github.com/Matthewacon/go-figure/config"
)

func NewSqlComponent(env config.IEnvironment) {
 cfg := env.GetConfig()
 var timeout SqlConfigValue
 var ok bool
 if timeout, ok = cfg.GetParameter(CONNECTION_TIMEOUT); !ok {
  timeout = SqlConfigValue(32)
  cfg.SetParameter(CONNECTION_TIMEOUT, timeout)
 }
 //listen for any changes to CONNECTION_TIMEOUT
 cfg.AddParameterListener(
  CONNECTION_TIMEOUT,
  config.PARAMETER_ACCESS_WRITE,
  func(prev config.IParameterValue, access config.ParameterAccess) error {
   //update something here
   newValue, ok := cfg.GetParameter(CONNECTION_TIMEOUT)
   if !ok {
    //entry was deleted, do something about it
   }
   timeout = newValue
   return nil
  },
 )
}
``` 

## Running the tests and benchmarks
Tests:
```sh
go test ./...
```

Benchmarks:
```sh
go test --bench=. ./...
```

## Performance
```
BenchmarkConfigKVInsertionPrebuilt-8                   	12965660	       91.6 ns/op
BenchmarkConfigKVInsertion-8                           	3994064	      506 ns/op
BenchmarkConfigKVLookup-8                              	1887486	      812 ns/op
BenchmarkConfigKVLookupPrebuilt-8                      	18696658	       64.3 ns/op
BenchmarkConfigKVRemoval-8                             	1426365	      755 ns/op
BenchmarkReadAccessCallback-8                          	12074638	       98.8 ns/op
BenchmarkWriteAccessCallback-8                         	8582522	      140 ns/op
BenchmarkAnyAccessCallbackOnRead-8                     	12217536	       97.6 ns/op
BenchmarkAnyAccessCallbackOnWrite-8                    	8502228	      141 ns/op
BenchmarkAnyAccessCallbackOnReadVerticallyScaled-8     	1000000	     1909 ns/op
BenchmarkAnyAccessCallbackOnReadHorizontallyScaled-8   	  34530	   119802 ns/op
```


## License
This project is licensed under the [M.I.T. License](https://github.com/Matthewacon/go-figure/blob/master/LICENSE).
