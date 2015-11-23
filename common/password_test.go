package common

import (
	"testing"
)

func TestSalt(t *testing.T) {
	salt, err := Salt()
	if err != nil {
		t.Fatal(err)
	}
	if len(salt) != saltSize {
		t.Fatalf("Salt size mismatched, expect %d, we get %d", saltSize, len(salt))
	}	
}

func TestEncrypt(t *testing.T) {
	salt, err := Salt()
	if err != nil {
		t.Fatal(err)
	}
	key1, err := Encrypt("abcde", salt)
	if err != nil {
		t.Fatal(err)
	}
	
	key2, err := Encrypt("abcde", salt)
	if err != nil {
		t.Fatal(err)
	}
	
	
	if len(key1) != len(key2) {
		t.Fatal("key1 and key2 should be same")
	}
	for i := 0; i < len(key1); i++ {
		if key1[i] != key2[i] {
			t.Fatal("key1 and key2 should be same")
		}
	}
}

func TestMd5(t *testing.T) {
	key1 := Md5Encrypt("abcde")
	if len(key1) != md5Size {
		t.Fatal("md5 failed")
	}
	
	key2 := Md5Encrypt("abcde")
	if len(key2) != md5Size {
		t.Fatal("md5 failed")
	}

	for i := 0; i < len(key1); i++ {
		if key1[i] != key2[i] {
			t.Fatal("key1 and key2 should be same")
		}
	}
}


func BenchmarkSalt(b *testing.B) {
	for i := 0; i < b.N; i++ { //use b.N for looping
		Salt()
    }
}

func BenchmarkEncrypt(b *testing.B) {
	for i := 0; i < b.N; i++ { //use b.N for looping
		salt, _ := Salt()
		Encrypt("abcdefghijklmn", salt)
    }	
}