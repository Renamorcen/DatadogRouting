//
//DatadogRouting.go
//
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
/*
*Location - failas
*Output = 2D array, kuris atspindetu excelyje atidaryta csv faila
*/
func getStrings(location string) [][]string{
	B_content, err := ioutil.ReadFile(location)
	reader := csv.NewReader(bytes.NewReader(B_content))
	strings, err := reader.ReadAll()
	if err!=nil{
		fmt.Println(err)
	}
	return strings
}
/*
*Visos parsing funkcijos panasios, gal butu imanoma padaryti koki mors polymorphism
*/
func parseBeers(location string) []beer {
	S_beer := getStrings(location)
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

func parseGeocodes(location string) []geocode {
	S_geocodes := getStrings(location)
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

func assignNames(in []geocode, location string) []geocode{
	S_breweries	:= getStrings(location)
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
/*
* geocodeSlice: araus daryklu array
* beerSlice : alaus array
* Bandziau daryti kazka panasaus kaip 
* FIND * WHERE geocode.ID == brewery.ID
* 
*/
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
/*
*Matematines funkcijos versti i radianus ir i laipsnius
*/
func toRadians(deg float64) float64{
	return deg*math.Pi/180
}

func toDegree(rad float64) float64{
	return rad * 180 / math.Pi
}
/*
*Pritaikau Inverse Haversine formule
*Sirdingas aciu https://www.movable-type.co.uk/scripts/latlong.html
*/
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
/*
*
*lat - pradine latitude
*lon - pradine longitude(?)
*boundSlice - basically "database"
*
*returns - array of geocodes kurie yra geriausias kelias
*	 - array of beers kuriuos surinko
*
*/
func greedyAlg(lat, lon float64, boundSlice []geocode) ([]geocode,[]beer) {
	fuel := 2000.0
	var optimalPath []geocode
	var beers []beer
	home :=geocode{-1, -1, lat, lon, "", nil, "Home"}
	optimalPath, beers = recursiveSoln(home, home, boundSlice, optimalPath, fuel,beers)
	optimalPath = append(optimalPath, home)
	return optimalPath, beers
}
/*
*Funkcija kuri panasiai printina kaip duotose skaidrese
*/
func printSoln(soln []geocode, beers []beer, lat, lon float64){
	if(len(soln)==0){
		fmt.Println("No nearby breweries found")
	}else{
		var distTravelled float64
		fmt.Printf("Went through %d breweries\n", len(soln))
		distTravelled = distTravelled + calcDist(lat, lon, soln[0].lat, soln[0].lon)
		for i:=range soln{
			if i!=0{
			step:=calcDist(soln[i-1].lat,soln[i-1].lon,soln[i].lat,soln[i].lon)
			distTravelled = distTravelled + step
				fmt.Printf("%s -[%f]->%s\n",soln[i-1].name,step, soln[i].name)
			}
		}
		distTravelled = distTravelled + calcDist(lat, lon, soln[len(soln)-1].lat, soln[len(soln)-1].lon)
		fmt.Printf("Total distnace travelled: %f\n",distTravelled)
		fmt.Printf("Found %d beers\n", len(beers))
		for i:=range beers{
			fmt.Printf("%s\n", beers[i].name)
		}
	}
}
/*
*Paties rasyta delete funkcija, nes kitaip neisimami reikiami elementai. Tad surandu juos pagal ID
*/
func remove(set []geocode, element geocode) []geocode{
	var index int
	for i:=range set{
		if set[i].brewery_id == element.brewery_id{
			index = i
		}
	}
	outputset:=append(set[:index], set[index+1:]...)
	return outputset
}
/*
*
*currentLoc - brewery, kuriame esame
*home - namu lokacija. Galetu buti ir global variable
*neighbourhood - neistyrinetu daryklu kolekcija
*optimalPathSF - kelias iki sitos daryklos
*fuelSF - degalai kuriuos turime sitoje darykloje
*beerSF - alus, kuri surinkome visose daryklose iki sitos
*
*grazina kelia, davusi daugiausia alaus is currentLoc i neighbourhood kaip geocode array
*Bei grazina alu kaip beer array
*
*Algoritmas yra basically brute force, nes neissiaiskinau kitokio budo surasti daugiausia alaus duodanti kelia
*laikas turetu buti daugmaz O(n!), mat su kiekviena aplankyta darykla lieka (n-1) daryklu dar tyrinejimui. Taip gaunasi kad is viso apieskome:
* n(n-1)(n-2)...*3*2
*
*/
func recursiveSoln(currentLoc, home geocode, neighbourhood, optimalPathSF []geocode, fuelSF float64, beersSF []beer) ([]geocode, []beer){
	var bestPath []geocode
	var bestBeer []beer
	if (fuelSF<0) || (fuelSF <calcDist(home.lat, home.lon, currentLoc.lat, currentLoc.lon)){
		return optimalPathSF, beersSF //base case, returnint su kuo atejau nieko nekeites, nes pasiektas dead end
	}
	optimalPathInHere := append(optimalPathSF, currentLoc) //nera dead end, prisegam sita prie optimal path
	beersInHere := beersSF
	for i:=range currentLoc.beers{
		beersInHere = append(beersInHere, currentLoc.beers[i]) //prisegame alus prie optimal path
	}
	for i:=range neighbourhood{//einam per visus pasiekiamus kaimynus ir pritaikome algoritma
		nextNeighbourhood := remove(neighbourhood, currentLoc) //isimame kita darykla is neistyrinetu sarasu
		fuelAfterFlight := fuelSF - calcDist(currentLoc.lat, currentLoc.lon, neighbourhood[i].lat, neighbourhood[i].lon)
		tempPath, tempBeers := recursiveSoln(neighbourhood[i], home, nextNeighbourhood, optimalPathInHere, fuelAfterFlight, beersInHere)
		if len(tempBeers) > len(bestBeer){//laikome atmintyje geriausia kelia, veliau ji returninsime atgal i virsu
			bestBeer = tempBeers
			bestPath = tempPath
		}
	}
	return bestPath, bestBeer
}


/*
* Pirmas arg yra Latitude
* Antras arg yra Longitude
*/
func main(){
	lat,err := strconv.ParseFloat(os.Args[1],32)
	lon,err := strconv.ParseFloat(os.Args[2],32)
	fmt.Printf("%f/%f\n",lat,lon)
	//parsinu beers.csv
	beerSlice := parseBeers("dumps/beers.csv")
	//Parsinu geocodes.csv
	geocodeSlice := parseGeocodes("dumps/geocodes.csv")
	//Duodu vardus visom daryklom
	geocodeSlice = assignNames(geocodeSlice,"dumps/breweries.csv")
	//Prie kiekvienos daryklos prikabinu alus
	boundSlice := bindBeersToBreweries(geocodeSlice,beerSlice)
	//Issprendziam uzduoti
	solution, beers := greedyAlg(lat, lon, boundSlice)
	printSoln(solution,beers,lat,lon)
	if err!= nil{
		fmt.Println(err)
	}
}
