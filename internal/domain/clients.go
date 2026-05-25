package domain

type Client struct {
	ID int64

	UUID string

	Name string

	Secret string

	Domain string

	Enabled bool

	Scopes []string
}
