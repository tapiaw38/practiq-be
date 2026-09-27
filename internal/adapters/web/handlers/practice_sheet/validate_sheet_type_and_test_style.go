package practicesheet

func validateSheetTypeAndTestStyle(sheetType, testStyle string) string {
	if sheetType != "" && sheetType != "practice" && sheetType != "level_test" {
		return "sheet_type must be 'practice' or 'level_test'"
	}
	if testStyle != "" && testStyle != "keyboard" && testStyle != "canvas" {
		return "test_style must be 'keyboard' or 'canvas'"
	}
	return ""
}
