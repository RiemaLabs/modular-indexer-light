package constant

var ApiState = StatusSync

type ApiStatus int

const (
	StateActive ApiStatus = iota + 1
	StatusSync
	StatusVerify
)

func (p ApiStatus) String() string {
	switch p {
	case StateActive:
		return "system is ready"
	case StatusSync:
		return "system is syncing"
	case StatusVerify:
		return "system is verifying"
	default:
		return ""
	}
}
