# mysql驱动

## 简介

本驱动在[github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)的基础上封装，提供类似于
python的第三方库MySQLdb的使用方式。

## 安装

    go get -u github.com/blurty/mylib/mysql

## 使用

### 插入数据

    m := mysql.NewMysql("127.0.0.1", 3306, "root", "111111", "test")
    m.Connect()
    sql := "insert user(username, password, token) values(\"wang\", \"aaaaaa\", \"aaa-aaaa-aaaaa\");"
    n, err := m.Insert(sql)
    if err != nil {
        fmt.Printf("got err:%v", err)
        return
    }
    fmt.Printf("insert operation affects %d rows\n", n)