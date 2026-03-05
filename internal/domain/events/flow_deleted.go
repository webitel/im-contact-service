package events





const (
	FlowSchemaDeleteTopic = "todo.flow_event.delete"
)

type FlowSchemaDeleted struct {
	FlowID string
}