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
)

type geocode struct{
	id int
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
	entryNumber:=len(geocodes)
	//Todo: implement linkedList, sukurti panasias strukturas kitiem data, tada gal processint pradesiu
	var geocodeArr[entryNumber] geocode
	for i:=range geocodes{
		geocodeArr[i]:=geocode(geocodes[i][0], geocodes[i][1], geocodes[i][2], geocodes[i][3], geocores[i][4])
		fmt.Println(geocodes[i][2])
	}
	if err!=nil{
		fmt.Println(err)
	}
}
