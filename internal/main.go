package main

import (
	"fmt"
	"insightful-intel/internal/domain"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	fmt.Println("Hello, World!")

	onapi := domain.NewOnapiDomain()
	dgii := domain.NewDgiiDomain()

	resp, err := onapi.SearchComercialName("NOVASCO")
	if err != nil {
		fmt.Println(err)
	}
	spew.Dump(resp)

	if len(resp) > 0 {
		details, err := onapi.GetDetails(resp[0].NumeroExpediente, resp[0].SerieExpediente)
		if err != nil {
			fmt.Println(err)
		}
		spew.Dump(details)
	}

	dgiiResp, err := dgii.GetRegister("NOVASCO")
	if err != nil {
		fmt.Println(err)
	}
	spew.Dump(dgiiResp)
}
