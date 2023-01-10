package utils

import (
	"context"
	"encoding/json"
	"google.golang.org/grpc/metadata"
)

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
