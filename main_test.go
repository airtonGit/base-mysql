package basemysql

import (
	"fmt"
	"testing"
)

func TestUpdate(t *testing.T) {
	var banco db
	banco.Connect("root", "eunaouso", "127.0.0.1", "3307", "projeto_connection")
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
	var banco db
	banco.Connect("root", "eunaouso", "127.0.0.1", "3307", "projeto_connection")
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
