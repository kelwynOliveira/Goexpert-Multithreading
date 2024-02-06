# Multithreading challenge

Goexpert postgraduation project

## Challenge Description

> In this challenge you'll have to use what we've learnt about multithreading and APIs to find the fastest result between two different APIs.
>
> The two requests will be made simultaneously to the following APIs:
>
> - https://brasilapi.com.br/api/cep/v1/ + cep
> - http://viacep.com.br/ws/ + cep + /json/
>
> The requirements for this challenge are:
>
> - Accept the API that delivers the fastest response and discard the slowest response.
> - The result of the request should be displayed in the command line with the address data, as well as which API sent it.
> - Limit the response time to 1 second. Otherwise, the timeout error should be displayed.

## How to execute

run the program with cep as parameter

e.g: `go run main.go 01153000`

e.g.: `.\main.exe 01153000`

## Challenge Solution Development

Since the project is not developed on the src go folder a module is needed.

To do this on cmd where the project is run `go mod init <module name>`. The parameter `<module name>` can be any string you like, it will act as an unique identifier to the module, but as good practice I will use the URL path where the code will be hosted.

e.g.: `go mod init github.com/kelwynOliveira/Goexpert-Multithreading`

### Verify the cep argument

To take the cep as argument the `os` package must be imported.

The if statement verify if the there is a second argument to use as cep.

```go
  if len(os.Args) != 2 {
		fmt.Println("Pass the CEP as arg")
		fmt.Print("e.g: .\\main.exe 01153000")
		return
	}
	cep := os.Args[1]
```

### Structs to receive the json from API

> used https://mholt.github.io/json-to-go/ to convert the JSON to Go struct

- https://brasilapi.com.br/api/cep/v1/ + cep

```go
type BrasilAPI struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}
```

- http://viacep.com.br/ws/ + cep + /json/

```go
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
```

### The main function

Since the program makes a request from two APIs and accepts the API that delivers the fastest response we will use two channels, start two go routines and use two "get functions".

The `select` statement will choose the valid option to be printed.

```go
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
```

### functions to get the response from APIs

Here we will have two "get functions" one for each API request

They receive `cep` and the `channel` as parameters.

```go
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
```

### The request function

The request function receives the API URL and returns a slice of bytes as result.

```go
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
```
