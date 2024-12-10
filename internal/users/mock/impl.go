package mock

import (
	"context"
	"fmt"

	"github.com/mipt-kp-2024-go-beer/loan-service/internal/fail"
	"github.com/mipt-kp-2024-go-beer/loan-service/internal/users"
)

func NewConn() users.Connection {
	return &implConn{}
}

type implConn struct {
}

func (*implConn) VerifyToken(ctx context.Context, authToken string) (*users.User, error) {
	switch {
	case authToken == "invalid":
		return nil, fmt.Errorf("%w: pretend invalid token", fail.ErrUserService)
	case authToken == "regular-user":
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
	case authToken == "librarian":
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
