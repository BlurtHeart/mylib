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
	ms.DB.Close()
}

// execute
// need to implement
func (ms *MySQLdb) Execute(sql string) error {
	_, err := ms.DB.Exec(sql)
	return err
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
