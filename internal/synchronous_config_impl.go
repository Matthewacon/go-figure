package internal

import (
	"fmt"
	"sync"
	"unsafe"

	"github.com/Matthewacon/gas"

	"github.com/Matthewacon/go-figure/config"
)

const concurrentModification = "Configuration bus detected a concurrent modification to the listeners!\n"

//TODO deadlock detection
type SynchronousConfigImpl struct {
	parameters      config.Parameters
	listeners       config.ParameterListeners
	currentListener *config.ParameterListener
	mutex           sync.Locker
	errorHandler    config.CallbackErrorHandler
	panicHandler    config.PanicHandler
}

//TODO cover nil panic (when go-away is ready)
func (cfg *SynchronousConfigImpl) detectPanic() {
	if p := recover(); p != nil {
		cfg.panicHandler(p)
	}
}

//TODO event loop checking (spawn a secondary goroutine that combs through the listener trace and looks for loops)
// O(P^2) for {P ∈ Z | 2 <= P <= N/2}, where N is the number of access listeners in the set
func (cfg *SynchronousConfigImpl) pushParameterEvent(key config.IParameterKey, access config.ParameterAccess, prevValue config.IParameterValue) {
	previousListener := cfg.currentListener
	accessListeners, ok := cfg.listeners[key]
	if ok {
		defer func() { cfg.currentListener = previousListener }()
		for _, listener := range *accessListeners {
			//don't recurse into the same listener if a change was made from that listener
			if listener.ParameterAccess & access != 0 && cfg.currentListener != listener.ParameterListener {
				cfg.currentListener = listener.ParameterListener
				invoke := *listener.ParameterListener
				if err := (invoke)(prevValue, access); err != nil {
					cfg.errorHandler(invoke, access, key, err)
				}
			}
		}
	}
}

func (cfg *SynchronousConfigImpl) AddParameterListener(key config.IParameterKey, access config.ParameterAccess, listener config.ParameterListener) {
	defer cfg.detectPanic()
	if cfg.currentListener != nil {
		panic(fmt.Errorf(concurrentModification))
	}
	cfg.mutex.Lock()
	defer cfg.mutex.Unlock()
	accessListeners, ok := cfg.listeners[key]
	if !ok {
		accessListeners = &[]config.ParameterListenerEntry{}
		cfg.listeners[key] = accessListeners
	}
	index := -1
	//Look for existing entry
	for i := 0; i < len(*accessListeners); i++ {
		if (*accessListeners)[i].ParameterListener == &listener {
			index = i
			break
		}
	}
	if index != -1 {
		//Update existing entry
		(*accessListeners)[index].ParameterAccess |= access
	} else {
		//Add new entry
		*accessListeners = append(
			*accessListeners,
			config.ParameterListenerEntry{
				access,
				&listener,
			},
		)
	}
}

func (cfg *SynchronousConfigImpl) RemoveParameterListener(key config.IParameterKey, access config.ParameterAccess, toRemove config.ParameterListener) {
	defer cfg.detectPanic()
	if cfg.currentListener != nil {
		panic(fmt.Errorf(concurrentModification))
	}
	cfg.mutex.Lock()
	defer cfg.mutex.Unlock()
	accessListeners, ok := cfg.listeners[key]
	if ok {
		for _, listener := range *accessListeners {
			if listener.ParameterListener == &toRemove {
				if listener.ParameterAccess & access != 0 {
					listener.ParameterAccess = ^(^listener.ParameterAccess | access)
					//remove listener if it no longer listens on any events
					if listener.ParameterAccess == 0 {
						//create new slice with all elements except this one
						previousListeners := *accessListeners
						*accessListeners = make([]config.ParameterListenerEntry, 0)
						for _, prev := range previousListeners {
							if prev.ParameterListener != &toRemove {
								*accessListeners = append(*accessListeners, prev)
							}
						}
					}
					return
				}
			}
		}
	}
}

func (cfg *SynchronousConfigImpl) GetParameterListeners(key config.IParameterKey) []config.ParameterListenerEntry {
	defer cfg.detectPanic()
	if listeners, ok := cfg.listeners[key]; ok {
		toReturn := make([]config.ParameterListenerEntry, len(*listeners))
		copy(toReturn, *listeners)
		return toReturn
	}
	return []config.ParameterListenerEntry{}
}

func (cfg *SynchronousConfigImpl) SetCallbackErrorHandler(h config.CallbackErrorHandler) config.CallbackErrorHandler {
	gas.AssertNonNil(h)
 prev := cfg.errorHandler
	cfg.errorHandler = h
	return prev
}

func (cfg *SynchronousConfigImpl) SetUnexpectedPanicHandler(h config.PanicHandler) config.PanicHandler {
	gas.AssertNonNil(h)
	prev := cfg.panicHandler
	cfg.panicHandler = h
	return prev
}

func (cfg *SynchronousConfigImpl) GetParameter(key config.IParameterKey) (config.IParameterValue, bool) {
	defer cfg.detectPanic()
	value, ok := cfg.parameters[key]
	cfg.pushParameterEvent(key, config.PARAMETER_ACCESS_READ, value)
	if !ok {
		return nil, false
	}
	return value, true
}

func (cfg *SynchronousConfigImpl) GetParameters() config.Parameters {
	defer cfg.detectPanic()
	//shallow copy parameters to prevent mutations
	params := config.Parameters{}
	for k, v := range cfg.parameters {
		params[k] = v
	}
	//fire read events
	for k, v := range cfg.parameters {
		cfg.pushParameterEvent(k, config.PARAMETER_ACCESS_READ, v)
	}
	return params
}

func (cfg *SynchronousConfigImpl) SetParameter(key config.IParameterKey, value config.IParameterValue) {
	defer cfg.detectPanic()
	prevValue, ok := cfg.parameters[key]
	if !ok {
		prevValue = nil
	}
	cfg.parameters[key] = value
	cfg.pushParameterEvent(key, config.PARAMETER_ACCESS_WRITE, prevValue)
}

func (cfg *SynchronousConfigImpl) SetParameters(params map[config.IParameterKey]config.IParameterValue) {
	defer cfg.detectPanic()
	for key, value := range params {
		prevValue, ok := cfg.parameters[key]
		if !ok {
			prevValue = nil
		}
		cfg.pushParameterEvent(key, config.PARAMETER_ACCESS_WRITE, prevValue)
		cfg.parameters[key] = value
	}
}

func (cfg *SynchronousConfigImpl) RemoveParameter(key config.IParameterKey) (config.IParameterValue, bool) {
	defer cfg.detectPanic()
	value, ok := cfg.parameters[key]
	if ok {
		delete(cfg.parameters, key)
		cfg.pushParameterEvent(key, config.PARAMETER_ACCESS_WRITE, value)
		return value, ok
	}
	return nil, false
}

func NewSynchronousConfigBus() config.IConfigBus {
	return &SynchronousConfigImpl{
		config.Parameters{},
		config.ParameterListeners{},
		nil,
		&sync.Mutex{},
		func(active config.ParameterListener, access config.ParameterAccess, key config.IParameterKey, err error) {
			panic(fmt.Errorf(
				"Encountered unexpected error from listener 0x%x on 0x%x access of [%v]:\n%s\n",
				unsafe.Pointer(&active),
				access,
				key.Key(),
				err.Error(),
			))
		},
		func(p interface{}) {
			panic(p)
		},
	}
}