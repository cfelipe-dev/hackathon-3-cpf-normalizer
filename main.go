package main

import (
	zendeskdata "cpf-normalizer/internal/zendeskdata"
	"encoding/json"
	"fmt"
	"log"
)

func zendeskData() {
	data, err := zendeskdata.SearchEndUser("19999999999")

	if err != nil {
		log.Fatalf("Erro ao buscar o usuário: %v", err)
	}

	fmt.Println(data)
}

func main() {
	data, err := zendeskdata.SearchEndUser("19999999999")

	if err != nil {
		log.Fatalf("Erro ao buscar o usuário: %v", err)
	}

	json, err := json.Marshal(data)

	if err != nil {
		log.Fatalf("Erro ao converter para JSON: %v", err)
	}

	fmt.Println("CPFs: ", string(json))
}
