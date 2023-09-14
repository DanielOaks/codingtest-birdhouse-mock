// Copyright (c) 2023 Daniel Oaks <daniel@danieloaks.net>
// released under the CC0 license

package main

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/DanielOaks/codingtest-birdhouse-mock/bh"
	"github.com/docopt/docopt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	usage := `birdhouse-mock
This acts as a fake backend for https://github.com/DanielOaks/codingtest-birdhouse-admin

Usage:
	birdhouse-mock run
	birdhouse-mock -h | --help
	birdhouse-mock --version

Options:
	-h --help          Show this screen.
	--version          Show version.`

	arguments, _ := docopt.ParseArgs(usage, nil, "0.0.1")

	if !arguments["run"].(bool) {
		return
	}

	config := bh.GetConfig()

	fmt.Println("Configuration:")
	fmt.Println(" ", arguments)
	fmt.Println(" ", config)

	fmt.Println("Generating mock registrations and data!")
	dataOrder, data := bh.GenerateData(config.Birdhouses)
	fmt.Println("Data:")
	fmt.Println(" ", data)

	fmt.Println("Starting server!")
	server := bh.NewServer(data, dataOrder)

	// gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "PUT", "PATCH"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// routes
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	router.GET("/registration", server.GetRegistration)
	router.GET("/registration/:ubid", server.GetSingleRegistration)
	router.GET("/house/:ubid/occupancy", server.GetOccupancy)

	address := net.JoinHostPort("0.0.0.0", strconv.Itoa(config.Server.Port))
	router.Run(address)
}
