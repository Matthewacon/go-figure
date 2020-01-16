package go_figure

import (
	"fmt"

	"github.com/Matthewacon/go-figure/config"
	"github.com/Matthewacon/go-figure/internal"
)

func NewSynchronousConfig(env config.IEnvironment) config.IConfigBus {
	cfg := internal.NewSynchronousConfigBus()
	env.SetConfig(cfg)
	return cfg
}

//TODO
func NewAsynchronousConfig(env config.IEnvironment) config.IConfigBus {
	panic(fmt.Errorf("unimplemented!\n"))
	return nil
}