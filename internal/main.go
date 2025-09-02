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

	resp, err := onapi.SearchComercialName("DM B&V")
	if err != nil {
		fmt.Println(err)
	}
	spew.Dump(resp)

	details, err := onapi.GetDetails(resp.Data[0].Id)
	if err != nil {
		fmt.Println(err)
	}
	spew.Dump(details)

	dgiiResp, err := dgii.GetRegister("DM")
	if err != nil {
		fmt.Println(err)
	}
	spew.Dump(dgiiResp)
}
