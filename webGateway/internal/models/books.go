package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	apis "github.com/ChangHChen/Reading-Copilot/webGateway/internal/APIs"
)

type BookMeta struct {
	GutenID  int           `json:"id"`
	Title    string        `json:"title"`
	Authors  []apis.Author `json:"authors"`
	ImageURL string        `json:"image_url"`
}

type BookModel struct {
	DB *sql.DB
}

func GetTopBooksList() ([]BookMeta, error) {
	url := "https://gutendex.com/books?languages=en&sort=popular"

	resp, err := http.Get(url)
	if err != nil {
		return nil, ErrFetchingData
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch books: status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResp apis.BookListAPIResponse
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return nil, err
	}
	var books []BookMeta
	for _, result := range apiResp.Results[:10] {
		book := BookMeta{
			GutenID:  result.GutenID,
			Title:    result.Title,
			Authors:  result.Authors,
			ImageURL: result.Formats["image/jpeg"],
		}
		books = append(books, book)
	}
	return books, nil
}
