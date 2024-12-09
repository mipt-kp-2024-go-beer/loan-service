package books

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/mipt-kp-2024-go-beer/loan-service/internal/fail"
)

func NewConn(url string) Connection {
	return &implConn{
		url: url,
		client: http.Client{
			Timeout: 10,
		},
	}
}

type implConn struct {
	url    string
	client http.Client
}

func (c *implConn) LookupBook(ID string) (Book, error) {
	request, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		fmt.Sprintf("%s/api/v1/books/%s", c.url, url.PathEscape(ID)),
		nil,
	)
	if err != nil {
		return Book{}, err
	}

	response, err := c.client.Do(request)
	if err != nil {
		return Book{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return Book{}, fmt.Errorf(
			"%w: %d %s %q",
			fail.ErrBookService,
			response.StatusCode,
			response.Status,
			string(body),
		)
	}

	var result struct {
		ID          string `json:"id"`
		Title       string `json:"title"`
		Author      string `json:"author"`
		Description string `json:"description"`
	}
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return Book{}, err
	}

	return Book{
		ID:          result.ID,
		Title:       result.Title,
		Author:      result.Author,
		Description: result.Description,
		TotalStock:  1, // TODO: Change when the book service implements this
	}, nil
}
