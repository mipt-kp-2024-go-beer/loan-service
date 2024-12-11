package books

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/mipt-kp-2024-go-beer/loan-service/internal/fail"
)

func NewConn(url string) Connection {
	return &implConn{
		url: url,
		client: http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type implConn struct {
	url    string
	client http.Client
}

func (c *implConn) LookupBook(ctx context.Context, ID string) (*Book, error) {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("http://%s/api/v1/books/%s", c.url, url.PathEscape(ID)),
		nil,
	)
	if err != nil {
		return nil, err
	}

	response, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf(
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
		Stock       string `json:"stock"`
	}
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	stock, err := strconv.ParseUint(result.Stock, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse stock: %w", fail.ErrBookService, err)
	}

	return &Book{
		ID:          result.ID,
		Title:       result.Title,
		Author:      result.Author,
		Description: result.Description,
		TotalStock:  uint(stock),
	}, nil
}
