package brute

func Search(input, substring string) []string {
	var matches []string

	for i := 0; i <= len(input)-len(substring); i++ {
		if input[i:i+len(substring)] == substring {
			matches = append(matches, input[i:i+len(substring)])
		}
	}

	return matches
}
