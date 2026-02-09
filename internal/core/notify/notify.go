package notify

type Service interface {
	ForAdmin(msg string)
}
