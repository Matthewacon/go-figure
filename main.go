package go_figure

import (
	"fmt"
	"go-figure/config"
	"go-figure/internal"
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