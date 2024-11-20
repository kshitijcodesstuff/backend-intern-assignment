package models

import (
	"testing"
)

func TestLoadStoreMaster(t *testing.T) {
	// Assuming StoreMaster.csv is in the correct location
	LoadStoreMaster("../StoreMaster.csv")
	if len(storeMaster) == 0 {
		t.Errorf("Expected storeMaster to be populated")
	}
}

func TestIsValidStore(t *testing.T) {
	LoadStoreMaster("../StoreMaster.csv")
	valid := IsValidStore("RP00001")
	if !valid {
		t.Errorf("Expected 'RP00001' to be a valid store ID")
	}
	invalid := IsValidStore("INVALID_ID")
	if invalid {
		t.Errorf("Expected 'INVALID_ID' to be invalid")
	}
}
