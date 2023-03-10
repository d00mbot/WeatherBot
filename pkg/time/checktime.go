package time

func CheckTime(t string) string {
	manyTime := []string{"01", "02", "03", "04", "05", "06", "07", "08", "09", "10", "11", "12",
		"13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23", "00"}

	if timeContains(manyTime, t) {
		return t
	} else {
		return ""
	}
}

func timeContains(tt []string, t string) bool {
	for _, v := range tt {
		if v == t {
			return true
		}
	}
	return false
}
