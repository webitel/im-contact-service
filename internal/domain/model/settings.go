package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/webitel/webitel-go-kit/pkg/errors"
)

const (
	RuleTypeBotMessages = "bot_messages"
	RuleTypeAddToGroup  = "add_to_group"
)
type SettingRule interface {
	GetType() string
	SetValue(any) error
	GetValue() any
}

type ContactSettings struct {
	ID         uuid.UUID           `json:"id" db:"id"`
	ContactID  uuid.UUID            `json:"contact_id" db:"contact_id"`
	UpdatedAt  time.Time            `json:"updated_at" db:"updated_at"`
	Rules    []SettingRule `json:"rules" db:"rules"`
}



func BuildRule(ruleName string) (SettingRule, error) {
	switch ruleName {
	case RuleTypeBotMessages:
		return newBotMessagesRule(), nil
	case RuleTypeAddToGroup:
		return newAddToGroupSetting(), nil
	default:
		return nil, errors.InvalidArgument(fmt.Sprintf("unknown rule type: %s", ruleName))
	}
}

func newBotMessagesRule() *botMessagesRule{
	return &botMessagesRule{
		name: RuleTypeBotMessages,
	}
}

func newAddToGroupSetting() *addToGroupRule {
	return &addToGroupRule{
		name: RuleTypeAddToGroup,
	}
}

type botMessagesRule struct {
	name string 
	allow     bool
}

func (s *botMessagesRule) GetType() string {
	return s.name
}

func (s *botMessagesRule) SetValue(value any) error {
	var (
		allow bool
		ok    bool
	)
	if allow, ok = value.(bool); !ok {
		return errors.InvalidArgument("invalid value type for botMessagesRule, expected bool")
	} 


	s.allow = allow
	return nil
}

func (s *botMessagesRule) GetValue() any {
	return s.allow
}


type addToGroupRule struct {
	name      string 
	allow bool
}

func (s *addToGroupRule) GetType() string {
	return s.name
}

func (s *addToGroupRule) SetValue(value any) error {
	var (
		allow bool
		ok    bool
	)
	if allow, ok = value.(bool); !ok {
		return errors.InvalidArgument("invalid value type for addToGroupRule, expected bool")
	}
	
	s.allow = allow
	return nil
}

func (s *addToGroupRule) GetValue() any {
	return s.allow
}

