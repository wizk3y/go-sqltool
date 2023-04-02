package sqltool

import "fmt"

// Begin -- start transaction
func (st *SQLTool) Begin() error {
	tx, err := st.db.Begin()
	if err != nil {
		return err
	}

	st.tx = tx
	st.isTransaction = true
	return nil
}

// Commit -- commit transaction
func (st *SQLTool) Commit() error {
	if !st.isTransaction {
		fmt.Printf("[sqltool] tx not found")
		return nil
	}

	err := st.tx.Commit()
	if err != nil {
		return err
	}

	st.isTransaction = false
	st.tx = nil
	return nil
}

// Rollback -- rollback transaction if transaction not commited
func (st *SQLTool) Rollback() error {
	if !st.isTransaction {
		return nil
	}

	err := st.tx.Rollback()

	st.isTransaction = false
	st.tx = nil
	return err
}
