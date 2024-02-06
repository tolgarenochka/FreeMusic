package app_errors

type FileNotFound struct {
	Message string
}

func (e *FileNotFound) Error() string {
	return e.Message
}
