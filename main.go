package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type CEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

const URL = "http://viacep.com.br/ws/"

func main() {

	http.HandleFunc("/", GetCepHandler)
	http.ListenAndServe(":8080", nil)
}

func GetCepHandler(writer http.ResponseWriter, reader *http.Request) {
	if reader.URL.Path != "/" {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	cepParam := reader.URL.Query().Get("cep")
	if cepParam == "" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	cep, errorCep := getCep(cepParam)
	if errorCep != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(cep)
}

func getCep(cep string) (*CEP, error) {
	response, errorResponse := http.Get(URL + cep + "/json/")
	if errorResponse != nil {
		return nil, errorResponse
	}
	defer response.Body.Close()

	body, errorBody := ioutil.ReadAll(response.Body)
	if errorBody != nil {
		return nil, errorBody
	}

	var c CEP
	errorC := json.Unmarshal(body, &c)
	if errorC != nil {
		return nil, errorC
	}

	return &c, nil
}
