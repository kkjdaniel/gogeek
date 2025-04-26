package main

import (
	"fmt"

	"github.com/kkjdaniel/gogeek/collection"
)

func main() {
	response, err := collection.Query("kkjdaniel", collection.WithOwned(true))
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("Response:", response)
}
