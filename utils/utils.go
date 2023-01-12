package utils

import (
	"context"
	"encoding/json"
	humanize "github.com/dustin/go-humanize"
	"google.golang.org/grpc/metadata"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// GenerateID based on current time
func GenerateID() int64 {
	return time.Now().UnixNano() + int64(rand.Intn(10000))
}

// ToByte converts any type to a byte slice.
func ToByte(i any) []byte {
	bt, _ := json.Marshal(i)
	return bt
}

// Dump to json using json marshal
func Dump(i any) string {
	return string(ToByte(i))
}

// DumpIncomingContext converts the metadata from the incoming context to a string representation using json marshal.
func DumpIncomingContext(c context.Context) string {
	md, _ := metadata.FromIncomingContext(c)
	return Dump(md)
}

// Offset to get offset from page and limit, min value for page = 1
func Offset(page, limit int) int {
	offset := (page - 1) * limit
	if offset < 0 {
		return 0
	}
	return offset
}

// StringToInt convert string to int
func StringToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

// StringToInt64 takes a string representation of a number and converts it to int64.
func StringToInt64(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return i
}

// Int64ToString takes an int64 number and converts it to a string representation
func Int64ToString(i int64) string {
	s := strconv.FormatInt(i, 10)
	return s
}

// CalculatePages function that takes in the total number of items and the number of items per page,
// and calculates the number of pages required to fit all the items.
// it returns the number of pages as an integer
func CalculatePages(total, size int) int {
	return int(math.Ceil(float64(total) / float64(size)))
}

// Int64ToRupiah takes an int64 amount and converts it to a string in Indonesian Rupiah format.
// This function uses humanize package to add a comma as thousand separator,
// and then replaces it with a dot to match the Rupiah format.
// It returns a string in the format "Rp<amount>"
func Int64ToRupiah(amount int64) string {
	humanizeValue := humanize.CommafWithDigits(float64(amount), 0)
	stringValue := strings.Replace(humanizeValue, ",", ".", -1)
	return "Rp" + stringValue
}

// FormatTimeRFC3339 Format time according to RFC3339Nano
func FormatTimeRFC3339(t *time.Time) (s string) {
	if t == nil {
		return
	}

	if t.Nanosecond() == 0 {
		return t.Format("2022-12-15T15:04:05.000000000Z07:00")
	}

	return t.Format(time.RFC3339Nano)
}
