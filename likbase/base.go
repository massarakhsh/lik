package likbase

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/massarakhsh/lik"
	"log"
	"strconv"
	"strings"
	"sync"
)

type DBase struct {
	DB   *sql.DB
	Sync sync.Mutex
}

type DBField struct {
	Name  string
	Proto string
}

type DBaser interface {
	Close()
	PrepareSql(what string, from string, where string, order string, limit ...int) string
	ControlTable(table string, fields []DBField) bool
	DropTable(table string)
	BuildOneMap(row *sql.Rows) *lik.DItemSet
	GetOneById(table string, id lik.IDB) lik.Seter
	GetOneBySql(stsql string) lik.Seter
	GetOneElm(what string, from string, where string, order string) lik.Seter
	GetListAll(table string) lik.Lister
	GetListBySql(stsql string) lik.Lister
	GetListElm(what string, from string, where string, order string, limit ...int) lik.Lister
	QueryRow(sql string) (*sql.Rows, bool)
	CalculeIDB(sql string) (lik.IDB, bool)
	CalculeInt(sql string) (int, bool)
	CalculeString(sql string) (string, bool)
	Execute(sql string, args ...interface{}) bool
	GetBinary(table string, id lik.IDB, field string) []byte
	SetBinary(table string, id lik.IDB, field string, val []byte)
	InsertElm(table string, sets lik.Seter) lik.IDB
	UpdateElm(table string, id lik.IDB, sets lik.Seter) bool
	DeleteElm(table string, id lik.IDB) bool
	LoadCountElm(table string, where string) int
}

var (
	FId = "id"
)

func SignId(table string, id lik.IDB) string {
	return fmt.Sprintf("%s%d", table, int(id))
}

func OpenDBase(driver string, logon string, connect string, dbname string) DBaser {
	dbs := &DBase{}
	if !dbs.openDBase(driver, logon, connect, dbname) {
		dbs = nil
	}
	return dbs
}

func (dbs *DBase) openDBase(driver string, logon string, connect string, dbname string) bool {
	strcon := connect
	if logon != "" {
		strcon = logon + "@" + connect
	}
	strcon += "/" + dbname
	db, err := sql.Open(driver, strcon)
	if err != nil {
		lik.SayError("Database " + dbname + " NOT opened")
		//log.Fatal(err)
		return false
	}
	dbs.DB = db
	return true
}

func (dbs *DBase) Close() {
	if dbs.DB != nil {
		_  = dbs.DB.Close()
		dbs.DB = nil
	}
}

func StrToIDB(str string) lik.IDB {
	return lik.IDB(lik.StrToInt(str))
}

func IDBToStr(id lik.IDB) string {
	return lik.IntToStr(int(id))
}

func PrepareSql(what string, from string, where string, order string, limit ...int) string {
	stsql := "select"
	if len(what) > 0 {
		stsql += " " + what
	} else {
		stsql += " *"
	}
	if len(from) > 0 {
		stsql += " from " + from
	}
	if len(where) > 0 {
		stsql += " where " + where
	}
	if len(order) > 0 {
		stsql += " order by " + order
	}
	count := 0
	first := 0
	if len(limit) > 1 {
		count, first = limit[0], limit[1]
	} else if len(limit) > 0 {
		count = limit[0]
	}
	if first > 0 {
		stsql += fmt.Sprintf(" limit %d,%d", count, first)
	} else if count > 0 {
		stsql += fmt.Sprintf(" limit %d", count)
	}
	return stsql
}

func (dbs *DBase) PrepareSql(what string, from string, where string, order string, limit ...int) string {
	if dbs.DB == nil {
		return ""
	}
	return PrepareSql(what, from, where, order, limit...)
}

func (dbs *DBase) DropTable(table string) {
	dbs.Execute(fmt.Sprintf("DROP TABLE IF EXISTS `%s`", table))
}

func (dbs *DBase) ControlTable(table string, fields []DBField) bool {
	for fase := 0; fase < 2; fase++ {
		if rows, err := dbs.DB.Query(fmt.Sprintf("SELECT * FROM `%s` LIMIT 1", table)); err == nil {
			columns, _ := rows.Columns()
			coltypes, _ := rows.ColumnTypes()
			for _, field := range fields {
				found := false
				right := false
				for nc := 0; nc < len(columns); nc++ {
					if field.Name == columns[nc] {
						found = true
						right = false
						ctp := coltypes[nc].DatabaseTypeName()
						if strings.Contains(field.Proto, "S") && strings.Contains(ctp, "CHAR") {
							right = true
						} else if strings.Contains(field.Proto, "T") && strings.Contains(ctp, "BLOB") {
							right = true
						} else if strings.Contains(field.Proto, "L") && strings.Contains(ctp, "INT") {
							right = true
						} else if strings.Contains(field.Proto, "I") && strings.Contains(ctp, "INT") {
							right = true
						} else if strings.Contains(field.Proto, "R") && strings.Contains(ctp, "REAL") {
							right = true
						} else if strings.Contains(field.Proto, "R") && strings.Contains(ctp, "DOUBLE") {
							right = true
						} else if strings.Contains(field.Proto, "D") && strings.Contains(ctp, "DATA") {
							right = true
						} else if strings.Contains(field.Proto, "D") && strings.Contains(ctp, "DATE") {
							right = true
						}
					}
				}
				if !found || !right {
					lik.SayWarning("ALTER " + table + ": " + field.Name + ", " + field.Proto)
					cmd := fmt.Sprintf("ALTER TABLE `%s`", table)
					if !found {
						cmd += fmt.Sprintf(" ADD COLUMN `%s`", field.Name)
					} else {
						cmd += fmt.Sprintf(" CHANGE COLUMN `%s` `%s`", field.Name, field.Name)
					}
					if strings.Index(field.Proto, "S") >= 0 {
						cmd += " VARCHAR(255) NOT NULL DEFAULT ''"
					} else if strings.Index(field.Proto, "I") >= 0 {
						cmd += " INTEGER NOT NULL DEFAULT 0"
					} else if strings.Index(field.Proto, "L") >= 0 {
						cmd += " BIGINT NOT NULL DEFAULT 0"
					} else if strings.Index(field.Proto, "R") >= 0 {
						cmd += " REAL NOT NULL DEFAULT 0.0"
					} else if strings.Index(field.Proto, "T") >= 0 {
						cmd += " LONGBLOB"
					} else if strings.Index(field.Proto, "D") >= 0 {
						cmd += " DATATIME"
					}
					lik.SayWarning(cmd)
					if _, err := dbs.DB.Exec(cmd); err != nil {
						log.Fatal(err)
					}
				}
			}
			break
		}
		sql := fmt.Sprintf("CREATE TABLE `%s` (`%s` BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY)", table, FId)
		if _, err := dbs.DB.Exec(sql); err != nil {
			log.Fatal(err)
		}
	}
	return true
}

func (dbs *DBase) QueryRow(sql string) (*sql.Rows, bool) {
	rows, err := dbs.DB.Query(sql)
	if err != nil {
		fmt.Println(err)
		fmt.Println("SQL: " + sql)
		return nil, false
	}
	return rows, true
}

func (dbs *DBase) CalculeIDB(sql string) (lik.IDB, bool) {
	id, ok := dbs.CalculeInt(sql)
	return lik.IDB(id), ok
}

func (dbs *DBase) CalculeInt(sql string) (int, bool) {
	val := 0
	ok := false
	if rows, isit := dbs.QueryRow(sql); isit {
		if rows.Next() {
			if rows.Scan(&val) == nil {
				ok = true
			}
		}
		if rows.Close() != nil {
		}
	}
	return val, ok
}

func (dbs *DBase) CalculeString(sql string) (string, bool) {
	val := ""
	ok := false
	if rows, isit := dbs.QueryRow(sql); isit {
		if rows.Next() {
			if rows.Scan(&val) == nil {
				ok = true
			}
		}
		if rows.Close() != nil {
		}
	}
	return val, ok
}

func (dbs *DBase) Execute(sql string, args ...interface{}) bool {
	_, err := dbs.DB.Exec(sql, args...)
	if err != nil {
		fmt.Println("MySql error " + sql)
		return false
	}
	return true
}

func (dbs *DBase) BuildOneMap(row *sql.Rows) *lik.DItemSet {
	elm := &lik.DItemSet{}
	types, _ := row.ColumnTypes()
	columns, _ := row.Columns()
	count := len(columns)
	values := make([]interface{}, count)
	scanArgs := make([]interface{}, count)
	for i := range values {
		scanArgs[i] = &values[i]
	}
	_ = row.Scan(scanArgs...)
	for i, v := range values {
		typ := types[i].DatabaseTypeName()
		key := columns[i]
		str := ""
		if v != nil {
			x := v.([]byte)
			str = string(x)
		}
		if typ == "INT" || typ == "INTEGER" {
			if val, ok := lik.StrToIntIf(str); ok {
				elm.SetValue(key, val)
			}
		} else if typ == "BIGINT" {
			if val, ok := lik.StrToIntIf(str); ok {
				elm.SetValue(key, val)
			}
		} else if typ == "REAL" {
			if rl, ok := strconv.ParseFloat(str, 64); ok == nil {
				_ = rl
			}
		} else /*if typ == "VARCHAR"*/ {
			elm.SetValue(key, str)
		}
	}

	return elm
}

func (dbs *DBase) GetOneBySql(stsql string) lik.Seter {
	if strings.Index(strings.ToLower(stsql), " limit") < 0 {
		stsql += " limit 1"
	}
	rows, err := dbs.DB.Query(stsql)
	if err != nil {
		lik.SayError("MySql error " + stsql)
		return nil
	}
	if !rows.Next() {
		rows.Close()
		return nil
	}
	elm := dbs.BuildOneMap(rows)
	rows.Close()
	return elm.ToSet()
}

func (dbs *DBase) GetOneById(table string, id lik.IDB) lik.Seter {
	whe := fmt.Sprintf("%s=%d", FId, int(id))
	stsql := dbs.PrepareSql("*", fmt.Sprintf("`%s`", table), whe, "", 1)
	return dbs.GetOneBySql(stsql)
}

func (dbs *DBase) GetOneElm(what string, from string, where string, order string) lik.Seter {
	stsql := dbs.PrepareSql(what, from, where, order, 1)
	return dbs.GetOneBySql(stsql)
}

func (dbs *DBase) GetListBySql(stsql string) lik.Lister {
	rows, err := dbs.DB.Query(stsql)
	if err != nil {
		lik.SayError("MySql error " + stsql)
		return nil
	}
	list := lik.BuildList()
	for rows.Next() {
		elm := dbs.BuildOneMap(rows)
		list.AddItems(elm)
	}
	rows.Close()
	return list
}

func (dbs *DBase) GetListElm(what string, from string, where string, order string, limit ...int) lik.Lister {
	stsql := dbs.PrepareSql(what, from, where, order, limit...)
	return dbs.GetListBySql(stsql)
}

func (dbs *DBase) GetListAll(table string) lik.Lister {
	return dbs.GetListElm("*", fmt.Sprintf("`%s`", table), "", FId)
}

func (dbs *DBase) GetBinary(table string, id lik.IDB, field string) []byte {
	var result []byte
	whe := fmt.Sprintf("%s=%d", FId, int(id))
	stsql := dbs.PrepareSql(field, fmt.Sprintf("`%s`", table), whe, "")
	if rows, ok := dbs.QueryRow(stsql); ok {
		if rows.Next() {
			values := make([]interface{}, 1)
			if err := rows.Scan(&values[0]); err == nil {
				if values[0] != nil {
					result = values[0].([]byte)
				}
			}
		}
		rows.Close()
	}
	return result
}

func (dbs *DBase) SetBinary(table string, id lik.IDB, field string, val []byte) {
	whe := fmt.Sprintf("%s=%d", FId, int(id))
	sql := fmt.Sprintf("update `%s` set `%s` =? where %s", table, field, whe)
	dbs.Execute(sql, val)
}

func (dbs *DBase) InsertElm(table string, sets lik.Seter) lik.IDB {
	return dbs.updateTableElm(table, 0, sets)
}

func (dbs *DBase) UpdateElm(table string, id lik.IDB, sets lik.Seter) bool {
	return dbs.updateTableElm(table, id, sets) == id
}

func (dbs *DBase) DeleteElm(table string, id lik.IDB) bool {
	ok := true
	whe := fmt.Sprintf("%s=%d", FId, int(id))
	if !dbs.Execute(fmt.Sprintf("DELETE FROM `%s` WHERE %s", table, whe)) {
		ok = false
	}
	return ok
}

func (dbs *DBase) LoadCountElm(table string, where string) int {
	stsql := dbs.PrepareSql("COUNT(*)", fmt.Sprintf("`%s`", table), where, "")
	count, _ := dbs.CalculeInt(stsql)
	return count
}

func (dbs *DBase) updateTableElm(table string, id lik.IDB, sets lik.Seter) lik.IDB {
	var id0 lik.IDB
	cmd := ""
	if sets != nil {
		for _, set := range sets.Values() {
			key := set.Key
			val := set.Val
			if key == "id" || key == "Id" {
				vali := val.ToInt()
				id0 = lik.IDB(vali)
				if id0 > 0 {
					if cmd != "" {
						cmd += ","
					}
					cmd += fmt.Sprintf("`%s`=%d", key, vali)
				}
			} else if val.IsInt() {
				vali := val.ToInt()
				if cmd != "" {
					cmd += ","
				}
				cmd += fmt.Sprintf("`%s`=%d", key, vali)
			} else if val.IsFloat() {
				valf := val.ToFloat()
				if cmd != "" {
					cmd += ","
				}
				cmd += fmt.Sprintf("`%s`=%f", key, valf)
			} else if val.IsString() {
				vals := val.ToString()
				if cmd != "" {
					cmd += ","
				}
				if vals == "CURRENT_TIMESTAMP" {
					cmd += fmt.Sprintf("`%s`=%s", key, vals)
				} else {
					cmd += fmt.Sprintf("`%s`=%s", key, lik.StrToQuotes(vals))
				}
			}
		}
	}
	if id <= 0 {
		var sqls string
		if cmd != "" {
			sqls = fmt.Sprintf("insert into `%s` set %s", table, cmd)
		} else {
			sqls = fmt.Sprintf("insert into `%s` () values ()", table)
		}
		if dbs.Execute(sqls) {
			if id0 > 0 {
				id = id0
			} else {
				id, _ = dbs.CalculeIDB(fmt.Sprintf("SELECT MAX(%s) FROM `%s`", FId, table))
			}
		}
	} else if cmd != "" {
		whe := fmt.Sprintf("%s=%d", FId, int(id))
		sqls := fmt.Sprintf("update `%s` set %s where %s", table, cmd, whe)
		if !dbs.Execute(sqls) {
			id = 0
		}
	}
	return id
}
