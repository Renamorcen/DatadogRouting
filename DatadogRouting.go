//DatadogRouting.go
//Vytenis pirma karta naudoja Go, nes panasi i C
//Vytenis Sliogeris
//Saltiniai sintaksei:
//https://www.tutorialspoint.com/go/
//https://gobyexample.com/command-line-arguments
//https://golang.org/

package main

import	(
	"os"
	"fmt"
	"io/ioutil"
	"encoding/csv"
	"bytes"
	"strconv"
)

//linkedlist 
type geocode struct{
	id	int
	brewery_id int
	lat float32
	lon float32
	accuracy string
}


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
	geocodeSlice :=  make([]geocode, len(geocodes))
	for i:=range geocodes{
		id := strconv
		geocodeC:=geocode{strconv.Atoi(geocodes[i][0]), geocodes[i][1], geocodes[i][2], geocodes[i][3], geocodes[i][4]}
		//fmt.Println(geocodes[i][4])
	}
	if err!=nil{
		fmt.Println(err)
	}
}
