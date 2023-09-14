// Copyright (c) 2023 Daniel Oaks <daniel@danieloaks.net>
// released under the CC0 license

package bh

import (
	_ "embed"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	perlin "github.com/aquilax/go-perlin"
	uuid "github.com/satori/go.uuid"
)

//go:embed moby-word-lists/names-f.txt
var namesF string

//go:embed moby-word-lists/names-m.txt
var namesM string

func getNameList() (names []string) {
	names = strings.Split(namesF, "\n")
	names = append(names, strings.Split(namesM, "\n")...)

	return names
}

//go:embed locations/World_Cities_Location_table.csv
var locationListCSV string

type Location struct {
	Latitude  float64
	Longitude float64
}

func getLocationList() (locations []Location) {
	r := csv.NewReader(strings.NewReader(locationListCSV))
	r.Comma = ';'

	for {
		record, err := r.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		// the fields are:
		//   id, country, city, latitude, longitude, altitude
		lat, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			log.Fatal(err)
		}

		long, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			log.Fatal(err)
		}

		locations = append(locations, Location{
			Latitude:  lat,
			Longitude: long,
		})

		// fmt.Println(record)
		// fmt.Println(" ", locations[len(locations)-1])
	}

	return locations
}

type OccupancyState struct {
	ID        string
	CreatedAt time.Time
	Eggs      int
	Birds     int
}

type Birdhouse struct {
	Name             string
	Location         Location
	OccupancyHistory []OccupancyState
}

type Data map[string]*Birdhouse

func GenerateData(conf BirdhousesConfig) ([]string, *Data) {
	names := getNameList()
	locations := getLocationList()

	// note: we use a slice here instead of just using an ordered list to
	//  simplify iterating when pagination is involved. it's a bit messy, but
	//  it's quick and should work well
	var dataOrder []string
	var data Data

	for {
		fmt.Println("Generating data")
		dataOrder = nil
		data = make(Data)
		var emptyGenerated int

		for i := 0; i < conf.Registrations; i++ {
			ubid := uuid.NewV4().String()
			dataOrder = append(dataOrder, ubid)

			if rand.Float64() < conf.EmptyRegistrationsPercentage {
				emptyGenerated++
				data[ubid] = nil
				continue
			}

			randomName := fmt.Sprintf("%s's Birdhouse", names[rand.Intn(len(names))])
			randomLocation := locations[rand.Intn(len(locations))]
			var occupancy []OccupancyState

			oneWeek := time.Hour * 24 * 7
			stepTime := time.Second * time.Duration(oneWeek.Seconds()/float64(conf.OccupancyUpdatesPerWeek))
			baseTime := time.Now()

			// generating occupancy values is interesting.
			currentEggs := float64(rand.Intn(50))
			currentEggModifier := float64(rand.Intn(8))
			eggSeed := int64(rand.Float64() * 1000)
			eggP := perlin.NewPerlin(2, 2, 4, eggSeed)

			currentBirds := float64(rand.Intn(30))
			currentBirdModifier := float64(rand.Intn(5))
			birdSeed := int64(rand.Float64() * 1000)
			birdP := perlin.NewPerlin(2, 2, 9, birdSeed)

			for j := 0; j < conf.OccupancyUpdatesPerWeek*conf.StandardOccupancyInWeeks; j++ {
				// random time to adjust the baseTime by for this update
				sleepTime := time.Second * time.Duration(rand.Intn(int(stepTime.Seconds()/2)))
				if rand.Float64() < .5 {
					sleepTime *= -1
				}
				occupancy = append(occupancy, OccupancyState{
					ID:        uuid.NewV4().String(),
					CreatedAt: baseTime.Add(sleepTime),
					Eggs:      int(math.Abs(currentEggs)),
					Birds:     int(math.Abs(currentBirds)),
				})
				currentEggs += eggP.Noise1D(float64(j)*.1) * currentEggModifier
				currentEggModifier += rand.Float64() * 2

				currentBirds += birdP.Noise1D(float64(j)*.15) * currentBirdModifier
				currentBirdModifier += rand.Float64() * 3

				baseTime = baseTime.Add(-stepTime)
			}

			data[ubid] = &Birdhouse{
				Name:             randomName,
				Location:         randomLocation,
				OccupancyHistory: occupancy,
			}
		}

		if conf.EmptyRegistrationsPercentage > 0 && emptyGenerated == 0 {
			fmt.Println("  Didn't generate any empty birdhouses, regenerating")
			continue
		}

		break
	}

	return dataOrder, &data
}
