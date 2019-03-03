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
	name		string
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
			geocodeObj:=geocode{id, brewery_id, lat, lon, accuracy, beerSlice, ""}
			geocodeSlice = append(geocodeSlice,geocodeObj)
		}
	}
	return geocodeSlice
}

func assignNames(in []geocode) []geocode{
	S_breweries	:= getStrings("dumps/breweries.csv")
	geocodeSlice	:= in
	for i:= range S_breweries{
		if i!=0{
			brew_id, err := strconv.Atoi(S_breweries[i][0])
			if err!=nil{
				fmt.Println(err)
			}
			for j:=range geocodeSlice{
				if geocodeSlice[j].brewery_id == brew_id{
					geocodeSlice[j].name = S_breweries[i][1]
				}
			}
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

func toDegree(rad float64) float64{
	return rad * 180 / math.Pi
}
//Pritaikau Inverse Haversine formule
func calcDist(lat1, lon1, lat2, lon2 float64) float64{
	var dist float64
	lon1R	:= toRadians(lon1)
	lon2R	:= toRadians(lon2)
	lat1R	:= toRadians(lat1)
	lat2R	:= toRadians(lat2)
	a	:= math.Sin((lat2R-lat1R)/2)*math.Sin((lat2R-lat1R)/2) + math.Cos(lat1R)*math.Cos(lat2R) * math.Sin((lon2R-lon1R)/2)*math.Sin((lon2R-lon1R)/2)
	c	:= 2 *  math.Atan2(math.Sqrt(a),math.Sqrt(1-a))
	R	:= 6371.0 //zemes spindulys kilometrais
	dist	= R*c
	return dist
}

func findNeighbourhood(lat,lon float64, boundSlice []geocode) []geocode{
	var d float64

	var neighbourhood []geocode
	d = 2000.0/6371.0

	lonR	:= toRadians(lon)
	latR	:= toRadians(lat)

	latMaxR :=math.Asin(math.Sin(latR)*math.Cos(d) + math.Cos(latR)*math.Sin(d)*math.Cos(0))
	latMinR :=math.Asin(math.Sin(latR)*math.Cos(d) + math.Cos(latR)*math.Sin(d)*math.Cos(math.Pi))
	lonMaxR :=lonR + math.Atan2(math.Sin(math.Pi/2)*math.Sin(d)*math.Cos(latR), math.Cos(d)-(math.Sin(latR)*math.Sin(latR)))
	lonMinR :=lonR + math.Atan2(math.Sin(math.Pi*3/2)*math.Sin(d)*math.Cos(latR), math.Cos(d)-(math.Sin(latR)*math.Sin(latR)))

	latMin	:= toDegree(latMinR)
	latMax	:= toDegree(latMaxR)
	lonMin	:= toDegree(lonMinR)
	lonMax	:= toDegree(lonMaxR)

	for i:=range boundSlice{
		breweryLat := boundSlice[i].lat
		breweryLon := boundSlice[i].lon
		if breweryLat > latMin && breweryLat < latMax && breweryLon > lonMin && breweryLon < lonMax{
			neighbourhood = append(neighbourhood, boundSlice[i])
		}
	}
	return neighbourhood
}

func getBeerCount(loc geocode)int{
	return len(loc.beers)
}

func greedyAlg(lat, lon float64, boundSlice []geocode) []geocode {
	neighbourhood := findNeighbourhood(lat,lon, boundSlice)
	fuel := 2000.0
	var optimalPath []geocode
	var beerCount int
	home :=geocode{-1, -1, lat, lon, "", nil, "Home"}
	optimalPath, beerCount = recursiveSoln(home, home, neighbourhood, optimalPath, fuel,0)
	fmt.Println(beerCount)
	return optimalPath
}
//SF (beercountSF, optimalpathSF) reiskia So Far, kad atskirti variables
func recursiveSoln(currentLoc, home geocode, neighbourhood, optimalPathSF []geocode, fuel float64, beercountSF int) ([]geocode, int){
	if calcDist(home.lat, home.lon, currentLoc.lat, currentLoc.lon) < fuel || len(neighbourhood)==0{
		return optimalPathSF, beercountSF //base case
	}

	var optimalPath []geocode
	beerCount := beercountSF

	for i:=range neighbourhood{
		optimalPathI, beerCountI := recursiveSoln(neighbourhood[i], home, append(neighbourhood[:i], neighbourhood[i+1:]...), append(optimalPathSF, neighbourhood[i]), fuel-calcDist(currentLoc.lat, currentLoc.lon, neighbourhood[i].lat, neighbourhood[i].lon), getBeerCount(neighbourhood[i]))
		if beerCountI > beerCount{
			optimalPath = optimalPathI
			beerCount = beerCountI
		}
	}
	return optimalPath, beerCount
}

func main(){
	lat,err := strconv.ParseFloat(os.Args[1],32)
	lon,err := strconv.ParseFloat(os.Args[2],32)
	fmt.Println(lat)
	fmt.Println(lon)
	//parsinu beers.csv
	beerSlice := parseBeers()
	//Parsinu geocodes.csv
	geocodeSlice := parseGeocodes()

	geocodeSlice = assignNames(geocodeSlice)
	boundSlice := bindBeersToBreweries(geocodeSlice,beerSlice)
	solution := greedyAlg(lat, lon, boundSlice)
	fmt.Println(solution)
	//greedyAlg(lon, lat, boundSlice)
	if err!= nil{
		fmt.Println(err)
	}
}
