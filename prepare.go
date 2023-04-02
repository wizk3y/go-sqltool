package sqltool

import (
	"encoding/json"
	"reflect"
	"time"

	"github.com/wizk3y/go-sqltool/internal"
)

// PrepareInsert -- parse model struct and values support INSERT INTO command
func (st *SQLTool) PrepareInsert(i interface{}, opts ...sqlToolOpt) {
	if st.actionType != insertAction {
		st.actionType = insertAction
	}

	var (
		iPkgPath   = reflect.TypeOf(i).PkgPath()
		iName      = reflect.TypeOf(i).String()
		needUpdate bool
	)
	if st.modelPkgPath != iPkgPath || st.modelName != iName {
		st.modelPkgPath = iPkgPath
		st.modelName = iName
		needUpdate = true
	}

	// apply opts
	for _, o := range opts {
		needUpdate = o.Apply(st) || needUpdate
	}

	// parse columns
	if needUpdate {
		st.parseColumns(i)
	}
	st.values = st.PrepareValues(i)
}

// PrepareSelect -- parse model struct and values support SELECT command
func (st *SQLTool) PrepareSelect(i interface{}, opts ...sqlToolOpt) {
	if st.actionType != selectAction {
		st.actionType = selectAction
	}

	var (
		iPkgPath   = reflect.TypeOf(i).PkgPath()
		iName      = reflect.TypeOf(i).String()
		needUpdate bool
	)
	if st.modelPkgPath != iPkgPath || st.modelName != iName {
		st.modelPkgPath = iPkgPath
		st.modelName = iName
		needUpdate = true
	}

	// apply opts
	for _, o := range opts {
		needUpdate = o.Apply(st) || needUpdate
	}

	// parse columns
	if needUpdate {
		st.parseColumns(i)
	}
}

// PrepareUpdate -- parse model struct and values support UPDATE command
func (st *SQLTool) PrepareUpdate(i interface{}, opts ...sqlToolOpt) {
	if st.actionType != updateAction {
		st.actionType = updateAction
	}

	var (
		iPkgPath   = reflect.TypeOf(i).PkgPath()
		iName      = reflect.TypeOf(i).String()
		needUpdate bool
	)
	if st.modelPkgPath != iPkgPath || st.modelName != iName {
		st.modelPkgPath = iPkgPath
		st.modelName = iName
		needUpdate = true
	}

	// apply opts
	for _, o := range opts {
		needUpdate = o.Apply(st) || needUpdate
	}

	// parse columns
	if needUpdate {
		st.parseColumns(i)
	}
	st.values = st.PrepareValues(i)
}

func (st *SQLTool) parseColumns(i interface{}) {
	st.columns = make([]string, 0)
	st.column2FieldName = make(map[string]string, 0)
	st.column2Type = make(map[string]reflect.Type, 0)

	t := reflect.TypeOf(i).Elem()
	for index := 0; index < t.NumField(); index++ {
		f := t.Field(index)
		jsonTags := internal.TrimedSpaceStringSlice(f.Tag.Get("json"), ",")
		if len(jsonTags) == 0 || jsonTags[0] == "-" {
			continue
		}
		st.addColumnFieldNameAndType(jsonTags[0], f.Name, f.Type)
	}
}

// Use to add column, FieldName, Type to SQLTool
func (st *SQLTool) addColumnFieldNameAndType(column string, name string, typeV reflect.Type) {
	if st.isIgnoreColumn(column) {
		return
	}
	st.columns = append(st.columns, column)
	st.column2FieldName[column] = name
	st.column2Type[column] = typeV
}

func (st *SQLTool) isIgnoreColumn(column string) bool {
	if _, ok := st.autoCreateDateTimeColumns[column]; ok && st.actionType == insertAction {
		return false
	}

	if _, ok := st.autoUpdateDateTimeColumns[column]; ok && (st.actionType == insertAction || st.actionType == updateAction) {
		return false
	}

	if _, ok := st.ignoreColumns[column]; ok {
		return true
	}

	if column == st.serialColumn && (st.actionType == insertAction || st.actionType == updateAction) {
		return true
	}

	return false
}

// GetColumns -- Use to get list columns when do SELECT command
func (st *SQLTool) GetColumns() []string {
	return st.columns
}

// PrepareValues -- help parse struct values for batch INSERT command
func (st *SQLTool) PrepareValues(i interface{}) []interface{} {
	if internal.IsZeroOfUnderlyingType(i) {
		return nil
	}

	values := make([]interface{}, 0)
	v := reflect.ValueOf(i).Elem()
	for _, column := range st.columns {
		if column == st.serialColumn {
			continue
		}

		fieldName, _ := st.column2FieldName[column]
		fieldValue := v.FieldByName(fieldName)
		fieldValueInterface := fieldValue.Interface()
		convertedValue := fieldValueInterface
		if internal.IsZeroOfUnderlyingType(fieldValueInterface) && internal.IsStringSliceContains(st.nullableColumns, column) {
			convertedValue = nil
		}

		// check if column is datetime
		if internal.IsStringSliceContains(st.dateTimeColumns, column) {
			if _, ok := st.autoCreateDateTimeColumns[column]; ok && st.actionType == insertAction {
				convertedValue = time.Now()
			} else if _, ok := st.autoUpdateDateTimeColumns[column]; ok && (st.actionType == insertAction || st.actionType == updateAction) {
				convertedValue = time.Now()
			} else {
				if internal.IsZeroOfUnderlyingType(fieldValueInterface) {
					convertedValue = nil
				} else {
					convertedValue = internal.GetTimeByUnit(fieldValue.Int(), st.dateTimeUnit)
				}
			}
		}

		vType, _ := st.column2Type[column]

		if vType.Kind() == reflect.Slice && vType.Elem().Kind() == reflect.Uint8 {
			values = append(values, convertedValue)
			continue
		}

		var errMarshal error
		switch vType.Kind() {
		case reflect.Slice, reflect.Struct, reflect.Ptr, reflect.Map:
			if internal.IsZeroOfUnderlyingType(fieldValueInterface) {
				convertedValue = nil
			} else {
				convertedValue, errMarshal = json.Marshal(convertedValue)
			}
		}
		if errMarshal != nil {
			panic(errMarshal)
		}

		values = append(values, convertedValue)
	}

	return values
}

// GetInsertValues -- Get values has been prepared by PrepareInsert
func (st *SQLTool) GetInsertValues() []interface{} {
	return st.values
}

// GetUpdateMap -- Get map field - values has been prepared by PrepareUpdate
func (st *SQLTool) GetUpdateMap() map[string]interface{} {
	m := make(map[string]interface{})

	for k, f := range st.columns {
		if f == st.serialColumn {
			continue
		}

		m[f] = st.values[k]
	}

	return m
}
