package course

func nullableUUID(value string) interface{} {
	if value == "" {
		return nil
	}
	return value
}
