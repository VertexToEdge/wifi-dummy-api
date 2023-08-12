package main

import (
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

var request_timestamp_log map[string][]int64

type OutBlock struct {
	Msg string `json:"MSG"`
}

type Result struct {
	Location string `json:"location"`
	Client   int    `json:"client"`
}

type WifiResponse struct {
	OutBlock []OutBlock `json:"OUT_BLOCK"`
	Result   []Result   `json:"RESULT"`
}

func constWifiNames() []string {
	return []string{
		"wifi1",
		"wifi2",
		"wifi3",
		"wifi4",
		"wifi5",
		"wifi6",
		"wifi7",
		"wifi8",
		"wifi9",
		"wifi10",
	}
}

func NewWifiResponse() WifiResponse {
	return WifiResponse{OutBlock: []OutBlock{}, Result: []Result{}}
}
func (r *WifiResponse) addOutBlock(msg string) {
	r.OutBlock = append(r.OutBlock, OutBlock{Msg: msg})
}
func (r *WifiResponse) addResult(location string, client int) {
	r.Result = append(r.Result, Result{Location: location, Client: client})
}

var mutex = &sync.Mutex{}

func main() {
	app := fiber.New()
	request_timestamp_log = make(map[string][]int64)

	app.Get("/svc/offcam/pub/WifiAllInfo", func(c *fiber.Ctx) error {
		auth_key := c.Query("AUTH_KEY", "test")
		var limit_seconds int64
		var limit_requests int
		if v, ok := strconv.Atoi(c.Query("limit_seconds", "30")); ok == nil {
			limit_seconds = int64(v)
		} else {
			error_wifi_response := NewWifiResponse()
			error_wifi_response.addOutBlock("limit_seconds is not integer")
			return c.JSON(error_wifi_response)
		}

		if v, ok := strconv.Atoi(c.Query("limit_requests", "5")); ok == nil {
			limit_requests = v
		} else {
			error_wifi_response := NewWifiResponse()
			error_wifi_response.addOutBlock("limit_requests is not integer")
			return c.JSON(error_wifi_response)
		}

		mutex.Lock() // Lock for request_timestamp_log
		// check if auth_key is not seen in request_timestamp_log
		if _, ok := request_timestamp_log[auth_key]; !ok {
			request_timestamp_log[auth_key] = []int64{time.Now().Unix()}
		}

		// trim request_timestamp_log that out of limit_seconds duration
		for _, timestamp := range request_timestamp_log[auth_key] {
			if timestamp < time.Now().Unix()-limit_seconds {
				request_timestamp_log[auth_key] = request_timestamp_log[auth_key][1:]
			}
		}

		// check if request_timestamp_log[auth_key] is out of limit_requests
		if len(request_timestamp_log[auth_key]) > limit_requests {
			error_wifi_response := NewWifiResponse()
			error_wifi_response.addOutBlock("exceed limit_requests per limit_seconds")
			return c.JSON(error_wifi_response)
		}

		// make random response data
		wifi_response := NewWifiResponse()
		for _, wifi_name := range constWifiNames() {
			wifi_response.addOutBlock("success")
			rand_client := rand.Intn(20)
			wifi_response.addResult("2학생회관-"+wifi_name, rand_client)
		}

		// append request_timestamp_log
		request_timestamp_log[auth_key] = append(request_timestamp_log[auth_key], time.Now().Unix())
		mutex.Unlock() // Unlock for request_timestamp_log
		return c.JSON(wifi_response)
	})

	log.Fatal(app.Listen(":8080"))
}
