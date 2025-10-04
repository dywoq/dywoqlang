package meta

import (
	"fmt"
	"math"
)

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

func Integral(num any) (Type[any], error) {
	switch v := num.(type) {

	case int8, int16, int32, int64, int:
		val := toInt64(v)
		switch {
		case val >= int64(Int8.Min) && val <= int64(Int8.Max):
			return Type[any]{Name: Int8.Name, Size: Int8.Size, Align: Int8.Align, Min: Int8.Min, Max: Int8.Max, Default: Int8.Default}, nil
		case val >= int64(Int16.Min) && val <= int64(Int16.Max):
			return Type[any]{Name: Int16.Name, Size: Int16.Size, Align: Int16.Align, Min: Int16.Min, Max: Int16.Max, Default: Int16.Default}, nil
		case val >= int64(Int32.Min) && val <= int64(Int32.Max):
			return Type[any]{Name: Int32.Name, Size: Int32.Size, Align: Int32.Align, Min: Int32.Min, Max: Int32.Max, Default: Int32.Default}, nil
		default:
			return Type[any]{Name: Int64.Name, Size: Int64.Size, Align: Int64.Align, Min: Int64.Min, Max: Int64.Max, Default: Int64.Default}, nil
		}

	case uint8, uint16, uint32, uint64, uint:
		val := toUint64(v)
		switch {
		case val <= uint64(Uint8.Max):
			return Type[any]{Name: Uint8.Name, Size: Uint8.Size, Align: Uint8.Align, Min: Uint8.Min, Max: Uint8.Max, Default: Uint8.Default}, nil
		case val <= uint64(Uint16.Max):
			return Type[any]{Name: Uint16.Name, Size: Uint16.Size, Align: Uint16.Align, Min: Uint16.Min, Max: Uint16.Max, Default: Uint16.Default}, nil
		case val <= uint64(Uint32.Max):
			return Type[any]{Name: Uint32.Name, Size: Uint32.Size, Align: Uint32.Align, Min: Uint32.Min, Max: Uint32.Max, Default: Uint32.Default}, nil
		default:
			return Type[any]{Name: Uint64.Name, Size: Uint64.Size, Align: Uint64.Align, Min: Uint64.Min, Max: Uint64.Max, Default: Uint64.Default}, nil
		}

	case float32, float64:
		val := toFloat64(v)
		if val >= float64(Float32.Min) && val <= float64(Float32.Max) {
			return Type[any]{Name: Float32.Name, Size: Float32.Size, Align: Float32.Align, Min: Float32.Min, Max: Float32.Max, Default: Float32.Default}, nil
		}
		return Type[any]{Name: Float64.Name, Size: Float64.Size, Align: Float64.Align, Min: Float64.Min, Max: Float64.Max, Default: Float64.Default}, nil

	default:
		return Type[any]{}, fmt.Errorf("unsupported type: %T", num)
	}
}

func toInt64(v any) int64 {
	switch x := v.(type) {
	case int8:
		return int64(x)
	case int16:
		return int64(x)
	case int32:
		return int64(x)
	case int64:
		return x
	case int:
		return int64(x)
	default:
		return 0
	}
}

func toUint64(v any) uint64 {
	switch x := v.(type) {
	case uint8:
		return uint64(x)
	case uint16:
		return uint64(x)
	case uint32:
		return uint64(x)
	case uint64:
		return x
	case uint:
		return uint64(x)
	default:
		return 0
	}
}

func toFloat64(v any) float64 {
	switch x := v.(type) {
	case float32:
		return float64(x)
	case float64:
		return x
	default:
		return 0
	}
}
