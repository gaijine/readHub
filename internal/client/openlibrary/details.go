package openlibrary

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"readHub/internal/domain"
)

func (c *Client) GetBookDetails(openLibraryID string) (domain.BookDetails, error) {
	url := c.baseURL + "/works/" + openLibraryID + ".json"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return domain.BookDetails{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return domain.BookDetails{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		status := strconv.Itoa(resp.StatusCode)
		return domain.BookDetails{}, fmt.Errorf("unexpected status code: %s", status)
	}

	var response workResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return domain.BookDetails{}, err
	}
	var coverURL string
	if len(response.Covers) == 0 {
		coverURL = ""
	} else {
		cover := strconv.FormatInt(response.Covers[0], 10)
		coverURL = "https://covers.openlibrary.org/b/id/" + cover + "-L.jpg"
	}

	book := domain.BookDetails{
		OpenLibraryID: openLibraryID,
		Title:         response.Title,
		CoverURL:      coverURL,
	}
	return book, nil
}
