package commands

type UsageError struct {
	Msg string
}

func (e UsageError) Error() string {
	return e.Msg
}
