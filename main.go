package main

import (
	"log"
	"simplewall-utils/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
