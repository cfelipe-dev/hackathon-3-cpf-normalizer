package normalizecpf

import (
	"context"
	"log"
	"os"
	"strings"

	godotenv "github.com/joho/godotenv"
	cpfcnpj "github.com/klassmann/cpfcnpj"
	openaisdk "github.com/sashabaranov/go-openai"
)

func SendRequest(CPF string, formatCPF bool) ([]string, error) {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apikey := os.Getenv("OPENAI_API_KEY")

	if apikey == "" {
		log.Fatal("OPENAI_API_KEY is not set")
	}

	client := openaisdk.NewClient(apikey)

	response, err := client.CreateChatCompletion(
		context.Background(),
		openaisdk.ChatCompletionRequest{
			Model: "gpt-3.5-turbo",
			Messages: []openaisdk.ChatCompletionMessage{
				{
					Role:    "system",
					Content: "Você é um assistente especializado em identificação e validação de CPFs. Dado um texto, identifique todos os possíveis CPFs, que pode ser nos seguintes formatos: apenas números, com pontos entre os números, com pontos e traços e com espaços. Valide se são corretos e formate-os no padrão XXXXXXXXXXX, retornando apenas o(s) CPF(s) formatado(s).",
				},
				{
					Role:    "user",
					Content: CPF,
				},
			},
		},
	)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	cpfs := strings.Split(response.Choices[0].Message.Content, "\n")

	var validCPFs []string

	for _, cpftext := range cpfs {
		cpf := cpfcnpj.NewCPF(cpftext)

		if cpf.IsValid() {
			if formatCPF {
				validCPFs = append(validCPFs, cpf.String())
			} else {
				validCPFs = append(validCPFs, cpftext)
			}
		}
	}

	return validCPFs, nil
}
