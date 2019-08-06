/**
* Created by Visual Studio Code.
* User: tuxer
* Created At: 2017-12-17 22:15:50
 */

package db

import (
	"reflect"
	"strconv"
	"time"
)

//Resultset ...
type Resultset map[string]interface{}

//GetInt ...
func (r Resultset) GetInt(name string) *int {
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

//GetIntD ...
func (r Resultset) GetIntD(name string) int {
	if val := r.GetInt(name); val != nil {
		return *val
	}
	return 0
}

//GetIntOr ...
func (r Resultset) GetIntOr(name string, defValue int) int {
	if val := r.GetInt(name); val != nil {
		return *val
	}
	return defValue
}

//GetTime ...
func (r Resultset) GetTime(name string) *time.Time {
	if val := r.getValue(name); val != nil {
		if val, ok := val.Interface().(**time.Time); ok {
			if *val == nil {
				return nil
			}
			return *val
		}
	}
	return nil
}

//GetTimeD ...
func (r Resultset) GetTimeD(name string) time.Time {
	if val := r.GetTime(name); val != nil {
		return *val
	}
	return time.Time{}
}

func (r Resultset) getValue(name string) *reflect.Value {
	if val, ok := r[name]; ok {
		val := reflect.ValueOf(val)
		return &val
	}
	return nil
}

//GetFloat ...
func (r Resultset) GetFloat(name string) *float64 {
	if val := r.getValue(name); val != nil {
		switch val := val.Interface().(type) {
		case **uint64:
			if *val == nil {
				return nil
			}
			floatVal := float64(int(**val))
			return &floatVal
		case **int64:
			if *val == nil {
				return nil
			}
			floatVal := float64(int(**val))
			return &floatVal
		case **float64:
			return *val
		}
	}
	return nil
}

//GetFloatD ...
func (r Resultset) GetFloatD(name string) float64 {
	if val := r.GetFloat(name); val != nil {
		return *val
	}
	return 0
}

//GetFloatOr ...
func (r Resultset) GetFloatOr(name string, defValue float64) float64 {
	if val := r.GetFloat(name); val != nil {
		return *val
	}
	return defValue
}

//GetString ...
func (r Resultset) GetString(name string) *string {
	if val := r.getValue(name); val != nil {
		if reflect.ValueOf(val.Interface()).Elem().IsNil() {
			return nil
		}
		str := ``
		switch val := val.Interface().(type) {
		case **string:
			str = **val
		case *[]byte:
			str = string(*val)
		case **uint64:
			str = strconv.FormatUint(**val, 10)
		case **int64:
			str = strconv.FormatInt(**val, 10)
		case **float64:
			str = strconv.FormatUint(uint64(**val), 10)
		case *time.Time:
			str = val.String()
		case **time.Time:
			v := *val
			str = v.String()
		default:
			println(`unable to parse string from ` + reflect.TypeOf(val).String())
		}
		return &str
	}
	return nil
}

//GetBytes ...
func (r Resultset) GetBytes(name string) []byte {
	if val := r.getValue(name); val != nil {
		if reflect.ValueOf(val.Interface()).Elem().IsNil() {
			return nil
		}
		switch val := val.Interface().(type) {
		case *[]byte:
			return *val
		case **uint64:
			return []byte(strconv.FormatUint(**val, 10))
		case **int64:
			return []byte(strconv.FormatInt(**val, 10))
		case **float64:
			return []byte(strconv.FormatUint(uint64(**val), 10))
		default:
			return []byte(``)
		}
	}
	return nil
}

//GetStringD ...
func (r Resultset) GetStringD(name string) string {
	if val := r.GetString(name); val != nil {
		return *val
	}
	return ``
}

//GetStringOr ...
func (r Resultset) GetStringOr(name string, defValue string) string {
	if val := r.GetString(name); val != nil {
		return *val
	}
	return defValue
}
