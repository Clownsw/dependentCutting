package util

// ReadAllLine 将content根据\n拆解
func ReadAllLine(content []byte) []string {
	var start = 0
	var result []string

	for i, c := range content {
		if c == '\n' {
			// start - i
			v := string(content[start:i])
			result = append(result, v)
			start = i
		}
	}

	return result
}
