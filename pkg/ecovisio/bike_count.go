package ecovisio

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type BikeCount struct {
	Date  time.Time `json:"date"`
	Count int       `json:"count"`
}

func GetBikesCountFromTotemInfo(info []TotemInfo) (bikeCount [][]BikeCount, err error) {

	for _, totem := range info {
		base, err := url.Parse(URL)
		if err != nil {
			return nil, err
		}
		// Path params
		base.Path += REQ_URI + "data/" + fmt.Sprint(totem.IdPdc)

		// Query params
		params := url.Values{}
		params.Add("idPdc", fmt.Sprint(totem.IdPdc))
		params.Add("fin", time.Now().Format("02/01/2006"))    // end date DD/MM/YYYY
		params.Add("debut", totem.Start.Format("02/01/2006")) // start date DD/MM/YYYY
		params.Add("interval", fmt.Sprint(4))                 // hardcoded, TODO: reverse it
		params.Add("idOrganisme", ID_ORGANISME)               // hardcoded, TODO: reverse it
		var flowIds string
		for _, id := range totem.Pratique {
			if id.Pratique == totem.MainPratique {
				if flowIds != "" {
					flowIds += ";"
				}
				flowIds += fmt.Sprint(id.Id)
			}
		}
		params.Add("flowIds", flowIds)
		base.RawQuery = params.Encode()

		log.Printf("Request: %v", base.String())
		resp, err := http.Get(base.String())
		if err != nil {
			return nil, err
		}

		log.Printf("received response: %v", resp.StatusCode)
		defer resp.Body.Close()
		var countJSONArray [][]string
		var totemCount []BikeCount
		err = json.NewDecoder(resp.Body).Decode(&countJSONArray)
		if err != nil {
			continue
		}
		found := false
		for _, entry := range countJSONArray {
			date, _ := time.Parse("01/02/2006", entry[0])
			count, _ := strconv.Atoi(entry[1])
			if !found && count == 0 {
				continue
			} else {
				found = true
			}
			totemCount = append(totemCount, BikeCount{Date: date, Count: count})
		}
		bikeCount = append(bikeCount, totemCount)
	}

	return bikeCount, err
}
