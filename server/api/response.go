package api

type Error struct {
	HTTPStatus int `json:"-"`

	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e Error) Error() string {
	return e.Message
}
