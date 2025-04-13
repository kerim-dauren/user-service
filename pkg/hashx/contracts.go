package hashx

type (
	Hasher interface {
		Hash(password string) (string, error)
	}
	Checker interface {
		CompareHashAndPassword(hashedPassword, password string) error
	}
)
