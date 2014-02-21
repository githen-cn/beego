package orm

import (
	"fmt"
	"reflect"
	"time"
)

func getExistPk(mi *modelInfo, ind reflect.Value) (column string, value interface{}, exist bool) {
	fi := mi.fields.pk

	v := ind.Field(fi.fieldIndex)
	if fi.fieldType&IsIntegerField > 0 {
		vu := v.Int()
		exist = vu > 0
		value = vu
	} else {
		vu := v.String()
		exist = vu != ""
		value = vu
	}

	column = fi.column
	return
}

func getFlatParams(fi *fieldInfo, args []interface{}, tz *time.Location) (params []interface{}) {

outFor:
	for _, arg := range args {
		val := reflect.ValueOf(arg)

		if arg == nil {
			params = append(params, arg)
			continue
		}

		switch v := arg.(type) {
		case []byte:
		case time.Time:
			if fi != nil && fi.fieldType == TypeDateField {
				arg = v.In(DefaultTimeLoc).Format(format_Date)
			} else {
				arg = v.In(tz).Format(format_DateTime)
			}
		default:
			kind := val.Kind()
			switch kind {
			case reflect.Slice, reflect.Array:

				var args []interface{}
				for i := 0; i < val.Len(); i++ {
					v := val.Index(i)

					var vu interface{}
					if v.CanInterface() {
						vu = v.Interface()
					}

					if vu == nil {
						continue
					}

					args = append(args, vu)
				}

				if len(args) > 0 {
					p := getFlatParams(fi, args, tz)
					params = append(params, p...)
				}
				continue outFor

			case reflect.Ptr, reflect.Struct:
				ind := reflect.Indirect(val)

				if ind.Kind() == reflect.Struct {
					typ := ind.Type()
					name := getFullName(typ)
					var value interface{}
					if mmi, ok := modelCache.getByFN(name); ok {
						if _, vu, exist := getExistPk(mmi, ind); exist {
							value = vu
						}
					}
					arg = value

					if arg == nil {
						panic(fmt.Errorf("need a valid args value, unknown table or value `%s`", name))
					}
				} else {
					arg = ind.Interface()
				}
			}
		}
		params = append(params, arg)
	}
	return
}
