package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
)

func genString(stringLength int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	b := make([]rune, stringLength)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func genFCMToken(installID string) string {
	return fmt.Sprintf("%s:APA91b%s", installID, genString(134))
}

func generateRequestData(referrer string) map[string]interface{} {
	installID := genString(11)
	return map[string]interface{}{
		"key":          "39a1b6f2bed0443b85f46b0daebe6f4a",
		"install_id":   installID,
		"referrer":     referrer,
		"warp_enabled": false,
		"tos":          nil,
		"type":         "Android",
		"locale":       "en_US",
		"fcm_token":    genFCMToken(installID),
	}
}

func makeRequest(referrer string) int {
	// Make a single request.
	jsonValue, _ := json.Marshal(generateRequestData(referrer))
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI("https://api.cloudflareclient.com/v0a745/reg")
	req.Header.SetContentType("application/json")
	req.SetBody(jsonValue)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	err := fasthttp.Do(req, resp)
	if err != nil {
		return -1
	}
	return resp.StatusCode()
}

func main() {
	rand.Seed(time.Now().UnixNano())
	var referrer string
	var numRequests int
	fmt.Print("Enter the referrer id: ")
	fmt.Scanln(&referrer)
	fmt.Print("Number of GB(s): ")
	fmt.Scanln(&numRequests)

	successes := 0
	failures := 0
	for i := 0; i < numRequests; i++ {
		fmt.Printf("\rgenerating GBs: [%s] %d%% (%d/%d)", strings.Repeat("=", int(float64(i+1)/float64(numRequests)*20)), int(float64(i+1)/float64(numRequests)*100), i+1, numRequests)
		time.Sleep(20 * time.Second)
		statusCode := makeRequest(referrer)
		if statusCode == 200 {
			successes++
		} else {
			failures++
			fmt.Printf("error on request %d/%d", i+1, numRequests)
		}
	}
	fmt.Printf("\nFinished with %d successes and %d failures.", successes, failures)
}
