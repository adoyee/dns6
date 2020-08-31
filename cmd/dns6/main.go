package main

import (
	"log"

	"github.com/adoyee/dns6"
)

func main() {
	log.Fatal(dns6.ListenAndServe("[::]:53"))
}
