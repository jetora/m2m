package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	_ "mysql"
	"os"
)

var myDB *sql.DB

type Hostinfo struct {
	DBUser,
	DBPassword,
	DBname,
	DBHost,
	DBPort,
	DBChar string
}

func connMysql(host *Hostinfo) (*sql.DB, error) {
	if host.DBHost != "" {
		host.DBHost = "tcp(" + host.DBHost + ":" + host.DBPort + ")"
	}
	db, err := sql.Open("mysql", host.DBUser+":"+host.DBPassword+"@"+host.DBHost+"/"+host.DBname+"?charset="+host.DBChar)
	return db, err
}
func SetDB(ip string) (myDB *sql.DB) {
	var server_info Hostinfo
	server_info.DBUser = "root"
	server_info.DBPassword = "rootpass"
	server_info.DBname = "test"
	server_info.DBHost = ip
	server_info.DBPort = "4306"
	server_info.DBChar = "utf8"
	myDB, _ = connMysql(&server_info)
	return myDB
}
func tab2txt(ip string) {
	myDB = SetDB(ip)
	defer myDB.Close()
	rows, err := myDB.Query("select * from jdorders_0.orders_0 limit 2")
	defer rows.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	tsql := bytes.Buffer{}
	tsql.WriteString("insert into jdorders_0.orders_0 (")
	columns, err := rows.Columns()
	for i, column := range columns {
		if i != len(columns)-1 {
			tsql.WriteString(column + ",")
		} else {
			tsql.WriteString(column)
		}
	}
	tsql.WriteString(") values \n")
	if err != nil {
		fmt.Println(err.Error())
	}
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			fmt.Println(err.Error())
		}
		tsql.WriteString("(")
		var value string
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			if i != len(columns)-1 {
				tsql.WriteString("'" + value + "'" + ",")
			} else {
				tsql.WriteString("'" + value + "'")
			}
		}
		tsql.WriteString("),\n")
	}
	if err = rows.Err(); err != nil {
		fmt.Println(err.Error())
	}
	outputFile, outputError := os.OpenFile("output.txt", os.O_WRONLY|os.O_CREATE, 0666)
	if outputError != nil {
		fmt.Printf("An error occurred with file opening or creation\n")
		return
	}
	defer outputFile.Close()
	outputWriter := bufio.NewWriter(outputFile)
	outputString := tsql.String()
	outputWriter.WriteString(outputString)
	outputWriter.Flush()
}
func main() {
	from_ip := "10.0.57.74"
	tab2txt(from_ip)
}
