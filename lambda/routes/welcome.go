package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	access_token := r.URL.Query().Get("access_token")

	fmt.Println("access_token")
	fmt.Println(access_token)

	w.Write([]byte(access_token))
	measurements := getMeasurements(r.Header.Get("access_token"))
	fmt.Println(measurements)

	pingDB()

	w.WriteHeader(http.StatusOK)
}

func getMeasurements(accessToken string) []Measurements {
	req, err := http.NewRequest(http.MethodGet, "https://wbsapi.withings.net/measure?action=getmeas&meastypes=1,6,8,88&category=1&startdate=1709149524&enddate=1716921924", nil)
	if err != nil {
		log.Fatalln("Failed to create req")
	}

	req.Header.Set("authorization", "Bearer "+accessToken)

	c := &http.Client{}
	response, err := c.Do(req)
	if err != nil {
		log.Fatalln("Failed to fetch measurement")
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseMeasurements responseMeasurements

	err = json.Unmarshal(responseData, &responseMeasurements)
	if err != nil {
		fmt.Println("--------", err)
	}

	measurements := parseMeasurements(responseMeasurements)

	return measurements
}

type Measurements struct {
	ID       int `json:"id"`
	Date     int `json:"date"`
	Weight   int `json:"weight"`   // 1
	FatMass  int `json:"fatmass"`  // 6
	BodyFat  int `json:"bodyfat"`  // 8
	BoneMass int `json:"bonemass"` //88
}

func parseMeasurements(oldslice responseMeasurements) []Measurements {

	var measurements []Measurements

	for _, m := range oldslice.Body.Measuregrps {
		date := m.Date
		id := m.Grpid

		var weight int
		var fatmass int
		var BodyFat int
		var BoneMass int
		for _, p := range m.Measures {

			switch p.Type {
			case 1:
				weight = p.Value
			case 6:
				fatmass = p.Value
			case 8:
				BodyFat = p.Value
			case 88:
				BoneMass = p.Value
			}
		}

		measurements = append(measurements, Measurements{ID: id, Date: date, Weight: weight, FatMass: fatmass, BodyFat: BodyFat, BoneMass: BoneMass})
	}

	return measurements
}

type responseMeasurements struct {
	Status int `json:"status"`
	Body   struct {
		Updatetime  int    `json:"updatetime"`
		Timezone    string `json:"timezone"`
		Measuregrps []struct {
			Grpid        int    `json:"grpid"`
			Attrib       int    `json:"attrib"`
			Date         int    `json:"date"`
			Created      int    `json:"created"`
			Modified     int    `json:"modified"`
			Category     int    `json:"category"`
			Deviceid     string `json:"deviceid"`
			HashDeviceid string `json:"hash_deviceid"`
			Measures     []struct {
				Value int `json:"value"`
				Type  int `json:"type"`
				Unit  int `json:"unit"`
				Algo  int `json:"algo,omitempty"`
				Fm    int `json:"fm,omitempty"`
			} `json:"measures"`
			Modelid int    `json:"modelid"`
			Model   string `json:"model"`
			Comment any    `json:"comment"`
		} `json:"measuregrps"`
	} `json:"body"`
}
