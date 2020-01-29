package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	//"strconv"
	//"time"
)

/* URLS:
/api/user/122434/ -- get user details (BloomSkyUser)
/api/skydevice/find_me/?name= -- search for device via name
/api/skydevice/find_me/?all=test -- returns all devices???
/api/user/122434/ -- show user details
*/
var bloomskyToken string = ""
var bloomUsername string = ""
var bloomPassword string = ""
var searchTerm string = ""

type BloomSkyDeviceData struct {
	Temperature           float32
	seaLevelPressure_inch string
	FWVersion2            string
	FWVersion1            string
	UV_level              int
	isOffline             string
	Dewpoint              float32
	TS                    int
	Rain                  bool
	night                 bool
	Pressure_inch         float32
	seaLevelPressure      string
	ImageURL              string
	Humidity              int
	Dewpoint_f            float32
	DeviceType            string
	Battery_level         int
	Luminance             int
	UV_level_max          int
	ThirdPartyUV          int
	Temperature_f         float32
	Voltage               int
	UVIndex               int
}

type BloomSkySearchForDeviceName struct {
	Count    int
	Previous string
	Next     string
	Results  []BloomSkyDeviceDetails
}

type BloomSkyDeviceUser struct {
	Username string
	Country  string
	Userid   int
	DeviceID string
	Nickname string
	Id       int
}

type BloomSkyDeviceDetails struct {
	C_or_F              bool
	Distance            int
	AccessoryType       string
	BatteryNotification bool
	CityName            string
	Following_User      []BloomSkyDeviceUser
	VideoList_C         []string
	PreviewImageList    []string
	DeviceName          string
	RainNotification    bool
	FullAddress         string
	RegisterTime        int
	forecast            int
	NumOfFollowers      int
	Data                BloomSkyDeviceData
	UTC                 int
	LON                 float64
	DST                 int
	Owner               string
	is_newzealand       string
	DeviceID            string
	BoundedStorm        string
	Searchable          bool
	LAT                 float64
	ALT                 float32
	StreetName          string
}

type BloomSkyUser struct {
	username               string
	email                  string
	cellphone_UDID         string
	cellphone_OS           string
	app_version            string
	cellphone_token        string
	cellphone_location_lon string
	cellphone_location_lat string
	cellphone_timezone     string
	cellphone_dst          string
	c_or_f                 string
	registration_time      string
	avatar                 string
	country                string
	nickname               string
	login_type             string
	id                     string
	b_or_h                 string
	m_or_km                string
	tl_notify              string
	followed_devices       string
	owned_devices          string
	location               string
	bio                    string
	android_token_44       string
}

type BloomSkySearchData struct {
	Temperature   string
	ImageURL      string
	TS            string
	WindSpeed     string
	Temperature_f string
	WindDirection string
}

type BloomSkySearchDevice struct {
	DeviceName  string
	UTC         string
	DST         string
	LON         string
	DeviceID    string
	LAT         string
	FullAddress string
	Data        BloomSkySearchData
}

func main() {

	if len(os.Args) < 4 {
		fmt.Printf(`bloomsearch <username> <password> <searchterm>
			username - Your bloomsky username
			password - Your bloomsky password
			searchTerm - Bloomsky you are searching for
		`)
		os.Exit(2)

	}

	bloomUsername = os.Args[1]
	bloomPassword = os.Args[2]
	searchTerm = os.Args[3]
	if loginAndGetToken() {
		fmt.Println("Got Auth Token, searching for", searchTerm)
		findDevices(searchTerm)
	} else {
		fmt.Println("[!] Invalid username and password")
	}

}

func findDevices(searchterm string) {

	url := "https://bskybackend.bloomsky.com/api/skydevice/find_me/?name=" + searchterm // this searches for a name
	//url := "https://bskybackend.bloomsky.com/api/skydevice/find_me/?all=test" // this returns all devices? wtf. Needs BloomSkySearchDevice struct
	fmt.Println("Making API request to", url)
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println("[!] Could not make API request", err)
		os.Exit(2)
	}

	req.Header.Set("Authorization", "Token "+bloomskyToken)
	fmt.Println("Token " + bloomskyToken)
	client := &http.Client{}
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)

	var bloomSkySearchResults BloomSkySearchForDeviceName
	err = json.Unmarshal(body, &bloomSkySearchResults)

	if err != nil {
		fmt.Println(err.Error())

	}

	bloomResults := bloomSkySearchResults.Results
	fmt.Printf("%-40s | %-25s | %s \n", "Device Name", "Device ID", "Street Name")
	fmt.Printf("%-40s | %-25s | %s \n", "-----------", "---------", "-----------")
	for k, _ := range bloomResults {
		fmt.Printf("%-40s | %-25s | %s \n", bloomResults[k].DeviceName, bloomResults[k].DeviceID, bloomResults[k].StreetName)
	}

}

func loginAndGetToken() bool {
	message := map[string]interface{}{
		"username": bloomUsername,
		"password": bloomPassword,
	}

	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := http.Post("https://bskybackend.bloomsky.com/auth/login/", "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		log.Fatalln(err)
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	tokenValue, ok := result["auth_token"]
	if ok {
		bloomskyToken = fmt.Sprintf("%v", tokenValue)
		return true
	}

	return false

}
