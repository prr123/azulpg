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
	"bufio"
	"strings"
	util "github.com/prr123/utility/utilLib"
	)

type cliObj struct {
	dbg bool
	user string
	pwd string
	db string
}

func procCli() (cli cliObj, err error) {
	numarg := len(os.Args)

    flags:=[]string{"dbg", "user", "pwd", "db"}

    useStr := " /user=username /pwd=password /db=database [/dbg]"
    helpStr := "pgsql"

    if numarg > len(flags) +1 {
        fmt.Println("too many arguments in cl!")
        fmt.Println("usage: %s %s\n", os.Args[0], useStr)
        os.Exit(-1)
    }

    if numarg == 1 || (numarg > 1 && os.Args[1] == "help") {
        fmt.Printf("help: %s\n", helpStr)
        fmt.Printf("usage is: %s %s\n", os.Args[0], useStr)
        os.Exit(1)
    }

    flagMap, err := util.ParseFlags(os.Args, flags)
    if err != nil {return cli, fmt.Errorf("util.ParseFlags: %v\n", err)}

    cli.dbg = false
    _, ok := flagMap["dbg"]
    if ok {cli.dbg = true}

    uval, ok := flagMap["user"]
    if !ok {
        return cli, fmt.Errorf("cli -- user not provided!")
    } else {
        if uval.(string) == "none" {return cli, fmt.Errorf("error: no user name provided!")}
        cli.user = uval.(string)
    }

    pval, ok := flagMap["pwd"]
    if !ok {
        return cli, fmt.Errorf("cli -- pwd not provided!")
    } else {
        if pval.(string) == "none" {return cli, fmt.Errorf("error: no password provided!")}
        cli.pwd = pval.(string)
    }

    dbval, ok := flagMap["db"]
    if !ok {
        return cli, fmt.Errorf("cli -- db not provided!")
    } else {
        if dbval.(string) == "none" {return cli, fmt.Errorf("error: no db name provided!")}
        cli.db = dbval.(string)
    }

	return cli, nil
}

func PrintCli(cli cliObj) {
	fmt.Println("********** CLI **********")
	fmt.Printf("dbg:  %t\n", cli.dbg)
	fmt.Printf("user: %s\n", cli.user)
	fmt.Printf("pwd:  %s\n", cli.pwd)
	fmt.Printf("db:   %s\n", cli.db)
	fmt.Println("*************************")

}

func PrintCli2(cli cliObj) {
	fmt.Printf("CLI>")
	fmt.Printf(" user: %s", cli.user)
	fmt.Printf(" pwd: %s", cli.pwd)
	fmt.Printf(" db: %s", cli.db)
	fmt.Printf(" dbg:  %t\n", cli.dbg)
}


func getSql() (sqlStr string, err error) {

	reader := bufio.NewReader(os.Stdin)

    for i:=0; i<5; i++ {
		if i== 0 {
			fmt.Printf("sql> ")
		} else {
			fmt.Printf("sql %d> ", i)
		}
		line, err := reader.ReadString('\n')
		if err != nil {log.Fatalf("error -- read line: %v\n", err)}
		if len(line) == 1 {break}
		trimLine := strings.TrimSuffix(line,"\n")
//		fmt.Printf("line>%s<\n", trimLine)
		if i>0 {trimLine = " " + trimLine}
		sqlStr += trimLine
	}
	sqlStr += ";"
	return sqlStr, nil
}

func main () {

//	var user, pwd string

	cli, err := procCli()
	if err != nil {log.Fatalf("error -- proCli: %v\n", err)}

	PrintCli2(cli)

	sqlStr :=""
	for i:=0; i< 5; i++ {
		var err error
		sqlStr, err = getSql()
		if err != nil {log.Fatalf("error -- sqlStr: %v\n", err)}

		idx := strings.Index(sqlStr, "end")
		if idx == 0 {os.Exit(0)}
		fmt.Printf("sql: %s\n", sqlStr)
		cmds := strings.Split(sqlStr, " ")
		fmt.Printf("%s: %s\n", cmds[0], sqlStr)
		switch cmds[0] {
		case "show":
			err := ProcShow(sqlStr)
			if err != nil {fmt.Printf("error -- show: %v\n", err)}
		case "end":
			break
		default:
			err:= ProcSql(sqlStr)
			if err != nil {fmt.Printf("error -- sql: %v\n", err)}
		}
	}
}

func ProcShow(sql string) (err error) {
	fmt.Printf("Proc show: %s\n", sql)
	return nil
}

func ProcSql(sql string) (err error) {
	fmt.Printf("Proc sql: %s\n", sql)
	return nil
}
