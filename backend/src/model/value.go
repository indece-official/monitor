package model

import (
	"fmt"
	"strconv"
	"time"
)

type ValueType string

const (
	ValueTypeNumber   ValueType = "number"
	ValueTypeText     ValueType = "text"
	ValueTypeDate     ValueType = "date"
	ValueTypeDatetime ValueType = "datetime"
	ValueTypeDuration ValueType = "duration"
)

type Value struct {
	valueType     ValueType
	value         string
	float64Value  float64
	timeValue     time.Time
	durationValue time.Duration
}

func (v *Value) LessThan(other *Value) (bool, error) {
	if other.valueType != v.valueType {
		return false, fmt.Errorf("value types don't match: %s vs %s", v.valueType, other.valueType)
	}

	switch v.valueType {
	case ValueTypeNumber:
		return v.float64Value < other.float64Value, nil
	case ValueTypeDate,
		ValueTypeDatetime:
		return v.timeValue.Compare(other.timeValue) < 0, nil
	case ValueTypeDuration:
		return v.durationValue < other.durationValue, nil
	default:
		return false, fmt.Errorf("less-than conmpareator not supported for type %s", v.valueType)
	}
}

func (v *Value) LessThanEqual(other *Value) (bool, error) {
	if other.valueType != v.valueType {
		return false, fmt.Errorf("value types don't match: %s vs %s", v.valueType, other.valueType)
	}

	switch v.valueType {
	case ValueTypeNumber:
		return v.float64Value <= other.float64Value, nil
	case ValueTypeDate,
		ValueTypeDatetime:
		return v.timeValue.Compare(other.timeValue) <= 0, nil
	case ValueTypeDuration:
		return v.durationValue <= other.durationValue, nil
	default:
		return false, fmt.Errorf("less-than-equal conmpareator not supported for type %s", v.valueType)
	}
}

func (v *Value) GreaterThan(other *Value) (bool, error) {
	if other.valueType != v.valueType {
		return false, fmt.Errorf("value types don't match: %s vs %s", v.valueType, other.valueType)
	}

	switch v.valueType {
	case ValueTypeNumber:
		return v.float64Value > other.float64Value, nil
	case ValueTypeDate,
		ValueTypeDatetime:
		return v.timeValue.Compare(other.timeValue) > 0, nil
	case ValueTypeDuration:
		return v.durationValue > other.durationValue, nil
	default:
		return false, fmt.Errorf("greater-than conmpareator not supported for type %s", v.valueType)
	}
}

func (v *Value) GreaterThanEqual(other *Value) (bool, error) {
	if other.valueType != v.valueType {
		return false, fmt.Errorf("value types don't match: %s vs %s", v.valueType, other.valueType)
	}

	switch v.valueType {
	case ValueTypeNumber:
		return v.float64Value >= other.float64Value, nil
	case ValueTypeDate,
		ValueTypeDatetime:
		return v.timeValue.Compare(other.timeValue) >= 0, nil
	case ValueTypeDuration:
		return v.durationValue >= other.durationValue, nil
	default:
		return false, fmt.Errorf("greater-than-equal conmpareator not supported for type %s", v.valueType)
	}
}

func (v *Value) Raw() interface{} {
	switch v.valueType {
	case ValueTypeNumber:
		return v.float64Value
	case ValueTypeDate,
		ValueTypeDatetime:
		return v.timeValue
	case ValueTypeDuration:
		return v.durationValue
	case ValueTypeText:
		return v.value
	default:
		return v.value
	}
}

func NewValue(valueType ValueType, value string) (*Value, error) {
	var err error

	val := &Value{}
	val.valueType = valueType
	val.value = value

	switch valueType {
	case ValueTypeNumber:
		val.float64Value, err = strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, err
		}
	case ValueTypeDate:
		val.timeValue, err = time.Parse(time.DateOnly, value)
		if err != nil {
			return nil, err
		}
	case ValueTypeDatetime:
		val.timeValue, err = time.Parse(time.RFC3339, value)
		if err != nil {
			return nil, err
		}
	case ValueTypeDuration:
		val.durationValue, err = time.ParseDuration(value)
		if err != nil {
			return nil, err
		}
	case ValueTypeText:
		break
	default:
		return nil, fmt.Errorf("unsupported value type %s", valueType)
	}

	return val, nil
}
