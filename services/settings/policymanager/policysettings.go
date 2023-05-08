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
	if result, err := policyMgr.validateWithModel(amap, "type", cond.cType); err == nil {
		cond.cType = result.(string)
		if _, ok := policyConditionTypeMap[cond.cType]; !ok {
			err := fmt.Errorf("found invalid condition type %s in PolicyCondition", cond.op)
			policyMgr.logger.Err("NewPolicyCondition error:", err)
			return nil, err
		}
	} else {
		policyMgr.logger.Err("NewPolicyCondition error:", err)
		return nil, err
	}
	if result, err := policyMgr.validateWithModel(amap, "op", cond.op); err == nil {
		cond.op = result.(string)
		if _, ok := policyConditionOpsMap[cond.op]; !ok {
			err := fmt.Errorf("found invalid op %s in PolicyCondition", cond.op)
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
	if result, err := policyMgr.validateWithModel(amap, "id", flow.id); err == nil {
		flow.id = result.(string)
	} else {
		policyMgr.logger.Err("NewPolicyFlowCategory id:", err)
		return nil, err
	}
	if result, err := policyMgr.validateWithModel(amap, "name", flow.name); err == nil {
		flow.name = result.(string)
	} else {
		policyMgr.logger.Err("NewPolicyFlowCategory with id: %s error on name:", flow.id, err)
		return nil, err
	}
	if result, err := policyMgr.validateWithModel(amap, "description", flow.description); err == nil {
		flow.description = result.(string)
	} else {
		policyMgr.logger.Err("NewPolicyFlowCategory with id: %s error on description:", flow.id, err)
		return nil, err
	}
	cmap := make([]interface{}, 1)
	if result, err := policyMgr.validateWithModel(amap, "conditions", cmap); err == nil {
		cmap = result.([]interface{})
		flow.conditions = make([]*PolicyCondition, len(cmap))
		for i, v := range cmap {
			var err error
			if flow.conditions[i], err = policyMgr.NewPolicyCondition(v); err != nil {
				policyMgr.logger.Err("Failed loading PolicyCondition %d for PolicyFlowCategory %s", i, flow.id)
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
	if result, err := policyMgr.validateWithModel(amap, "id", conf.id); err == nil {
		conf.id = result.(string)
	} else {
		policyMgr.logger.Err("NewPolicyConiguration with invalid id", err)
		return nil, err
	}
	if result, err := policyMgr.validateWithModel(amap, "name", conf.name); err == nil {
		conf.name = result.(string)
	} else {
		policyMgr.logger.Err("NewPolicyConiguration with id %s has invalid\\ name", conf.id, err)
		return nil, err
	}
	if result, err := policyMgr.validateWithModel(amap, "description", conf.description); err == nil {
		conf.description = result.(string)
	} else {
		policyMgr.logger.Err("NewPolicyConiguration with id %s has invalid description", conf.id, err)
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

	if result, err := policyMgr.validateWithModel(amap, "id", policy.id); err == nil {
		policy.id = result.(string)
	} else {
		policyMgr.logger.Err("NewPolicy with has invalid id", err)
		return nil, err
	}
	if result, err := policyMgr.validateWithModel(amap, "name", policy.name); err == nil {
		policy.name = result.(string)
	} else {
		policyMgr.logger.Err("NewPolicy with id %s has invalid name", policy.id, err)
		return nil, err
	}
	if result, err := policyMgr.validateWithModel(amap, "description", policy.description); err == nil {
		policy.description = result.(string)
	} else {
		policyMgr.logger.Err("NewPolicy with id %s has invalid description", policy.id, err)
		return nil, err
	}
	if result, err := policyMgr.validateWithModel(amap, "enabled", policy.enabled); err == nil {
		policy.enabled = result.(bool)
	} else {
		policyMgr.logger.Err("NewPolicy with id %s has invalid enabled", policy.id, err)
		return nil, err
	}
	policy.configurations = make(map[string]*PolicyConfiguration)
	configMap := make([]interface{}, 1)
	if result, err := policyMgr.validateWithModel(amap, "configurations", configMap); err == nil {
		configMap = result.([]interface{})
		for _, v := range configMap {
			id := v.(string)
			if _, ok := policyMgr.configurations[id]; ok {
				policy.configurations[id] = policyMgr.configurations[id]
			} else {
				err := fmt.Errorf("invalid Policy %s(%s) with unknown configuration: %s",
					policy.id, policy.name, id)
				policyMgr.logger.Err("NewPolicy error", err)
				return nil, err
			}
		}
	} else {
		err := fmt.Errorf("invalid Policy %s(%s) with nil configuration",
			policy.id, policy.name)
		policyMgr.logger.Err("NewPolicy error", err)
		return nil, err
	}
	policy.flows = make(map[string]*PolicyFlowCategory)
	flowMap := make([]interface{}, 1)
	if result, err := policyMgr.validateWithModel(amap, "flows", flowMap); err == nil {
		pfmap := result.([]interface{})
		for _, v := range pfmap {
			id := v.(string)
			if _, ok := policyMgr.flowCategories[id]; ok {
				policy.flows[id] = policyMgr.flowCategories[id]
			} else {
				err := fmt.Errorf("invalid Policy %s(%s) with unknown flow: %s",
					policy.id, policy.name, id)
				policyMgr.logger.Err("NewPolicy error", err)
				return nil, err
			}
		}
	} else {
		err := fmt.Errorf("invalid Policy %s(%s) with no flows", policy.id, policy.name)
		policyMgr.logger.Err("NewPolicy error", err)
		return nil, err
	}
	return &policy, nil
}
