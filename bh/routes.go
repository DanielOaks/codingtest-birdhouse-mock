// Copyright (c) 2023 Daniel Oaks <daniel@danieloaks.net>
// released under the CC0 license

package bh

import (
	"math"
	"net/http"
	"strconv"

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
	if (*s.data)[itemKey] != nil {
		item.Birdhouse = &birdhouseEntry{
			UbidValue:           itemKey,
			Name:                (*s.data)[itemKey].Name,
			Latitude:            (*s.data)[itemKey].Location.Latitude,
			Longitude:           (*s.data)[itemKey].Location.Longitude,
			LastOccupancyUpdate: (*s.data)[itemKey].OccupancyHistory[0].CreatedAt.UTC().Format("2006-01-02T15:04:05.000Z"),
		}
	}

	return item
}

func (s *Server) GetRegistration(c *gin.Context) {
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

	c.JSON(http.StatusOK, s.getReg(ubid))
}

func (s *Server) GetOccupancy(c *gin.Context) {
	ubid := c.Param("ubid")
	page, limit := getPageAndLimit(c)
	totalEntries := 0
	if (*s.data)[ubid] != nil {
		totalEntries = len((*s.data)[ubid].OccupancyHistory)
	}

	// return all data if no limit is set
	if limit == -1 {
		limit = totalEntries
	}

	var entries []occupancyEntry

	baseIndex := (page - 1) * limit
	for i := 0; i < limit; i++ {
		if baseIndex+i >= totalEntries {
			break
		}
		entries = append(entries, occupancyEntry{
			ID:        (*s.data)[ubid].OccupancyHistory[baseIndex+i].ID,
			Eggs:      (*s.data)[ubid].OccupancyHistory[baseIndex+i].Eggs,
			Birds:     (*s.data)[ubid].OccupancyHistory[baseIndex+i].Birds,
			CreatedAt: (*s.data)[ubid].OccupancyHistory[baseIndex+i].CreatedAt.Format("2006-01-02T15:04:05.000Z"),
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
