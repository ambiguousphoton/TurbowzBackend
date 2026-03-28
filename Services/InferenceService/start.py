import subprocess
import platform

path = "/Users/vyoam/Desktop/Personal/GoServers"

files = [
    "services/ImageReturnService/Imagereturner.go",
    "services/ServerDataSearch/ServerDataSearch.go",
    "Services/VideoMetaDataService/GetUpdateVideoMD.go",
    "Services/UserData/UserDataService.go",
    "Services/ServerDataStream/ServerDataStream.go",
    "Services/CommentService/CommentService.go",
    "Services/ServerDataReceive/ServerDataReceive.go",
    "Services/SocketConnectionService/SocketConnection.go",
    "Services/RecommendationService/Recommendation.go",
    "Services/ActivityService/Activity.go",
    "Services/EcoDataGetUpdate/EcoDataService.go",
    "Services/CommunicationService/Communicate.go",
    "Services/FollowUserService/FollowUserService.go",
    "Services/AdsAndRevenueService/Ads.go",
    "Services/TrendingService/TrendingTrigger.go",
]

for f in files:
    subprocess.Popen([
        "open",
        "-a", "Terminal",
        "--args",
        "bash", "-c",
        f"cd '{path}' && go run '{f}'; exec bash"
    ])
