package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type viaCEP struct {
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

func getCEPFromViaCEP(cep string, ch chan viaCEP) {
	res, err := http.Get(fmt.Sprintf("https://viacep.com.br/ws/%s/json", cep))
	if err != nil {
		slog.Error("Viacep:", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		slog.Error("Viacep:", err)
	}

	var c viaCEP
	err = json.Unmarshal(body, &c)
	if err != nil {
		slog.Error("Viacep:", err)
	}

	ch <- c
}

type brasilAPICEP struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

func getCEPFromBrasilAPI(cep string, ch chan<- brasilAPICEP) {
	res, err := http.Get(fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep))
	if err != nil {
		slog.Error("Brasilapi:", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		slog.Error("Brasilapi:", err)
	}

	var c brasilAPICEP
	err = json.Unmarshal(body, &c)
	if err != nil {
		slog.Error("Brasilapi:", err)
	}

	ch <- c
}

func main() {
	brasilapi, viacep := make(chan brasilAPICEP), make(chan viaCEP)

	zipcode := flag.String("cep", "01001000", "Seach info for a brazilian zipcode")
	flag.Parse()

	go getCEPFromViaCEP(*zipcode, viacep)
	go getCEPFromBrasilAPI(*zipcode, brasilapi)

	select {
	case res := <-viacep:
		fmt.Printf("From viacep.com.br:\n%v\n", res)
	case res := <-brasilapi:
		fmt.Printf("From brasilapi.com.br:\n%v\n", res)
	case <-time.After(time.Second):
		slog.Error("Timeout was exceeded")
	}
}
