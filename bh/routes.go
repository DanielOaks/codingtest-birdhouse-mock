// Copyright (c) 2023 Daniel Oaks <daniel@danieloaks.net>
// released under the CC0 license

package bh

import (
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func getPageAndLimit(c *gin.Context) (page, limit int) {
	pageString := c.Query("page")
	page, err := strconv.Atoi(pageString)
	if err != nil {
		page = 1
	}

	limitString := c.Query("limit")
	limit, err = strconv.Atoi(limitString)
	if err != nil {
		limit = -1
	}

	return page, limit
}

type Server struct {
	data  *Data
	order []string
}

func NewServer(data *Data, dataOrder []string) *Server {
	server := Server{
		data:  data,
		order: dataOrder,
	}

	return &server
}

type birdhouseEntry struct {
	UbidValue           string  `json:"ubidValue"`
	Name                string  `json:"name"`
	Latitude            float64 `json:"latitude"`
	Longitude           float64 `json:"longitude"`
	LastOccupancyUpdate string  `json:"lastOccupancyUpdate"`
}

type registrationEntry struct {
	Value     string          `json:"value"`
	Birdhouse *birdhouseEntry `json:"birdhouse,omitempty"`
}

type occupancyEntry struct {
	ID        string `json:"id"`
	Eggs      int    `json:"eggs"`
	Birds     int    `json:"birds"`
	CreatedAt string `json:"created_at"`
}

func (s *Server) getReg(itemKey string) registrationEntry {
	item := registrationEntry{
		Value: itemKey,
	}
	theData := (*s.data)[itemKey]
	if theData != nil {
		item.Birdhouse = &birdhouseEntry{
			UbidValue:           itemKey,
			Name:                theData.Name,
			Latitude:            theData.Location.Latitude,
			Longitude:           theData.Location.Longitude,
			LastOccupancyUpdate: theData.OccupancyHistory[0].CreatedAt.UTC().Format("2006-01-02T15:04:05.000Z"),
		}
	}

	return item
}

func (s *Server) GetRegistrations(c *gin.Context) {
	page, limit := getPageAndLimit(c)
	totalItems := len(*s.data)

	// return all data if no limit is set
	if limit == -1 {
		limit = totalItems
	}

	var items []registrationEntry

	baseIndex := (page - 1) * limit
	for i := 0; i < limit; i++ {
		if baseIndex+i >= totalItems {
			break
		}
		itemKey := s.order[baseIndex+i]

		items = append(items, s.getReg(itemKey))
	}

	c.JSON(http.StatusOK, gin.H{
		"items": items,
		"meta": map[string]int{
			"totalItems":   totalItems,
			"itemCount":    len(items),
			"itemsPerPage": limit,
			"totalPages":   int(math.Ceil(float64(totalItems) / float64(limit))),
			"currentPage":  page,
		},
	})
}

func (s *Server) GetSingleRegistration(c *gin.Context) {
	ubid := c.Param("ubid")

	if (*s.data)[ubid] == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "The given registration was not found",
			"params": map[string]any{
				"value": ubid,
			},
		})
		return
	}

	c.JSON(http.StatusOK, s.getReg(ubid))
}

func (s *Server) GetOccupancy(c *gin.Context) {
	ubid := c.Param("ubid")
	if ubid != c.GetHeader("X-UBID") {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Authentication is required for this endpoint",
		})
		return
	}
	theData := (*s.data)[ubid]
	if theData == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "The given birdhouse was not found",
			"params": map[string]any{
				"ubid": ubid,
			},
		})
		return
	}

	page, limit := getPageAndLimit(c)
	order := strings.ToLower(c.Query("order"))
	if order == "" {
		order = "desc"
	}

	totalEntries := len(theData.OccupancyHistory)

	// return all data if no limit is set
	if limit == -1 {
		limit = totalEntries
	}

	var entries []occupancyEntry

	baseIndex := (page - 1) * limit
	if order == "asc" {
		baseIndex = (totalEntries - 1) - ((page - 1) * limit)
	}

	for i := 0; i < limit; i++ {
		thisIndex := baseIndex + i
		if order == "asc" {
			thisIndex = baseIndex - i
		}

		if thisIndex < 0 || thisIndex >= totalEntries {
			break
		}
		theOccupancy := theData.OccupancyHistory[thisIndex]
		entries = append(entries, occupancyEntry{
			ID:        theOccupancy.ID,
			Eggs:      theOccupancy.Eggs,
			Birds:     theOccupancy.Birds,
			CreatedAt: theOccupancy.CreatedAt.Format("2006-01-02T15:04:05.000Z"),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"items": entries,
		"meta": map[string]int{
			"totalItems":   totalEntries,
			"itemCount":    len(entries),
			"itemsPerPage": limit,
			"totalPages":   int(math.Ceil(float64(totalEntries) / float64(limit))),
			"currentPage":  page,
		},
	})
}
