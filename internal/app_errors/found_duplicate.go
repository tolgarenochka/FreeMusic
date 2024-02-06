package app_errors

type DuplicateFound struct {
	Message string
}

func (e *DuplicateFound) Error() string {
	return e.Message
}
