package utils

import (
	"context"
	"encoding/json"
	"google.golang.org/grpc/metadata"
	"math"
	"math/rand"
	"strconv"
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

// CalculatePages function that takes in the total number of items and the number of items per page,
// and calculates the number of pages required to fit all the items.
// it returns the number of pages as an integer
func CalculatePages(total, size int) int {
	return int(math.Ceil(float64(total) / float64(size)))
}
