#!/bin/bash

DATE=$(date +%Y-%m-%d)
mkdir -p logs/$DATE
find logs -type d -mtime +7 -exec rm -rf {} + 2>/dev/null

go run Services/ImageReturnService/ImageReturner.go         >> logs/$DATE/ImageReturnService.log 2>&1 &
go run Services/ServerDataSearch/ServerDataSearch.go        >> logs/$DATE/ServerDataSearch.log 2>&1 &
go run Services/VideoMetaDataService/GetUpdateVideoMD.go    >> logs/$DATE/VideoMetaData.log 2>&1 &
go run Services/UserData/UserDataService.go                 >> logs/$DATE/UserData.log 2>&1 &
go run Services/ServerDataStream/ServerDataStream.go        >> logs/$DATE/ServerDataStream.log 2>&1 &
go run Services/CommentService/CommentService.go            >> logs/$DATE/CommentService.log 2>&1 &
go run Services/ServerDataReceive/ServerDataReceive.go      >> logs/$DATE/ServerDataReceive.log 2>&1 &
go run Services/SocketConnectionService/SocketConnection.go >> logs/$DATE/SocketConnection.log 2>&1 &
go run Services/RecommendationService/Recommendation.go     >> logs/$DATE/RecommendationService.log 2>&1 &
go run Services/ActivityService/Activity.go                 >> logs/$DATE/ActivityService.log 2>&1 &
go run Services/EcoDataGetUpdate/EcoDataService.go          >> logs/$DATE/EcoDataService.log 2>&1 &
go run Services/CommunicationService/Communicate.go         >> logs/$DATE/CommunicationService.log 2>&1 &
go run Services/FollowUserService/FollowUserService.go      >> logs/$DATE/FollowUserService.log 2>&1 &
go run Services/AdsAndRevenueService/Ads.go                 >> logs/$DATE/AdsAndRevenueService.log 2>&1 &
go run Services/TrendingService/TrendingTrigger.go          >> logs/$DATE/TrendingService.log 2>&1 &
go run Services/EventService/EventService.go                >> logs/$DATE/EventService.log 2>&1 &

wait
