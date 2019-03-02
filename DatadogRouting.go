//
//DatadogRouting.go
//
//Vytenis pirma karta naudoja Go, nes panasi i C, tad sansas kazka ismokti
//Vytenis Sliogeris
//Saltiniai sintaksei:
//https://www.tutorialspoint.com/go/
//https://gobyexample.com/command-line-arguments
//https://golang.org/
//Problema panasi i knapsack problema, tad kazka bandysim

package main

import	(
	"os"
	"fmt"
	"io/ioutil"
	"encoding/csv"
	"bytes"
	"strconv"
	"math"
)

//Struktura laikyti viena lokacijos entry
type geocode struct{
	id		int
	brewery_id	int
	lat		float64
	lon		float64
	accuracy	string
	beers		[]beer
}
//Struktura alaus databazei
type beer struct{
	id		int
	brewery_id	int
	name		string
	cat_id		int
	style_id	int
}
//Dummy code algoritmui

func getStrings(location string) [][]string{
	B_content, err := ioutil.ReadFile(location)
	reader := csv.NewReader(bytes.NewReader(B_content))
	strings, err := reader.ReadAll()
	if err!=nil{
		fmt.Println(err)
	}
	return strings
}

func parseBeers() []beer {
	S_beer := getStrings("dumps/beers.csv")
	var beerSlice []beer
	for i:=range S_beer{
		if i!=0{
			id, err		:= strconv.Atoi(S_beer[i][0])
			brewery_id, err := strconv.Atoi(S_beer[i][1])
			name		:= S_beer[i][2]
			cat_id, err	:= strconv.Atoi(S_beer[i][3])
			style_id, err	:= strconv.Atoi(S_beer[i][4])
			if err!=nil{
				fmt.Println(err)
			}
			beerObj:=beer{id, brewery_id, name, cat_id, style_id}
			beerSlice = append(beerSlice,beerObj)
		}
	}
	return beerSlice
}

func parseGeocodes() []geocode {
	S_geocodes := getStrings("dumps/geocodes.csv")
	var geocodeSlice []geocode
	for i:=range S_geocodes{
		if i!=0{
			id, err		:= strconv.Atoi(S_geocodes[i][0])
			brewery_id, err := strconv.Atoi(S_geocodes[i][1])
			lat, err	:= strconv.ParseFloat(S_geocodes[i][2], 32)
			lon, err	:= strconv.ParseFloat(S_geocodes[i][3], 32)
			accuracy	:= S_geocodes[i][4]
			var beerSlice	[]beer
			if err!=nil{
				fmt.Println(err)
			}
			geocodeObj:=geocode{id, brewery_id, lat, lon, accuracy, beerSlice}
			geocodeSlice = append(geocodeSlice,geocodeObj)
		}
	}
	return geocodeSlice
}

func bindBeersToBreweries(geocodeSlice []geocode, beerSlice []beer) []geocode{
	outputGeocodeSlice := geocodeSlice
	for i:=range beerSlice{
		for j:=range geocodeSlice{
			if beerSlice[i].brewery_id ==geocodeSlice[j].brewery_id{
				outputBeerSlice := append(geocodeSlice[j].beers,beerSlice[i])
				geocodeSlice[j].beers = outputBeerSlice
			}
		}
	}
	return outputGeocodeSlice
}


func PrintBeers(slice geocode){
	beerSlice := slice.beers
	for i:=range beerSlice{
		fmt.Println(beerSlice[i].name)
	}
}
func toRadians(deg float64) float64{
	return deg*math.Pi/180
}
//Pritaikau Inverse Haversine formule
func calcDist(lon1, lat1, lon2, lat2 float64) float64{
	var dist float64
	lon1R := toRadians(lon1)
	lon2R := toRadians(lon2)
	lat1R := toRadians(lat1)
	lat2R := toRadians(lat2)
	R := 6371.0 //zemes spindulys kilometrais
	dist = 2.0 * R * math.Asin(math.Sqrt(math.Sin((lat2R-lat1R)/2)*math.Sin((lat2R-lat1R)/2) + math.Cos(lat1R)*math.Cos(lat2R)*math.Sin((lon2R-lon1R)/2)*math.Sin((lon2R-lon1R)/2)))
	return dist
}

func greedyAlg(lon, lat float64, boundSlice []geocode) {
	fmt.Println("ur mum gay")
}

func main(){
	lon,err := strconv.ParseFloat(os.Args[1],32)
	lat,err := strconv.ParseFloat(os.Args[2],32)
	fmt.Println(lon)
	//parsinu beers.csv
	beerSlice := parseBeers()
	//Parsinu geocodes.csv
	geocodeSlice := parseGeocodes()

	boundSlice := bindBeersToBreweries(geocodeSlice,beerSlice)
	PrintBeers(boundSlice[2])
	greedyAlg(lon, lat, boundSlice)
	if err!= nil{
		fmt.Println(err)
	}
}
