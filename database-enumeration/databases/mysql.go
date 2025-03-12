package databases

import (
	db_instance "database-enumeration/db-instance"
	"database/sql"
	"fmt"
	"log"
)

type MySQL struct {
	Host string
	Db   sql.DB
}

func NewMySQL(host string) (*MySQL, error) {
	m := MySQL{Host: host}
	err := m.connect()
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// connect just connect to the host
func (m *MySQL) connect() error {

	db, err := sql.Open("mysql", fmt.Sprintf("root:password@tcp(%s:3306)/information_schema", m.Host))
	if err != nil {
		log.Panicln(err)
	}
	m.Db = *db
	return nil
}

// GetSchema extract the schema from the db_instance.Database
func (m *MySQL) GetSchema() (*db_instance.Schema, error) {
	var s = new(db_instance.Schema)

	sql := `SELECT TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME FROM columns
	WHERE TABLE_SCHEMA NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys')
	ORDER BY TABLE_SCHEMA, TABLE_NAME`
	schemarows, err := m.Db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer schemarows.Close()

	var prevschema, prevtable string
	var db db_instance.Database
	var table db_instance.Table
	for schemarows.Next() {
		var currschema, currtable, currcol string
		if err := schemarows.Scan(&currschema, &currtable, &currcol); err != nil {
			return nil, err
		}

		if currschema != prevschema {
			if prevschema != "" {
				db.Tables = append(db.Tables, table)
				s.Databases = append(s.Databases, db)
			}
			db = db_instance.Database{Name: currschema, Tables: []db_instance.Table{}}
			prevschema = currschema
			prevtable = ""
		}

		if currtable != prevtable {
			if prevtable != "" {
				db.Tables = append(db.Tables, table)
			}
			table = db_instance.Table{Name: currtable, Columns: []string{}}
			prevtable = currtable
		}
		table.Columns = append(table.Columns, currcol)
	}
	db.Tables = append(db.Tables, table)
	s.Databases = append(s.Databases, db)
	if err := schemarows.Err(); err != nil {
		return nil, err
	}

	return s, nil
}
