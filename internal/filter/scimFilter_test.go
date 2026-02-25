// Copyright © 2026 Ping Identity Corporation

package filter

import (
	"testing"
)

func TestBuildScimFilter_SimpleFilterSet_Success(t *testing.T) {

	expectedSCIMFilter := `((name1 eq "value1") OR (name1 eq "value2") OR (name1 eq "VALUE3"))`

	filterSet := make([]interface{}, 0)

	filterSet = append(filterSet, map[string]interface{}{
		"name":   "name1",
		"values": []string{"value1", "value2", "VALUE3"},
	})

	if actualSCIMFilter := BuildScimFilter(filterSet, map[string]string{}); actualSCIMFilter != expectedSCIMFilter {
		t.Fatalf("\nExpected: \t%s\ngot:\t\t%s", expectedSCIMFilter, actualSCIMFilter)
	}

}

func TestBuildScimFilter_MultipleFilterSet1_Success(t *testing.T) {

	expectedSCIMFilter := `((name1 eq "value1") OR (name1 eq "value2") OR (name1 eq "VALUE3")) AND ((name2 eq "value1"))`

	filterSet := make([]interface{}, 0)

	filterSet = append(filterSet, map[string]interface{}{
		"name":   "name1",
		"values": []string{"value1", "value2", "VALUE3"},
	})

	filterSet = append(filterSet, map[string]interface{}{
		"name":   "name2",
		"values": []string{"value1"},
	})

	if actualSCIMFilter := BuildScimFilter(filterSet, map[string]string{}); actualSCIMFilter != expectedSCIMFilter {
		t.Fatalf("\nExpected: \t%s\ngot:\t\t%s", expectedSCIMFilter, actualSCIMFilter)
	}

}

func TestBuildScimFilter_MultipleFilterSet2_Success(t *testing.T) {

	expectedSCIMFilter := `((name1 eq "value1") OR (name1 eq "value2") OR (name1 eq "VALUE3")) AND ((name2 eq "value1") OR (name2 eq "value2"))`

	filterSet := make([]interface{}, 0)

	filterSet = append(filterSet, map[string]interface{}{
		"name":   "name1",
		"values": []string{"value1", "value2", "VALUE3"},
	})

	filterSet = append(filterSet, map[string]interface{}{
		"name":   "name2",
		"values": []string{"value1", "value2"},
	})

	if actualSCIMFilter := BuildScimFilter(filterSet, map[string]string{}); actualSCIMFilter != expectedSCIMFilter {
		t.Fatalf("\nExpected: \t%s\ngot:\t\t%s", expectedSCIMFilter, actualSCIMFilter)
	}

}

func TestBuildScimFilter_CustomFilterMap_Success(t *testing.T) {

	expectedSCIMFilter := `((name1 eq "value1") OR (name1 eq "value2") OR (name1 eq "VALUE3")) AND ((name2[id eq "value1"]) OR (name2[id eq "value2"]))`

	filterSet := make([]interface{}, 0)

	filterSet = append(filterSet, map[string]interface{}{
		"name":   "name1",
		"values": []string{"value1", "value2", "VALUE3"},
	})

	filterSet = append(filterSet, map[string]interface{}{
		"name":   "name2.id",
		"values": []string{"value1", "value2"},
	})

	customFilterMap := map[string]string{
		"name2.id": `name2[id eq "%s"]`,
	}

	if actualSCIMFilter := BuildScimFilter(filterSet, customFilterMap); actualSCIMFilter != expectedSCIMFilter {
		t.Fatalf("\nExpected: \t%s\ngot:\t\t%s", expectedSCIMFilter, actualSCIMFilter)
	}

}
