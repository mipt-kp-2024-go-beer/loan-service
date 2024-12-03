package users

// User stores the information about a user
type User struct {
	// ID is the unique identifier of the user
	ID string
	// Login is the mnemonic identifier of the user
	Login string
	// Name is the first name of the user
	Name string
	// Surname is the last name of the user
	Surname string
	// Permissions is the bitmask for the user's permissions
	Permissions Permission
}

// Permission is the flag enum for the permissions a user might have
type Permission uint64

const (
	// PermManageBooks allows the user to add, edit and delete books from the library
	PermManageBooks Permission = 1 << 0
	// PermQueryTotalStock allows the user to get the total stored book count
	PermQueryTotalStock Permission = 1 << 1
	// PermChangeTotalStock allows the user to register updates to the total stored book count.
	// Requires PermGetTotalStock as a prerequisite.
	PermChangeTotalStock Permission = 1 << 2
	// PermQueryUsers allows the user to get information about other users, including their permissions.
	// Not required to get information about oneself, other rules apply.
	PermQueryUsers Permission = 1 << 3
	// PermManageUsers allows the user to add, edit and delete other users.
	// Not required to manage oneself, other rules apply.
	// Requires PermQueryUsers as a prerequisite.
	PermManageUsers Permission = 1 << 4
	// PermGrantPermissions allows the user to grant permissions to other users.
	// Only a subset of own permissions may be granted.
	// Requires PermQueryUsers as a prerequisite.
	PermGrantPermissions Permission = 1 << 5
	// PermLoanBooks allows the user to register book takeouts and returns.
	PermLoanBooks Permission = 1 << 6
	// PermQueryAvailableStock allows the user to get the number of available (not lent out) copies of a book.
	PermQueryAvailableStock Permission = 1 << 7
	// PermQueryReservations allows the user to get information related to book reservations.
	PermQueryReservations Permission = 1 << 8
)

// Connection is the interface for the private API of the users microservice
type Connection interface {
	// VerifyToken cheks the authentication token and returns the information about the associated user if it is valid
	VerifyToken(authToken string) (User, error)
}

// HasPerm returns true if the user has the given permission
func (u User) HasPerm(p Permission) bool {
	return u.Permissions&p != 0
}
