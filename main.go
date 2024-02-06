package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type BrasilAPI struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

type ViaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Pass the CEP as arg")
		fmt.Print("e.g: .\\main.exe 01153000")
		return
	}
	cep := os.Args[1]

	channelBrasilAPI := make(chan BrasilAPI)
	channelViaCep := make(chan ViaCEP)

	go GetBrasilAPI(cep, channelBrasilAPI)
	go GetViaCEP(cep, channelViaCep)

	select {
	case brasilCep := <-channelBrasilAPI:
		fmt.Printf("From BrasilAPI: %+v\n", brasilCep)

	case viaCep := <-channelViaCep:
		fmt.Printf("From ViaCEP: %+v\n", viaCep)

	case <-time.After(time.Second):
		fmt.Printf("TimeOut")
	}

}

func GetBrasilAPI(cep string, chBrasil chan BrasilAPI) {
	var address BrasilAPI
	brasilAPIURL := "https://brasilapi.com.br/api/cep/v1/" + cep

	result, err := requestCEP(brasilAPIURL)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(result, &address)
	if err != nil {
		panic(err)
	}

	chBrasil <- address
}

func GetViaCEP(cep string, chViaCep chan ViaCEP) {
	var address ViaCEP
	viaCEPURL := "http://viacep.com.br/ws/" + cep + "/json/"

	result, err := requestCEP(viaCEPURL)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(result, &address)
	if err != nil {
		panic(err)
	}

	chViaCep <- address
}

func requestCEP(APIURL string) ([]byte, error) {
	request, err := http.Get(APIURL)
	if err != nil {
		return nil, err
	}

	result, err := io.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}

	return result, nil
}
