package dataplex_test

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	dataplex "github.com/hashicorp/terraform-provider-google/google/services/dataplex"
)

func TestNumberOfAspectsValidation(t *testing.T) {
	fieldName := "aspects"
	numbers_100 := make([]interface{}, 100)
	for i := 0; i < 100; i++ {
		numbers_100[i] = i
	}
	numbers_99 := make([]interface{}, 99)
	for i := 0; i < 99; i++ {
		numbers_99[i] = i
	}
	numbers_empty := make([]interface{}, 0)
	map_100 := make(map[string]interface{}, 100)
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("key%d", i)
		map_100[key] = i
	}
	map_99 := make(map[string]interface{}, 99)
	for i := 0; i < 99; i++ {
		key := fmt.Sprintf("key%d", i)
		map_99[key] = i
	}
	map_empty := make(map[string]interface{}, 0)

	testCases := []struct {
		name        string
		input       interface{}
		expectError bool
		errorMsg    string
	}{
		{"too many aspects in a slice", numbers_100, true, "The maximal number of aspects is 99."},
		{"max number of aspects in a slice", numbers_99, false, ""},
		{"min number of aspects in a slice", numbers_empty, false, ""},
		{"too many aspects in a map", map_100, true, "The maximal number of aspects is 99."},
		{"max number of aspects in a map", map_99, false, ""},
		{"min number of aspects in a map", map_empty, false, ""},
		{"a string is not a valid input", "xelpatad", true, "to be array"},
		{"nil is not a valid input", nil, true, "to be array"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, errors := dataplex.NumberOfAspectsValidation(tc.input, fieldName)
			hasError := len(errors) > 0

			if hasError != tc.expectError {
				t.Fatalf("%s: NumberOfAspectsValidation() error expectation mismatch: got error = %v (%v), want error = %v", tc.name, hasError, errors, tc.expectError)
			}

			if tc.expectError && tc.errorMsg != "" {
				found := false
				for _, err := range errors {
					if strings.Contains(err.Error(), tc.errorMsg) { // Check if error message contains the expected substring
						found = true
						break
					}
				}
				if !found {
					t.Errorf("%s: NumberOfAspectsValidation() expected error containing %q, but got: %v", tc.name, tc.errorMsg, errors)
				}
			}
		})
	}
}

func TestProjectNumberValidation(t *testing.T) {
	fieldName := "some_field"
	testCases := []struct {
		name        string
		input       interface{}
		expectError bool
		errorMsg    string
	}{
		{"valid input", "projects/1234567890/locations/us-central1", false, ""},
		{"valid input with only number", "projects/987/stuff", false, ""},
		{"valid input with trailing slash content", "projects/1/a/b/c", false, ""},
		{"valid input minimal", "projects/1/a", false, ""},
		{"invalid input trailing slash only", "projects/555/", true, "has an invalid format"},
		{"invalid type - int", 123, true, `to be string, but got int`},
		{"invalid type - nil", nil, true, `to be string, but got <nil>`},
		{"invalid format - missing 'projects/' prefix", "12345/locations/us", true, "has an invalid format"},
		{"invalid format - project number starts with 0", "projects/0123/data", true, "has an invalid format"},
		{"invalid format - no project number", "projects//data", true, "has an invalid format"},
		{"invalid format - letters instead of number", "projects/abc/data", true, "has an invalid format"},
		{"invalid format - missing content after number/", "projects/123", true, "has an invalid format"},
		{"invalid format - empty string", "", true, "has an invalid format"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, errors := dataplex.ProjectNumberValidation(tc.input, fieldName)
			hasError := len(errors) > 0

			if hasError != tc.expectError {
				t.Fatalf("%s: ProjectNumberValidation() error expectation mismatch: got error = %v (%v), want error = %v", tc.name, hasError, errors, tc.expectError)
			}

			if tc.expectError && tc.errorMsg != "" {
				found := false
				for _, err := range errors {
					if strings.Contains(err.Error(), tc.errorMsg) { // Check if error message contains the expected substring
						found = true
						break
					}
				}
				if !found {
					t.Errorf("%s: ProjectNumberValidation() expected error containing %q, but got: %v", tc.name, tc.errorMsg, errors)
				}
			}
		})
	}
}

func TestAspectProjectNumberValidation(t *testing.T) {
	fieldName := "some_field"
	testCases := []struct {
		name        string
		input       interface{}
		expectError bool
		errorMsg    string
	}{
		{"valid input", "1234567890.compute.googleapis.com/Disk", false, ""},
		{"valid input minimal", "1.a", false, ""},
		{"invalid input trailing dot only", "987.", true, "has an invalid format"},
		{"invalid type - int", 456, true, `to be string, but got int`},
		{"invalid type - nil", nil, true, `to be string, but got <nil>`},
		{"invalid format - missing number", ".compute.googleapis.com/Disk", true, "has an invalid format"},
		{"invalid format - number starts with 0", "0123.compute.googleapis.com/Disk", true, "has an invalid format"},
		{"invalid format - missing dot", "12345compute", true, "has an invalid format"},
		{"invalid format - letters instead of number", "abc.compute.googleapis.com/Disk", true, "has an invalid format"},
		{"invalid format - missing content after dot", "12345", true, "has an invalid format"},
		{"invalid format - empty string", "", true, "has an invalid format"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, errors := dataplex.AspectProjectNumberValidation(tc.input, fieldName)
			hasError := len(errors) > 0

			if hasError != tc.expectError {
				t.Fatalf("%s: AspectProjectNumberValidation() error expectation mismatch: got error = %v (%v), want error = %v", tc.name, hasError, errors, tc.expectError)
			}

			if tc.expectError && tc.errorMsg != "" {
				found := false
				for _, err := range errors {
					if strings.Contains(err.Error(), tc.errorMsg) { // Check if error message contains the expected substring
						found = true
						break
					}
				}
				if !found {
					t.Errorf("%s: AspectProjectNumberValidation() expected error containing %q, but got: %v", tc.name, tc.errorMsg, errors)
				}
			}
		})
	}
}

func TestFilterAspects(t *testing.T) {
	testCases := []struct {
		name            string
		aspectKeySet    map[string]struct{}
		resInput        map[string]interface{}
		expectedAspects map[string]interface{}
		expectError     bool
		errorMsg        string
	}{
		{"aspects key is absent", map[string]struct{}{"keep": {}}, map[string]interface{}{"otherKey": "value"}, nil, false, ""},
		{"aspects value is nil", map[string]struct{}{"keep": {}}, map[string]interface{}{"aspects": nil}, nil, false, ""},
		{"empty aspectKeySet", map[string]struct{}{}, map[string]interface{}{"aspects": map[string]interface{}{"one": map[string]interface{}{"data": 1}, "two": map[string]interface{}{"data": 2}}}, map[string]interface{}{}, false, ""},
		{"keep all aspects", map[string]struct{}{"one": {}, "two": {}}, map[string]interface{}{"aspects": map[string]interface{}{"one": map[string]interface{}{"data": 1}, "two": map[string]interface{}{"data": 2}}}, map[string]interface{}{"one": map[string]interface{}{"data": 1}, "two": map[string]interface{}{"data": 2}}, false, ""},
		{"keep some aspects", map[string]struct{}{"two": {}, "three_not_present": {}}, map[string]interface{}{"aspects": map[string]interface{}{"one": map[string]interface{}{"data": 1}, "two": map[string]interface{}{"data": 2}}}, map[string]interface{}{"two": map[string]interface{}{"data": 2}}, false, ""},
		{"input aspects map is empty", map[string]struct{}{"keep": {}}, map[string]interface{}{"aspects": map[string]interface{}{}}, map[string]interface{}{}, false, ""},
		{"aspects is wrong type", map[string]struct{}{"keep": {}}, map[string]interface{}{"aspects": "not a map"}, nil, true, "FilterAspects: 'aspects' field is not a map[string]interface{}, got string"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resCopy := deepCopyMap(tc.resInput)
			originalAspectsBeforeCall := deepCopyValue(resCopy["aspects"])

			err := dataplex.FilterAspects(tc.aspectKeySet, resCopy)

			if tc.expectError {
				if err == nil {
					t.Fatalf("%s: Expected an error, but got nil", tc.name)
				}
				if tc.errorMsg != "" && !strings.Contains(err.Error(), tc.errorMsg) {
					t.Errorf("%s: Expected error message containing %q, got %q", tc.name, tc.errorMsg, err.Error())
				}
				if !reflect.DeepEqual(resCopy["aspects"], originalAspectsBeforeCall) {
					t.Errorf("%s: resCopy['aspects'] was modified during error case.\nBefore: %#v\nAfter: %#v", tc.name, originalAspectsBeforeCall, resCopy["aspects"])
				}
				return
			}

			if err != nil {
				t.Fatalf("%s: Did not expect an error, but got: %v", tc.name, err)
			}

			actualAspectsRaw, aspectsKeyExists := resCopy["aspects"]

			if tc.expectedAspects == nil {
				if aspectsKeyExists && actualAspectsRaw != nil {
					if tc.name == "aspects key is absent" {
						if aspectsKeyExists {
							t.Errorf("%s: Expected 'aspects' key to be absent, but it exists with value: %v", tc.name, actualAspectsRaw)
						}
					} else {
						t.Errorf("%s: Expected 'aspects' value to be nil, but got: %v", tc.name, actualAspectsRaw)
					}
				}
				return
			}

			if !aspectsKeyExists {
				t.Fatalf("%s: Expected 'aspects' key to exist, but it was absent. Expected value: %#v", tc.name, tc.expectedAspects)
			}

			actualAspects, ok := actualAspectsRaw.(map[string]interface{})
			if !ok {
				t.Fatalf("%s: Expected 'aspects' to be a map[string]interface{}, but got %T. Value: %#v", tc.name, actualAspectsRaw, actualAspectsRaw)
			}

			if !reflect.DeepEqual(actualAspects, tc.expectedAspects) {
				t.Errorf("%s: FilterAspects() result mismatch:\ngot:  %#v\nwant: %#v", tc.name, actualAspects, tc.expectedAspects)
			}
		})
	}
}

func TestAddAspectsToSet(t *testing.T) {
	testCases := []struct {
		name         string
		initialSet   map[string]struct{}
		aspectsInput interface{}
		expectedSet  map[string]struct{}
		expectError  bool
		errorMsg     string
	}{
		{"add to empty set", map[string]struct{}{}, []interface{}{map[string]interface{}{"aspect_key": "key1"}, map[string]interface{}{"aspect_key": "key2"}}, map[string]struct{}{"key1": {}, "key2": {}}, false, ""},
		{"add to existing set", map[string]struct{}{"existing": {}}, []interface{}{map[string]interface{}{"aspect_key": "key1"}}, map[string]struct{}{"existing": {}, "key1": {}}, false, ""},
		{"add duplicate keys", map[string]struct{}{}, []interface{}{map[string]interface{}{"aspect_key": "key1"}, map[string]interface{}{"aspect_key": "key1"}, map[string]interface{}{"aspect_key": "key2"}}, map[string]struct{}{"key1": {}, "key2": {}}, false, ""},
		{"input aspects is empty slice", map[string]struct{}{"existing": {}}, []interface{}{}, map[string]struct{}{"existing": {}}, false, ""},
		{"input aspects is nil", map[string]struct{}{"original": {}}, nil, map[string]struct{}{"original": {}}, false, ""},
		{"input aspects is wrong type", map[string]struct{}{}, "not a slice", map[string]struct{}{}, true, "AddAspectsToSet: input 'aspects' is not a []interface{}, got string"},
		{"item in slice is not a map", map[string]struct{}{}, []interface{}{"not a map"}, map[string]struct{}{}, true, "AddAspectsToSet: item at index 0 is not a map[string]interface{}, got string"},
		{"item map missing aspect_key", map[string]struct{}{}, []interface{}{map[string]interface{}{"wrong_key": "key1"}}, map[string]struct{}{}, true, "AddAspectsToSet: 'aspect_key' not found in aspect item at index 0"},
		{"aspect_key is not a string", map[string]struct{}{}, []interface{}{map[string]interface{}{"aspect_key": 123}}, map[string]struct{}{}, true, "AddAspectsToSet: 'aspect_key' in item at index 0 is not a string, got int"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			currentSet := make(map[string]struct{})
			for k, v := range tc.initialSet {
				currentSet[k] = v
			}

			err := dataplex.AddAspectsToSet(currentSet, tc.aspectsInput)

			if tc.expectError {
				if err == nil {
					t.Fatalf("%s: Expected an error, but got nil", tc.name)
				}
				if tc.errorMsg != "" && !strings.Contains(err.Error(), tc.errorMsg) {
					t.Errorf("%s: Expected error message containing %q, got %q", tc.name, tc.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("%s: Did not expect an error, but got: %v", tc.name, err)
				}
				if !reflect.DeepEqual(currentSet, tc.expectedSet) {
					t.Errorf("%s: AddAspectsToSet() result mismatch:\ngot:  %v\nwant: %v", tc.name, currentSet, tc.expectedSet)
				}
			}
		})
	}
}

func TestInverseTransformAspects(t *testing.T) {
	testCases := []struct {
		name             string
		resInput         map[string]interface{}
		expectedAspects  []interface{}
		expectNilAspects bool
		expectError      bool
		errorMsg         string
	}{
		{"aspects key is absent", map[string]interface{}{"otherKey": "value"}, nil, true, false, ""},
		{"aspects value is nil", map[string]interface{}{"aspects": nil}, nil, true, false, ""},
		{"aspects is empty map", map[string]interface{}{"aspects": map[string]interface{}{}}, []interface{}{}, false, false, ""},
		{"aspects with one entry", map[string]interface{}{"aspects": map[string]interface{}{"key1": map[string]interface{}{"data": "value1"}}}, []interface{}{map[string]interface{}{"aspectKey": "key1", "aspect": map[string]interface{}{"data": "value1"}}}, false, false, ""},
		{"aspects with multiple entries", map[string]interface{}{"aspects": map[string]interface{}{"key2": map[string]interface{}{"data": "value2"}, "key1": map[string]interface{}{"data": "value1"}}}, []interface{}{map[string]interface{}{"aspectKey": "key1", "aspect": map[string]interface{}{"data": "value1"}}, map[string]interface{}{"aspectKey": "key2", "aspect": map[string]interface{}{"data": "value2"}}}, false, false, ""},
		{"aspects is wrong type (not map)", map[string]interface{}{"aspects": "not a map"}, nil, false, true, "InverseTransformAspects: 'aspects' field is not a map[string]interface{}, got string"},
		{"aspect value is not a map", map[string]interface{}{"aspects": map[string]interface{}{"key1": "not a map value"}}, nil, false, true, "InverseTransformAspects: value for key 'key1' is not a map[string]interface{}, got string"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resCopy := deepCopyMap(tc.resInput)
			originalAspectsBeforeCall := deepCopyValue(resCopy["aspects"])

			err := dataplex.InverseTransformAspects(resCopy)

			if tc.expectError {
				if err == nil {
					t.Fatalf("%s: Expected an error, but got nil", tc.name)
				}
				if tc.errorMsg != "" && !strings.Contains(err.Error(), tc.errorMsg) {
					t.Errorf("%s: Expected error message containing %q, got %q", tc.name, tc.errorMsg, err.Error())
				}
				if !reflect.DeepEqual(resCopy["aspects"], originalAspectsBeforeCall) {
					t.Errorf("%s: resCopy['aspects'] was modified during error case.\nBefore: %#v\nAfter: %#v", tc.name, originalAspectsBeforeCall, resCopy["aspects"])
				}
				return
			}

			if err != nil {
				t.Fatalf("%s: Did not expect an error, but got: %v", tc.name, err)
			}

			actualAspectsRaw, aspectsKeyExists := resCopy["aspects"]

			if tc.expectNilAspects {
				if aspectsKeyExists && actualAspectsRaw != nil {
					t.Errorf("%s: Expected 'aspects' to be nil or absent, but got: %#v", tc.name, actualAspectsRaw)
				}
				return
			}

			if !aspectsKeyExists {
				t.Fatalf("%s: Expected 'aspects' key in result map, but it was missing. Expected value: %#v", tc.name, tc.expectedAspects)
			}
			if actualAspectsRaw == nil && tc.expectedAspects != nil {
				t.Fatalf("%s: Expected 'aspects' to be non-nil, but got nil. Expected value: %#v", tc.name, tc.expectedAspects)
			}

			actualAspectsSlice, ok := actualAspectsRaw.([]interface{})
			if !ok {
				if tc.expectedAspects != nil || actualAspectsRaw != nil {
					t.Fatalf("%s: Expected 'aspects' to be []interface{}, but got %T. Value: %#v", tc.name, actualAspectsRaw, actualAspectsRaw)
				}
			}

			if actualAspectsSlice != nil {
				sortAspectSlice(actualAspectsSlice)
			}
			if tc.expectedAspects != nil {
				sortAspectSlice(tc.expectedAspects)
			}

			if !reflect.DeepEqual(actualAspectsSlice, tc.expectedAspects) {
				t.Errorf("%s: InverseTransformAspects() result mismatch:\ngot:  %#v\nwant: %#v", tc.name, actualAspectsSlice, tc.expectedAspects)
			}
		})
	}
}

func TestTransformAspects(t *testing.T) {
	testCases := []struct {
		name             string
		objInput         map[string]interface{}
		expectedAspects  map[string]interface{}
		expectNilAspects bool
		expectError      bool
		errorMsg         string
	}{
		{"aspects key is absent", map[string]interface{}{"otherKey": "value"}, nil, true, false, ""},
		{"aspects value is nil", map[string]interface{}{"aspects": nil}, nil, true, false, ""},
		{"aspects is empty slice", map[string]interface{}{"aspects": []interface{}{}}, map[string]interface{}{}, false, false, ""},
		{"aspects with one item", map[string]interface{}{"aspects": []interface{}{map[string]interface{}{"aspectKey": "key1", "aspect": map[string]interface{}{"data": "value1"}}}}, map[string]interface{}{"key1": map[string]interface{}{"data": "value1"}}, false, false, ""},
		{"aspects with one item that has no aspect", map[string]interface{}{"aspects": []interface{}{map[string]interface{}{"aspectKey": "key1"}}}, map[string]interface{}{"key1": map[string]interface{}{"data": map[string]interface{}{}}}, false, false, ""},
		{"aspects with multiple items", map[string]interface{}{"aspects": []interface{}{map[string]interface{}{"aspectKey": "key1", "aspect": map[string]interface{}{"data": "value1"}}, map[string]interface{}{"aspectKey": "key2", "aspect": map[string]interface{}{"data": "value2"}}}}, map[string]interface{}{"key1": map[string]interface{}{"data": "value1"}, "key2": map[string]interface{}{"data": "value2"}}, false, false, ""},
		{"aspects with duplicate aspectKey", map[string]interface{}{"aspects": []interface{}{map[string]interface{}{"aspectKey": "key1", "aspect": map[string]interface{}{"data": "value_first"}}, map[string]interface{}{"aspectKey": "key2", "aspect": map[string]interface{}{"data": "value2"}}, map[string]interface{}{"aspectKey": "key1", "aspect": map[string]interface{}{"data": "value_last"}}}}, map[string]interface{}{"key1": map[string]interface{}{"data": "value_last"}, "key2": map[string]interface{}{"data": "value2"}}, false, false, ""},
		{"aspects is wrong type (not slice)", map[string]interface{}{"aspects": "not a slice"}, nil, false, true, "TransformAspects: 'aspects' field is not a []interface{}, got string"},
		{"item in slice is not a map", map[string]interface{}{"aspects": []interface{}{"not a map"}}, nil, false, true, "TransformAspects: item in 'aspects' slice at index 0 is not a map[string]interface{}, got string"},
		{"item map missing aspectKey", map[string]interface{}{"aspects": []interface{}{map[string]interface{}{"wrongKey": "k1", "aspect": map[string]interface{}{}}}}, nil, false, true, "TransformAspects: 'aspectKey' not found in aspect item at index 0"},
		{"aspectKey is not a string", map[string]interface{}{"aspects": []interface{}{map[string]interface{}{"aspectKey": 123, "aspect": map[string]interface{}{}}}}, nil, false, true, "TransformAspects: 'aspectKey' in item at index 0 is not a string, got int"},
		{"aspect is present but wrong type", map[string]interface{}{"aspects": []interface{}{map[string]interface{}{"aspectKey": "key1", "aspect": "not a map"}}}, map[string]interface{}{"key1": map[string]interface{}{"data": map[string]interface{}{}}}, false, false, ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			objCopy := deepCopyMap(tc.objInput)
			originalAspectsBeforeCall := deepCopyValue(objCopy["aspects"])

			err := dataplex.TransformAspects(objCopy)

			if tc.expectError {
				if err == nil {
					t.Fatalf("%s: Expected an error, but got nil", tc.name)
				}
				if tc.errorMsg != "" && !strings.Contains(err.Error(), tc.errorMsg) {
					t.Errorf("%s: Expected error message containing %q, got %q", tc.name, tc.errorMsg, err.Error())
				}
				if !reflect.DeepEqual(objCopy["aspects"], originalAspectsBeforeCall) {
					t.Errorf("%s: objCopy['aspects'] was modified during error case.\nBefore: %#v\nAfter: %#v", tc.name, originalAspectsBeforeCall, objCopy["aspects"])
				}
				return
			}

			if err != nil {
				t.Fatalf("%s: Did not expect an error, but got: %v", tc.name, err)
			}

			actualAspectsRaw, aspectsKeyExists := objCopy["aspects"]

			if tc.expectNilAspects {
				if aspectsKeyExists && actualAspectsRaw != nil {
					t.Errorf("%s: Expected 'aspects' to be nil or absent, but got: %#v", tc.name, actualAspectsRaw)
				}
				return
			}

			if !aspectsKeyExists {
				t.Fatalf("%s: Expected 'aspects' key in result map, but it was missing. Expected value: %#v", tc.name, tc.expectedAspects)
			}
			if actualAspectsRaw == nil && tc.expectedAspects != nil {
				t.Fatalf("%s: Expected 'aspects' to be non-nil, but got nil. Expected value: %#v", tc.name, tc.expectedAspects)
			}

			actualAspectsMap, ok := actualAspectsRaw.(map[string]interface{})
			if !ok {
				if tc.expectedAspects != nil || actualAspectsRaw != nil {
					t.Fatalf("%s: Expected 'aspects' to be map[string]interface{}, but got %T. Value: %#v", tc.name, actualAspectsRaw, actualAspectsRaw)
				}
			}

			if !reflect.DeepEqual(actualAspectsMap, tc.expectedAspects) {
				t.Errorf("%s: TransformAspects() result mismatch:\ngot:  %#v\nwant: %#v", tc.name, actualAspectsMap, tc.expectedAspects)
			}
		})
	}
}

func deepCopyMap(original map[string]interface{}) map[string]interface{} {
	if original == nil {
		return nil
	}
	copyMap := make(map[string]interface{}, len(original))
	for key, value := range original {
		copyMap[key] = deepCopyValue(value)
	}
	return copyMap
}

func deepCopySlice(original []interface{}) []interface{} {
	if original == nil {
		return nil
	}
	copySlice := make([]interface{}, len(original))
	for i, value := range original {
		copySlice[i] = deepCopyValue(value)
	}
	return copySlice
}

func deepCopyValue(value interface{}) interface{} {
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case map[string]interface{}:
		return deepCopyMap(v)
	case []interface{}:
		return deepCopySlice(v)
	default:
		return v
	}
}

func sortAspectSlice(slice []interface{}) {
	if slice == nil {
		return
	}
	sort.SliceStable(slice, func(i, j int) bool {
		mapI, okI := slice[i].(map[string]interface{})
		mapJ, okJ := slice[j].(map[string]interface{})
		if !okI || !okJ {
			return false
		}

		keyIRaw, keyIExists := mapI["aspectKey"]
		keyJRaw, keyJExists := mapJ["aspectKey"]

		if !keyIExists || !keyJExists {
			return false
		}

		keyI, okI := keyIRaw.(string)
		keyJ, okJ := keyJRaw.(string)
		if !okI || !okJ {
			return false
		}
		return keyI < keyJ
	})
}

func TestAccDataplexEntry_dataplexEntryUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_number": envvar.GetTestProjectNumberFromEnv(),
		"random_suffix":  acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataplexEntryDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataplexEntry_dataplexEntryFullUpdatePrepare(context),
			},
			{
				ResourceName:            "google_dataplex_entry.test_entry_full",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"aspects", "entry_group_id", "entry_id", "location"},
			},
			{
				Config: testAccDataplexEntry_dataplexEntryUpdate(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_dataplex_entry.test_entry_full", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_dataplex_entry.test_entry_full",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"aspects", "entry_group_id", "entry_id", "location"},
			},
		},
	})
}

func testAccDataplexEntry_dataplexEntryFullUpdatePrepare(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dataplex_aspect_type" "tf-test-aspect-type-full-one" {
  aspect_type_id         = "tf-test-aspect-type-full%{random_suffix}-one"
  location     = "us-central1"
  project      = "%{project_number}"

  metadata_template = <<EOF
{
  "name": "tf-test-template",
  "type": "record",
  "recordFields": [
    {
      "name": "type",
      "type": "enum",
      "annotations": {
        "displayName": "Type",
        "description": "Specifies the type of view represented by the entry."
      },
      "index": 1,
      "constraints": {
        "required": true
      },
      "enumValues": [
        {
          "name": "VIEW",
          "index": 1
        }
      ]
    }
  ]
}
EOF
}

resource "google_dataplex_aspect_type" "tf-test-aspect-type-full-two" {
  aspect_type_id         = "tf-test-aspect-type-full%{random_suffix}-two"
  location     = "us-central1"
  project      = "%{project_number}"

  metadata_template = <<EOF
{
  "name": "tf-test-template",
  "type": "record",
  "recordFields": [
    {
      "name": "story",
      "type": "enum",
      "annotations": {
        "displayName": "Story",
        "description": "Specifies the story of an entry."
      },
      "index": 1,
      "constraints": {
        "required": true
      },
      "enumValues": [
        {
          "name": "SEQUENCE",
          "index": 1
        },
        {
          "name": "DESERT_ISLAND",
          "index": 2
        }
      ]
    }
  ]
}
EOF
}

resource "google_dataplex_entry_group" "tf-test-entry-group-full" {
  entry_group_id = "tf-test-entry-group-full%{random_suffix}"
  project = "%{project_number}"
  location = "us-central1"
}

resource "google_dataplex_entry_type" "tf-test-entry-type-full" {
  entry_type_id = "tf-test-entry-type-full%{random_suffix}"
  project = "%{project_number}"
  location = "us-central1"

  required_aspects {
    type = google_dataplex_aspect_type.tf-test-aspect-type-full-one.name
  }
}

resource "google_dataplex_entry" "test_entry_full" {
  entry_group_id = google_dataplex_entry_group.tf-test-entry-group-full.entry_group_id
  project = "%{project_number}"
  location = "us-central1"
  entry_id = "tf-test-entry-full%{random_suffix}"
  entry_type = google_dataplex_entry_type.tf-test-entry-type-full.name
  fully_qualified_name = "bigquery:%{project_number}.test-dataset"
  parent_entry = "projects/%{project_number}/locations/us-central1/entryGroups/tf-test-entry-group-full%{random_suffix}/entries/some-other-entry"
  entry_source {
    resource = "bigquery:%{project_number}.test-dataset"
    system = "System III"
    platform = "BigQuery"
    display_name = "Human readable name"
    description = "Description from source system"
    labels = {
      "old-label": "old-value"
      "some-label": "some-value"
    }

    ancestors {
      name = "ancestor-one"
      type = "type-one"
    }

    ancestors {
      name = "ancestor-two"
      type = "type-two"
    }

    create_time = "2023-08-03T19:19:00.094Z"
    update_time = "2023-08-03T20:19:00.094Z"
  }

  aspects {
    aspect_key = "%{project_number}.us-central1.tf-test-aspect-type-full%{random_suffix}-one"
    aspect {
      data = <<EOF
          {"type": "VIEW"    }
        EOF
    }
  }

  aspects {
    aspect_key = "%{project_number}.us-central1.tf-test-aspect-type-full%{random_suffix}-two"
    aspect {
      data = <<EOF
          {"story": "SEQUENCE"    }
        EOF
    }
  }
 depends_on = [google_dataplex_aspect_type.tf-test-aspect-type-full-two, google_dataplex_aspect_type.tf-test-aspect-type-full-one]
}
`, context)
}

func testAccDataplexEntry_dataplexEntryUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dataplex_aspect_type" "tf-test-aspect-type-full-one" {
  aspect_type_id         = "tf-test-aspect-type-full%{random_suffix}-one"
  location     = "us-central1"
  project      = "%{project_number}"

  metadata_template = <<EOF
{
  "name": "tf-test-template",
  "type": "record",
  "recordFields": [
    {
      "name": "type",
      "type": "enum",
      "annotations": {
        "displayName": "Type",
        "description": "Specifies the type of view represented by the entry."
      },
      "index": 1,
      "constraints": {
        "required": true
      },
      "enumValues": [
        {
          "name": "VIEW",
          "index": 1
        }
      ]
    }
  ]
}
EOF
}

resource "google_dataplex_aspect_type" "tf-test-aspect-type-full-two" {
  aspect_type_id         = "tf-test-aspect-type-full%{random_suffix}-two"
  location     = "us-central1"
  project      = "%{project_number}"

  metadata_template = <<EOF
{
  "name": "tf-test-template",
  "type": "record",
  "recordFields": [
    {
      "name": "story",
      "type": "enum",
      "annotations": {
        "displayName": "Story",
        "description": "Specifies the story of an entry."
      },
      "index": 1,
      "constraints": {
        "required": true
      },
      "enumValues": [
        {
          "name": "SEQUENCE",
          "index": 1
        },
        {
          "name": "DESERT_ISLAND",
          "index": 2
        }
      ]
    }
  ]
}
EOF
}

resource "google_dataplex_entry_group" "tf-test-entry-group-full" {
  entry_group_id = "tf-test-entry-group-full%{random_suffix}"
  project = "%{project_number}"
  location = "us-central1"
}

resource "google_dataplex_entry_type" "tf-test-entry-type-full" {
  entry_type_id = "tf-test-entry-type-full%{random_suffix}"
  project = "%{project_number}"
  location = "us-central1"

  labels = { "tag": "test-tf" }
  display_name = "terraform entry type"
  description = "entry type created by Terraform"

  type_aliases = ["TABLE", "DATABASE"]
  platform = "GCS"
  system = "CloudSQL"

  required_aspects {
    type = google_dataplex_aspect_type.tf-test-aspect-type-full-one.name
  }
}

resource "google_dataplex_entry" "test_entry_full" {
  entry_group_id = google_dataplex_entry_group.tf-test-entry-group-full.entry_group_id
  project = "%{project_number}"
  location = "us-central1"
  entry_id = "tf-test-entry-full%{random_suffix}"
  entry_type = google_dataplex_entry_type.tf-test-entry-type-full.name
  fully_qualified_name = "bigquery:%{project_number}.test-dataset-modified"
  parent_entry = "projects/%{project_number}/locations/us-central1/entryGroups/tf-test-entry-group-full%{random_suffix}/entries/some-other-entry"
  entry_source {
    resource = "bigquery:%{project_number}.test-dataset-modified"
    system = "System III - modified"
    platform = "BigQuery-modified"
    display_name = "Human readable name-modified"
    description = "Description from source system-modified"
    labels = {
      "some-label": "some-value-modified"
      "new-label": "new-value"
    }

    ancestors {
      name = "ancestor-one"
      type = "type-one"
    }

    ancestors {
      name = "ancestor-two"
      type = "type-two"
    }

    create_time = "2024-08-03T19:19:00.094Z"
    update_time = "2024-08-03T20:19:00.094Z"
  }

  aspects {
    aspect_key = "%{project_number}.us-central1.tf-test-aspect-type-full%{random_suffix}-one"
    aspect {
      data = <<EOF
     {"type": "VIEW"    }
        EOF
    }
  }
 depends_on = [google_dataplex_aspect_type.tf-test-aspect-type-full-two, google_dataplex_aspect_type.tf-test-aspect-type-full-one]
}
`, context)
}
