package config

//Interface type for environments
type IEnvironment interface {
	SetConfig(cfg IConfigBus)
	GetConfig() IConfigBus
	String() string
}

//Interface type for all config parameter keys
type IParameterKey interface {
	Key() interface{}
	String() string
}

//Interface type for all config parameter values
type IParameterValue interface {
	Value() interface{}
	String() string
}

type Parameters map[IParameterKey]IParameterValue

type ParameterAccess uint8
const (
	PARAMETER_ACCESS_RESERVED,
	PARAMETER_ACCESS_READ,
	PARAMETER_ACCESS_WRITE,
	PARAMETER_ACCESS_ANY ParameterAccess =
		0,
		1 << 0,
		1 << 1,
		PARAMETER_ACCESS_READ | PARAMETER_ACCESS_WRITE
)

type ParameterListener func(prev IParameterValue, access ParameterAccess) error
type ParameterListenerEntry struct {
	ParameterAccess
	*ParameterListener
}
type ParameterListeners map[IParameterKey]*[]ParameterListenerEntry

type CallbackErrorHandler func(active ParameterListener, access ParameterAccess, key IParameterKey, err error)
//The panic handler should cover all functions EXCEPT: SetCallbackErrorHandler, SetUnexpectedPanicHandler
type PanicHandler func(p interface{})

type IConfigBus interface {
	AddParameterListener(key IParameterKey, access ParameterAccess, listener ParameterListener)
	RemoveParameterListener(key IParameterKey, access ParameterAccess, listener ParameterListener)
	GetParameterListeners(key IParameterKey) []ParameterListenerEntry
 SetCallbackErrorHandler(handler CallbackErrorHandler) CallbackErrorHandler
	SetUnexpectedPanicHandler(handler PanicHandler) PanicHandler
	GetParameter(key IParameterKey) (IParameterValue, bool)
	GetParameters() Parameters
	SetParameter(key IParameterKey, value IParameterValue)
	SetParameters(params map[IParameterKey]IParameterValue)
	RemoveParameter(key IParameterKey) (IParameterValue, bool)
}
