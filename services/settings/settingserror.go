package settings

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/r3labs/diff/v2"
)

/*
The purpose of the structs/functions in this file is to handle the case when a CONFIRM object comes back from sync-settings.
This would happen when trying to enable a rule with an associated disable policy or disable a policy with an associated enabled rule.

In the structs below,
the policy = parent
rule = child
interface can be parent to policy

Steps of execution :
1. Sync-settings returns an error with the CONFIRM object as a string in it.
2. We determine the current settings.
3. We determine the changeset between the current settings the newly created settings that failed.
4. From the changeset we put each change into an enabled, deleted, or disabled array.
5. For each of the three arrays:
	1. For each change in the array:
		1. find the affected item given the change. If enabling a rule, the affected item would be the rule.
		2. Based on the affected item id, collect the parents/children invalid items based on the type of change from the sync-settings error
			enable          = build based on parents
			disable/deleted = build based on children
	2. From all invalid items for change, put into a UI message
6. Put all three arrays UI messages into one message and send to UI
*/

// SetSettingsError represents the top structure from sync-settings with CONFIRM
// so we can pick up the error
type SetSettingsError struct {
	Confirm Confirmation `json:"CONFIRM"`
}

// Confirmation represents the array of invalid items from sync-settings
type Confirmation struct {
	InvalidItems map[string]InvalidItem `json:"invalidItems"`
}

// InvalidItem represents the invalid item from sync-settings with child/parent id,
// reason, value, and type
type InvalidItem struct {
	Reason   string `json:"reason"`
	Type     string `json:"type"`
	Value    string `json:"value"`
	ParentID string `json:"parentId"`
}

// SetSettingsErrorUI represents the structure sent to the UI for reading the error
type SetSettingsErrorUI struct {
	MainTranslationString string          `json:"mainTranslationString"`
	InvalidReason         string          `json:"invalidReason"`
	AffectedValues        []AffectedValue `json:"affectedValues"`
}

// AffectedValue represents each affected item that will be read by the UI
type AffectedValue struct {
	AffectedType  string `json:"affectedType"`
	AffectedValue string `json:"affectedValue"`
}

// determineSetSettingsError determines the error to send back after sync-settings
// @param origErr error - original error from sync-settings
// @param output string - original output from sync-settings
// @param settingsFile string - location of settings file
// @param jsonSettings map[string]interface{} - new json settings that did not succeed in saving
// @return interface{} - map of error and output of the newly built errorMessage and the original output from sync-settings
// @return error - error from building the new message or the origErr which means success in building the message
func determineSetSettingsError(origErr error, output string, settingsFile string, jsonSettings map[string]interface{}) (interface{}, error) {
	var err error = nil

	// get old settings to compare with new settings to determine changes
	jsonSettingsOld, oldSettingsErr := readSettingsFileJSON(settingsFile)
	if oldSettingsErr == nil {
		// build messages for delete, disable, and enable
		errorMessage, buildMessageErr := buildMessage(jsonSettingsOld, jsonSettings, origErr)
		if buildMessageErr == nil {
			return map[string]interface{}{"error": errorMessage, "output": output}, origErr
		}
		err = buildMessageErr
	} else {
		err = oldSettingsErr
	}

	// make sure err has something
	if err == nil {
		err = errors.New("Failed to determine set settings")
	}

	logger.Warn("Failed to find the set settings error: %s\n", err.Error())
	logger.Warn("Original error: %s\n", origErr.Error())
	return map[string]interface{}{"error": "failed_sync_settings", "output": output}, err
}

// buildMessage builds the message that will given to the UI to display
// @param jsonSettingsOld map[string]interface{} - current settings on system to compare with newly created settings
// @param jsonSettings map[string]interface{} - newly created settings that failed to save
// @param origErr error - original error from sync-settings
// @return string - build message to be sent to the UI
// @return error - any error in building the message, nil if none
func buildMessage(jsonSettingsOld map[string]interface{}, jsonSettings map[string]interface{}, origErr error) (string, error) {
	// determine changeset between new settings and old
	deletedChanges, disableChanges, enableChanges, changeSetErr := determineChangeSet(jsonSettingsOld, jsonSettings)
	if changeSetErr != nil {
		return "", changeSetErr
	}

	// create the SetSettingsError golang struct from original error for easier handling
	settingsError, settingsErrorStructErr := getSettingsErrorStruct(origErr)
	if settingsErrorStructErr != nil {
		return "", settingsErrorStructErr
	}

	// create array of SetSettingsErrorUI and process deleted, disabled, and enabled changes
	messages := make([]SetSettingsErrorUI, 0)
	if len(deletedChanges) > 0 {
		deleteMessage, deleteMsgErr := determineMessage(settingsError, "deleted", "child", deletedChanges, jsonSettingsOld)
		if deleteMsgErr != nil {
			return "", deleteMsgErr
		}
		messages = append(messages, deleteMessage)
	}

	if len(disableChanges) > 0 {
		disableMessage, disableMsgErr := determineMessage(settingsError, "disabled", "child", disableChanges, jsonSettingsOld)
		if disableMsgErr != nil {
			return "", disableMsgErr
		}
		messages = append(messages, disableMessage)
	}

	if len(enableChanges) > 0 {
		enableMessage, enableMsgErr := determineMessage(settingsError, "enabled", "parent", enableChanges, jsonSettingsOld)
		if enableMsgErr != nil {
			return "", enableMsgErr
		}
		messages = append(messages, enableMessage)
	}

	// convert structs into string to send to UI
	bytes, bytesErr := json.Marshal(messages)
	if bytesErr != nil {
		logger.Warn("Failed to unmarshal messages: %s\n", bytesErr.Error())
		return "", bytesErr
	}

	return "CONFIRM: " + string(bytes), nil
}

// determineChangeSet determines changes that occured from new/old so we can show relevant information for each change
// @param jsonSettingsOld map[string]interface{} - current settings on system to compare with newly created settings
// @param jsonSettings map[string]interface{} - newly created settings that failed to save
// @return []diff.Change - array of changes that were 'deleted' changes
// @return []diff.Change - array of changes that were 'disabled' changes
// @return []diff.Change - array of changes that were 'enabled' changes
// @return error - any error in getting arrays, nil if none
func determineChangeSet(jsonSettingsOld map[string]interface{}, jsonSettings map[string]interface{}) ([]diff.Change, []diff.Change, []diff.Change, error) {
	deletedChanges := make([]diff.Change, 0)
	disableChanges := make([]diff.Change, 0)
	enableChanges := make([]diff.Change, 0)

	// use the diff package to get the diff changes
	changes, diffErr := diff.Diff(jsonSettingsOld, jsonSettings)
	if diffErr != nil {
		logger.Warn("Failed to diff the json settings")
		return nil, nil, nil, diffErr
	}

	// loop through each change and based on the type and the changes from/to, put into appropriate array
	for _, changeRaw := range changes {
		// delete changes
		if changeRaw.To == nil {
			deletedChanges = append(deletedChanges, changeRaw)
		} else {
			// get from/to for enable/disable
			change, ok := changeRaw.To.(bool)
			if !ok {
				logger.Warn("Unsupported change type, ignorning\n")
				continue
			}
			oldVal, oldValOk := changeRaw.From.(bool)
			if !oldValOk {
				logger.Warn("Unsupported change type, ignoring\n")
				continue
			}
			// disable changes
			if changeRaw.Type == "update" && change == false && oldVal == true {
				disableChanges = append(disableChanges, changeRaw)
			}
			// enable changes
			if changeRaw.Type == "update" && change == true && oldVal == false {
				enableChanges = append(enableChanges, changeRaw)
			}
		}
	}

	// if deletes are accompanied by disable/enable changes, then something went wrong or the UI was changed
	if len(deletedChanges) > 0 && len(disableChanges) > 0 && len(enableChanges) > 0 {
		logger.Warn("This should not happen unless we start doing a soft delete\n")
		return nil, nil, nil, errors.New("Need to rethink deletes")
	}

	return deletedChanges, disableChanges, enableChanges, nil
}

// getSettingsErrorStruct converts the origErr error into a golang struct we can manipulate easily
// @param origErr error - original error from sync-settings
// @return *SetSettingsError - pointer to the SetSettingsError struct
// @return error - any error in getting struct, nil if none
func getSettingsErrorStruct(origErr error) (*SetSettingsError, error) {
	settingsError := &SetSettingsError{
		Confirm: Confirmation{},
	}
	err := json.Unmarshal([]byte(origErr.Error()), &settingsError)
	if err != nil {
		logger.Err("Couldn't get error: %s\n", err.Error())
		return nil, err
	}

	return settingsError, nil
}

// determineMessage determines the message for a single change set found - enable/disable/deleted changes
// @param settingsError *SetSettingsError - struct of respones from sync-settings
// @param invalidReason string - why sync-settings returned an error, could be enabled, deleted, or disabled
// @param buildFrom string - how to build the message. Build based on parents of items (for enabled changes) or child (for deleted/disabled)
// @param changes []diff.Change - array of changes to build message for
// @param jsonSettingsOld map[string]interface{} - current settings on system to get current values from
// @return SetSettingsErrorUI - UI message for the changes
// @return error - error from determining message, nil if none
func determineMessage(settingsError *SetSettingsError, invalidReason string, buildFrom string, changes []diff.Change, jsonSettingsOld map[string]interface{}) (SetSettingsErrorUI, error) {
	// create SetSettingsErrorUI and set the translation string
	newErr := SetSettingsErrorUI{}
	newErr.MainTranslationString = "affected_item_disabled_or_deleted"
	if invalidReason == "enabled" {
		newErr.MainTranslationString = "affected_item_on_enable"
	}
	newErr.InvalidReason = invalidReason
	newErr.AffectedValues = make([]AffectedValue, 0)

	for _, change := range changes {
		// goal from getAffectedItemId is to get the item object's id, like a policy or a rule
		pathForAffectedItem := change.Path // use path for deletes only
		// use all except last item for disable/enable
		if invalidReason == "disabled" || invalidReason == "enabled" {
			pathForAffectedItem = pathForAffectedItem[:len(pathForAffectedItem)-1]
		}
		affectedItemID, err := getAffectedItemID(jsonSettingsOld, pathForAffectedItem)
		if err != nil {
			return newErr, errors.New("Could not find affected item id")
		}

		// for each change build a message
		messages, buildErr := buildIndividualMessage(affectedItemID, invalidReason, buildFrom, settingsError)
		if buildErr != nil {
			logger.Warn("Failed to create whole individual message: %s\n", buildErr.Error())
			return newErr, buildErr
		}
		newErr.AffectedValues = append(newErr.AffectedValues, messages...)
	}

	// if no values found, sync-settings did not error because of something in this change, which is wrong
	if len(newErr.AffectedValues) == 0 {
		logger.Warn("No valid affected values found, erroring\n")
		return newErr, errors.New("No affected values found")
	}

	return newErr, nil
}

// getAffectedItemId gets the item object Id, like a policyId or a ruleId
// @param jsonSettingsOld map[string]interface{} - current settings on system to get current values from
// @param path []string - path of settings to get from
// @return string - id of the affected item
// @return error - error from getting item, nil if none
func getAffectedItemID(jsonSettingsOld map[string]interface{}, path []string) (string, error) {
	valueRaw, err := getSettingsFromJSON(jsonSettingsOld, path)
	if err != nil {
		logger.Warn("Failed to get settings items\n")
		return "", errors.New("Failed to get affected item")
	}

	value, ok := valueRaw.(map[string]interface{})
	if !ok {
		logger.Warn("Failed to get rule as could not convert to right type\n")
		return "", errors.New("Failed to get affected item")
	}

	return determineID(path, value)
}

// buildIndividualMessage builds the message's affected values from a single change i.e. disabling interface would return enabled policies/rules
// @param id string - id of object of the change
// @param invalidReason string - what the change is attempting to do, like enable/delete/disabled
// @param buildFrom string - how to build the message. Build based on parents of items (for enabled changes) or child (for deleted/disabled)
// @param settingsError *SetSettingsError - settings error struct from sync settings
func buildIndividualMessage(id string, invalidReason string, buildFrom string, settingsError *SetSettingsError) ([]AffectedValue, error) {
	affectedValues := make([]AffectedValue, 0)

	// determine the nextId to use, the child/parent
	nextIds := determineNextID(settingsError.Confirm.InvalidItems, id, buildFrom)

	for _, nextID := range nextIds {
		if len(nextID) <= 0 {
			continue
		}

		// get the child/parent to add next to affectedValues
		invalidItemToAdd, found := settingsError.Confirm.InvalidItems[nextID]
		if !found {
			logger.Warn("Could not find invalid id: %s\n", nextID)
			return affectedValues, errors.New("Could not find invalid item id")
		}

		affectedValue := AffectedValue{
			AffectedType:  invalidItemToAdd.Type,
			AffectedValue: invalidItemToAdd.Value,
		}
		affectedValues = append(affectedValues, affectedValue)

		// looks for children/parents of the newly added affected value and adds them
		moreValues, buildErr := buildIndividualMessage(nextID, invalidReason, buildFrom, settingsError)
		if buildErr != nil {
			logger.Warn("Failed to create whole individual message: %s\n", buildErr.Error())
			return affectedValues, buildErr
		}
		affectedValues = append(affectedValues, moreValues...)
	}

	return affectedValues, nil
}

// determineNextID determines the next ids to look at, whether the childIds or parentId
// @param allInvalidItems map[string]InvalidItem - map [invalidItemId] invalidItem - contains all elements affected by the current operation
// @param currentId string - the ID of the element we are currently processing
// @param buildFrom string - how to build the message. Build based on parents of items (for enabled changes) or child (for deleted/disabled)
// @return []string - next ids of invalidItems to be processed
func determineNextID(allInvalidItems map[string]InvalidItem, currentId string, buildFrom string) []string {
	invalidItem, ok := allInvalidItems[currentId]
	if !ok {
		return make([]string, 0, 0)
	}

	if buildFrom == "child" { //disable/delete
		// return []string{invalidItem.ChildID}
		childIds := []string{}
		for key, val := range allInvalidItems {
			if val.ParentID == currentId {
				childIds = append(childIds, key)
			}
		}
		return childIds
	}
	if buildFrom == "parent" { // enable
		return []string{invalidItem.ParentID}
	}

	return make([]string, 0, 0)
}

// determineID determines the id of a given affectedItem
// @param path []string - path of change to look at, determines what type of object we're changing
// @param value map[string]interface{} - object of affectedItem
// @return string - id of affected item
// @return error - any error, nil if none
func determineID(path []string, value map[string]interface{}) (string, error) {
	var rawID interface{}
	var found bool = false
	if path[0] == "wan" && path[1] == "policy_chains" && path[3] == "rules" { // rules
		rawID, found = value["ruleId"]
	} else if path[0] == "wan" && path[1] == "policies" { // policies
		rawID, found = value["policyId"]
	} else if path[0] == "network" && path[1] == "interfaces" { // interfaces
		rawID, found = value["interfaceId"]
	}

	if !found {
		logger.Warn("Failed to find given changed id: %v\n", rawID)
		return "", errors.New("Failed to find changed id")
	}

	// string for policies/rules
	id, ok := rawID.(string)
	if !ok {
		// float64 for interfaces
		idFloat, ok := rawID.(float64)
		if !ok {
			logger.Warn("Failed to get id\n")
			return "", errors.New("Failed to get id")
		}
		id = strconv.FormatFloat(idFloat, 'f', -1, 64)
	}

	return id, nil
}
