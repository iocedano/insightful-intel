package main

import (
	"fmt"
	"insightful-intel/internal/domain"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	fmt.Println("Hello, World!")

	scj := domain.NewScjDomain()
	onapi := domain.NewOnapiDomain()
	dgii := domain.NewDgiiDomain()
	pgr := domain.NewPgrDomain()

	companyTarget := "NOVASCO"
	var holder string
	var targer string
	var manager string

	resp, err := onapi.SearchComercialName(companyTarget)
	if err != nil {
		fmt.Println(err)
	}
	spew.Dump("-------ONAPI-----\n", resp)

	if len(resp) > 0 {
		details, err := onapi.GetDetails(resp[0].NumeroExpediente, resp[0].SerieExpediente)
		if err != nil {
			fmt.Println(err)
		}
		spew.Dump("-------------\n", details)
		holder = details.Titular
		targer = details.Texto
		manager = details.Gestor

	}

	dgiiResp, err := dgii.GetRegister(companyTarget)
	if err != nil {
		fmt.Println(err)
	}
	spew.Dump("-----DGII--------\n", dgiiResp)

	pgrHolderNews, err := pgr.Search(holder)
	if err != nil {
		fmt.Println(err)
	}
	spew.Dump("------PGR----holder---\n", pgrHolderNews)

	pgrTargerNews, err := pgr.Search(targer)
	if err != nil {
		fmt.Println(err)
	}
	spew.Dump("-----PGR----targer----\n", pgrTargerNews)

	pgrManagerNews, err := pgr.Search(manager)
	if err != nil {
		fmt.Println(err)
	}
	spew.Dump("-----PGR----manager----\n", pgrManagerNews)

	scjCases, err := scj.Search(holder)
	if err != nil {
		fmt.Println(err)
	}
	spew.Dump("-------SCJ---holder---\n", scjCases)

	scjCases, err = scj.Search(targer)
	if err != nil {
		fmt.Println(err)
	}

	spew.Dump("-------SCJ---targer---\n", scjCases)

	scjCases, err = scj.Search(manager)
	if err != nil {
		fmt.Println(err)
	}

	spew.Dump("-------SCJ---manager---\n", scjCases)
}
