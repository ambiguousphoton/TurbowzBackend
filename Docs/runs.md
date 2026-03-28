Activating RecommendationGenratorModelEnv in RecommendationVectors folder
<!-- python -m venv inv -->
source inv//bin/activate
uvicorn Inference:app --host 0.0.0.0 --port 9000
uvicorn AudioToText:app --host 0.0.0.0 --port 9018

Frontend 

- npx expo start


Backend Go

go run services/ImageReturnService/Imagereturner.go         (8088)
go run services/ServerDataSearch/ServerDataSearch.go        (8082)
go run Services/VideoMetaDataService/GetUpdateVideoMD.go    (7999)
go run Services/UserData/UserDataService.go                 (8100) 
go run Services/ServerDataStream/ServerDataStream.go        (8091)
go run Services/CommentService/CommentService.go            (7200)
go run Services/ServerDataReceive/ServerDataReceive.go      (8080)
go run Services/SocketConnectionService/SocketConnection.go (8181)
go run Services/RecommendationService/Recommendation.go     (8007)
go run Services/ActivityService/Activity.go                 (7992)
go run Services/EcoDataGetUpdate/EcoDataService.go          (7011)
go run Services/CommunicationService/Communicate.go         (8280)
go run Services/FollowUserService/FollowUserService.go      (8010)
go run Services/AdsAndRevenueService/Ads.go                 (8991)
go run Services/TrendingService/TrendingTrigger.go          (9090)
go run Services/EventService/EventService.go                (7002)





Task Executer
go run Services/TaskExecuterService/TaskExecuter.go         (7110)

Trigger Service
go run Services/TrendingService/TrendingTrigger.go             
