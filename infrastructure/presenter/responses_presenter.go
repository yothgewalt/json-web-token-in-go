package presenter

type ResponseAny[T any] struct {
	Any []T `json:"any"`
}

type ResponseOnlyMessage struct {
	Message string `json:"message"`
}

type ResponseOnlyError struct {
	Error error `json:"error"`
}

type ResponseCollection struct {
	ID           uint   `json:"id"`
	Username     string `json:"username"`
	Firstname    string `json:"firstname"`
	Lastname     string `json:"lastname"`
	PhoneNumber  string `json:"phone_number"`
	EmailAddress string `json:"email_address"`
	FacebookLink string `json:"facebook_link"`
}
