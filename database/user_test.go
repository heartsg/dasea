package database

import (
	"testing"
)

func TestUser(t *testing.T) {
	InitUsers()
	user, err := CreateUser("abc@abc.com", "def", "manager")
	if err != nil {
		t.Fatal(err)
	}
	if user.Id != 2 {
		t.Fatal("User id not properly set after insert")
	}
	user2, err := GetUser(user.Id)
	if err != nil {
		t.Fatal(err)
	}
	if user2 == nil {
		t.Fatal("GetUser failed")
	}
	if user.Email != user2.Email {
		t.Fatalf("Email comparison failed: %s != %s", user.Email, user2.Email)
	}
	
	if user2.IsVerified {
		t.Fatal("Created user verified, which should not happen")
	}
	
	user2.Verified() 
	
	if !user2.IsVerified {
		t.Fatal("User should be verified")
	}
	
	user3, err := GetUserByEmail("abc@abc.com")
	if err != nil {
		t.Fatal(err)
	}
	if user3 == nil {
		t.Fatal("GetUserByEmail failed")
	}
	if user.Email != user3.Email {
		t.Fatalf("Email comparison failed: %s != %s", user.Email, user2.Email)
	}
	if !user3.IsVerified {
		t.Fatal("User should be verified")
	}
	
	if !user3.Validate("def") {
		t.Fatal("User password validation failed")
	}
	
	
	user4, err := GetUserByToken(user3.Token)
	if err != nil {
		t.Fatal(err)
	}
	if user4 == nil {
		t.Fatal("GetUserByToken failed")
	}
	if user.Email != user4.Email {
		t.Fatalf("Email comparison failed: %s != %s", user.Email, user2.Email)
	}
	if !user4.IsVerified {
		t.Fatal("User should be verified")
	}
	
	
	user.HardDelete()
	user5, err := GetUser(user.Id)
	if err == nil || user5 != nil {
		t.Fatal("user not deleted")
	}
}