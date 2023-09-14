// Copyright (c) 2023 Daniel Oaks <daniel@danieloaks.net>
// released under the CC0 license

package bh

import (
	"os"
	"strconv"
)

type BirdhousesConfig struct {
	// how many registrations to generate
	Registrations int

	// how many of those registrations have no birdhouse, e.g. 0.2 == 20%
	EmptyRegistrationsPercentage float64

	// ow many weeks to generate occupancy data for
	StandardOccupancyInWeeks int

	// how many occupancy updates to have per week for most birdhouses
	OccupancyUpdatesPerWeek int

	// how many birdhouses have broken, garbage occupancy data.
	// this includes duplicated entries, entries with the same ID, and more!
	BrokenBirdhousePercentage float64
}

type ServeConfig struct {
	// which port to host this server on
	Port int
}

type Config struct {
	Birdhouses BirdhousesConfig
	Server     ServeConfig
}

func GetConfig() (config Config) {
	// birdhouses
	//
	registrations, err := strconv.Atoi(os.Getenv("BH_REGISTRATIONS"))
	if err != nil {
		registrations = 20
	}
	config.Birdhouses.Registrations = registrations

	emptyRegistrations, err := strconv.ParseFloat(os.Getenv("BH_EMPTY_REGISTRATIONS"), 64)
	if err != nil {
		emptyRegistrations = 0.1
	}
	config.Birdhouses.EmptyRegistrationsPercentage = emptyRegistrations

	weeks, err := strconv.Atoi(os.Getenv("BH_OCCUPANCY_WEEKS"))
	if err != nil {
		weeks = 25
	}
	config.Birdhouses.StandardOccupancyInWeeks = weeks

	updatesPerWeek, err := strconv.Atoi(os.Getenv("BH_UPDATES_PER_WEEK"))
	if err != nil {
		updatesPerWeek = 14
	}
	config.Birdhouses.OccupancyUpdatesPerWeek = updatesPerWeek

	brokenBirdhouses, err := strconv.ParseFloat(os.Getenv("BH_BROKEN_BIRDHOUSES"), 64)
	if err != nil {
		brokenBirdhouses = 0.1
	}
	config.Birdhouses.BrokenBirdhousePercentage = brokenBirdhouses

	// server
	//
	port, err := strconv.Atoi(os.Getenv("BH_SERVE_PORT"))
	if err != nil {
		port = 7000
	}
	config.Server.Port = port

	return config
}
