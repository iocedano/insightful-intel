package main

import (
	"fmt"
	"insightful-intel/internal/domain"
)

func main() {
	fmt.Println("Hello, World!")

	// onapi := domain.NewOnapiDomain()
	dgii := domain.NewDgiiDomain()

	// resp, err := onapi.SearchComercialName("DM B&V")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// spew.Dump(resp)

	// details, err := onapi.GetDetails(resp.Data[0].Id)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// spew.Dump(details)

	dgii.GetRegister("DM")
}
