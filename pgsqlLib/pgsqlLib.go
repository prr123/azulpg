package pgsql


import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"context"
    "github.com/jackc/pgx/v5"
    util "github.com/prr123/utility/utilLib"
)

type dbCon struct {
	dburl string
	dbcon *pgx.Conn
	ctx context.Context
	dbg bool
}

func (db dbCon)InitDb() (err error){
    ctx := context.Background()
//    dburl :="postgresql://dbuser:dbtest@/testdb"

    dbcon, err := pgx.Connect(ctx, db.dburl)
    if err != nil {return fmt.Errorf("Unable to create db connection: %v\n", err)}
    defer dbcon.Close(ctx)
	db.ctx = ctx
	db.dbcon = dbcon
	return nil
}

func ProcCli(cliStr []string) (db dbCon, err error) {

	var dbinp struct {
		user string
		pwd string
		db string
	}
//    numarg := len(os.Args)
    numarg := len(cliStr)

    flags:=[]string{"dbg", "user", "pwd", "db"}

    useStr := " /user=username /pwd=password /db=database [/dbg]"
    helpStr := "pgsql"

    if numarg > len(flags) +1 {
        fmt.Println("too many arguments in cl!")
        fmt.Println("usage: %s %s\n", cliStr[0], useStr)
        os.Exit(-1)
    }

    if numarg == 1 || (numarg > 1 && cliStr[1] == "help") {
        fmt.Printf("help: %s\n", helpStr)
        fmt.Printf("usage is: %s %s\n", cliStr[0], useStr)
        os.Exit(1)
    }

    flagMap, err := util.ParseFlags(cliStr, flags)
    if err != nil {return db, fmt.Errorf("util.ParseFlags: %v\n", err)}

    db.dbg = false
    _, ok := flagMap["dbg"]
    if ok {db.dbg = true}

    uval, ok := flagMap["user"]
    if !ok {
        return db, fmt.Errorf("cli -- user not provided!")
    } else {
        if uval.(string) == "none" {return db, fmt.Errorf("error: no user name provided!")}
        dbinp.user = uval.(string)
    }

    pval, ok := flagMap["pwd"]
    if !ok {
        return db, fmt.Errorf("cli -- pwd not provided!")
    } else {
        if pval.(string) == "none" {return db, fmt.Errorf("error: no password provided!")}
        dbinp.pwd = pval.(string)
    }

    dbval, ok := flagMap["db"]
    if !ok {
        return db, fmt.Errorf("cli -- db not provided!")
    } else {
        if dbval.(string) == "none" {return db, fmt.Errorf("error: no db name provided!")}
        dbinp.db = dbval.(string)
    }
// postgresql://dbuser:dbtest@/testdb
//    dburl :="postgresql://dbuser:dbtest@/testdb"
	db.dburl = "postgresql://" + dbinp.user + ":" + dbinp.db + "@/" + dbinp.pwd

    return db, nil
}

func GetSql() (sqlStr string, err error) {

    reader := bufio.NewReader(os.Stdin)

//    multi := false
    for i:=0; i<5; i++ {
        if i== 0 {
            fmt.Printf("sql> ")
        } else {
            fmt.Printf("sql %d> ", i)
        }
        line, err := reader.ReadString('\n')
        if err != nil {return sqlStr, fmt.Errorf("reading line: %v\n", err)}

        if len(line) == 1 {break}
		if edx:=strings.Index(line, "end"); edx>-1 {return "end", nil}

        trimLine := strings.TrimSuffix(line,"\n")
//      fmt.Printf("line>%s<\n", trimLine)
        if i>0 {
			sqlStr = sqlStr + " " + trimLine
//			multi = true
		} else {
			sqlStr = trimLine
		}
        idx := strings.Index(sqlStr, ";")
        if idx > -1 {break}
    }
	if len(sqlStr) > 5 {
        idx := strings.Index(sqlStr, ";")
        if idx == -1 {sqlStr += ";"}
	}
	return sqlStr, nil
}

func ProcShow(sql string) (err error) {
    fmt.Printf("Proc show: %s\n", sql)
    return nil
}

func ProcSql(sql string) (err error) {
    fmt.Printf("Proc sql: %s\n", sql)
    return nil
}

func PrintDburl(db dbCon) {
	fmt.Printf("dburl: %s\n",db.dburl)
}
