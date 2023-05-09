package policy

import (
	"errors"
	"fmt"
	"reflect"
)

var policyConditionTypeMap = map[string]int{
	"CLIENT_ADDRESS": 1,
	"CLIENT_PORT":    1,
	"DAY_OF_WEEK":    1,
	"DEST_ADDRESS":   1,
	"INTERFACE":      1,
	"SERVER_ADDRESS": 1,
	"SERVER_PORT":    1,
	"SOURCE_ADDRESS": 1,
	"PROTOCOL_TYPE":  1,
	"TIME_OF_DAY":    1,
	"VLAN_ID":        1,
}

// Valid PolicyCondition Ops - there may be more at some point
// == implies an OR operation between the different entries in the value arrray
// != implies an AND operation between the different entries in the value array
// all other operations assume a single entry in the value array (or string)
var policyConditionOpsMap = map[string]int{"==": 1, "!=": 1, "<": 1, ">": 1, "<=": 1, ">=": 1}

func (policyMgr *PolicyManager) validateWithModel(amap map[string]interface{}, key string, model interface{}) (interface{}, error) {
	if _, ok := amap[key]; ok {
		if reflect.TypeOf(amap[key]) == reflect.TypeOf(model) {
			switch model.(type) {
			case bool:
				return amap[key].(bool), nil
			case string:
				return amap[key].(string), nil
			case []interface{}:
				return amap[key].([]interface{}), nil
			default:
				// Should not happen
				return nil, errors.New("unexpected type in validateWithModel")
			}
		} else {
			return nil, fmt.Errorf("invalid type %T expected %T for key %s", amap[key], model, key)
		}
	} else {
		return nil, fmt.Errorf("invalid nil type expected %T for key %s", model, key)
	}
}

func (policyMgr *PolicyManager) NewPolicyCondition(value interface{}) (*PolicyCondition, error) {
	amap := value.(map[string]interface{})
	cond := PolicyCondition{}
	if result, err := policyMgr.validateWithModel(amap, "type", cond.CType); err == nil {
		cond.CType = result.(string)
		if _, ok := policyConditionTypeMap[cond.CType]; !ok {
			err := fmt.Errorf("found invalid condition type %s in PolicyCondition", cond.Op)
			policyMgr.logger.Err("NewPolicyCondition error:", err)
			return nil, err
		}
	} else {
		policyMgr.logger.Err("NewPolicyCondition error:", err)
		return nil, err
	}
	if result, err := policyMgr.validateWithModel(amap, "op", cond.Op); err == nil {
		cond.Op = result.(string)
		if _, ok := policyConditionOpsMap[cond.Op]; !ok {
			err := fmt.Errorf("found invalid op %s in PolicyCondition", cond.Op)
			policyMgr.logger.Err("NewPolicyCondition error:", err)
			return nil, err
		}
	} else {
		policyMgr.logger.Err("NewPolicyCondition error:", err)
		return nil, err
	}
	if _, ok := amap["value"]; ok {
		switch amap["value"].(type) {
		case string:
			cond.value = make([]string, 1)
			cond.value[0] = amap["value"].(string)
		case []interface{}:
			carray := amap["value"].([]interface{})
			cond.value = make([]string, len(carray))
			for i, v := range carray {
				cond.value[i] = v.(string)
			}
		default:
			err := fmt.Errorf("invalid type %T for value in PolicyCondition: ", amap["type"])
			policyMgr.logger.Err("NewPolicyCondition error:", err)
			return nil, err
		}
	} else {
		err := fmt.Errorf("invalid PolicyCondition with nil value")
		policyMgr.logger.Err("NewPolicyCondition error:", err)
		return nil, err
	}
	return &cond, nil
}

func (policyMgr *PolicyManager) NewPolicyFlowCategory(value interface{}) (*PolicyFlowCategory, error) {
	amap := value.(map[string]interface{})
	flow := PolicyFlowCategory{}
	if result, err := policyMgr.validateWithModel(amap, "id", flow.Id); err == nil {
		flow.Id = result.(string)
	} else {
		policyMgr.logger.Err("NewPolicyFlowCategory id:", err)
		return nil, err
	}
	if result, err := policyMgr.validateWithModel(amap, "name", flow.Name); err == nil {
		flow.Name = result.(string)
	} else {
		policyMgr.logger.Err("NewPolicyFlowCategory with id: %s error on name:", flow.Id, err)
		return nil, err
	}
	if result, err := policyMgr.validateWithModel(amap, "description", flow.Description); err == nil {
		flow.Description = result.(string)
	} else {
		policyMgr.logger.Err("NewPolicyFlowCategory with id: %s error on description:", flow.Id, err)
		return nil, err
	}
	cmap := make([]interface{}, 1)
	if result, err := policyMgr.validateWithModel(amap, "conditions", cmap); err == nil {
		cmap = result.([]interface{})
		flow.conditions = make([]*PolicyCondition, len(cmap))
		for i, v := range cmap {
			var err error
			if flow.conditions[i], err = policyMgr.NewPolicyCondition(v); err != nil {
				policyMgr.logger.Err("Failed loading PolicyCondition %d for PolicyFlowCategory %s", i, flow.Id)
				return &flow, err
			}
		}
	} else {
		policyMgr.logger.Err("NewPolicyFlowCategory errror:", err)
		return nil, err
	}
	return &flow, nil
}

func (policyMgr *PolicyManager) NewPolicyConfiguration(value interface{}) (*PolicyConfiguration, error) {
	amap := value.(map[string]interface{})
	conf := PolicyConfiguration{}
	if result, err := policyMgr.validateWithModel(amap, "id", conf.Id); err == nil {
		conf.Id = result.(string)
	} else {
		policyMgr.logger.Err("NewPolicyConiguration with invalid id", err)
		return nil, err
	}
	if result, err := policyMgr.validateWithModel(amap, "name", conf.Name); err == nil {
		conf.Name = result.(string)
	} else {
		policyMgr.logger.Err("NewPolicyConiguration with id %s has invalid\\ name", conf.Id, err)
		return nil, err
	}
	if result, err := policyMgr.validateWithModel(amap, "description", conf.Description); err == nil {
		conf.Description = result.(string)
	} else {
		policyMgr.logger.Err("NewPolicyConiguration with id %s has invalid description", conf.Id, err)
		return nil, err
	}
	/* 	for i, v := range amap["plugins"].([]string) {
	   		flow.plugins[i] =
	   	}
	*/
	return &conf, nil
}

func (policyMgr *PolicyManager) NewPolicy(value interface{}) (*Policy, error) {
	amap := value.(map[string]interface{})
	policy := Policy{}

	if result, err := policyMgr.validateWithModel(amap, "id", policy.Id); err == nil {
		policy.Id = result.(string)
	} else {
		policyMgr.logger.Err("NewPolicy with has invalid id", err)
		return nil, err
	}
	if result, err := policyMgr.validateWithModel(amap, "name", policy.Name); err == nil {
		policy.Name = result.(string)
	} else {
		policyMgr.logger.Err("NewPolicy with id %s has invalid name", policy.Id, err)
		return nil, err
	}
	if result, err := policyMgr.validateWithModel(amap, "description", policy.Description); err == nil {
		policy.Description = result.(string)
	} else {
		policyMgr.logger.Err("NewPolicy with id %s has invalid description", policy.Id, err)
		return nil, err
	}
	if result, err := policyMgr.validateWithModel(amap, "enabled", policy.Enabled); err == nil {
		policy.Enabled = result.(bool)
	} else {
		policyMgr.logger.Err("NewPolicy with id %s has invalid enabled", policy.Id, err)
		return nil, err
	}
	return &policy, nil
}
