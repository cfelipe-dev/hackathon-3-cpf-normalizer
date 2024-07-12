package main

import (
	zendeskdata "cpf-normalizer/internal/zendeskdata"
	"encoding/json"
	"fmt"
	"log"
)

func GetDataAndFormatCPFs(userPhoneNumber string, formatCPF bool) string {
	data, err := zendeskdata.SearchEndUser(userPhoneNumber, formatCPF)

	if err != nil {
		log.Fatalf("Erro ao buscar o usuário: %v", err)
	}

	json, err := json.Marshal(data)

	if err != nil {
		log.Fatalf("Erro ao converter para JSON: %v", err)
	}

	fmt.Println(string(json))

	return string(json)
}
