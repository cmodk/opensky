package opensky

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/cmodk/go-simplehttp"
	"github.com/sirupsen/logrus"
)

type Opensky struct {
	lg    *logrus.Logger
	sh    simplehttp.SimpleHttp
	debug bool
}

func New(u string, p string, logger *logrus.Logger) *Opensky {

	return &Opensky{
		sh: simplehttp.New(fmt.Sprintf("https://%s:%s@opensky-network.org/api", u, p), logger),
		lg: logger,
	}
}

func (fs *Opensky) SetDebug(d bool) {
	fs.debug = d
	fs.sh.SetDebug(d)
}

type State struct {
	ICAO         string
	Callsign     string
	TimePosition time.Time
	Longitude    float64
	Latitude     float64
	Velocity     float64
	GeoAltitude  float64
	TrueTrack    float64
}

func (os *Opensky) StatesAll(lamin float64, lomin float64, lamax float64, lomax float64) ([]State, error) {

	url := fmt.Sprintf("/states/all?lamin=%f&lomin=%f&lamax=%f&lomax=%f", lamin, lomin, lamax, lomax)

	resp, err := os.sh.Get(url)
	if err != nil {
		return []State{}, err
	}

	var osss []State
	raw := struct {
		States [][]interface{} `json:"states"`
	}{}

	if err := json.Unmarshal([]byte(resp), &raw); err != nil {
		return []State{}, err
	}

	for _, state := range raw.States {
		oss := State{}
		oss.ICAO = strings.Trim(state[0].(string), " ")
		oss.Callsign = strings.Trim(state[1].(string), " ")
		oss.TimePosition = parseUnix(state[3])
		oss.Longitude = parseFloat(state[5])
		oss.Latitude = parseFloat(state[6])
		oss.GeoAltitude = parseFloat(state[7])
		oss.Velocity = parseFloat(state[9])
		oss.TrueTrack = parseFloat(state[10])

		osss = append(osss, oss)

	}

	return osss, nil
}

type Flight struct {
	ICAO                string   `json:"icao24"`
	FirstSeen           unixTime `json:"firstSeen"`
	EstDepartureAirport string   `json:"estDepartureAirport"`
	LastSeen            unixTime `json:"lastSeen"`
	EstArrivalAirport   string   `json:"estArrivalAirport"`
	Callsign            string   `json:"callsign"`
}

func (os *Opensky) FlightGet(icao string) (Flight, error) {
	now := time.Now()
	url := fmt.Sprintf("/flights/aircraft?icao24=%s&begin=%d&end=%d",
		icao,
		now.Add(time.Duration(-1)*time.Hour).Unix(),
		now.Unix())

	resp, err := os.sh.Get(url)
	if err != nil {
		return Flight{}, err
	}

	log.Printf("Resp: %s\n", resp)
	flights := []Flight{}
	if err := json.Unmarshal([]byte(resp), &flights); err != nil {
		return Flight{}, err
	}

	log.Printf("Flights: %v\n", flights)
	if len(flights) == 0 {
		return Flight{}, errors.New("No flights in response")
	}

	flight := flights[0]
	flight.Callsign = strings.Trim(flight.Callsign, " ")
	return flight, nil
}

func parseUnix(i interface{}) (t time.Time) {
	if i == nil {
		return
	}
	unix_time := int64(i.(float64))
	t = time.Unix(unix_time, 0)
	return
}

func parseFloat(i interface{}) (f float64) {
	if i == nil {
		return
	}
	f = i.(float64)
	return
}

func parseInt(i interface{}) (v int) {
	if i == nil {
		return
	}
	v = i.(int)
	return
}
