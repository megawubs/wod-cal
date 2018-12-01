package wod

import (
	"encoding/json"
	"fmt"
	"github.com/megawubs/calendar"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type WOD struct {
	Id        int64  `json:"id_appointment"`
	DateStart string `json:"date_start"`
	DateEnd   string `json:"date_end"`
	Name      string `json:"name"`
}

type WODS []WOD

type result struct {
	ResultSet WODS
}

type Gym struct {
	Id   int64  `json:"id_gym"`
	Name string `json:"name"`
	City string `json:"city"`
}

type Gyms []Gym

type gymResult struct {
	ResultSet Gyms
}

func All(apiKey string, date time.Time) (WODS, error) {

	gyms, err := fetchGyms(apiKey)
	if err != nil {
		return nil, fmt.Errorf("could not fetch gyms: %s", err)
	}
	wods := make(WODS, 0)
	for _, gym := range gyms {

		data := url.Values{
			"data[service]":                                        {"agenda"},
			"data[method]":                                         {"ownAppointmentsUser"},
			"data[idc]":                                            {"2520"},
			"data[reset]":                                          {"1"},
			"data[date]":                                           {date.Format("2006-01-02")},
			"data[app]":                                            {"wodapp"},
			"data[language]":                                       {"nl_NL"},
			"data[version]":                                        {"10.0"},
			"data[clientUserAgent]":                                {"browser"},
			"data[token]":                                          {apiKey},
			"data[customApp]":                                      {"0"},
			"data[id_gym]":                                         {strconv.FormatInt(gym.Id, 10)},
			"data[companySyncVersions][accesslevels][version]":     {"12"},
			"data[companySyncVersions][accesslevels][version_pnp]": {"26"},
			"data[companySyncVersions][styles][version]":           {"12"},
			"data[companySyncVersions][styles][version_pnp]":       {"1"},
			"data[id_appuser_li]":                                  {"129038"},
		}

		b, err := request(data)
		if err != nil {
			return nil, fmt.Errorf("could not make request: %s", err)
		}

		var all result

		if err := json.Unmarshal(b, &all); err != nil {
			return nil, fmt.Errorf("go-wod/wod: reading response failed %s", err)
		}

		wods = append(wods, all.ResultSet...)
	}
	return wods, nil
}

func (workouts WODS) MarshallICalendar(c *calendar.Calendar, location *time.Location) error {
	for _, wod := range workouts {
		layout := "02-01-2006 15:04"

		start, err := time.ParseInLocation(layout, wod.DateStart, location)

		if err != nil {
			return fmt.Errorf("unable to parse start time %s", err)
		}

		end, err := time.ParseInLocation(layout, wod.DateEnd, location)
		if err != nil {
			return fmt.Errorf("unable to parse end time %s", err)
		}

		if end.Weekday() == time.Saturday {
			end = end.Add(time.Duration(30) * time.Minute)
		}

		if wod.Name == "WOD" {
			end = end.Add(30 * time.Minute)
		}

		event := calendar.NewEvent(wod.Id, "", start, end, "Crossfit - "+wod.Name)

		c.Add(event)
	}
	return nil
}

func fetchGyms(apiKey string) (Gyms, error) {
	params := url.Values{
		"data[service]":         {"user"},
		"data[method]":          {"gymsUser"},
		"data[app]":             {"wodapp"},
		"data[language]":        {"nl_NL"},
		"data[version]":         {"13.0"},
		"data[clientUserAgent]": {"browser"},
		"data[token]":           {apiKey},
		"data[customApp]":       {"0"},
		"data[id_appuser_li]":   {"129038"},
	}

	b, err := request(params)
	if err != nil {
		return nil, fmt.Errorf("could not fetch gyms: %s", err)
	}

	var gyms gymResult
	if err := json.Unmarshal(b, &gyms); err != nil {
		return nil, fmt.Errorf("go-wod/wod: reading response failed %s", err)
	}

	return gyms.ResultSet, nil
}

func request(values url.Values) ([]byte, error) {
	client := &http.Client{}

	reader := strings.NewReader(values.Encode())

	req, err := http.NewRequest("POST", "https://ws.personaltrainerapp.nl", reader)
	if err != nil {
		return nil, fmt.Errorf("go-wod/wod: creating request failed %s", err)
	}

	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("go-wod/wod: POST request failed %s", err)
	}

	return ioutil.ReadAll(res.Body)
}
