package helpers

import (
	"crypto/rand"
	"io"
)

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func GenerateNumber(max int) string {
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}
func GetEndRangeDate(month string) string {
	end := ""
	switch month {
	case "JAN":
		end = "31"
	case "FEB":
		end = "28"
	case "MAR":
		end = "31"
	case "APR":
		end = "30"
	case "MAY":
		end = "31"
	case "JUN":
		end = "30"
	case "JUL":
		end = "31"
	case "AUG":
		end = "31"
	case "SEP":
		end = "30"
	case "OCT":
		end = "31"
	case "NOV":
		end = "30"
	case "DEC":
		end = "31"
	}
	return end
}
func GetEndRangeDateTwo(month string) string {
	end := ""
	switch month {
	case "01":
		end = "31"
	case "02":
		end = "28"
	case "03":
		end = "31"
	case "04":
		end = "30"
	case "05":
		end = "31"
	case "06":
		end = "30"
	case "07":
		end = "31"
	case "08":
		end = "31"
	case "09":
		end = "30"
	case "10":
		end = "31"
	case "11":
		end = "30"
	case "12":
		end = "31"
	}
	return end
}
