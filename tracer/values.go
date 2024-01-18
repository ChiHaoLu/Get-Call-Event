// Ref. https://github.com/stackup-wallet/stackup-bundler/blob/main/pkg/tracer/values.go

package tracer

var (
	// Loaded JS tracers for simulating various EntryPoint methods using debug_traceCall.
	Loaded, _ = NewTracers()
)