// pgsql
// a program that is an alternative to psql
//
// author: prr, azul software
// date: 27 Sept 2024
// copyright 2024 prr, azul software

package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	pgdb "db/azuldb/pgsqlLib"
	)

func main () {

	db, err := pgdb.ProcCli(os.Args)
	if err != nil {log.Fatalf("error -- proCli: %v\n", err)}

	pgdb.PrintDburl(db)

	err = db.InitDb()
	if err != nil {log.Fatalf("error -- initdb: %v\n", err)}

	log.Printf("db connected!")

	sqlStr :=""
	for i:=0; i< 5; i++ {
		var err error
		sqlStr, err = pgdb.GetSql(db.DbName)
		if err != nil {log.Fatalf("error -- sqlStr: %v\n", err)}
//		fmt.Printf("sql: %s\n", sqlStr)
		if len(sqlStr) == 0 {continue}
		idx := strings.Index(sqlStr, "exit")
		if idx >= 0 {
			fmt.Println("*** exiting ***")
			os.Exit(0)
		}
		cmds := strings.Split(sqlStr, " ")
//		fmt.Printf("%s: %s\n", cmds[0], sqlStr)
		switch cmds[0] {
		case "select":
			fields, values, err := db.ProcSelect(sqlStr)
			if err != nil {fmt.Printf("error -- select: %v\n", err)}
			db.PrintSelect(fields, values)
		case "show":
			err := db.ProcShow(sqlStr)
			if err != nil {fmt.Printf("error -- show: %v\n", err)}
		case "exit":
			break
		default:
			err:= db.ProcSql(sqlStr)
			if err != nil {fmt.Printf("error -- sql: %v\n", err)}
		}
	}
	db.CloseDb()
	fmt.Println("*** success exiting ***")
}

