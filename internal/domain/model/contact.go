package model

type Contact struct {
	BaseModel
	IssuerId      string `json:"issuer_id" db:"issuer_id"`
	SubjectId     string `json:"subject_id" db:"subject_id"`
	ApplicationId string `json:"application_id" db:"application_id"`
	Type          string `json:"type" db:"type"`

	Name     string            `json:"name" db:"name"`
	Username string            `json:"username" db:"username"`
	Metadata map[string]string `json:"metadata" db:"metadata"`
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

func ContactAllowedFields() []string {
	return []string{"issuer_id", "application_id", "type", "name", "username", "metadata",
		"id", "domain_id", "created_by", "updated_by", "created_at", "updated_at"}
}
