package tests

import (
	"fmt"
	"go-figure/config"
	"testing"

	common "go-figure/internal/metrics"
)

func TestKVInsertionAndRetrieval(t *testing.T) {
	defer common.CatchUnexpectedPanic(t)
	env := common.DefaultEnvAndConfig()
	cfg := env.GetConfig()
	count := 100000
	//insert count kv config pairs
	for i := 0; i < count; i++ {
		cfg.SetParameter(
			common.IntKeyValue(i),
			common.IntKeyValue(-i),
		)
	}
	//verify total entry count
	if num := len(cfg.GetParameters()); num != count {
		t.Errorf(
			"Failed to insert all KV parameter pairs, inserted: %d, retrieved: %d\n",
			count,
			num,
		)
		return
	}
	//verify config pairs
	for i := 0; i < count; i++ {
		value, ok := cfg.GetParameter(common.IntKeyValue(i))
		if !ok {
			t.Errorf(
				"Missing kv pair [%d,%d]!\n",
				i, -i,
			)
			return
		}
		if value != common.IntKeyValue(-i) {
			t.Errorf(
				"Failed to verify kv config pair, expected: [%d,%d] received: [%d,%d]\n",
				i, -i,
				i, value.Value(),
			)
			return
		}
	}
}

func addAccessCallback(t *testing.T, cfg config.IConfigBus, kv common.IntKeyValue) {
	cfg.AddParameterListener(
		kv,
		config.PARAMETER_ACCESS_WRITE,
		func(value config.IParameterValue, access config.ParameterAccess) error {
			return nil
		},
	)
	if len(cfg.GetParameterListeners(kv)) != 1 {
		t.Errorf(
			"Failed to add paramter listener for access 0x%x on [%v]\n",
			config.PARAMETER_ACCESS_WRITE,
			kv.Key(),
		)
	}
}

func TestAddAccessCallback(t *testing.T) {
	defer common.CatchUnexpectedPanic(t)
	cfg := common.DefaultEnvAndConfig().GetConfig()
	kv := common.IntKeyValue(0)
	addAccessCallback(t, cfg, kv)
}

func TestRemoveAccessCallback(t *testing.T) {
	defer common.CatchUnexpectedPanic(t)
	cfg := common.DefaultEnvAndConfig().GetConfig()
	kv := common.IntKeyValue(0)
	addAccessCallback(t, cfg, kv)
	listener := cfg.GetParameterListeners(kv)[0]
	cfg.RemoveParameterListener(kv, listener.ParameterAccess, *listener.ParameterListener)
}

func addCheckedListener(t *testing.T, cfg config.IConfigBus, key config.IParameterKey, value config.IParameterValue, access config.ParameterAccess) {
	cfg.AddParameterListener(
		key,
		access,
		func(val config.IParameterValue, acc config.ParameterAccess) error {
			if acc != access {
				t.Errorf(
					"Callback for key '%d' expected access type: 0x%x, received type: 0x%x\n",
					key.Key(),
					access,
					acc,
				)
				return nil
			}
			if val != value {
				t.Errorf(
					"Callback for key '%d' expected value: '%d', received value: '%d'\n",
					key.Key(),
					value.Value(),
					val.Value(),
				)
				return nil
			}
			return nil
		},
	)
	//Check that callback was registered
	if len(cfg.GetParameterListeners(key)) != 1 {
		t.Errorf(
			"Failed to register parameter access listener for key '%d'\n",
			key.Key(),
		)
		return
	}
}

func TestKeyReadAccessCallback(t *testing.T) {
	defer common.CatchUnexpectedPanic(t)
	cfg := common.DefaultEnvAndConfig().GetConfig()
	count := 100
	for i := 0; i < count; i++ {
		key, value := common.IntKeyValue(i), common.IntKeyValue(i)
		cfg.SetParameter(key, value)
		addCheckedListener(t, cfg, key, value, config.PARAMETER_ACCESS_READ)
	}
	for i := 0; i < count; i++ {
		kv := common.IntKeyValue(i)
		//Invoke read access callbacks
		_, _ = cfg.GetParameter(common.IntKeyValue(i))
		//Should not invoke read access callbacks
		cfg.SetParameter(kv, kv)
	}
}

func TestKeyWriteAccessCallback(t *testing.T) {
	defer common.CatchUnexpectedPanic(t)
	cfg := common.DefaultEnvAndConfig().GetConfig()
	count := 100
	for i := 0; i < count; i++ {
		kv := common.IntKeyValue(i)
		cfg.SetParameter(kv, kv)
		addCheckedListener(t, cfg, kv, kv, config.PARAMETER_ACCESS_WRITE)
	}
 for i := 0; i < count; i++ {
 	kv := common.IntKeyValue(i)
		//Invoke write access callbacks
 	cfg.SetParameter(kv, kv)
		//Should not invoke read access callbacks
		_, _ = cfg.GetParameter(kv)
 }
}

func TestKeyAnyAccessCallback(t *testing.T) {
	defer common.CatchUnexpectedPanic(t)
	cfg := common.DefaultEnvAndConfig().GetConfig()
	count := 100
	calls := map[config.IParameterKey]int{}
	for i := 0; i < count; i++ {
		kv := common.IntKeyValue(i)
		cfg.SetParameter(kv, kv)
		cfg.AddParameterListener(
			kv,
			config.PARAMETER_ACCESS_ANY,
			func(value config.IParameterValue, access config.ParameterAccess) error {
				if value != kv {
					t.Errorf(
						"Callback for key '%d' expected value: '%d', received value: '%d'\n",
						kv,
						kv,
						value,
					)
				}
				calls[kv] += 1
				return nil
			},
		)
	}
	for i := 0; i < count; i++ {
		kv := common.IntKeyValue(i)
		//Invoke read access callbacks
		_, _ = cfg.GetParameter(kv)
		//Invoke write access callbacks
		cfg.SetParameter(kv, kv)
		if nCalls := calls[kv]; nCalls != 2 {
			t.Errorf(
				"Parameter callback registered with access 0x%x did not fire for every access event: expected: %d, fired: %d\n",
				config.PARAMETER_ACCESS_ANY,
				2,
				nCalls,
			)
			return
		}
	}
}

func TestDefaultCallbackErrorHandler(t *testing.T) {
	defer common.CatchExpectedPanic(t)
	cfg := common.DefaultEnvAndConfig().GetConfig()
	kv := common.IntKeyValue(0)
	cfg.AddParameterListener(
		kv,
		config.PARAMETER_ACCESS_ANY,
		func(value config.IParameterValue, access config.ParameterAccess) error {
			return fmt.Errorf("")
		},
	)
	_, _ = cfg.GetParameter(kv)
}

func TestCustomCallbackErrorHandler(t *testing.T) {
	defer common.CatchUnexpectedPanic(t)
	cfg := common.DefaultEnvAndConfig().GetConfig()
	kv := common.IntKeyValue(0)
	cfg.AddParameterListener(
		kv,
		config.PARAMETER_ACCESS_ANY,
		func(value config.IParameterValue, access config.ParameterAccess) error {
			return fmt.Errorf("")
		},
	)
	cfg.SetCallbackErrorHandler(
		func(active config.ParameterListener, access config.ParameterAccess, key config.IParameterKey, err error) {
			//silently ignore listener errors
		},
	)
	_, _ = cfg.GetParameter(kv)
}

func TestDefaultUnexpectedPanicHandler(t *testing.T) {
	defer common.CatchExpectedPanic(t)
	cfg := common.DefaultEnvAndConfig().GetConfig()
	kv := common.IntKeyValue(0)
	cfg.AddParameterListener(
		kv,
		config.PARAMETER_ACCESS_ANY,
		func(value config.IParameterValue, access config.ParameterAccess) error {
			panic(fmt.Errorf(""))
		},
	)
	_, _ = cfg.GetParameter(kv)
}

func TestCustomUnexpectedPanicHandler(t *testing.T) {
	defer common.CatchUnexpectedPanic(t)
	cfg := common.DefaultEnvAndConfig().GetConfig()
	kv := common.IntKeyValue(0)
	cfg.AddParameterListener(
		kv,
		config.PARAMETER_ACCESS_ANY,
		func(value config.IParameterValue, access config.ParameterAccess) error {
			return fmt.Errorf("")
		},
	)
	cfg.SetUnexpectedPanicHandler(func(p interface{}) {
		//silently ignore panics
	})
	_, _ = cfg.GetParameter(kv)
}

//TODO
//func TestEventLoopCheck(t *testing.T) {
//	defer common.CatchExpectedPanic(t)
//	cfg := common.DefaultEnvAndConfig().GetConfig()
//	kv := common.IntKeyValue(0)
//	for i := 0; i < 2; i++ {
//		cfg.AddParameterListener(
//			kv,
//			config.PARAMETER_ACCESS_READ,
//			func(value config.IParameterValue, access config.ParameterAccess) error {
//				_ = i
//				cfg.GetParameter(kv)
//				return nil
//			},
//		)
//	}
//	_ = cfg.GetParameter(kv)
//}

func TestConcurrentModificationCheck(t *testing.T) {
	defer common.CatchExpectedPanic(t)
	cfg := common.DefaultEnvAndConfig().GetConfig()
	kv := common.IntKeyValue(0)
	cfg.AddParameterListener(
		kv,
		config.PARAMETER_ACCESS_READ,
		func(value config.IParameterValue, access config.ParameterAccess) error {
			cfg.AddParameterListener(
				kv,
				config.PARAMETER_ACCESS_READ,
				nil,
			)
			return nil
		},
	)
	_, _ = cfg.GetParameter(kv)
}
