package util

func GetData(str string) string {
	state := 0
	start := 0
	end := len(str) - 1
	for i := 0; i < len(str); i++ {
		if str[i] == '{' {
			if state == 1 {
				start = i
				break
			}
			state++
		} else if str[i] == '}' {
			break
		}
	}
	return str[start:end]
}
