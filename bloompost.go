package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

func postImageFromCam(deviceID string, imagePath string) {

	var POSTData []byte
	fileData, err := ioutil.ReadFile(string(imagePath))
	if err != nil {
		fmt.Println("[!] Could not find imagePath, please check the path and try again")
		os.Exit(2)
	}
	POSTData = fileData

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	url := "https://bskybackend.bloomsky.com/devc/skydevice/" + deviceID
	querystring := fmt.Sprintf(`{"DeviceID":"%s","FWVersion1":"1.4.2","FWVersion2":"1.2.4","HWVersion":"1.0.1","DeviceType":"SKY1","Temperature":20.58,"Humidity":66,"Voltage":2621,"UVIndex":1448,"Luminance":17829742,"Rain":0,"Pressure":1004,"ChargerStatus":1,"TS":%s}`, deviceID, timestamp)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(POSTData))
	if err != nil {
		fmt.Println("[!] Could not make request to bloomsky, returned:\n %v", err)
	}

	q := req.URL.Query()
	q.Add("Info", querystring)
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var result map[string]interface{}
	json.Unmarshal([]byte(body), &result)

	resultVal := result["ResponseValue"].(float64)

	if int(resultVal) == 200 {
		fmt.Println("Updated Bloomsky!")
	} else {
		fmt.Println("Response value is:", resultVal)
		fmt.Printf("Invalid ResponseValue, likely your DeviceID is incorrect, full response: %s", body)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf(`bloompost <DeviceID> <imagePath>
			DeviceID - Device ID of your Bloomsky camera, use bloomsearch to find this or look at the traffic from your device
			imagePath - full path of the image to upload
		`)
		os.Exit(2)

	}
	var bloomDeviceID = os.Args[1]
	var imageToUpload = os.Args[2]

	postImageFromCam(bloomDeviceID, imageToUpload)
}
