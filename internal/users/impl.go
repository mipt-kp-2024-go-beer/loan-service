package users

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/mipt-kp-2024-go-beer/loan-service/internal/fail"
)

func NewConn(url string) Connection {
	return &implConn{
		url: url,
		client: http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

type implConn struct {
	url    string
	client http.Client
}

func (c *implConn) VerifyToken(ctx context.Context, authToken string) (*User, error) {
	responseID, err := c.makeRequest(ctx, "/user/id", authToken)
	if err != nil {
		return nil, err
	}
	defer responseID.Body.Close()

	var resultID struct {
		ID string `json:"ID"`
	}
	err = json.NewDecoder(responseID.Body).Decode(&resultID)
	if err != nil {
		return nil, err
	}

	responsePermissions, err := c.makeRequest(ctx, "/user/permissions", authToken)
	if err != nil {
		return nil, err
	}
	defer responsePermissions.Body.Close()

	var resultPermissions struct {
		// Yes, the typo is on their side, unfortunately...
		Permissions string `json:"permissios"`
	}
	err = json.NewDecoder(responsePermissions.Body).Decode(&resultPermissions)
	if err != nil {
		return nil, err
	}

	permissions, err := strconv.ParseUint(resultPermissions.Permissions, 10, 64)
	if err != nil {
		return nil, err
	}

	return &User{
		ID: resultID.ID,
		// TODO: Populate if the user service starts providing these
		Login:       "",
		Name:        "",
		Surname:     "",
		Permissions: Permission(permissions),
	}, nil
}

func (c *implConn) makeRequest(ctx context.Context, endpoint string, authToken string) (*http.Response, error) {
	packedToken, err := json.Marshal(struct {
		Token string `json:"token"`
	}{
		Token: authToken,
	})
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("http://%s%s", c.url, endpoint),
		bytes.NewReader(packedToken),
	)
	if err != nil {
		return nil, err
	}

	response, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		response.Body.Close()
		return nil, fmt.Errorf(
			"%w: %d %s %q",
			fail.ErrUserService,
			response.StatusCode,
			response.Status,
			string(body),
		)
	}

	return response, nil
}
