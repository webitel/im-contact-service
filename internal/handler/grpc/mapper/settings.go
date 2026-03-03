package mapper

import (
	"time"

	"github.com/google/uuid"
	"github.com/webitel/im-contact-service/gen/go/contact/v1"
	"github.com/webitel/im-contact-service/internal/domain/model"
	"github.com/webitel/im-contact-service/internal/service/dto"
	"github.com/webitel/webitel-go-kit/pkg/errors"
)


var (
	RuleNames = map[string]contact.RuleType{
		model.RuleTypeAddToGroup: contact.RuleType_ALLOW_ADD_TO_GROUPS,
		model.RuleTypeBotMessages: contact.RuleType_ALLOW_BOT_MESSAGES,
	}
)

//go:generate goverter gen github.com/webitel/im-contact-service/internal/handler/grpc/mapper

// goverter:converter
// goverter:matchIgnoreCase
// goverter:extend github.com/google/uuid:Parse
// goverter:extend time:UnixMilli
// goverter:extend ConvertRuleFromProto
type SettingsInConverter interface {
	ConvertGetSettingsRequest(*contact.GetContactSettingsRequest) (*dto.GetContactSettingsRequest, error)
	// goverter:useZeroValueOnPointerInconsistency
	ConvertUpdateSettingsRequest(*contact.UpdateContactSettingsRequest) (*dto.UpdateContactSettingsRequest, error)
}

// goverter:converter
// goverter:matchIgnoreCase
// goverter:ignoreUnexported
// goverter:extend ConvertUUID
// goverter:extend ConvertTimeToInt64
// goverter:extend ConvertRuleToProto
type SettingsOutConverter interface {
	ConvertSettings(*model.ContactSettings) (*contact.Settings, error)
}


func ConvertUUID(id uuid.UUID) string {
	return id.String()
}


func ConvertTimeToInt64(in time.Time) int64 {
	return in.UnixMilli()
}


func ConvertRuleFromProto(rule *contact.Rule) (model.SettingRule, error) {
	var ruleType string
	for serviceName, protoName := range RuleNames {
		if protoName == rule.GetType() {
			ruleType = serviceName
		}
	}
	if ruleType == "" {
		return nil, errors.InvalidArgument("unknown rule type")
	}
	internalRule, err := model.BuildRule(ruleType)
	if err != nil {
		return nil, err
	}

	var ruleValue any
	switch val := rule.GetValue().(type) {
	case *contact.Rule_Switcher:
		ruleValue = val.Switcher
	case *contact.Rule_List:
		ruleValue = val.List.GetValues()
	}

	if err := internalRule.SetValue(ruleValue); err != nil {
		return nil, err
	}
	return internalRule, nil
}


func ConvertRuleToProto(rule model.SettingRule) (*contact.Rule, error) {
	var (
		ok bool
		protoRule contact.Rule
	) 
	protoRule.Type, ok = RuleNames[rule.GetType()]
	if !ok {
		return nil, errors.InvalidArgument("unknown respective proto type for such rule", errors.WithValue("rule_type", rule.GetType()))
	}

	switch val := rule.GetValue().(type)  {
	case bool:
		protoRule.Value = &contact.Rule_Switcher{Switcher: val}
	case []int32:
		protoRule.Value = &contact.Rule_List{List: &contact.ListRuleValue{Values: val}}
	case []int:
		var list []int32
		for _, elem := range val {
			list = append(list, int32(elem))
		}
		protoRule.Value = &contact.Rule_List{List: &contact.ListRuleValue{Values:list}}
	default:
		return nil, errors.InvalidArgument("value not supported")
	}

	return &protoRule, nil

}
