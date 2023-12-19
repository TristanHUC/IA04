package main

import (
	"fmt"
	"github.com/go-faker/faker/v4"
)

func main() {
	fmt.Println(faker.FirstName() + " " + faker.LastName())
}
