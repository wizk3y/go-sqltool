package sqltool

import "reflect"

type sqlToolOpt interface {
	Apply(conf *SQLTool) bool
}

type serialColumnOpt string

// SerialColumnOpt -- ignore serial/auto-increment column when call GetFillValues or GetUpdateMap
func SerialColumnOpt(column string) sqlToolOpt {
	return serialColumnOpt(column)
}

func (o serialColumnOpt) Apply(st *SQLTool) bool {
	st.serialColumn = string(o)

	return false
}

type nullableColumnsOpt []string

// NullableColumnsOpt -- column not set value to nil instead of zero value of type
func NullableColumnsOpt(columns []string) sqlToolOpt {
	return nullableColumnsOpt(columns)
}

func (o nullableColumnsOpt) Apply(st *SQLTool) bool {
	st.nullableColumns = o

	return false
}

type dateTimeColumnsOpt []string

// DateTimeColumnsOpt -- set date/time columns help parse to value type
func DateTimeColumnsOpt(columns []string) sqlToolOpt {
	return dateTimeColumnsOpt(columns)
}

func (o dateTimeColumnsOpt) Apply(st *SQLTool) bool {
	st.dateTimeColumns = o

	return false
}

type dateTimeUnitOpt string

// DateTimeUnitOpt -- parse/prepare date time column value with time unit
func DateTimeUnitOpt(column string) sqlToolOpt {
	return dateTimeUnitOpt(column)
}

func (o dateTimeUnitOpt) Apply(st *SQLTool) bool {
	st.dateTimeUnit = string(o)

	return false
}

type autoCreateDateTimeColumnsOpt []string

// AutoCreateDateTimeColumnsOpt -- set time.Now() to columns as default value for Insert query, ignore any input value
func AutoCreateDateTimeColumnsOpt(columns []string) sqlToolOpt {
	return autoCreateDateTimeColumnsOpt(columns)
}

func (o autoCreateDateTimeColumnsOpt) Apply(st *SQLTool) bool {
	mapColumns := make(map[string]bool)

	for _, c := range o {
		mapColumns[c] = true
	}

	if reflect.DeepEqual(mapColumns, st.autoCreateDateTimeColumns) {
		return false
	}

	st.autoCreateDateTimeColumns = mapColumns
	return true
}

type autoUpdateDateTimeColumnsOpt []string

// AutoUpdateDateTimeColumnsOpt -- set time.Now() to columns as default value for Insert/Update query, ignore any input value
func AutoUpdateDateTimeColumnsOpt(columns []string) sqlToolOpt {
	return autoUpdateDateTimeColumnsOpt(columns)
}

func (o autoUpdateDateTimeColumnsOpt) Apply(st *SQLTool) bool {
	mapColumns := make(map[string]bool)

	for _, c := range o {
		mapColumns[c] = true
	}

	if reflect.DeepEqual(mapColumns, st.autoUpdateDateTimeColumns) {
		return false
	}

	st.autoUpdateDateTimeColumns = mapColumns
	return true
}

type ignoreColumnsOpt []string

// IgnoreColumnsOpt -- ignore columns when do select query
func IgnoreColumnsOpt(columns []string) sqlToolOpt {
	return ignoreColumnsOpt(columns)
}

func (o ignoreColumnsOpt) Apply(st *SQLTool) bool {
	mapColumns := make(map[string]bool)

	for _, c := range o {
		mapColumns[c] = true
	}

	if reflect.DeepEqual(mapColumns, st.ignoreColumns) {
		return false
	}

	st.ignoreColumns = mapColumns
	return true
}
