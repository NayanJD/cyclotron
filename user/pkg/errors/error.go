package customerrors

type ConstError string

func (err ConstError) Error() string {
	return string(err)
}
