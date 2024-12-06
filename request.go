package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func HandleRequest(args []string) {
	// read config and parse
	rawConfig, err := ParseRawConfig("rc.toml")
	if err != nil {
		log.Fatal(err)
	}
	config, err := NewConfigFromRaw(rawConfig)
	if err != nil {
		log.Fatal(err)
	}

	// parse args
	if len(args) != 1 {
		log.Fatal("Invalid number of arguments. Should be 1: rc <request_name>")
	}
	request_name := args[0]

	// do request
	request_config, ok := config.Requests[request_name]
	if !ok {
		log.Fatalf("Unknown request: %s", request_name)
	}
	resp, err := DoRequest(&request_config)
	if err != nil {
		log.Fatal(err)
	}
	err = HandleResponse(&request_config, resp)
	if err != nil {
		log.Fatal(err)
	}
}

// makes the http request, main function where all logics are
func DoRequest(config *RequestConfig) (*http.Response, error) {
	// create body
	var bodyReader io.Reader
	body, ok := config.Body.GetValue()
	if ok {
		bodyReader = strings.NewReader(string(body))
	} else {
		bodyReader = nil
	}

	// create request
	req, err := http.NewRequest(config.Method.String(), config.Url.String(), bodyReader)
	if err != nil {
		return nil, err
	}
	for _, value := range config.Headers {
		req.Header.Add(value.Key, value.Value)
	}

	// create client
	client := &http.Client{}

	// do request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func HandleResponse(config *RequestConfig, resp *http.Response) error {
	switch config.BodyType {
	case JSON:
		var json_body map[string]any
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(data, &json_body); err != nil {
			return err
		}
		for key, value := range json_body {
			fmt.Printf("%s: %v\n", key, value)
		}
		return nil
	case TEXT:
		return fmt.Errorf("not implemented")
	}
	return fmt.Errorf("invalid body format %s", config.BodyType)
}
