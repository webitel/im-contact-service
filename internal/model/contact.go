package model

type Contact struct {
	BaseModel
	IssuerId      string `json:"issuer_id"`
	ApplicationId string `json:"application_id"`
	Type          string `json:"type"`

	Name     string            `json:"name"`
	Username string            `json:"username"`
	Metadata map[string]string `json:"metadata"`
}

func (c *Contact) Equal(compare *Contact) bool {
	if c == nil && compare == nil {
		return true
	}

	return c.ApplicationId == compare.ApplicationId &&
		c.Id == compare.Id &&
		c.DomainId == compare.DomainId &&
		c.IssuerId == compare.IssuerId &&
		c.Name == compare.Name &&
		c.Type == compare.Type &&
		c.Username == compare.Username
}
