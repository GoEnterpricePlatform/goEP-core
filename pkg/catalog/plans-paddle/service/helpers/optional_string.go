package helpers

func OptionalString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
