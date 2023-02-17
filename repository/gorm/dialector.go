package gorm

import (
	"database/sql/driver"
	"reflect"

	"github.com/google/uuid"
	"gorm.io/driver/mysql"
)

type MySQLDialector struct {
	mysql.Dialector
}

func (x MySQLDialector) Explain(sql string, vars ...interface{}) string {
	var convertParams func(interface{}, int)
	var newVars = make([]interface{}, len(vars))
	convertParams = func(v interface{}, idx int) {
		switch v := v.(type) {
		case driver.Valuer:
			reflectValue := reflect.ValueOf(v)
			if v != nil && reflectValue.IsValid() && ((reflectValue.Kind() == reflect.Ptr && !reflectValue.IsNil()) || reflectValue.Kind() != reflect.Ptr) {
				r, _ := v.Value()
				convertParams(r, idx)
			} else {
				newVars[idx] = v
			}
		case []byte:
			id, err := uuid.FromBytes(v)
			if err != nil {
				newVars[idx] = v
			} else {
				newVars[idx] = "UUID:" + id.String()
			}
		case uuid.UUID:
			newVars[idx] = "UUID:" + v.String()
		default:
			newVars[idx] = v
		}
	}
	for idx, v := range vars {
		convertParams(v, idx)
	}
	return x.Dialector.Explain(sql, newVars...)
}
