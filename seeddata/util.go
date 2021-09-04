package seeddata

func randomString(n int) string {
	letters := []rune("bcdefghjklmnopqrstvwxyz0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[randomGenerator.Intn(len(letters))]
	}
	return string(s)
}

func has(value string, list []string) bool {
	for _, v := range list {
		if value == v {
			return true
		}
	}
	return false
}
