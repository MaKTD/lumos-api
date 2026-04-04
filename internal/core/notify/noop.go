package notify

type noopService struct{}

func NewNoop() Service {
	return &noopService{}
}

func (n *noopService) ForAdmin(_ string) {}
