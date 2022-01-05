package ecovisio

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type TotemInfoDate struct {
	time.Time
}

func (d *TotemInfoDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("02/01/2006", s)
	if err != nil {
		return err
	}
	d.Time = t

	return nil
}

func (d *TotemInfoDate) MarshallJSON() ([]byte, error) {
	return []byte(d.Time.String()), nil
}

type PhotoLink struct {
	Link string `json:"lien"`
}

type Pratique struct {
	Pratique int `json:"pratique"`
	Id       int `json:"id"`
}

type TotemInfo struct {
	Token              *string        `json:"token"`
	IdPdcImg           int            `json:"id_pdc_img"`
	IdPdc              int            `json:"idPdc"`
	Latitude           float32        `json:"lat"`
	Longitude          float32        `json:"lon"`
	Name               string         `json:"nom"`
	PublicMessage      *string        `json:"publicMessage"`
	Photo              []PhotoLink    `json:"photo"`
	PublicLink         *string        `json:"lienPublic"`
	Pratique           []Pratique     `json:"pratique"`
	MainPratique       int            `json:"mainPratique"`
	Filter             *string        `json:"filtre"`
	SiteForm           *string        `json:"formule_site"`
	End                *TotemInfoDate `json:"fin"`
	Start              *TotemInfoDate `json:"debut"`
	EndPeriod          *TotemInfoDate `json:"finPeriode"`
	StartPeriod        *TotemInfoDate `json:"debutPeriode"`
	CurrentYearDefault int            `json:"current_year_default"`
	ExternalURL        string         `json:"externalUrl"`
	OrgName            string         `json:"nomOrganisme"`
	Logo               string         `json:"logo"`
	Country            string         `json:"pays"`
	Sig                int            `json:"sig"`
	PublicPicto        *string        `json:"pictoPublic"`
	Today              *TotemInfoDate `json:"today"`
	Total              int            `json:"total"`
	LastDay            int            `json:"lastDay"`
	DailyAverage       int            `json:"moyD"`
	TotalLY            *int           `json:"totalLY"`
	LastDay8           *int           `json:"lastDay8"`
	DLYAverage         *int           `json:"moyDLY"`
	NumberDays         *int           `json:"nbDays"`
	NumberDaysLY       *int           `json:"nbDaysLY"`
}

func GetTotemList() (TotemInfo []TotemInfo, err error) {
	base, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}

	// Path params
	base.Path += REQ_URI + ID_ORGANISME // harcoded (organisme ID)

	// Query params
	params := url.Values{}
	params.Add("WithNull", "true")
	base.RawQuery = params.Encode()

	log.Printf("Request: %v", base.String())
	resp, err := http.Get(base.String())
	if err != nil {
		return nil, err
	}

	log.Printf("received response: %v", resp.StatusCode)
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&TotemInfo)

	return TotemInfo, err
}

func FindTotems(totemsIn []TotemInfo, name string) ([]TotemInfo, error) {
	var totemsOut []TotemInfo
	re, err := regexp.Compile("(?i)" + name)
	if err != nil {
		return nil, err
	}
	for _, totem := range totemsIn {
		if re.MatchString(totem.Name) {
			totemsOut = append(totemsOut, totem)
		}
	}
	return totemsOut, nil
}
