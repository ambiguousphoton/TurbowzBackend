#!/bin/bash


go run Services/ImageReturnService/ImageReturner.go         &
go run Services/ServerDataSearch/ServerDataSearch.go     &
go run Services/VideoMetaDataService/GetUpdateVideoMD.go    &
go run Services/UserData/UserDataService.go                 &
go run Services/ServerDataStream/ServerDataStream.go        &
go run Services/CommentService/CommentService.go            &
go run Services/ServerDataReceive/ServerDataReceive.go      &
go run Services/SocketConnectionService/SocketConnection.go &
go run Services/RecommendationService/Recommendation.go     &
go run Services/ActivityService/Activity.go                 &
go run Services/EcoDataGetUpdate/EcoDataService.go          &
go run Services/CommunicationService/Communicate.go         &
go run Services/FollowUserService/FollowUserService.go      &
go run Services/AdsAndRevenueService/Ads.go                 &
go run Services/TrendingService/TrendingTrigger.go          &
go run Services/EventService/EventService.go                &

wait
