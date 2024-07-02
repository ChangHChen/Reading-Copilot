package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Book struct {
	GutenID  int      `json:"id"`
	Title    string   `json:"title"`
	Authors  []Author `json:"authors"`
	ImageURL string   `json:"image_url"`
}

type Author struct {
	Name string `json:"name"`
}

type APIResponse struct {
	Results []struct {
		GutenID int               `json:"id"`
		Title   string            `json:"title"`
		Authors []Author          `json:"authors"`
		Formats map[string]string `json:"formats"`
	} `json:"results"`
}

func main() {
	url := "https://gutendex.com/books?languages=en&sort=popular"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching data:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	var apiResp APIResponse
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	var books []Book
	for _, result := range apiResp.Results[:10] {
		book := Book{
			GutenID:  result.GutenID,
			Title:    result.Title,
			Authors:  result.Authors,
			ImageURL: result.Formats["image/jpeg"],
		}
		books = append(books, book)
	}

	for _, book := range books {
		fmt.Printf("GutenID: %d, Title: %s, Image URL: %s\n", book.GutenID, book.Title, book.ImageURL)
		for _, author := range book.Authors {
			fmt.Printf("Author: %s\n", author.Name)
		}
	}
}
