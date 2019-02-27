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

//Struktura laikyti viena lokacijos entry
type geocode struct{
	id		int
	brewery_id	int
	lat		float64
	lon		float64
	accuracy	string
}

//Dummy code algoritmui

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
	geocodeSlice :=  make([]geocode, len(geocodes))
	for i:=range geocodes{
		if i!=0{
			id, err		:= strconv.Atoi(geocodes[i][0])
			brewery_id, err := strconv.Atoi(geocodes[i][1])
			lat, err	:= strconv.ParseFloat(geocodes[i][2], 32)
			lon, err	:= strconv.ParseFloat(geocodes[i][3], 32)
			accuracy	:= geocodes[i][4]
			if err!=nil{
				fmt.Println(err)
			}
			geocodeC:=geocode{id, brewery_id, lat, lon, accuracy}
			geocodeSlice[i-1] = geocodeC
		}
	}
	fmt.Println(geocodeSlice[1].lat)
	if err!=nil{
		fmt.Println(err)
	}
}
