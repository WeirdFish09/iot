package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type UpdateRequest struct {
	State int `json:"state"`
	Mode  int `json:"mode"`
}

type InfoResp map[string]map[string]int

func main() {
	url := flag.String("url", "http://localhost:8080", "URL to connect to")
	pin := flag.Int("pin", -1, "PIN")
	mode := flag.Bool("input", true, "is input")
	state := flag.Bool("state", true, "is high")
	info := flag.Bool("info", false, "just get info")
	flag.Parse()
	if *info {
		fullUrl := fmt.Sprintf("%s/info", *url)
		resp, err := http.Get(fullUrl)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(-1)
		}
		defer resp.Body.Close()
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		fmt.Println(bodyString)
		return
	}

	if *pin == -1 {
		fmt.Println("PIN is required")
		os.Exit(-1)
	}

	var req UpdateRequest
	if *state {
		req.State = 1
		fmt.Println("State - high")
	} else {
		fmt.Println("State - low")
		req.State = 0
	}

	if *mode {
		req.Mode = 0
		fmt.Println("Mode - input")
	} else {
		fmt.Println("Mode - output")
		req.Mode = 1
	}
	body, err := json.Marshal(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(-2)
	}
	fullUrl := fmt.Sprintf("%s/pinstate/%d", *url, *pin)
	fmt.Println("Request URL:", fullUrl)
	res, err := http.Post(fullUrl, "application/json", bytes.NewBuffer(body))
	fmt.Println("Response code:%d", res.StatusCode)
	if err != nil {
		fmt.Println(err)
		os.Exit(-3)
	}
	if res.StatusCode != 200 {
		fmt.Println(res.Status)
		os.Exit(-4)
	}
}
