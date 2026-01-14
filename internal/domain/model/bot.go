package model

type WebitelBot struct {
	BaseModel

	FlowId      int    `json:"flow_id"`
	DisplayName string `json:"display_name"`
}

func BotAllowedFields() []string {
	return []string{"id", "created_at", "updated_at", "domain_id", "flow_id", "display_name"}
}