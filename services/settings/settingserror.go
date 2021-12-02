package settings

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/r3labs/diff/v2"
	"github.com/untangle/golang-shared/services/logger"
)

// SetSettingsError TODO
type SetSettingsError struct {
	Confirm Confirmation `json:"CONFIRM"`
}

// Confirmation TODO
type Confirmation struct {
	InvalidItems map[string]InvalidItem `json:"invalidItems"`
}

// InvalidItem TODO
type InvalidItem struct {
	ChildType  string `json:"childType"`
	ChildID    string `json:"childId"`
	Reason     string `json:"reason"`
	Type       string `json:"type"`
	Value      string `json:"value"`
	Link       string `json:"link"`
	ParentID   string `json:"parentId"`
	ParentType string `json:"parentType"`
}

// SetSettingsErrorUI todo
type SetSettingsErrorUI struct {
	MainTranslationString string          `json:"mainTranslationString"`
	InvalidReason         string          `json:"invalidReason"`
	AffectedValues        []AffectedValue `json:"affectedValues"`
}

// AffectedValue TODO
type AffectedValue struct {
	AffectedType  string `json:"affectedType"`
	AffectedValue string `json:"affectedValue"`
	AffectedLink  string `json:"affectedLink"`
}

// determineSetSettingsError determines the error to send back after sync-settings
func determineSetSettingsError(origErr error, output string, settingsFile string, jsonSettings map[string]interface{}) (interface{}, error) {
	var err error = nil
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

// todo
func buildMessage(jsonSettingsOld map[string]interface{}, jsonSettings map[string]interface{}, origErr error) (string, error) {
	deletedChanges, disableChanges, enableChanges, changeSetErr := determineChangeSet(jsonSettingsOld, jsonSettings, origErr)
	if changeSetErr != nil {
		return "", changeSetErr
	}

	settingsError, settingsErrorStructErr := getSettingsErrorStruct(origErr)
	if settingsErrorStructErr != nil {
		return "", settingsErrorStructErr
	}

	messages := make([]SetSettingsErrorUI, 0)
	if len(deletedChanges) > 0 {
		deleteMessage, deleteMsgErr := determineMessage(settingsError, "deleted", "child", deletedChanges, jsonSettingsOld)
		if deleteMsgErr != nil {
			return "", deleteMsgErr
		}
		messages = append(messages, deleteMessage)
	}

	if len(disableChanges) > 0 {
		disableMessage, deleteMsgErr := determineMessage(settingsError, "disabled", "child", disableChanges, jsonSettingsOld)
		if deleteMsgErr != nil {
			return "", deleteMsgErr
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

	bytes, bytesErr := json.Marshal(messages)
	if bytesErr != nil {
		logger.Warn("Failed to unmarshal messages: %s\n", bytesErr.Error())
		return "", bytesErr
	}

	return "CONFIRM: " + string(bytes), nil
}

// todo
func determineChangeSet(jsonSettingsOld map[string]interface{}, jsonSettings map[string]interface{}, origErr error) ([]diff.Change, []diff.Change, []diff.Change, error) {
	deletedChanges := make([]diff.Change, 0)
	disableChanges := make([]diff.Change, 0)
	enableChanges := make([]diff.Change, 0)

	changes, diffErr := diff.Diff(jsonSettingsOld, jsonSettings)
	if diffErr != nil {
		logger.Warn("Failed to diff the json settings")
		return nil, nil, nil, diffErr
	}

	for _, changeRaw := range changes {
		if changeRaw.To == nil {
			deletedChanges = append(deletedChanges, changeRaw)
		} else {
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
			if changeRaw.Type == "update" && change == false && oldVal == true {
				disableChanges = append(disableChanges, changeRaw)
			}
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

// todo
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

func determineMessage(settingsError *SetSettingsError, invalidReason string, buildFrom string, changes []diff.Change, jsonSettingsOld map[string]interface{}) (SetSettingsErrorUI, error) {
	// create delete message, don't need to correlate so just use the output
	newErr := SetSettingsErrorUI{}
	newErr.MainTranslationString = "affected_item_disabled_or_deleted"
	if invalidReason == "enabled" {
		newErr.MainTranslationString = "affected_item_on_enable"
	}
	newErr.InvalidReason = invalidReason
	newErr.AffectedValues = make([]AffectedValue, 0)

	for _, change := range changes {
		pathForAffectedItem := change.Path
		if invalidReason == "disabled" || invalidReason == "enabled" {
			pathForAffectedItem = pathForAffectedItem[:len(pathForAffectedItem)-1]
		}
		affectedItemID, err := getAffectedItemID(jsonSettingsOld, pathForAffectedItem)
		if err != nil {
			return newErr, errors.New("Could not find affected item id")
		}

		messages, buildErr := buildIndividualMessage(affectedItemID, invalidReason, buildFrom, settingsError)
		if buildErr != nil {
			logger.Warn("Failed to create whole individual message: %s\n", buildErr.Error())
			return newErr, buildErr
		}
		newErr.AffectedValues = append(newErr.AffectedValues, messages...)
	}

	if len(newErr.AffectedValues) == 0 {
		logger.Warn("No valid affected values found, erroring\n")
		return newErr, errors.New("No affected values found")
	}

	return newErr, nil
}

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

func buildIndividualMessage(id string, invalidReason string, buildFrom string, settingsError *SetSettingsError) ([]AffectedValue, error) {
	affectedValues := make([]AffectedValue, 0)

	invalidItem, found := settingsError.Confirm.InvalidItems[id]
	if !found {
		logger.Warn("Could not find changed id: %s\n", id)
		return affectedValues, errors.New("Could not find changed id")
	}

	nextID := determineNextID(invalidItem, buildFrom)

	if len(nextID) <= 0 {
		// done
		return affectedValues, nil
	}

	invalidItemToAdd, found := settingsError.Confirm.InvalidItems[nextID]
	if !found {
		logger.Warn("Could not find invalid id: %s\n", nextID)
		return affectedValues, errors.New("Could not find invalid item id")
	}

	affectedLink := invalidItemToAdd.Link
	if invalidReason == "deleted" || invalidReason == "disabled" {
		affectedLink = ""
	}
	affectedValue := AffectedValue{
		AffectedType:  invalidItemToAdd.Type,
		AffectedValue: invalidItemToAdd.Value,
		AffectedLink:  affectedLink,
	}
	affectedValues = append(affectedValues, affectedValue)

	moreValues, buildErr := buildIndividualMessage(nextID, invalidReason, buildFrom, settingsError)
	if buildErr != nil {
		logger.Warn("Failed to create whole individual message: %s\n", buildErr.Error())
		return affectedValues, buildErr
	}
	affectedValues = append(affectedValues, moreValues...)

	return affectedValues, nil
}

// determine
func determineNextID(invalidItem InvalidItem, buildFrom string) string {
	if buildFrom == "child" { //disable/delete
		return invalidItem.ChildID
	} else if buildFrom == "parent" { // enable
		return invalidItem.ParentID
	}
	return ""
}

// todo
func determineID(path []string, value map[string]interface{}) (string, error) {
	var rawID interface{}
	var found bool = false
	if path[0] == "wan" && path[1] == "policy_chains" && path[3] == "rules" {
		rawID, found = value["ruleId"]
	} else if path[0] == "wan" && path[1] == "policies" {
		rawID, found = value["policyId"]
	} else if path[0] == "network" && path[1] == "interfaces" {
		rawID, found = value["interfaceId"]
	}

	if !found {
		logger.Warn("Failed to find given changed id: %v\n", rawID)
		return "", errors.New("Failed to find changed id")
	}

	id, ok := rawID.(string)
	if !ok {
		idFloat, ok := rawID.(float64)
		if !ok {
			logger.Warn("Failed to get id\n")
			return "", errors.New("Failed to get id")
		}
		id = strconv.FormatFloat(idFloat, 'f', -1, 64)
	}

	return id, nil
}
