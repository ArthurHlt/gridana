package converters

var ErrProbeNotFound errConverter = "No probe found"

type errConverter string

func (e errConverter) Error() string {
	return string(e)
}
