package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"

	apis "github.com/ChangHChen/Reading-Copilot/webGateway/internal/APIs"
)

type BookMeta struct {
	GutenID       int           `json:"id"`
	Title         string        `json:"title"`
	Authors       []apis.Author `json:"authors"`
	ImageURL      string        `json:"image_url"`
	LocalImageURL string
	LocalTextURL  string
}

type BookModel struct {
	DB *sql.DB
}

func (m *BookModel) GetTopBooksList() ([]BookMeta, error) {
	var books []BookMeta
	var jsonData string

	err := m.DB.QueryRow("SELECT cache_value FROM gutendex_cache WHERE cache_key = 'topBooks' AND last_updated > NOW() - INTERVAL 1 DAY").Scan(&jsonData)
	if err != nil {
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

		jsonData = string(body)
		_, err = m.DB.Exec("REPLACE INTO gutendex_cache (cache_key, cache_value) VALUES ('topBooks', ?)", jsonData)
		if err != nil {
			return nil, err
		}
	}

	var apiResp apis.BookListAPIResponse
	if err := json.Unmarshal([]byte(jsonData), &apiResp); err != nil {
		return nil, err
	}
	for _, result := range apiResp.Results[:10] {
		localImageURL, _ := cacheCoverImage(result.GutenID, result.Formats["image/jpeg"])
		book := BookMeta{
			GutenID:       result.GutenID,
			Title:         result.Title,
			Authors:       result.Authors,
			ImageURL:      result.Formats["image/jpeg"],
			LocalImageURL: localImageURL,
		}
		books = append(books, book)
	}
	return books, nil
}

func Search(keyword string) ([]BookMeta, error) {
	var books []BookMeta
	var jsonData string
	keyword = url.QueryEscape(keyword)
	url := fmt.Sprintf("https://gutendex.com//books?search=%s&&sort=popular&&language=en", keyword)
	resp, err := http.Get(url)
	if err != nil {
		return nil, ErrFetchingData
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	jsonData = string(body)

	var apiResp apis.BookListAPIResponse
	if err := json.Unmarshal([]byte(jsonData), &apiResp); err != nil {
		return nil, err
	}
	length := len(apiResp.Results)
	if length == 0 {
		return nil, ErrNoSearchResult
	}
	if length > 10 {
		length = 10
	}
	for _, result := range apiResp.Results[:length] {
		localImageURL, _ := cacheCoverImage(result.GutenID, result.Formats["image/jpeg"])
		book := BookMeta{
			GutenID:       result.GutenID,
			Title:         result.Title,
			Authors:       result.Authors,
			ImageURL:      result.Formats["image/jpeg"],
			LocalImageURL: localImageURL,
		}
		books = append(books, book)
	}
	return books, nil
}

func cacheCoverImage(gutenID int, imageURL string) (string, error) {
	cachePath := filepath.Join("cache", strconv.Itoa(gutenID), "cover.jpg")
	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(cachePath), 0755); err != nil {
			return "", err
		}
		resp, err := http.Get(imageURL)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()
		out, err := os.Create(cachePath)
		if err != nil {
			return "", err
		}
		defer out.Close()
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return "", err
		}
	}
	return cachePath, nil
}

func GetBookCache(gutenID int) (BookMeta, error) {
	
}