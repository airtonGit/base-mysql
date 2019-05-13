package basemysql

import (
	"database/sql"
	"fmt"
	"strings"

	//Driver mysql para manter conexao com banco
	_ "github.com/go-sql-driver/mysql"
)

//SQLExecError diferente de nil quando alguma execução
//SQL falha
// type SQLExecError struct {
// 	When time.Time
// 	What string
// }

// func (e *SQLExecError) Error() string {
// 	return fmt.Sprintf("As %v, %s", e.When, e.What)
// }

//BaseMysql possui principais metodos para execução select, insert update e delete em
//banco mysql, mantém 1 conexão ativa sem pool
type BaseMysql interface {
	//conexao *sql.DB
	connect(hostname string, dbname string, username string, password string) *sql.DB
	FetchLines(query string, values []string, options []string)
}

type db struct {
	con *sql.DB //conexão banco
}

//Connect mantém conexao ativa com banco
func (db *db) Connect(dbUser, dbPassword, dbHostname, dbPort, database string) error {
	con, err := sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		dbUser,
		dbPassword,
		dbHostname,
		dbPort, database))

	if err != nil {
		return fmt.Errorf("Erro ao conectar banco: %v", err)
	}
	db.con = con
	// TODO defer db.Close()
	return nil
}

func (db *db) Close() {
	db.con.Close()
}

func (db *db) FetchLines(table string, selectFields []string, where string, valuesWhere []interface{}) (result *[]interface{}, er error) { //options []string
	//log::info(__METHOD__ . " option is " . $option . " padrao é:" . PDO::FETCH_BOTH);
	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s",
		strings.Join(selectFields, ", "), table, where)
	rows, err := db.con.Query(query, valuesWhere...)
	if err != nil {
		return nil, fmt.Errorf("FAIL TO PREPARE SELECT: %v", err)
	}

	scannedRows := make([]interface{}, 0)

	readCols := make([]interface{}, len(selectFields))
	writeCols := make([]string, len(selectFields))
	for i := range writeCols {
		readCols[i] = &writeCols[i]
	}

	for rows.Next() {
		err := rows.Scan(readCols...)
		if err != nil {
			return nil, fmt.Errorf("Falha ao scan campos: %v", err)
		}
		scannedRows = append(scannedRows, writeCols)
	}

	fmt.Println("Linhas scanedas ", scannedRows)
	return &scannedRows, nil
}

// func (db *db) Insert(table string, row interface{}) { //columns map[string]string) {
// 	var sqlColumns, sqlValues, values []string

// 	for i, v := range columns {
// 		sqlColumns = append(sqlColumns, fmt.Sprintf("`%s`", i)) //preciso do key - nome da coluna e nao seu valor
// 		sqlValues = append(sqlValues, "?")
// 		values = append(values, v)
// 	}
// 	query := fmt.Sprintf(
// 		"INSERT INTO `%s` (%s) VALUES (%s);",
// 		table,
// 		strings.Join(sqlColumns, ", "),
// 		strings.Join(sqlValues, ", "))
// 	stmt, err := db.con.Prepare(query)
// 	if err != nil {
// 		panic("FAIL TO PREPARE INSERT")
// 	}
// 	args := make([]interface{}, len(values))
// 	for i, v := range values {
// 		args[i] = v
// 	}
// 	result, err := stmt.Exec(args...)
// 	if err != nil {
// 		panic("UNABLE TO INSERT OBJECT")
// 	} else {
// 		//return self::$conn->lastInsertId();
// 	}
// }

func (db *db) Update(table string, id uint, columns map[string]string) (sql.Result, error) {
	var sqlColumns, values []string
	for i, v := range columns {
		sqlColumns = append(sqlColumns, fmt.Sprintf("`%s` = ?", i))
		values = append(values, v)
	}

	query := fmt.Sprintf(
		"UPDATE `%s` SET %s WHERE `id` = %d;",
		table,
		strings.Join(sqlColumns, ", "),
		id)
	stmt, err := db.con.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("FAIL TO PREPARE UPDATE QUERY: %v", err)
	}
	args := make([]interface{}, len(values))
	for i, v := range values {
		args[i] = v
	}
	result, err := stmt.Exec(args...)
	if err != nil {
		return nil, fmt.Errorf("UNABLE TO UPDATE OBJECT: %v", err)
	}
	return result, nil
}

func (db *db) delete(table, id string) (sql.Result, error) {
	stmt, err := db.con.Prepare(fmt.Sprintf("DELETE FROM `%s` WHERE `ID` = ?", table))
	result, err := stmt.Exec(id)
	if err != nil {
		return nil, fmt.Errorf("UNABLE TO DELETE OBJECT: %v", err)
	}
	return result, nil
}

func (db *db) startTransaction() (sql.Result, error) {
	stmt, err := db.con.Prepare("START TRANSACTION")
	if err != nil {
		panic("ERRO AO Prepare query")
	}
	result, err := stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("UNABLE TO START TRANSACTION: %v", err)
	}
	return result, nil
}

func (db *db) commit() (sql.Result, error) {
	stmt, err := db.con.Prepare("COMMIT")
	if err != nil {
		return nil, fmt.Errorf("UNABLE TO PREPARE COMMIT: %v", err)
	}
	result, err := stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("UNABLE TO COMMIT TRANSACTION: %v", err)
	}
	return result, nil
}

func (db *db) rollback() (sql.Result, error) {
	stmt, err := db.con.Prepare("ROLLBACK")
	if err != nil {
		return nil, fmt.Errorf("UNABLE TO PREPARE ROLLBACK: %v", err)
	}
	result, err := stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("UNABLE TO ROLLBACK TRANSACTION: %v", err)
	}
	return result, nil
}
