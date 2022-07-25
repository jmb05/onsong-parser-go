package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/jmb05/Onsong-Parser-go/onsong"
)

func main() {
	content, err := ioutil.ReadFile("test-song.onsong")

	if err != nil {
		log.Fatal(err)
	}

	para := onsong.SplitParagraphs(string(content))
	for _, p := range para {
		for _, l := range p {
			fmt.Println(l)
		}
		fmt.Println("--------------------------------")
	}
}
