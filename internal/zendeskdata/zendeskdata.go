package zendeskdata

import (
	"bytes"
	"cpf-normalizer/internal/normalizecpf"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type FormattedCPF struct {
	CPF string `json:"cpf"`
}

type SearchEndUserType struct {
	Data struct {
		Search struct {
			Edges []struct {
				Node struct {
					Notes string `json:"notes"`
				} `json:"node"`
			} `json:"edges"`
		} `json:"search"`
	} `json:"data"`
	Errors []string `json:"errors"`
}

type QueryVariables struct {
	Limit      int    `json:"limit"`
	PageNumber int    `json:"pageNumber"`
	Query      string `json:"query"`
}

type RequestBody struct {
	Query     string         `json:"query"`
	Variables QueryVariables `json:"variables"`
}

const searchEndUserQuery = `query searchEndUsers($query: String!, $limit: Int, $pageNumber: Int) {
  search(type: USER, limit: $limit, query: $query, pageNumber: $pageNumber) {
    edges {
      node {
        ...EndUserFragment
        __typename
      }
      __typename
    }
    pageInfo {
      hasNextPage
      hasPreviousPage
      __typename
    }
    count
    __typename
  }
}

fragment EndUserFragment on User {
  id
  name
  email
  phone
  notes
}`

func SearchEndUser(inputPhone string, formatCPF bool) ([]FormattedCPF, error) {
	err := godotenv.Load()

	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	zendeskapiurl := os.Getenv("ZENDESK_API_URL")

	if zendeskapiurl == "" {
		panic("ZENDESK_API_URL is not set")
	}

	data := QueryVariables{
		Limit:      25,
		PageNumber: 1,
		Query:      fmt.Sprintf("phone:%s", inputPhone),
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": os.Getenv("ZENDESK_BASIC_AUTH"),
	}

	body := RequestBody{
		Query:     searchEndUserQuery,
		Variables: data,
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	req, err := http.NewRequest("POST", zendeskapiurl, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("non-OK HTTP status: %s, response body: %s", resp.Status, string(bodyBytes))
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result SearchEndUserType
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("graphql errors: %s", result.Errors)
	}

	var notes []string
	for _, edge := range result.Data.Search.Edges {
		notes = append(notes, edge.Node.Notes)
	}

	var normalizedTexts []FormattedCPF

	for _, note := range notes {
		normalizedCPFs, err := normalizecpf.SendRequest(note, formatCPF)
		if err != nil {
			return nil, fmt.Errorf("failed to normalize CPF: %w", err)
		}

		for _, cpf := range normalizedCPFs {
			normalizedTexts = append(normalizedTexts, FormattedCPF{CPF: cpf})
		}
	}

	return normalizedTexts, nil
}
