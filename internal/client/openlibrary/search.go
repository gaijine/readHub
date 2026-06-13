package openlibrary

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"readHub/internal/domain"

	"github.com/k0kubun/pp"
)

const searchPath = "/search.json"

func (c *Client) SearchBooks(query string) ([]domain.SearchBook, error) {
	params := url.Values{} // создали пустую коробку(контейнер параметров, обычно хранит данные map[string][]string)
	params.Set("q", query) // добавили ключ/значение, ключ=q, значение=query
	params.Set("limit", "5")
	encodedQuery := params.Encode() // метод Encode кодирует параметры, и получается limit=5&q=harry+potter

	url := c.baseURL + searchPath + "?" + encodedQuery

	req, err := http.NewRequest(http.MethodGet, url, nil) // создали запрос но не отправили
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req) // отправлен запрос клиенту, и в Response: Status, StatusCode, Header, Body
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		status := strconv.Itoa(resp.StatusCode)
		return nil, fmt.Errorf("unexpected status code: %s", status)
	}

	var response SearchResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	var resultList []domain.SearchBook
	for _, v := range response.Docs {
		openLibID := strings.TrimPrefix(v.Key, "/works/")
		cover := strconv.FormatInt(v.CoverID, 10)
		coverURL := "https://covers.openlibrary.org/b/id/" + cover + "-L.jpg"

		book := domain.SearchBook{
			OpenLibraryID: openLibID,
			Title:         v.Title,
			Author:        v.AuthorName,
			CoverURL:      coverURL,
		}
		resultList = append(resultList, book)
	}
	pp.Println(resultList)
	return resultList, nil
}

// func (c *Client) GetByOpenLibraryID(openLibraryID string) (domain.SearchBook, error){

// }
