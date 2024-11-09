package pgsql


import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"context"
    "github.com/jackc/pgx/v5/pgconn"
    "github.com/jackc/pgx/v5"
    util "github.com/prr123/utility/utilLib"
)

type dbCon struct {
	dburl string
	DbName string
	dbcon *pgx.Conn
	ctx context.Context
	dbg bool
}

type dbrowval []any


func (db *dbCon)InitDb() (err error){
    ctx := context.Background()
//    dburl :="postgresql://dbuser:dbtest@/testdb"

    dbcon, err := pgx.Connect(ctx, db.dburl)
    if err != nil {return fmt.Errorf("Unable to create db connection: %v\n", err)}
//    defer dbcon.Close(ctx)
	(*db).ctx = ctx
	(*db).dbcon = dbcon
	return nil
}

func ProcCli(cliStr []string) (dbpt *dbCon, err error) {

	var dbinp struct {
		user string
		pwd string
		db string
	}

	var db dbCon
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
    if err != nil {return nil, fmt.Errorf("util.ParseFlags: %v\n", err)}

    db.dbg = false
    _, ok := flagMap["dbg"]
    if ok {db.dbg = true}

    uval, ok := flagMap["user"]
    if !ok {
        return nil, fmt.Errorf("cli -- user not provided!")
    } else {
        if uval.(string) == "none" {return nil, fmt.Errorf("error: no user name provided!")}
        dbinp.user = uval.(string)
    }

    pval, ok := flagMap["pwd"]
    if !ok {
        return nil, fmt.Errorf("cli -- pwd not provided!")
    } else {
        if pval.(string) == "none" {return nil, fmt.Errorf("error: no password provided!")}
        dbinp.pwd = pval.(string)
    }

    dbval, ok := flagMap["db"]
    if !ok {
        return nil, fmt.Errorf("cli -- db not provided!")
    } else {
        if dbval.(string) == "none" {return nil, fmt.Errorf("error: no db name provided!")}
        dbinp.db = dbval.(string)
    }
// postgresql://dbuser:dbtest@/testdb
//    dburl :="postgresql://dbuser:dbtest@/testdb"
	db.dburl = "postgresql://" + dbinp.user + ":" + dbinp.db + "@/" + dbinp.pwd
	db.DbName = dbinp.db
    return &db, nil
}

func GetSql(dbnam string) (sqlStr string, err error) {

    reader := bufio.NewReader(os.Stdin)

//    multi := false
    for i:=0; i<5; i++ {
        if i== 0 {
            fmt.Printf("%s> ", dbnam)
        } else {
            fmt.Printf("%s sql %d> ", dbnam, i)
        }
        line, err := reader.ReadString('\n')
        if err != nil {return sqlStr, fmt.Errorf("reading line: %v\n", err)}

        if len(line) == 1 {break}
		if edx:=strings.Index(line, "exit"); edx>-1 {return "exit", nil}

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

func (db *dbCon) ProcShow(sql string) (err error) {
    fmt.Printf("Proc show: %s\n", sql)

    return nil
}

func (db *dbCon) ProcSql(sql string) (err error) {
    fmt.Printf("Proc sql: %s\n", sql)


    return nil
}


func (db *dbCon) ProcSelect(query string) (fields []pgconn.FieldDescription, valList []dbrowval, err error) {
    fmt.Printf("Proc select: %s\n", query)
	ctx := (*db).ctx
	pgcon := (*db).dbcon
	rows, err := pgcon.Query(ctx, query)
    if err != nil {
        return fields, valList, fmt.Errorf("select query failed: %v\n", err)
    }
    defer rows.Close()

    fields = rows.FieldDescriptions()
	valList = make([]dbrowval,0,50)
	count:=0
    for rows.Next() {
        val, err := rows.Values()
        if err != nil {
            return fields, valList, fmt.Errorf("row[%d]: get row values : %v\n", count, err)
        }
		valList = append(valList, val)
 		count++
	}

    return fields, valList, nil
}

func (db *dbCon)CloseDb() {
	ctx := (*db).ctx
	pgcon := (*db).dbcon
	pgcon.Close(ctx)
}

func PrintDburl(db *dbCon) {
	fmt.Printf("dburl: %s\n",db.dburl)
}

func (db *dbCon)PrintSelect(fields []pgconn.FieldDescription, values []dbrowval) {

	pgm := (*db).dbcon.TypeMap()
	for i:=0; i<len(fields); i++ {
		fmt.Printf("%-20s|", fields[i].Name)
	}
	fmt.Printf("\n")

	for i:=0; i<len(fields); i++ {
		field := fields[i]
       ftyp, res := pgm.TypeForOID(field.DataTypeOID)
        if !res {
            fmt.Printf("%-20s|","unkown")
        } else {
            fmt.Printf("%-20s|", ftyp.Name)
        }
	}
	fmt.Printf("\n")



	if len(values) == 0 {return}
	val := values[0]
	for i:=0; i<len(fields); i++ {
		fmt.Printf("%-20T|", val[i])
	}
	fmt.Printf("\n")
	for i:=0; i<len(fields); i++ {
		fmt.Printf("====================|")
	}
	fmt.Printf("\n")


	for i:=0; i<len(values); i++ {
		valrow:= values[i]
		for j:=0; j<len(fields); j++ {
			fmt.Printf("%-20v|", valrow[j])
		}
		fmt.Printf("\n")
	}
}
