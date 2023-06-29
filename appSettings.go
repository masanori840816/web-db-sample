package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type AppSettings struct {
	BaseURL          string `json:"baseUrl"`
	Host             string `json:"host"`
	Port             int    `json:"port"`
	ConnectionString string `json:"connectionString"`
	JWSMessage       string `json:"jwsMessage"`
}

func NewAppSettings() (settings AppSettings, err error) {
	result := &AppSettings{}
	cur, _ := os.Getwd()
	file, err := os.Open(fmt.Sprintf("%s/appSettings.json", cur))
	if err != nil {
		return *result, err
	}
	fileinfo, err := file.Stat()
	if err != nil {
		return *result, err
	}
	// Read file
	fileData := make([]byte, fileinfo.Size())
	_, err = file.Read(fileData)
	if err != nil {
		return *result, err
	}
	err = json.Unmarshal(fileData, &result)
	if err != nil {
		return AppSettings{}, err
	}
	return *result, nil
}
