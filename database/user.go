package database

import (
	"time"
	"errors"
	"dasea/common"
)

var userTokenExpireDuration time.Duration
var userVerificationTokenExpireDuration time.Duration

func init() {
	userTokenExpireDuration = time.Duration(24*7) * time.Hour
	userVerificationTokenExpireDuration = time.Duration(24*7) * time.Hour
}

func GenerateUserToken() (string, error) {
	//retry 5 times before it fail
	for i := 0; i < 5; i++ {
		token, err := common.Token(16)
		if err != nil {
			return "", err
		}
		has, err := mysqlEngine.Get(&User{Token:token})
		if err != nil {
			return "", err
		}
		if (!has) {
			return token, nil
		}
	}
	return "", errors.New("Generate token failed")
}

func CreateUser(email string, password string, level string ) (*User, error) {
	salt, err := common.Salt()
	if err != nil {
		return nil, err
	}
	key, err := common.Encrypt(password, salt)
	if err != nil {
		return nil, err
	}
	token, err := GenerateUserToken()
	if err != nil {
		return nil, err
	}
	expireAt := time.Now().Add(userTokenExpireDuration)
	
	user := User{
		Email: email,
		Password: key,
		Salt: salt,
		Token: token,
		Level: level,
		IsVerified: false,
		TokenExpireAt: expireAt,
	}
	_, err = mysqlEngine.Insert(&user)
	if err != nil {
		return nil, err
	}
	
	verificationToken, err := common.Token(16)
	verificationExpireAt := time.Now().Add(userVerificationTokenExpireDuration)
	if err == nil {
		mysqlEngine.Insert(&UserVerification{UserId:user.Id, Token:verificationToken, TokenExpireAt:verificationExpireAt})
	}
	
	return &user, nil
}

func GetUser(id int64) (*User, error) {
	user := new(User)
	has, err := mysqlEngine.Id(id).Get(user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("Invalid user Id")
	}
	return user, nil
}

func GetUserByEmail(email string) (*User, error) {
	user := &User{Email:email}
	has, err := mysqlEngine.Get(user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("Invalid email")
	}
	return user, nil
}

func GetUserByToken(token string) (*User, error) {
	user := &User{Token:token}
	has, err := mysqlEngine.Get(user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("Token not found")
	}
	if time.Now().After(user.TokenExpireAt) {
		return nil, errors.New("Token expired")
	}
	return user, nil
}

func DeleteUser(id int64) {
	mysqlEngine.Delete(&User{Id:id})
}

//Mainly used to clear user verification table
func (user* User) Verified() {
	user.IsVerified = true;
	mysqlEngine.Id(user.Id).Cols("is_verified").Update(user)
	mysqlEngine.Delete(&UserVerification{UserId:user.Id})
}

func (user* User) Delete() {
	mysqlEngine.Id(user.Id).Delete(&User{})
	mysqlEngine.Delete(&UserVerification{UserId:user.Id})
	mysqlEngine.Delete(&UserPasswordReset{UserId:user.Id})
}

func (user *User) HardDelete() {
	mysqlEngine.Id(user.Id).Delete(&User{})
	mysqlEngine.Delete(&UserVerification{UserId:user.Id})
	mysqlEngine.Delete(&UserPasswordReset{UserId:user.Id})
	mysqlEngine.Id(user.Id).Unscoped().Delete(user)
}

func (user *User) Validate(password string) bool {
	key, err := common.Encrypt(password, user.Salt)
	if err != nil {
		return false
	}
	if len(key) != len(user.Password) {
		return false
	}
	for i := 0; i < len(key); i++ {
		if key[i] != user.Password[i] {
			return false
		}
	}
	return true
}
