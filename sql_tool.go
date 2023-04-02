package sqltool

import (
	"context"
	"database/sql"
	"reflect"
)

type actionType string

const (
	insertAction actionType = "insert"
	selectAction actionType = "select"
	updateAction actionType = "update"
	deleteAction actionType = "delete"
)

// SQLTool --
type SQLTool struct {
	ctx           context.Context
	db            *sql.DB
	isTransaction bool
	tx            *sql.Tx

	actionType actionType
	// related to struct
	modelPkgPath     string
	modelName        string
	columns          []string
	column2FieldName map[string]string
	column2Type      map[string]reflect.Type
	values           []interface{}
	// related to opt
	serialColumn              string
	nullableColumns           []string
	dateTimeColumns           []string
	dateTimeUnit              string
	autoCreateDateTimeColumns map[string]bool
	autoUpdateDateTimeColumns map[string]bool
	ignoreColumns             map[string]bool
}

// NewTool -- generic sql tool
func NewTool(ctx context.Context, db *sql.DB) (st SQLTool) {
	st.ctx = ctx
	st.db = db
	// opt default
	st.serialColumn = "id"
	st.dateTimeUnit = "ns"
	return
}
