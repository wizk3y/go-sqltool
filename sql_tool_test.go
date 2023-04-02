package sqltool_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Masterminds/squirrel"
	"github.com/wizk3y/go-sqltool"
)

func Test_SQLTool_Insert(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error when open mock database connection, details: %v", err)
	}
	defer db.Close()

	query := `INSERT INTO user (created_at,updated_at,username,pass) VALUES (?,?,?,?)`

	mock.ExpectPrepare(query).
		ExpectExec().
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "sample", "sample").
		WillReturnResult(sqlmock.NewResult(1, 0))

	// real code
	sqlTool := sqltool.NewTool(context.Background(), db)
	req := struct {
		ID        int64  `json:"id"`
		CreatedAt int64  `json:"created_at"`
		UpdatedAt int64  `json:"updated_at"`
		Username  string `json:"username"`
		Pass      string `json:"pass"`
	}{
		Username: "sample",
		Pass:     "sample",
	}
	sqlTool.PrepareInsert(&req,
		sqltool.DateTimeColumnsOpt([]string{"created_at", "updated_at"}),
		sqltool.AutoCreateDateTimeColumnsOpt([]string{"created_at"}),
		sqltool.AutoUpdateDateTimeColumnsOpt([]string{"updated_at"}),
	)

	query, args, err := squirrel.Insert("user").
		Columns(sqlTool.GetColumns()...).
		Values(sqlTool.GetInsertValues()...).
		ToSql()
	if err != nil {
		t.Fatalf("error when build query")
	}

	_, err = sqlTool.Exec(query, args...)
	if err != nil {
		t.Fatalf("error when execute insert query, details: %v", err)
	}
}

func Test_SQLTool_Select(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error when open mock database connection, details: %v", err)
	}
	defer db.Close()

	query := "SELECT id, created_at, updated_at, username, pass, obj, ptr, map, slice FROM user WHERE id = ?"

	mock.ExpectPrepare(query).
		ExpectQuery().
		WithArgs(1).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "created_at", "updated_at", "username", "pass", "obj", "ptr", "map", "slice"}).
				AddRow(1, "2023-03-30 23:57:48.348585", "2023-03-30 23:57:48.348585", "sample", "sample", `{"sub_field":"abc"}`, "true", `{"a":"b"}`, `[1,3]`),
		)

	// real code
	sqlTool := sqltool.NewTool(context.Background(), db)
	res := struct {
		ID        int64  `json:"id"`
		CreatedAt int64  `json:"created_at"`
		UpdatedAt int64  `json:"updated_at"`
		Username  string `json:"username"`
		Pass      string `json:"pass"`
		Obj       struct {
			SubField string `json:"sub_field"`
		} `json:"obj"`
		Ptr   *bool              `json:"ptr"`
		Map   map[string]*string `json:"map"`
		Slice []*int64           `json:"slice"`
	}{}
	sqlTool.PrepareSelect(&res,
		sqltool.DateTimeColumnsOpt([]string{"created_at", "updated_at"}),
	)
	query, args, err := squirrel.Select(sqlTool.GetColumns()...).
		From("user").
		Where(squirrel.Eq{"id": 1}).
		ToSql()
	if err != nil {
		t.Fatalf("error when build query")
	}

	err = sqlTool.SelectOne(&res, query, args...)
	if err != nil {
		t.Fatalf("error when execute select one, details: %v", err)
	}
}

func Test_SQLTool_Update(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error when open mock database connection, details: %v", err)
	}
	defer db.Close()

	query := `UPDATE user SET pass = ?, updated_at = ? WHERE id = ?`

	mock.ExpectPrepare(query).
		ExpectExec().
		WithArgs("sample", sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 0))

	// real code
	sqlTool := sqltool.NewTool(context.Background(), db)
	req := struct {
		ID        int64  `json:"id"`
		CreatedAt int64  `json:"created_at"`
		UpdatedAt int64  `json:"updated_at"`
		Username  string `json:"username"`
		Pass      string `json:"pass"`
	}{
		Pass: "sample",
	}
	sqlTool.PrepareUpdate(&req,
		sqltool.DateTimeColumnsOpt([]string{"created_at", "updated_at"}),
		sqltool.AutoCreateDateTimeColumnsOpt([]string{"created_at"}),
		sqltool.AutoUpdateDateTimeColumnsOpt([]string{"updated_at"}),
		sqltool.IgnoreColumnsOpt([]string{"created_at", "updated_at", "username"}),
	)

	query, args, err := squirrel.Update("user").
		SetMap(sqlTool.GetUpdateMap()).
		Where(squirrel.Eq{"id": 1}).
		ToSql()
	if err != nil {
		t.Fatalf("error when build query")
	}

	_, err = sqlTool.Exec(query, args...)
	if err != nil {
		t.Fatalf("error when execute update query, details: %v", err)
	}
}

func Test_SQLTool_Delete(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error when open mock database connection, details: %v", err)
	}
	defer db.Close()

	mock.ExpectPrepare("DELETE FROM user WHERE id = ?").
		ExpectExec().
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 0))

	// real code
	sqlTool := sqltool.NewTool(context.Background(), db)

	query, args, err := squirrel.Delete("user").
		Where(squirrel.Eq{"id": 1}).
		ToSql()
	if err != nil {
		t.Fatalf("error when build query")
	}

	_, err = sqlTool.Exec(query, args...)
	if err != nil {
		t.Fatalf("error when execute delete query, details: %v", err)
	}
}
