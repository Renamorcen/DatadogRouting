//DatadogRouting.go
//Vytenis pirma karta naudoja Go, nes panasi i C
//Vytenis Sliogeris
//Saltiniai sintaksei:
//https://www.tutorialspoint.com/go/
//https://gobyexample.com/command-line-arguments
//https://golang.org/pkg/encoding/csv

package main

import	(
	"os"
	"fmt"
	"io/ioutil"
	"encoding/csv"
	"bytes"
)

func greedyAlg(long, lat int) int{
	return 0
}

func main(){
	lon := os.Args[1]
	lat := os.Args[2]
	B_geocodes, err := ioutil.ReadFile("dumps/geocodes.csv")
	reader := csv.NewReader(bytes.NewReader(B_geocodes))
	geocodes, err := reader.ReadAll()
	fmt.Println(lat+"/"+lon)
	fmt.Printf("%s",geocodes[2][1])
	if err!=nil{
		fmt.Println(err)
	}
}
