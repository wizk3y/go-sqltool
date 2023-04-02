package sqltool

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/wizk3y/go-sqltool/internal"
)

// SelectOne -- do select one row
func (st *SQLTool) SelectOne(dest interface{}, query string, args ...interface{}) (err error) {
	rows, err := st.queryContext(st.ctx, query, args...)
	if err != nil {
		return
	}
	defer rows.Close()
	var vp reflect.Value

	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer, not a value, to StructScan destination")
	}
	if v.IsNil() {
		return errors.New("nil pointer passed to StructScan destination")
	}
	direct := reflect.Indirect(v)

	base := internal.Deref(v.Type())

	if rows.Next() {
		vp = reflect.New(base)
		err = st.scanAndFill(rows, vp.Interface())
		if err == nil {
			direct.Set(vp.Elem())
		}
	} else {
		err = sql.ErrNoRows
	}

	return
}

// Select -- do select, same as SelectOne but return list results
func (st *SQLTool) Select(dest interface{}, query string, args ...interface{}) error {
	rows, err := st.queryContext(st.ctx, query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	var vp reflect.Value

	value := reflect.ValueOf(dest)

	// json.Unmarshal returns errors for these
	if value.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer, not a value, to StructScan destination")
	}
	if value.IsNil() {
		return errors.New("nil pointer passed to StructScan destination")
	}
	direct := reflect.Indirect(value)

	slice, err := internal.GetBaseType(value.Type(), reflect.Slice)
	if err != nil {
		return err
	}

	isPtr := slice.Elem().Kind() == reflect.Ptr
	base := internal.Deref(slice.Elem())
	empty := true

	for rows.Next() {
		vp = reflect.New(base)
		err = st.scanAndFill(rows, vp.Interface())
		if err != nil {
			fmt.Printf("[sqltool] error while scan and fill values, details: %v", err)
			continue
		}

		empty = false

		// append
		if isPtr {
			direct.Set(reflect.Append(direct, vp))
		} else {
			direct.Set(reflect.Append(direct, reflect.Indirect(vp)))
		}
	}

	if empty {
		err = sql.ErrNoRows
	}

	return err
}

func (st *SQLTool) queryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	var (
		stmt *sql.Stmt
		err  error
	)

	if st.isTransaction {
		stmt, err = st.tx.PrepareContext(ctx, query)
	} else {
		stmt, err = st.db.PrepareContext(ctx, query)
	}
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.QueryContext(ctx, args...)
}

// scanAndFill -- scan row then fill to dest
func (st *SQLTool) scanAndFill(rows *sql.Rows, dest interface{}) (err error) {
	values := make([]interface{}, 0)
	for _, column := range st.columns {
		vType, _ := st.column2Type[column]
		switch vType.Kind() {
		case reflect.String:
			var nstr = &sql.NullString{}
			values = append(values, nstr)
		case reflect.Bool:
			var nbool = &sql.NullBool{}
			values = append(values, nbool)
		case reflect.Float32, reflect.Float64:
			var nfloat64 = &sql.NullFloat64{}
			values = append(values, nfloat64)
		case reflect.Int64:
			if internal.IsStringSliceContains(st.dateTimeColumns, column) {
				var ntime = &sql.NullTime{}
				values = append(values, ntime)
			} else {
				var nint64 = &sql.NullInt64{}
				values = append(values, nint64)
			}
		case reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int:
			var nint64 = &sql.NullInt64{}
			values = append(values, nint64)
		case reflect.Slice, reflect.Struct, reflect.Map, reflect.Ptr:
			var nbytes = &sql.RawBytes{}
			values = append(values, nbytes)
		}
	}
	err = rows.Scan(values...)
	if err != nil {
		fmt.Printf("[sqltool] error while scan sql.Rows, fields: %v, details: %v", st.columns, err)
		return
	}

	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer, not a value, to StructScan destination")
	}
	if v.IsNil() {
		return errors.New("nil pointer passed to StructScan destination")
	}

	ve := v.Elem()
	for index, column := range st.columns {
		// get field name
		fieldName, _ := st.column2FieldName[column]
		vType, _ := st.column2Type[column]
		value := values[index]

		st.fillValueBySQLType(ve, column, fieldName, vType, value, false)
	}

	return
}

func (st *SQLTool) fillValueBySQLType(ve reflect.Value, column, fieldName string, vType reflect.Type, value interface{}, ptr bool) {
	switch vType.Kind() {
	case reflect.String:
		v := value.(*sql.NullString).String

		if ptr {
			ve.FieldByName(fieldName).Set(reflect.ValueOf(internal.StringPtr(v)))
		} else {
			ve.FieldByName(fieldName).SetString(v)
		}
	case reflect.Bool:
		v := value.(*sql.NullBool).Bool

		if ptr {
			ve.FieldByName(fieldName).Set(reflect.ValueOf(internal.BoolPtr(v)))
		} else {
			ve.FieldByName(fieldName).SetBool(v)
		}
	case reflect.Float32, reflect.Float64:
		v := value.(*sql.NullFloat64).Float64

		if ptr {
			ve.FieldByName(fieldName).Set(reflect.ValueOf(internal.Float64Ptr(v)))
		} else {
			ve.FieldByName(fieldName).SetFloat(v)
		}
	case reflect.Int64:
		if internal.IsStringSliceContains(st.dateTimeColumns, column) {
			var vtime = value.(*sql.NullTime)
			if vtime.Valid {
				ve.FieldByName(fieldName).SetInt(internal.TimestampByUnit(vtime.Time, st.dateTimeUnit))
			}
			break
		}
		fallthrough
	case reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int:
		v := value.(*sql.NullInt64).Int64

		if ptr {
			ve.FieldByName(fieldName).Set(reflect.ValueOf(internal.Int64Ptr(v)))
		} else {
			ve.FieldByName(fieldName).SetInt(v)
		}
	case reflect.Ptr:
		val := value.(*sql.RawBytes)
		if len(*val) < 1 {
			break
		}

		fillValueByType(ve, fieldName, internal.Deref(vType), string(*val), true)
	case reflect.Slice:
		val := value.(*sql.RawBytes)

		// If the slice is slice of bytes
		if ve.FieldByName(fieldName).Type().Elem().Kind() == reflect.Uint8 {
			var tmp = make([]byte, len(*val))
			copy(tmp, *val)
			ve.FieldByName(fieldName).SetBytes(tmp)
			break
		}

		fallthrough
	case reflect.Struct, reflect.Map:
		val := value.(*sql.RawBytes)
		if len(*val) < 1 {
			break
		}

		fillValueByType(ve, fieldName, vType, string(*val), false)
	}
}

func fillValueByType(ve reflect.Value, fieldName string, vType reflect.Type, valueStr string, ptr bool) {
	switch vType.Kind() {
	case reflect.String, reflect.Bool, reflect.Float64, reflect.Float32, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value := internal.CastValueTo(valueStr, vType, ptr)

		ve.FieldByName(fieldName).Set(reflect.ValueOf(value))
	case reflect.Slice, reflect.Struct, reflect.Map:
		var dataValue reflect.Value
		dataValue = reflect.New(vType)
		err := json.Unmarshal([]byte(valueStr), dataValue.Interface())
		if err != nil {
			fmt.Printf("[sqltool] error while parse value to slice/struct/map, field: %s, details: %v", fieldName, err)
			break
		}

		if ptr {
			ve.FieldByName(fieldName).Set(dataValue)
		} else {
			ve.FieldByName(fieldName).Set(dataValue.Elem())
		}
	}
}
