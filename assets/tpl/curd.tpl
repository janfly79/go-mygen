
const (
    {{.StructTableName}}_DB = "comic"
)






// 添加
func Add(ctx context.Context, value {{.StructTableName}})(lastId int64, err error) {
    conn := db.Get(ctx, "comic")
    sql := "INSERT INTO {{.TableName}} ({{.InsertFieldList}}) " +
        "VALUES ({{.InsertMark}})"
    q := db.SQLInsert("{{.TableName}}", sql)
    res, err := conn.ExecContext(ctx, q,
        {{range .InsertInfo}}value.{{.HumpName}},// {{.Comment}}
        {{end}})
    if err != nil {
        return
    }
    lastId, _ = res.LastInsertId()
    return
}

// 删除单条记录
func Del(ctx context.Context, where string, args []interface{}) (err error) {
    conn := db.Get(ctx, "comic")
	sql := "delete from {{.TableName}} " + where
	q := db.SQLDelete("{{.TableName}}", sql)

	_, err = conn.ExecContext(ctx, q, args...)
    return
}

// 获取单条记录
func Get(ctx context.Context, sqlText string, args []interface{})(row {{.StructTableName}}, err error){
    conn := db.Get(ctx, "comic")
    q := db.SQLSelect("{{.TableName}}", sqlText)
    err = conn.QueryRowContext(ctx, q, args...).Scan(
            		{{range .NullFieldsInfo}}&row.{{.HumpName}},// {{.Comment}}
            		{{end}})
    return
}

// 更新
func Update(ctx context.Context, sqlText string, args []interface{})(err error) {
    conn := db.Get(ctx, "comic")
    q := db.SQLUpdate("{{.TableName}}", sqlText)
    _, err = conn.ExecContext(ctx, q, args...)
    return
}

// 列表
func List(ctx context.Context, sqlText string, args []interface{})(rowsResult []{{.StructTableName}}, err error) {
    conn := db.Get(ctx, "comic")
    q := db.SQLSelect("{{.TableName}}", sqlText)
    rows, err := conn.QueryContext(ctx, q, args...)
    if err != nil {
    		return
    	}
    defer rows.Close()

    for rows.Next() {
        row := {{.StructTableName}}{}
        if err = rows.Scan(
            {{range .NullFieldsInfo}}&row.{{.HumpName}},// {{.Comment}}
                        		{{end}}
        ); err != nil {
            return
        }
        rowsResult = append(rowsResult, row)
    }

    err = rows.Err()

    return
}

// 总数
func Count(ctx context.Context, where string, args []interface{}) (total int32, err error){
    conn := db.Get(ctx, "comic")
    sqlText := "select count(*) from {{.TableName}} where " + where
    q := db.SQLSelect("{{.TableName}}", sqlText)
    err = conn.QueryRowContext(ctx, q, args...).Scan(&total)

    return
}

