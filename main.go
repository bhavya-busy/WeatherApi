package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"unicode"
)

type Pages struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

func find(v string) {
	var numStr string
	for _, char := range v {
		if unicode.IsDigit(char) {
			numStr += string(char)
		} else if numStr != "" {
			printNumber(numStr)
			numStr = ""
		}
	}
	if numStr != "" {
		printNumber(numStr)
	}
}

func printNumber(numStr string) {
	if num, err := strconv.Atoi(numStr); err == nil {
		fmt.Println(num)
	} else {
		fmt.Println("Error:", err)
	}
}

func main() {
	fmt.Println("weather")

	url := "https://jsonmock.hackerrank.com/api/weather/search?name="
	var name string

	fmt.Print("Enter the characters of name: ")
	fmt.Scanln(&name)
	url = url + name

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	defer response.Body.Close()
	body, readErr := io.ReadAll(response.Body)
	if readErr != nil {
		fmt.Println(readErr)
	}
	pages := Pages{}
	jsonErr := json.Unmarshal(body, &pages)
	if jsonErr != nil {
		fmt.Println(jsonErr)
	}

	for i := 1; i <= pages.TotalPages; i++ {
		nurl := url + fmt.Sprintf("&page=%v", i)
		res, err := http.Get(nurl)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return
		}
		var data struct {
			Data []struct {
				Name    string
				Weather string
				Status  []string
			}
		}

		jsonErr := json.Unmarshal(body, &data)
		if jsonErr != nil {
			fmt.Println(jsonErr)
		}

		for _, v := range data.Data {
			fmt.Println(v.Name)
			find(v.Weather)
			for _, val := range v.Status {
				find(val)
			}

		}

	}

}
