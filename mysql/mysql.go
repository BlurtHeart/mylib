// Go MySQL Driver - A MySQL-Driver similar to python's MySQLdb
// by using open source project github.com/go-sql-driver/mysql
//
// Copyright 2017 Blurt Authors. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at http://mozilla.org/MPL/2.0/.

package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLdb struct {
	Host     string
	Port     int
	UserName string
	Password string
	DBName   string
	Charset  string
	DB       *sql.DB
	Rows     *sql.Rows
}

func NewMysql(host string, port int, username, password, dbname string) *MySQLdb {
	return &MySQLdb{
		Host:     host,
		Port:     port,
		UserName: username,
		Password: password,
		DBName:   dbname,
		Charset:  "utf8",
	}
}

func (ms *MySQLdb) SetCharset(charset string) {
	ms.Charset = charset
}

func (ms *MySQLdb) Connect() (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", ms.UserName, ms.Password, ms.Host, ms.Port, ms.DBName, ms.Charset)
	ms.DB, err = sql.Open("mysql", dsn)
	return
}

func (ms *MySQLdb) Close() {
	if ms.Rows != nil {
		ms.Rows.Close()
	}
	if ms.DB != nil {
		ms.DB.Close()
	}
}

// execute
func (ms *MySQLdb) Execute(sql string) error {
	_, err := ms.DB.Exec(sql)
	return err
}

func (ms *MySQLdb) Query(sql string, args ...interface{}) (err error) {
	rows, err := ms.DB.Query(sql)
	if err != nil {
		return
	}
	ms.Rows = rows
	return
}

// when fetch all results, you can loop yourself by using this method
func (ms *MySQLdb) FetchOne(args ...interface{}) error {
	if ms.Rows.Next() {
		return ms.Rows.Scan(args)
	}
	return nil
}

// timestamp:v["create_time"].(string)
// int:strconv.Atoi(v["id"].(string))
// string:v["name"].(string)
func (ms *MySQLdb) FetchAll() ([]map[string]interface{}, error) {
	colums, err := ms.Rows.Columns()
	if err != nil {
		return nil, err
	}
	count := len(colums)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for ms.Rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		ms.Rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range colums {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	return tableData, nil
}

// return number of affect rows by insert operation
func (ms *MySQLdb) Insert(sql string) (n int64, err error) {
	res, err := ms.DB.Exec(sql)
	if err != nil {
		return
	}
	n, err = res.RowsAffected()
	return
}

// create table
func (ms *MySQLdb) CreateTable(sql string) error {
	_, err := ms.DB.Exec(sql)
	return err
}

// need to be implement
// this method aim to improve delete efficiency
func (ms *MySQLdb) Delete(sql string) error {
	panic("need to implement")
	return nil
}
