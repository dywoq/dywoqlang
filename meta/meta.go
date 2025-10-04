package meta

import "math"

// Type represents a metadata for DywoqLang type.
type Type[I any] struct {
	Name    string `json:"name"`
	Size    int    `json:"size"`
	Align   int    `json:"align"`
	Min     I      `json:"min"`
	Max     I      `json:"max"`
	Default I      `json:"default"`
}

// Metadata info of type.
var (
	Int8  = Type[int8]{Name: "i8", Size: 1, Align: 1, Min: math.MinInt8, Max: math.MaxInt8, Default: 0}
	Int16 = Type[int16]{Name: "i16", Size: 2, Align: 2, Min: math.MinInt16, Max: math.MaxInt16, Default: 0}
	Int32 = Type[int32]{Name: "i32", Size: 4, Align: 4, Min: math.MinInt32, Max: math.MaxInt32, Default: 0}
	Int64 = Type[int64]{Name: "i64", Size: 8, Align: 8, Min: math.MinInt64, Max: math.MaxInt64, Default: 0}

	Uint8  = Type[uint8]{Name: "u8", Size: 1, Align: 1, Min: 0, Max: math.MaxUint8, Default: 0}
	Uint16 = Type[uint16]{Name: "u16", Size: 2, Align: 2, Min: 0, Max: math.MaxUint16, Default: 0}
	Uint32 = Type[uint32]{Name: "u32", Size: 4, Align: 4, Min: 0, Max: math.MaxUint32, Default: 0}
	Uint64 = Type[uint64]{Name: "u64", Size: 8, Align: 8, Min: 0, Max: math.MaxUint64, Default: 0}

	Float32 = Type[float32]{
		Name: "f32", Size: 4, Align: 4,
		Min: -math.MaxFloat32, Max: math.MaxFloat32,
		Default: 0,
	}

	Float64 = Type[float64]{
		Name: "f64", Size: 8, Align: 8,
		Min: -math.MaxFloat64, Max: math.MaxFloat64,
		Default: 0,
	}
)
