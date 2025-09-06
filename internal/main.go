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
	pgr := domain.NewPgrDomain()

	resp, err := onapi.SearchComercialName("NOVASCO")
	if err != nil {
		fmt.Println(err)
	}

	var holder string
	var targer string

	if len(resp) > 0 {
		details, err := onapi.GetDetails(resp[0].NumeroExpediente, resp[0].SerieExpediente)
		if err != nil {
			fmt.Println(err)
		}
		spew.Dump(details)
		holder = details.Titular
		targer = details.Texto

	}

	dgiiResp, err := dgii.GetRegister("NOVASCO")
	if err != nil {
		fmt.Println(err)
	}
	spew.Dump(dgiiResp)

	pgrHolderNews, err := pgr.Search(holder)
	if err != nil {
		fmt.Println(err)
	}
	spew.Dump(pgrHolderNews)

	pgrTargerNews, err := pgr.Search(targer)
	if err != nil {
		fmt.Println(err)
	}
	spew.Dump(pgrTargerNews)

}
