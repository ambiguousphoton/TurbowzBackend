package models

type FollowData struct{
	FollowerCount     int64;
	FolloweeCount     int64;
	AlreadyFollowed   bool;
}