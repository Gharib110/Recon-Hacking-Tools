package databases

import (
	db_instance "database-enumeration/db-instance"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Mongo struct {
	Host    string
	session *mgo.Session
}

func NewMongo(host string) (*Mongo, error) {
	m := Mongo{Host: host}
	err := m.connect()
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (m *Mongo) connect() error {
	s, err := mgo.Dial(m.Host)
	if err != nil {
		return err
	}
	m.session = s
	return nil
}

// GetSchema extract the schema from the db_instance.Database
func (m *Mongo) GetSchema() (*db_instance.Schema, error) {
	var s = new(db_instance.Schema)

	dbNames, err := m.session.DatabaseNames()
	if err != nil {
		return nil, err
	}

	for _, dbname := range dbNames {
		db := db_instance.Database{Name: dbname, Tables: []db_instance.Table{}}
		collections, err := m.session.DB(dbname).CollectionNames()
		if err != nil {
			return nil, err
		}

		for _, collection := range collections {
			table := db_instance.Table{Name: collection, Columns: []string{}}

			var docRaw bson.Raw
			err := m.session.DB(dbname).C(collection).Find(nil).One(&docRaw)
			if err != nil {
				return nil, err
			}

			var doc bson.RawD
			if err := docRaw.Unmarshal(&doc); err != nil {
				if err != nil {
					return nil, err
				}
			}

			for _, f := range doc {
				table.Columns = append(table.Columns, f.Name)
			}
			db.Tables = append(db.Tables, table)
		}
		s.Databases = append(s.Databases, db)
	}
	return s, nil
}
