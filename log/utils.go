package log

import (
	"git.garena.com/shopee/common/ulog"
	"git.garena.com/shopee/platform/splog"
	"git.garena.com/shopee/platform/splog/core"
)

func toSplogTypedField(field ulog.TypedField) splog.TypedField {
	return splog.TypedField{
		Key:       field.Key(),
		Type:      toSplogFieldType(field.FieldType()),
		Integer:   field.NumericValue(),
		String:    field.StringValue(),
		Interface: field.InterfaceValue(),
	}
}

func toSplogLevel(level ulog.Level) (splog.Level, error) {
	lvl, err := core.ParseLevel(level.Name)
	return splog.Level(lvl), err
}

func toSplogFieldType(fieldType ulog.FieldType) core.FieldType {
	switch fieldType {
	case ulog.StringType:
		return core.StringType
	case ulog.BoolType:
		return core.BoolType
	case ulog.Int8Type:
		return core.Int8Type
	case ulog.Int16Type:
		return core.Int16Type
	case ulog.Int32Type:
		return core.Int32Type
	case ulog.Int64Type:
		return core.Int64Type
	case ulog.Uint8Type:
		return core.Uint8Type
	case ulog.Uint16Type:
		return core.Uint16Type
	case ulog.Uint32Type:
		return core.Uint32Type
	case ulog.Uint64Type:
		return core.Uint64Type
	case ulog.Float32Type:
		return core.Float32Type
	case ulog.Float64Type:
		return core.Float64Type
	case ulog.ErrorType:
		return core.ErrorType
	case ulog.ReflectType:
		return core.ReflectType
	case ulog.StringerType:
		return core.StringerType
	case ulog.SkipType:
		return core.SkipType
	default:
		return core.UnknownType
	}
}
