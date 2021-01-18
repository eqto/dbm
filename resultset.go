package db

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

//Resultset ...
type Resultset map[string]interface{}

//IntNil ...
func (r Resultset) IntNil(name string) *int {
	if val := r.getValue(name); val != nil {
		switch val := val.Interface().(type) {
		case **uint64:
			if *val == nil {
				return nil
			}
			intVal := int(**val)
			return &intVal
		case **int64:
			if *val == nil {
				return nil
			}
			intVal := int(**val)
			return &intVal
		case **float64:
			if *val == nil {
				return nil
			}
			intVal := int(**val)
			return &intVal
		case *[]uint8:
			if intVal, e := strconv.Atoi(string(*val)); e == nil {
				return &intVal
			}
			return nil
		case **string:
			if *val == nil {
				return nil
			}
			if intVal, e := strconv.Atoi(**val); e == nil {
				return &intVal
			}
			return nil
		default:

		}
	}
	return nil
}

//Int ...
func (r Resultset) Int(name string) int {
	if val := r.IntNil(name); val != nil {
		return *val
	}
	return 0
}

//IntOr ...
func (r Resultset) IntOr(name string, defValue int) int {
	if val := r.IntNil(name); val != nil {
		return *val
	}
	return defValue
}

//Time ...
func (r Resultset) Time(name string) time.Time {
	if val := r.TimeNil(name); val != nil {
		return *val
	}
	return time.Time{}
}

//TimeNil ...
func (r Resultset) TimeNil(name string) *time.Time {
	if val, ok := r[name]; ok && val != nil {
		switch val := val.(type) {
		case *sql.NullTime:
			if val.Valid {
				return &val.Time
			}
		default:
			println(fmt.Sprintf(
				`unable to parse time from field '%s' with type '%v'`,
				name, reflect.TypeOf(val)))
		}
	}
	return nil
}

//Float ...
func (r Resultset) Float(name string) float64 {
	return r.FloatOr(name, 0)
}

//FloatOr ...
func (r Resultset) FloatOr(name string, defValue float64) float64 {
	if val := r.FloatNil(name); val != nil {
		return *val
	}
	return defValue
}

//FloatNil ...
func (r Resultset) FloatNil(name string) *float64 {
	if val, ok := r[name]; ok && val != nil {
		float := 0.0
		switch val := val.(type) {
		case *uint8:
			float = float64(*val)
		case *int8:
			float = float64(*val)
		case *uint16:
			float = float64(*val)
		case *int16:
			float = float64(*val)
		case *uint32:
			float = float64(*val)
		case *int32:
			float = float64(*val)
		case *uint64:
			float = float64(*val)
		case *int64:
			float = float64(*val)
		case *float32:
			float = float64(*val)
		case *float64:
			float = *val
		case *sql.NullInt64:
			if !val.Valid {
				return nil
			}
			float = float64(val.Int64)
		case *sql.NullFloat64:
			if !val.Valid {
				return nil
			}
			float = val.Float64
		default:
			println(fmt.Sprintf(
				`unable to parse float from field '%s' with type '%v'`,
				name, reflect.TypeOf(val)))
			return nil
		}
		return &float

	}
	return nil
}

//String ...
func (r Resultset) String(name string) string {
	return r.StringOr(name, ``)
}

//StringOr ...
func (r Resultset) StringOr(name string, defValue string) string {
	if val := r.StringNil(name); val != nil {
		return *val
	}
	return defValue
}

//StringNil ...
func (r Resultset) StringNil(name string) *string {
	if val, ok := r[name]; ok && val != nil {
		str := ``
		switch val := val.(type) {
		case *[]uint8:
			str = string(*val)
		case *uint8:
			str = strconv.FormatUint(uint64(*val), 10)
		case *int8:
			str = strconv.FormatInt(int64(*val), 10)
		case *uint16:
			str = strconv.FormatUint(uint64(*val), 10)
		case *int16:
			str = strconv.FormatInt(int64(*val), 10)
		case *uint32:
			str = strconv.FormatUint(uint64(*val), 10)
		case *int32:
			str = strconv.FormatInt(int64(*val), 10)
		case *uint64:
			str = strconv.FormatUint(uint64(*val), 10)
		case *int64:
			str = strconv.FormatInt(int64(*val), 10)
		case *float32:
			str = strconv.FormatFloat(float64(*val), 'f', -1, 32)
		case *float64:
			str = strconv.FormatFloat(float64(*val), 'f', -1, 64)
		case *sql.NullInt64:
			if !val.Valid {
				return nil
			}
			str = strconv.FormatInt(val.Int64, 10)
		case *sql.NullFloat64:
			if !val.Valid {
				return nil
			}
			str = strconv.FormatFloat(val.Float64, 'f', -1, 64)
		case *sql.NullTime:
			if !val.Valid {
				return nil
			}
			str = val.Time.String()
		default:
			println(fmt.Sprintf(
				`unable to parse string from field '%s' with type '%v'`,
				name, reflect.TypeOf(val)))
			return nil
		}
		return &str
	}
	return nil
}

//Bytes ...
func (r Resultset) Bytes(name string) []byte {
	if val, ok := r[name]; ok && val != nil {
		switch val := val.(type) {
		case *[]uint8:
			return *val
		default:
			println(fmt.Sprintf(
				`unable to parse bytes from field '%s' with type '%v'`,
				name, reflect.TypeOf(val)))
		}

	}
	return nil
}

//Interface ...
func (r Resultset) Interface(name string) interface{} {
	if val := r.getValue(name); val != nil {
		return val.Interface()
	}
	return nil
}

func (r Resultset) getValue(name string) *reflect.Value {
	if val, ok := r[name]; ok {
		val := reflect.ValueOf(val)
		return &val
	}
	return nil
}
