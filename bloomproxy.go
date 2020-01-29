package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/nlopes/slack"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var sendSlack = 1
var gifFrameNum int = 1
var slackToken = "SLACKTOKENGOESHERE"
var api = slack.New(slackToken, slack.OptionDebug(false))

func SpecificNameServer(ctx context.Context, network, address string) (net.Conn, error) {
	d := net.Dialer{}
	return d.DialContext(ctx, "udp", "1.1.1.1:53")
}

func MakeHTTPRequest(requestType string, requestHost string, requestHeaders map[string][]string, requestDestURL string, POSTData []byte) string {

	DNSResolver := net.Resolver{
		PreferGo: true,
		Dial:     SpecificNameServer,
	}

	ctx := context.Background()
	ipaddr, err := DNSResolver.LookupIPAddr(ctx, requestHost)
	if err != nil {
		panic(err)
	}

	gifFramePath := "nyan/" + strconv.Itoa(gifFrameNum) + ".jpg"
	slackSendMessage("Attempting to send a different image")

	gifData, err := ioutil.ReadFile(string(gifFramePath))
	if err == nil {
		POSTData = gifData
		slackUploadFile(gifFramePath)
		gifFrameNum++
		if gifFrameNum == 12 {
			gifFrameNum = 1
		}

	} else {
		fmt.Println(err.Error())
		slackSendMessage("Failed to read gif frame :(")
	}

	req, err := http.NewRequest(requestType, "http://"+ipaddr[0].String()+requestDestURL, bytes.NewBuffer(POSTData))

	for k, v := range requestHeaders {
		if k != "Content-Length" {
			req.Header.Set(string(k), string(v[0]))
		}
	}
	req.Host = requestDestURL // cant set this in the headers? yeah, I dunno man.

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("```PANIC!!! \n" + err.Error() + "```")
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	return string(body)
}
func slackSendMessage(message string) {
	if sendSlack == 0 {
		return
	}
	api.PostMessage("#bloomsky", slack.MsgOptionText(message, false))

}
func slackUploadFile(fileName string) {
	if sendSlack == 0 {
		return
	}
	params := slack.FileUploadParameters{
		Title:    "Image Test",
		Filetype: "jpg",
		File:     fileName,
		Channels: []string{"#bloomsky"},
	}
	file, err := api.UploadFile(params)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("Name: %s, URL: %s\n", file.Name, file.URL)

}
func main() {

	fmt.Printf("Webserver Started\n")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		var requestHeaders = r.Header

		var requestHost = r.Host
		var requestDestURL = r.URL.String()
		var requestType = r.Method
		requestBody, err := ioutil.ReadAll(r.Body)
		path := r.URL.Path[1:]
		if strings.Contains(path, "jpg") {
			fmt.Println("Serving JPG! " + path)
			slackUploadFile(string(path))
			data, err := ioutil.ReadFile(string(path))
			if err == nil {
				w.Write(data)
			} else {
				w.WriteHeader(404)
				w.Write([]byte("404 Something went wrong - " + http.StatusText(404)))
			}
			return
		}

		if err != nil {

		}

		stringBody := string(requestBody)

		slackSendMessage(fmt.Sprintf("New Incoming HTTP Request of type `%s` to `%s` looking for `%s` ", requestType, requestHost, requestDestURL))
		t := time.Now()
		currTime := t.Format("2006-01-02 15:04:05")

		fmt.Printf("%s -- Incoming Request to %s \n", currTime, requestDestURL)
		var headersString = "*Headers*\n```"
		for k, v := range requestHeaders {
			headersString += k + " : " + v[0] + "\n"
		}
		headersString += "```"
		slackSendMessage(headersString)

		if strings.Contains(stringBody, "JFIF") {

			//Lets save the image
			fileName := t.Format("200601021504") + ".jpg"
			fileData := requestBody
			ioutil.WriteFile(fileName, fileData, 0644)
			slackSendMessage("*Image being uploaded:*")
			slackUploadFile(string(fileName))
			os.Remove(fileName)

		} else {
			slackSendMessage("*Non Image Request Data:* ```" + stringBody + "```")
		}

		var realRequestResponse = MakeHTTPRequest(requestType, requestHost, requestHeaders, requestDestURL, requestBody)
		slackSendMessage("Request Response: ```" + realRequestResponse + "```")
		slackSendMessage("────────────────────────────────────────────────────────────────────")
	})

	http.ListenAndServe(":80", nil)
}
