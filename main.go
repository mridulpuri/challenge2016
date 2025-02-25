package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

type Distributor struct {
	Name    string
	Include []string
	Exclude []string
	Parent  *Distributor
}

type City struct {
	City    string
	State   string
	Country string
}

var cityMap map[string]City
var distributors map[string]*Distributor

func normalize(s string) string {
	return strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(s), " ", ""))
}

func loadCities(path string) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	cityMap = make(map[string]City)
	for _, rec := range records[1:] {
		cityKey := normalize(rec[3] + "-" + rec[4] + "-" + rec[5])
		cityMap[cityKey] = City{
			City:    normalize(rec[3]),
			State:   normalize(rec[4]),
			Country: normalize(rec[5]),
		}
	}
}

func loadDistributors(path string) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	distributors = make(map[string]*Distributor)
	for _, rec := range records[1:] {
		name := normalize(rec[0])
		includes := strings.Split(rec[1], "|")
		excludes := strings.Split(rec[2], "|")
		parentName := normalize(rec[3])
		var parent *Distributor
		if parentName != "" {
			parent = distributors[parentName]
		}
		for i := range includes {
			includes[i] = normalize(includes[i])
		}
		for i := range excludes {
			excludes[i] = normalize(excludes[i])
		}
		distributors[name] = &Distributor{
			Name:    name,
			Include: includes,
			Exclude: excludes,
			Parent:  parent,
		}
	}
}

func isAuthorized(dist *Distributor, region string) (bool, string) {
	region = normalize(region)

	if !regionExists(region) {
		return false, "Region does not exist"
	}

	if dist.Parent != nil {
		parentAuthorized, _ := isAuthorized(dist.Parent, region)
		if !parentAuthorized {
			return false, "Denied by parent distributor: " + dist.Parent.Name
		}
	}

	for _, ex := range dist.Exclude {
		if matchesRegion(region, ex) {
			return false, "Explicitly excluded by distributor: " + dist.Name
		}
	}

	for _, inc := range dist.Include {
		if matchesRegion(region, inc) {
			return true, "Explicitly included by distributor: " + dist.Name
		}
	}

	return false, "No matching inclusion found"
}


func matchesRegion(query, perm string) bool {
	qParts := strings.Split(query, "-")
	pParts := strings.Split(perm, "-")

	if len(qParts) < len(pParts) {
		return false
	}

	for i := 1; i <= len(pParts); i++ {
		if qParts[len(qParts)-i] != pParts[len(pParts)-i] {
			return false
		}
	}
	return true
}

func regionExists(region string) bool {
	_, cityOk := cityMap[region]
	if cityOk {
		return true
	}
	for city := range cityMap {
		if strings.HasSuffix(city, region) {
			return true
		}
	}
	return false
}

func main() {
	loadCities("cities.csv")
	loadDistributors("distributors.csv")

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter distributor name: ")
	scanner.Scan()
	distName := normalize(scanner.Text())

	dist, exists := distributors[distName]
	if !exists {
		fmt.Printf("Distributor %s not found\n", distName)
		return
	}

	fmt.Print("Enter region (format CITY-STATE-COUNTRY): ")
	scanner.Scan()
	region := normalize(scanner.Text())

	authorized, reason := isAuthorized(dist, region)
	if authorized {
		fmt.Printf("%s is authorized to distribute in %s. Reason: %s\n", dist.Name, region, reason)
	} else {
		fmt.Printf("%s is NOT authorized to distribute in %s. Reason: %s\n", dist.Name, region, reason)
	}
}