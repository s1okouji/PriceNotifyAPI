package util

import "time"

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

func RegularRequest(f func()) {
	y, m, d := time.Now().AddDate(0, 0, 1).Date()
	date := time.Date(y, m, d, 10, 0, 0, 0, time.UTC)
	time.AfterFunc(time.Until(date), f)
}
