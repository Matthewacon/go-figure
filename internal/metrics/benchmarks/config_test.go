package benchmarks

import (
	"testing"

	"github.com/Matthewacon/go-figure/config"
	"github.com/Matthewacon/go-figure/internal/metrics"
)

func BenchmarkConfigKVInsertionPrebuilt(b *testing.B) {
	cfg := metrics.DefaultEnvAndConfig().GetConfig()
	kv := metrics.IntKeyValue(0)
	for i := 0; i < b.N; i++ {
		cfg.SetParameter(kv, kv)
	}
}

var insertionCfg = metrics.DefaultEnvAndConfig().GetConfig()
func BenchmarkConfigKVInsertion(b *testing.B) {
	for i := 0; i < b.N; i++ {
		kv := metrics.IntKeyValue(i)
		insertionCfg.SetParameter(kv, kv)
	}
}

func BenchmarkConfigKVLookup(b *testing.B) {
	cfg := metrics.DefaultEnvAndConfig().GetConfig()
	for i := 0; i < b.N; i++ {
		kv := metrics.IntKeyValue(i)
		cfg.SetParameter(kv, kv)
	}
	for i := 0; i < b.N; i++ {
		_, _ = cfg.GetParameter(metrics.IntKeyValue(i))
	}
}

func BenchmarkConfigKVLookupPrebuilt(b *testing.B) {
	cfg := metrics.DefaultEnvAndConfig().GetConfig()
	kv := metrics.IntKeyValue(0)
	cfg.SetParameter(kv, kv)
	for i := 0; i < b.N; i++ {
		_, _ = cfg.GetParameter(kv)
	}
}

func BenchmarkConfigKVRemoval(b *testing.B) {
	//setup for benchmark
	cfg := metrics.DefaultEnvAndConfig().GetConfig()
	for i := 0; i < b.N; i++ {
		kv := metrics.IntKeyValue(i)
		cfg.SetParameter(kv, kv)
	}
	for i := 0; i < b.N; i++ {
		kv := metrics.IntKeyValue(i)
		_, _ = cfg.RemoveParameter(kv)
	}
}

var zeroParameterListener = func(value config.IParameterValue, access config.ParameterAccess) error {
	return nil
}

func BenchmarkReadAccessCallback(b *testing.B) {
	cfg := metrics.DefaultEnvAndConfig().GetConfig()
	key := metrics.IntKeyValue(0)
	cfg.AddParameterListener(
		key,
		config.PARAMETER_ACCESS_READ,
		zeroParameterListener,
	)
	for i := 0; i < b.N; i++ {
		//Invoke read access callback
		_, _ = cfg.GetParameter(key)
	}
}

func BenchmarkWriteAccessCallback(b *testing.B) {
	cfg := metrics.DefaultEnvAndConfig().GetConfig()
	cfg.AddParameterListener(
		metrics.IntKeyValue(0),
		config.PARAMETER_ACCESS_WRITE,
		zeroParameterListener,
	)
	key := metrics.IntKeyValue(0)
	for i := 0; i < b.N; i++ {
		//Invoke write access callback
		cfg.SetParameter(key, key)
	}
}

func BenchmarkAnyAccessCallbackOnRead(b *testing.B) {
	cfg := metrics.DefaultEnvAndConfig().GetConfig()
	cfg.AddParameterListener(
		metrics.IntKeyValue(0),
		config.PARAMETER_ACCESS_ANY,
		zeroParameterListener,
	)
	key := metrics.IntKeyValue(0)
	for i := 0; i < b.N; i++ {
		//Invoke read access callback
		_, _ = cfg.GetParameter(key)
	}
}


func BenchmarkAnyAccessCallbackOnWrite(b *testing.B) {
	cfg := metrics.DefaultEnvAndConfig().GetConfig()
	key := metrics.IntKeyValue(0)
	cfg.AddParameterListener(
		key,
		config.PARAMETER_ACCESS_ANY,
		zeroParameterListener,
	)
	for i := 0; i < b.N; i++ {
		//Invoke write access callback
		cfg.SetParameter(key, key)
	}
}

func BenchmarkAnyAccessCallbackOnReadVerticallyScaled(b *testing.B) {
 cfg := metrics.DefaultEnvAndConfig().GetConfig()
	for i := 0; i < b.N; i++ {
		kv := metrics.IntKeyValue(i)
		cfg.SetParameter(kv, kv)
		cfg.AddParameterListener(
			kv,
			config.PARAMETER_ACCESS_ANY,
			func(value config.IParameterValue, access config.ParameterAccess) error {
				_ = i
				return nil
			},
		)
	}
	for i := 0; i < b.N; i++ {
		//Invoke read callbacks
		_, _ = cfg.GetParameter(metrics.IntKeyValue(i))
	}
}

func BenchmarkAnyAccessCallbackOnReadHorizontallyScaled(b *testing.B) {
	cfg := metrics.DefaultEnvAndConfig().GetConfig()
	kv := metrics.IntKeyValue(0)
	cfg.SetParameter(kv, kv)
	for i := 0; i < b.N; i++ {
		cfg.AddParameterListener(
			kv,
			config.PARAMETER_ACCESS_ANY,
			func(value config.IParameterValue, access config.ParameterAccess) error {
				_ = i
				return nil
			},
		)
	}
	for i := 0; i < b.N; i++ {
		//Invoke read access
		_, _ = cfg.GetParameter(kv)
	}
}