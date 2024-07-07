package models

import (
	"bufio"
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
	TotalPageNum  int
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
	localTextFileURL, _ := cacheBookText(apiResp.GutenID, apiResp.Formats["text/plain; charset=us-ascii"])
	localTextURL, totalPageNum, err := splitTextIntoPages(localTextFileURL)
	if err != nil {
		return BookMeta{}, err
	}

	book = BookMeta{
		GutenID:       apiResp.GutenID,
		Title:         apiResp.Title,
		Authors:       apiResp.Authors,
		ImageURL:      apiResp.Formats["image/jpeg"],
		TextURL:       apiResp.Formats["text/plain; charset=us-ascii"],
		LocalImageURL: localImageURL,
		LocalTextURL:  localTextURL,
		TotalPageNum:  totalPageNum,
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

func (m *BookModel) GetLatest10History(userID int) ([]BookMeta, error) {
	var bookIDList []int
	stmt := `SELECT book_id FROM reading_progress WHERE user_id=? ORDER BY last_updated DESC LIMIT 10`
	rows, err := m.DB.Query(stmt, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var bookID int
		if err := rows.Scan(&bookID); err != nil {
			return nil, err
		}
		bookIDList = append(bookIDList, bookID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(bookIDList) == 0 {
		return nil, ErrNoSearchResult
	}

	books := make([]BookMeta, len(bookIDList))

	for i, bookID := range bookIDList {
		books[i], err = m.GetBookByID(bookID)
		if err != nil {
			return nil, err
		}
	}
	return books, nil
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

func splitTextIntoPages(filePath string) (string, int, error) {
	baseDir := filepath.Dir(filePath)
	pagesDir := filepath.Join(baseDir, "pages")

	if _, err := os.Stat(pagesDir); !os.IsNotExist(err) {
		totalPageNum, _ := countTxtFiles(pagesDir)
		return pagesDir, totalPageNum, nil
	}
	err := os.Mkdir(pagesDir, 0755)
	if err != nil {
		return "", 0, err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return "", 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var paragraphs []string
	var currentParagraph strings.Builder
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || line == "\n" {
			if currentParagraph.Len() > 0 {
				paragraphs = append(paragraphs, currentParagraph.String())
				currentParagraph.Reset()
			}
		} else {
			if currentParagraph.Len() > 0 {
				currentParagraph.WriteString(" ")
			}
			currentParagraph.WriteString(line)
		}
	}
	if currentParagraph.Len() > 0 {
		paragraphs = append(paragraphs, currentParagraph.String())
	}

	pageCount := 1
	var wordsInPage int
	var pageContent strings.Builder

	for _, paragraph := range paragraphs {
		wordCount := len(strings.Fields(paragraph))
		if wordCount+wordsInPage > 300 && wordsInPage != 0 {
			err := savePage(pagesDir, pageCount, pageContent.String())
			if err != nil {
				return "", 0, err
			}
			pageContent.Reset()
			wordsInPage = 0
			pageCount++
		}
		pageContent.WriteString(paragraph + "\n\n")
		wordsInPage += wordCount
	}

	if pageContent.Len() > 0 {
		err := savePage(pagesDir, pageCount, pageContent.String())
		if err != nil {
			return "", 0, err
		}
	}
	totalPageNum, _ := countTxtFiles(pagesDir)

	return pagesDir, totalPageNum, nil
}

func savePage(dir string, pageNumber int, content string) error {
	filename := fmt.Sprintf("page_%d.txt", pageNumber)
	filePath := filepath.Join(dir, filename)
	return os.WriteFile(filePath, []byte(content), 0644)
}

func countTxtFiles(dir string) (int, error) {
	var count int
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".txt" {
			count++
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return count, nil
}
