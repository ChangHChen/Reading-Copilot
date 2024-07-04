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
	"strings"

	apis "github.com/ChangHChen/Reading-Copilot/webGateway/internal/APIs"
)

type BookMeta struct {
	GutenID       int
	Title         string
	Authors       []apis.Author
	ImageURL      string
	TextURL       string
	LocalImageURL string
	LocalTextURL  string
}

type BookModel struct {
	DB *sql.DB
}

func fetchFromGutendex(url string) ([]byte, error) {
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
	return body, nil
}

func parseBookList(jsonbody []byte) ([]BookMeta, error) {
	var books []BookMeta
	var apiResp apis.BookListAPIResponse
	if err := json.Unmarshal([]byte(jsonbody), &apiResp); err != nil {
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

func parseBook(jsonbody []byte) (BookMeta, error) {
	var book BookMeta
	var apiResp apis.BookAPIResponse
	if err := json.Unmarshal([]byte(jsonbody), &apiResp); err != nil {
		return BookMeta{}, err
	}
	localImageURL, _ := cacheCoverImage(apiResp.GutenID, apiResp.Formats["image/jpeg"])
	localTextURL, _ := cacheBookText(apiResp.GutenID, apiResp.Formats["text/plain; charset=us-ascii"])

	book = BookMeta{
		GutenID:       apiResp.GutenID,
		Title:         apiResp.Title,
		Authors:       apiResp.Authors,
		ImageURL:      apiResp.Formats["image/jpeg"],
		TextURL:       apiResp.Formats["text/plain; charset=us-ascii"],
		LocalImageURL: localImageURL,
		LocalTextURL:  localTextURL,
	}

	return book, nil
}

func (m *BookModel) GetTopBooksList() ([]BookMeta, error) {
	var books []BookMeta
	var jsonString string
	var jsonData []byte

	err := m.DB.QueryRow("SELECT cache_value FROM gutendex_cache WHERE cache_key = 'topBooks' AND last_updated > NOW() - INTERVAL 1 DAY").Scan(&jsonString)
	jsonData = []byte(jsonString)
	if err != nil {
		url := "https://gutendex.com/books?languages=en&sort=popular"
		jsonData, err = fetchFromGutendex(url)
		if err != nil {
			return nil, err
		}
		_, err = m.DB.Exec("REPLACE INTO gutendex_cache (cache_key, cache_value) VALUES ('topBooks', ?)", jsonData)
		if err != nil {
			return nil, err
		}
	}
	books, err = parseBookList(jsonData)
	if err != nil {
		return nil, err
	}
	return books, nil
}

func (m *BookModel) GetBookByID(id int) (BookMeta, error) {
	var book BookMeta
	var jsonString string
	var jsonData []byte

	err := m.DB.QueryRow("SELECT cache_value FROM gutendex_cache WHERE cache_key = ? AND last_updated > NOW() - INTERVAL 15 DAY", fmt.Sprintf("GutenIDX%d", id)).Scan(&jsonString)
	jsonData = []byte(jsonString)
	if err != nil {
		url := fmt.Sprintf("https://gutendex.com/books/%d", id)
		jsonData, err = fetchFromGutendex(url)
		if err != nil {
			return BookMeta{}, err
		}
		_, err = m.DB.Exec("REPLACE INTO gutendex_cache (cache_key, cache_value) VALUES (?, ?)", fmt.Sprintf("GutenIDX%d", id), jsonData)
		if err != nil {
			return BookMeta{}, err
		}
	}
	book, err = parseBook(jsonData)
	if err != nil {
		return BookMeta{}, err
	}
	return book, nil
}

func Search(keyword string) ([]BookMeta, error) {
	var books []BookMeta
	keyword = url.QueryEscape(keyword)
	url := fmt.Sprintf("https://gutendex.com//books?search=%s&&sort=popular&&language=en", keyword)
	jsonData, err := fetchFromGutendex(url)
	if err != nil {
		return nil, err
	}
	books, err = parseBookList(jsonData)
	if err != nil {
		return nil, err
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

func cacheBookText(gutenID int, textURL string) (string, error) {
	cachePath := filepath.Join("cache", strconv.Itoa(gutenID), "plain_text.txt")
	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(cachePath), 0755); err != nil {
			return "", err
		}
		resp, err := http.Get(textURL)
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
		err = readAndProcessTextFile(cachePath)
		if err != nil {
			return "", err
		}
	}

	return cachePath, nil
}

func processText(text string) string {
	lines := strings.Split(text, "\n")
	var result strings.Builder

	for i := 0; i < len(lines); i++ {
		currentLine := strings.TrimSpace(lines[i])

		if currentLine == "" {
			if i+1 < len(lines) && strings.TrimSpace(lines[i+1]) != "" {
				result.WriteString("\n\n")
			}
		} else {
			if result.Len() > 0 {
				result.WriteString(" ")
			}
			result.WriteString(currentLine)
		}
	}

	return result.String()
}

func readAndProcessTextFile(filepath string) error {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	processedContent := processText(string(content))
	err = os.WriteFile(filepath, []byte(processedContent), 0644)
	if err != nil {
		return err
	}
	return nil
}
