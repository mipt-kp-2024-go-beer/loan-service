package mock

import (
	"context"
	"fmt"

	"github.com/mipt-kp-2024-go-beer/loan-service/internal/fail"
	"github.com/mipt-kp-2024-go-beer/loan-service/internal/users"
)

func NewUsersConn() users.Connection {
	return &implUsersConn{}
}

type implUsersConn struct {
}

func (*implUsersConn) VerifyToken(ctx context.Context, authToken string) (*users.User, error) {
	switch {
	case authToken == "token-invalid":
		return nil, fmt.Errorf("%w: pretend invalid token", fail.ErrUserService)
	case authToken == "token-regular-user":
		return &users.User{
			ID:      "vasya-pupkin",
			Login:   "",
			Name:    "",
			Surname: "",
			Permissions: users.Permission(
				users.PermQueryAvailableStock |
					users.PermQueryReservations |
					users.PermQueryTotalStock,
			),
		}, nil
	case authToken == "token-librarian":
		return &users.User{
			ID:      "yuuko-shirakawa",
			Login:   "",
			Name:    "",
			Surname: "",
			// I.e. all permissions
			Permissions: users.Permission(0xffffffffffffffff),
		}, nil
	}
	panic("Unexpected request to mock user service!")
}
