package models



type UserAuth struct {
	AuthID             int64
	UserID             int64
	UserLoginAccount   string
	UserPhoneNumber    string
	UserHashedPassword string
	AccountCreatedAt   string
}