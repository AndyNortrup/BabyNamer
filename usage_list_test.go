package main

import "testing"

func TestCombineUsageLists(t *testing.T) {

	origin1 := NameOrigin{Code: "code1", Plain: "Code1", Description: "First Code"}
	origin2 := NameOrigin{Code: "code2", Plain: "Code2", Description: "Second Code"}

	settingUsage := make(map[string]*SettingUsage)
	settingUsage[origin1.Code] = &SettingUsage{
		User:       "user@email.com",
		Enabled:    true,
		NameOrigin: origin1,
	}

	allUsages := make(map[string]NameOrigin)
	allUsages[origin1.Code] = origin1
	allUsages[origin2.Code] = origin2

	result := combineUsageLists(settingUsage, allUsages, "user@email.com")

	if result[0].Enabled == false {
		t.Log("Failed to take setting from user settings.")
		t.Fail()
	}

	if result[1].Enabled == true {
		t.Logf("Failed to generate correct user setting from allUsages:\n%v", result[1])
		t.Fail()
	}
}
