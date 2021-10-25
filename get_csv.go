package main

import (
	"net/http"
	"encoding/csv"
	"io"
)

func getCSVFromUrl(url string) ([][]string, string) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", "Chrome/56.0.2924.76")
	response, err := client.Do(req)
	if err != nil {
		return [][]string{}, "getCSV() Failed"
	}

	defer response.Body.Close()
	reader := csv.NewReader(response.Body)
	list := make([][]string, 0)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			break
			// should or do something useful log.Fatal(err)
		}
		list = append(list, record)
	}

	return list, "getCSV() Success"
}