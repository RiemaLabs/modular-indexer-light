package utils

import "encoding/json"

type RawValue interface{ RawMap | RawSlice }

// RawMap is an any-typed map that is compatible with JS-side object representation.
//
// See: js.ValueOf
type RawMap = map[string]any

// RawSlice like RawMap, an any-typed slice compatible with JS-side array representation.
type RawSlice = []any

// Raw marshals a struct into a raw value to feed the JS (WebAssembly) side.
func Raw[T RawValue](v any) T {
	data, _ := json.Marshal(v)
	var out T
	_ = json.Unmarshal(data, &out) // stupid but works
	return out
}
