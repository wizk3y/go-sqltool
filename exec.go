package sqltool

import "database/sql"

// Exec -- do insert/update/delete or execute procedure
func (st *SQLTool) Exec(query string, args ...interface{}) (sql.Result, error) {
	var (
		stmt *sql.Stmt
		err  error
	)

	if st.isTransaction {
		stmt, err = st.tx.PrepareContext(st.ctx, query)
	} else {
		stmt, err = st.db.PrepareContext(st.ctx, query)
	}
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.ExecContext(st.ctx, args...)
}
