package common

import (
	"testing"
)

func TestToken(t *testing.T) {
	token, err := Token(32)
	if err != nil {
		t.Fatal(err)
	}
	
	if len(token) != 32 {
		t.Fatal("Length should be 32")
	}
}

func BenchmarkToken(b *testing.B) {
	for i := 0; i < b.N; i++ { //use b.N for looping
		Token(32)
    }
}