package opensky

import (
	"encoding/json"
	"fmt"
	"log"
)

type Route struct {
	Callsign     string   `json:"callsign"`
	OperatorIata string   `json:"operatorIata"`
	FlightNumber int      `json:"flightNumber"`
	Route        []string `json:"route"`
}

func (os *Opensky) RouteGet(icao string) (Route, error) {

	url := fmt.Sprintf("/routes?callsign=%s", icao)

	resp, err := os.sh.Get(url)
	if err != nil {
		return Route{}, err
	}

	log.Printf("Resp: %s\n", resp)
	r := Route{}
	if err := json.Unmarshal([]byte(resp), &r); err != nil {
		return Route{}, err
	}

	return r, nil

}
