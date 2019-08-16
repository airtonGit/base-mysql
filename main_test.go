package basemysql

import (
	"fmt"
	"testing"
	"time"
)

func TestUpdate(t *testing.T) {
	banco := Db{User: "root", Password: "eunaouso", Host: "127.0.0.1", Port: "3307", Database: "projeto_connection"}
	banco.Connect()
	defer banco.Close()

	var m map[string]string

	m = make(map[string]string, 2)

	m["name"] = "HELIO SAVITSKI"
	m["cpf"] = "02620527929"

	_, err := banco.Update("instructor", 1, m)
	if err != nil {
		t.Error(`Update falhou`)
	}
}

func TestFetchLines(t *testing.T) {
	// var banco db
	// banco.Connect("root", "eunaouso", "127.0.0.1", "3307", "projeto_connection")
	banco := Db{User: "root", Password: "eunaouso", Host: "127.0.0.1", Port: "3307", Database: "projeto_connection"}
	banco.Connect()
	defer banco.Close()

	campos := []string{"name", "cpf"}
	where := "id = ?"
	var valuesWhere []interface{}
	valuesWhere = make([]interface{}, 0)
	valuesWhere = append(valuesWhere, 1)

	_, err := banco.FetchLines("instructor", campos, where, valuesWhere)
	if err != nil {
		fmt.Println("Falha ao select: ", err)
	}
}

func TestInsert(t *testing.T) {
	banco := Db{User: "root", Password: "eunaouso", Host: "127.0.0.1", Port: "3307", Database: "projeto_connection"}
	banco.Connect()
	defer banco.Close()

	agora := time.Now()
	agoraMysql := agora.Format("2006-01-02 15:04:05")
	fmt.Println("Agora", agoraMysql)
	var tests = []struct {
		input []interface{}
		want  int64
	}{
		/*caso de teste*/ {
			/* input */ []interface{}{nil, "Airton Teste1", agoraMysql, agoraMysql, "Desc1"}, // region id, date_created, date_modified, name, description
			/* want param */ 1},
		{
			/* input */ []interface{}{nil, "Airton Teste2", agoraMysql, agoraMysql, "Desc2"},
			/* want param */ 1}}

	for _, test := range tests {
		if got, err := banco.Insert("region", test.input); got < 0 || err != nil {
			t.Errorf("Insert(%s, %s) = %d", "region", test.input, got)
		}
	}

}
