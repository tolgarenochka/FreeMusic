package app_errors

type NoData struct {
	Message string
}

func (e *NoData) Error() string {
	return e.Message
}
