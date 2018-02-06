/**
* Created by Visual Studio Code.
* User: tuxer
* Created At: 2017-12-17 22:15:50
 */

package db

import (
	"gitlab.com/tuxer/go-logger"
	"reflect"
	"strconv"
	"time"
)

//Resultset ...
type Resultset map[string]interface{}


//GetInt ...
func (r Resultset) GetInt(name string) *int {
	if val := r.getValue(name); val != nil {
		if i, ok := val.Interface().(**uint64); ok {
			if *i == nil {
				return nil
			}
			intVal := int(**i)
			return &intVal
		} else if str, ok := val.Interface().(**string); ok {
			if *str == nil {
				return nil
			}
			intVal, e := strconv.Atoi(**str)
			if e != nil {
				return nil
			}
			return &intVal
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
		case **float64:
			str = strconv.FormatUint(uint64(**val), 10)
		case *time.Time:
			str = val.String()
		default:
			log.D(`masuk2`)
			println(reflect.TypeOf(val).String())
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
