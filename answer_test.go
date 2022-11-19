package main

import (
	"testing"

	fuzz "github.com/AdaLogics/go-fuzz-headers"
)

var counter int

func init() {
	counter = 0
}

func Fuzz(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		fuzzConsumer := fuzz.NewConsumer(data)
		targetStruct := &CreateCompletionResponse{}
		err := fuzzConsumer.GenerateStruct(targetStruct)
		if err != nil {
			return
		}
		if targetStruct.Choices != nil {
			return
		}
		counter++
		if counter == 10000 {
			t.Fatalf("%+v\n", targetStruct)
		}
	})
}
