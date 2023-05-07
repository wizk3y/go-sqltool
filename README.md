# go-sqltool

go-sqltool is a Go toolkit help you working with SQL database (without using ORM). This project recommend using [squirrel](https://github.com/Masterminds/squirrel) to build query, which help prevent some issue related to SQL Injection.

## Install
```shell
go get github.com/wizk3y/go-sqltool
```

**Note:** go-sqltool uses [Go Modules](https://github.com/golang/go/wiki/Modules) to manage dependencies.

## Usage
- Create new instance of `SQLTool`
```go
st := sqltool.NewTool(ctx, m.db)
```
- If your query is insert/update or select, use `PrepareInsert`/`PrepareUpdate` or `PrepareSelect` equivalent
```go
var data = struct {
    ID       int64  `json:"id"`
    Username string `json:"username"`
}{}
st.PrepareSelect(&data)
```
- Using `Select`, `SelectOne`, `Exec` to execute your query
```go
query := "SELECT id, username FROM user WHERE id = ? LIMIT 1"
args := []interface{1}

err = st.SelectOne(&data, query, args...)
if err != nil {
    return nil, err
}
```

## Advance usage
- [Transaction](https://github.com/wizk3y/go-sqltool-doc/tree/master/transaction.md)
- [Batch insert](https://github.com/wizk3y/go-sqltool-doc/tree/master/batch_insert.md)
- [Prepare opt](https://github.com/wizk3y/go-sqltool-doc/tree/master/prepare_opt.md)

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)