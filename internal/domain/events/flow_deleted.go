package events

const (
	FlowSchemaDeleteTopic = "flow_schema.delete.#"
)

type FlowSchemaDeleted struct {
	FlowID string `json:"flow_id"`
}
