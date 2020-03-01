package basemysql

import (
	"database/sql"
	"fmt"
	"strings"

	//Driver mysql para manter conexao com banco
	_ "github.com/go-sql-driver/mysql"
)

//Db mantem conexao com banco MySQL
//Provê métodos para update/insert/select
type Db struct {
	con *sql.DB //conexão banco
	//Conn deve substituir con em futura versão major
	Conn                                 *sql.DB //conexão banco
	User, Password, Host, Port, Database string
}

//Connect mantém conexao ativa com banco
func (db *Db) Connect() error {
	con, err := sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		db.User,
		db.Password,
		db.Host,
		db.Port, db.Database))

	if err != nil {
		return fmt.Errorf("Erro ao conectar banco: %v", err)
	}
	db.con = con
	db.Conn = con
	// TODO defer db.Close()
	return nil
}

//Close encerra conexao com banco MySQL
func (db *Db) Close() {
	db.con.Close()
}

//FetchLines retorna select de campos fornecidos em slice de strings
func (db *Db) FetchLines(table string, selectFields []string, where string, valuesWhere []interface{}) (result *[]interface{}, er error) { //options []string
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

//Insert - Adiciona nova linha na tabela informada, deve-se informar todas as colunas na mesma ordem dos campos na declaração da tabela.
func (db *Db) Insert(table string, columnsValue []interface{}) (int64, error) {
	var sqlAnchors []string

	for range columnsValue {
		sqlAnchors = append(sqlAnchors, "?")
	}

	query := fmt.Sprintf(
		"INSERT INTO `%s` VALUES (%s);",
		table,
		strings.Join(sqlAnchors, ", "))
	stmt, err := db.con.Prepare(query)
	if err != nil {
		return -1, fmt.Errorf("FAIL TO PREPARE INSERT error: %s", err.Error())
	}
	result, err := stmt.Exec(columnsValue...)
	if err != nil {
		return -1, fmt.Errorf("UNABLE TO INSERT OBJECT error: %s", err.Error())
	}
	return result.LastInsertId()
}

//Update realiza update de linhas selecionadas por where
func (db *Db) Update(table string, id uint, columns map[string]string) (sql.Result, error) {
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

func (db *Db) delete(table, id string) (sql.Result, error) {
	stmt, err := db.con.Prepare(fmt.Sprintf("DELETE FROM `%s` WHERE `ID` = ?", table))
	result, err := stmt.Exec(id)
	if err != nil {
		return nil, fmt.Errorf("UNABLE TO DELETE OBJECT: %v", err)
	}
	return result, nil
}

func (db *Db) startTransaction() (sql.Result, error) {
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

func (db *Db) commit() (sql.Result, error) {
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

func (db *Db) rollback() (sql.Result, error) {
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
