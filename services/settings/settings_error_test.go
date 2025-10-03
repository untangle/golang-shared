package settings

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetermineNextIds(t *testing.T) {
	type testCase struct {
		name        string
		allItems    map[string]InvalidItem
		crtId       string
		buildFrom   string
		expectedIds []string // order is irrelevant, both expected and result arrays are sorted before asserts
	}

	testCases := []testCase{
		{
			name:        "from child id not found",
			allItems:    map[string]InvalidItem{},
			crtId:       "1",
			buildFrom:   "child",
			expectedIds: []string{},
		},
		{
			name:        "from parent id not found",
			allItems:    map[string]InvalidItem{},
			crtId:       "2",
			buildFrom:   "parent",
			expectedIds: []string{},
		},
		{
			name: "id present but build from something we did not expect",
			allItems: map[string]InvalidItem{
				"2": {},
				"3": {ParentID: "2"},
				"4": {ParentID: "3"},
			},
			crtId:       "3",
			buildFrom:   "ohLookATypoHowUnfortunate",
			expectedIds: []string{},
		},
		{
			name: "from parent id here but nothing else related",
			allItems: map[string]InvalidItem{
				"qwe": {},
				"asd": {},
			},
			crtId:       "qwe",
			buildFrom:   "parent",
			expectedIds: []string{},
		},
		{
			name: "from child id here but nothing else related",
			allItems: map[string]InvalidItem{
				"qwe": {},
				"asd": {},
			},
			crtId:       "asd",
			buildFrom:   "child",
			expectedIds: []string{},
		},
		{
			name: "from parent, crtId between parent and children",
			allItems: map[string]InvalidItem{
				"qwe": {},
				"asd": {ParentID: "qwe"},
				"zxc": {ParentID: "asd"},
				"cxz": {ParentID: "asd"},
			},
			crtId:       "asd",
			buildFrom:   "parent",
			expectedIds: []string{"qwe"},
		},
		{
			name: "from child, crtId between parent and children",
			allItems: map[string]InvalidItem{
				"qwe": {},
				"asd": {ParentID: "qwe"},
				"zxc": {ParentID: "asd"},
				"cxz": {ParentID: "asd"},
			},
			crtId:       "asd",
			buildFrom:   "child",
			expectedIds: []string{"zxc", "cxz"},
		},
		{
			name: "from parent, child has comma seperated parentIDs",
			allItems: map[string]InvalidItem{
				"childID1":  {ParentID: "parentID1,parentID2"},
				"parentID1": {ParentID: ""},
				"parentID2": {ParentID: ""},
			},
			crtId:       "childID1",
			buildFrom:   "parent",
			expectedIds: []string{"parentID1", "parentID2"},
		},
	}

	Startup()

	for _, test := range testCases {
		ids := determineNextIds(test.allItems, test.crtId, test.buildFrom)

		assert.Equal(t, len(test.expectedIds), len(ids), test.name)

		sort.Strings(ids)
		sort.Strings(test.expectedIds)

		for i := range ids {
			assert.Equal(t, test.expectedIds[i], ids[i], fmt.Sprintf("%s: %d", test.name, i))
		}
	}
}

func TestDetermineSettingsError(t *testing.T) {
	type testCase struct {
		Id             string            `json:"id"`
		InvalidReason  string            `json:"invalidReason"`
		BuildFrom      string            `json:"buildFrom"`
		SettingsError  *SetSettingsError `json:"settingsError"`
		ExpectError    bool              `json:"expectError,omitempty"`
		ExpectedResult []AffectedValue   `json:"expectedResult,omitempty"`
	}

	testFiles := []string{
		"./testdata/error_tests/sync-settings-error-1.json",
		"./testdata/error_tests/sync-settings-error-2.json",
		"./testdata/error_tests/sync-settings-error-3.json",
	}

	for _, testFile := range testFiles {
		raw, err := os.ReadFile(testFile)
		assert.Nil(t, err, fmt.Sprintf("Cannot read testFile %s", testFile))

		test := testCase{}

		err = json.Unmarshal(raw, &test)
		assert.Nil(t, err, fmt.Sprintf("Cannot unmarshal testFile %s", testFile))

		result, err := buildIndividualMessage(test.Id, test.InvalidReason, test.BuildFrom, test.SettingsError)

		if test.ExpectError {
			assert.NotNil(t, err, testFile)
			continue
		}

		assert.Nil(t, err, testFile)
		assert.Equal(t, len(test.ExpectedResult), len(result), testFile)

		getLessFunc := func(arr []AffectedValue) func(i, j int) bool {
			return func(i, j int) bool {
				if arr[i].AffectedType == arr[j].AffectedType {
					return arr[i].AffectedValue > arr[j].AffectedValue
				}
				return arr[i].AffectedType > arr[j].AffectedType
			}
		}

		sort.Slice(result, getLessFunc(result))
		sort.Slice(test.ExpectedResult, getLessFunc(test.ExpectedResult))

		for i := range test.ExpectedResult {
			assert.Equal(t, test.ExpectedResult[i].AffectedType, result[i].AffectedType, testFile)
			assert.Equal(t, test.ExpectedResult[i].AffectedValue, result[i].AffectedValue, testFile)
		}
	}
}

func TestDetermineChangeSet(t *testing.T) {
	// Sample old and new settings with multiple key-value pairs
	oldSettings := map[string]interface{}{
		"key1": true,
		"key2": false,
		"key3": "value",
		"key4": 42,
		"key5": []string{"a", "b", "c"},
	}

	newSettings := map[string]interface{}{
		"key2": true,
		"key3": "new_value",
		"key4": 100,
		"key5": []string{"x", "y", "z"},
		"key6": "new_key",
	}

	// Call the determineChangeSet function with the sample settings
	deletedChanges, disableChanges, enableChanges, err := determineChangeSet(oldSettings, newSettings)

	// Check the behavior and output
	assert.NoError(t, err, "Expected no error")

	// Validate the deleted changes
	assert.Equal(t, 1, len(deletedChanges), "Expected one deleted change")
	assert.Equal(t, "delete", deletedChanges[0].Type, "Expected deleted change type")
	assert.Equal(t, "key1", deletedChanges[0].Path[0], "Expected deleted field 'key6'")

	// Validate the disable changes
	assert.Equal(t, 0, len(disableChanges), "Expected one disable change")

	// Validate the enable changes
	assert.Equal(t, 1, len(enableChanges), "Expected one enable change")
	assert.Equal(t, "update", enableChanges[0].Type, "Expected enable change type")
	assert.Equal(t, "key2", enableChanges[0].Path[0], "Expected enable field 'key2'")
	assert.True(t, enableChanges[0].To.(bool), "Expected enable change 'To' value to be true")
	assert.False(t, enableChanges[0].From.(bool), "Expected enable change 'From' value to be false")
}
