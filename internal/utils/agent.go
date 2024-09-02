package utils

type Agent struct {
	id    uint
	admin bool
}

func NewAgent(id uint, admin bool) *Agent {
	return &Agent{
		id:    id,
		admin: admin,
	}
}

func (a *Agent) ID() uint {
	return a.id
}

func (a *Agent) IsAdmin() bool {
	return a.admin
}
