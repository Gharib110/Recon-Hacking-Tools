package database_enumeration

import (
	"fmt"
	"regexp"
)

type DBInstance interface {
	GetSchema() (*Schema, error)
}

type Schema struct {
	Databases []Database
}

type Database struct {
	Name   string
	Tables []Table
}

type Table struct {
	Name    string
	Columns []string
}

func Search(m DBInstance) error {
	s, err := m.GetSchema()
	if err != nil {
		return err
	}

	re := getRegex()
	for _, database := range s.Databases {
		for _, table := range database.Tables {
			for _, field := range table.Columns {
				for _, r := range re {
					if r.MatchString(field) {
						fmt.Println(database)
						fmt.Printf("[+] HIT: %s\n", field)
					}
				}
			}
		}
	}
	return nil
}

// getRegex sample important regex of a database
// should determine the structures of a DB then create regex for data inside of it
func getRegex() []*regexp.Regexp {
	return []*regexp.Regexp{
		regexp.MustCompile(`(?i)social`),
		regexp.MustCompile(`(?i)ssn`),
		regexp.MustCompile(`(?i)pass(word)?`),
		regexp.MustCompile(`(?i)hash`),
		regexp.MustCompile(`(?i)security`),
		regexp.MustCompile(`(?i)key`),
	}
}

func (s *Schema) String() string {
	var ret string
	for _, database := range s.Databases {
		ret += fmt.Sprint(database.String() + "\n")
	}
	return ret
}

func (d *Database) String() string {
	ret := fmt.Sprintf("[DB] = %+s\n", d.Name)
	for _, table := range d.Tables {
		ret += table.String()
	}
	return ret
}

func (t *Table) String() string {
	ret := fmt.Sprintf("    [TABLE] = %+s\n", t.Name)
	for _, field := range t.Columns {
		ret += fmt.Sprintf("       [COL] = %+s\n", field)
	}
	return ret
}
