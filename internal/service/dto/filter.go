package dto

type (
	BaseFilter struct {
		DomainId int
		Page int
		Size int
		Sort string
		Q string
	}
)

func (b *BaseFilter) GetQ() *string {
	if b.Q == "" {
		return nil
	}

	return &b.Q
}