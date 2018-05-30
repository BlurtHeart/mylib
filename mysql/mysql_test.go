// Go MySQL Driver - A MySQL-Driver similar to python's MySQLdb
// by using open source project github.com/go-sql-driver/mysql
//
// Copyright 2017 Blurt Authors. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at http://mozilla.org/MPL/2.0/.

package mysql

//func TestMySql_Insert(t *testing.T) {
//	m := NewMysql("127.0.0.1", 3306, "root", "111111", "test")
//	m.Connect()
//	defer m.Close()
//	_, err := m.Execute("insert user(username, password, token) values(\"wang\", \"aaaaaa\", \"aaa-aaaa-aaaaa\");")
//	if err != nil {
//		t.Errorf("got err:%v", err)
//	}
//}

//func TestMySQLdb_CreateTable(t *testing.T) {
//	m := NewMysql("127.0.0.1", 3306, "root", "111111", "test")
//	m.Connect()
//	defer m.Close()
//	err := m.CreateTable("create table if not exists user_test(id int unsigned auto_increment, name varchar(20) not null, phone varchar(13), primary key(id))ENGINE=InnoDB DEFAULT CHARSET=utf8;")
//	if err != nil {
//		t.Errorf("got err:%v", err)
//	}
//}
