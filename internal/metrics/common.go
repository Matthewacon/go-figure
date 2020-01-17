package metrics

import (
 "fmt"
 "testing"

 "github.com/Matthewacon/go-figure"
 "github.com/Matthewacon/go-figure/config"
)

type Environment struct {
 string
 config.IConfigBus
}

//IEnvironment
func (e *Environment) SetConfig(cfg config.IConfigBus) { e.IConfigBus = cfg }
func (e *Environment) GetConfig() config.IConfigBus    { return e.IConfigBus }
func (e *Environment) String() string                  { return e.string }
func (e *Environment) IsLive() bool                    { return e.IConfigBus != nil }

func DefaultEnvAndConfig() config.IEnvironment {
 env := &Environment{ string: "env" }
 go_figure.NewSynchronousConfig(env)
 return env
}

func EnvWithNameAndConfig(name string) config.IEnvironment {
 env := &Environment{ string: name }
 env.SetConfig(env)
 return env
}

func CatchUnexpectedPanic(t *testing.T) {
 if err := recover(); err != nil {
  t.Errorf(fmt.Sprintf("%v", err))
 }
}

func CatchExpectedPanic(t *testing.T) {
 if err := recover(); err == nil {
  t.Errorf("Expected panic!\n")
 }
}

type IntKeyValue int
//IParameterKey and IParameterValue
func (i IntKeyValue) String() string { return fmt.Sprintf("%d", i) }
func (i IntKeyValue) Key() interface{} { return i }
func (i IntKeyValue) Value() interface{} { return i }
