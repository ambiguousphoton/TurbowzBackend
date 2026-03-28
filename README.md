# Turbowz Backend

A microservices-based backend for a short-form video and social content platform. Built with Go and Python, using PostgreSQL for data storage and JWT for authentication.

## Overview

Turbowz is a platform for sharing:
- Short informational videos (up to 5 min)
- Ecos (image + text posts for sharing opinions)
- Events
- Real-time messaging via WebSockets

The backend is composed of independent Go microservices, each running on its own port, along with Python-based AI/ML inference services.

## Tech Stack

- **Go 1.24** — Core backend services
- **Python** — AI inference (CLIP embeddings, Whisper audio-to-text)
- **PostgreSQL** — Database
- **JWT** — Authentication
- **WebSockets** — Real-time messaging
- **FFmpeg** — Video/audio processing
- **HLS** — Video streaming

### Go Dependencies
- `github.com/golang-jwt/jwt/v5` — JWT auth
- `github.com/google/uuid` — UUID generation
- `github.com/gorilla/websocket` — WebSocket support
- `github.com/lib/pq` — PostgreSQL driver
- `golang.org/x/crypto` — Password hashing

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                     Client (Frontend)                   │
└──────────────┬──────────────────────────────────────────┘
               │
    ┌──────────▼──────────────────────────────────────┐
    │              Go Microservices                     │
    │                                                  │
    │  UserData ─ Auth, Profiles, Saves       (8100)   │
    │  ServerDataReceive ─ Uploads            (8080)   │
    │  ServerDataStream ─ Video Streaming     (8091)   │
    │  ServerDataSearch ─ Search              (8082)   │
    │  VideoMetaData ─ Video CRUD             (7999)   │
    │  EcoDataService ─ Eco CRUD              (7011)   │
    │  EventService ─ Events                  (7002)   │
    │  CommentService ─ Comments              (7200)   │
    │  FollowUserService ─ Follow/Unfollow    (8010)   │
    │  RecommendationService ─ Recs           (8007)   │
    │  ImageReturnService ─ Images/Thumbs     (8088)   │
    │  CommunicationService ─ WebSocket Chat  (8280)   │
    │  ActivityService ─ History/Analytics    (7992)   │
    │  AdsAndRevenueService ─ Banner Ads      (8991)   │
    │  TrendingService ─ Trending Trigger     (9090)   │
    │  TaskExecuterService ─ Background Tasks (7110)   │
    │  ConnectionService ─ User Connections   (8001)   │
    └──────────────┬──────────────────────────────────┘
                   │
    ┌──────────────▼──────────────────────────────────┐
    │           Python Inference Services               │
    │                                                  │
    │  CLIP Embedding (vectorize video/user)  (9000)   │
    │  Whisper Audio-to-Text                  (9018)   │
    └──────────────┬──────────────────────────────────┘
                   │
    ┌──────────────▼──────────┐
    │       PostgreSQL        │
    └─────────────────────────┘
```

## Services

| Service | Port | Description |
|---|---|---|
| UserData | 8100 | Account creation, authentication, profile management, saves |
| ServerDataReceive | 8080 | Video, eco, event, and profile photo uploads |
| ServerDataStream | 8091 | HLS video streaming |
| ServerDataSearch | 8082 | Search videos, ecos, and users by keyword |
| VideoMetaData | 7999 | Video metadata, views, likes, scores, trending, saved videos |
| EcoDataService | 7011 | Eco metadata, likes, scores, trending |
| EventService | 7002 | Event metadata and view counts |
| CommentService | 7200 | Push and get comments on videos |
| FollowUserService | 8010 | Follow/unfollow, follower/following lists |
| RecommendationService | 8007 | Video recommendations (content-based and user-based) |
| ImageReturnService | 8088 | Serve thumbnails, eco images, and profile photos |
| CommunicationService | 8280 | Real-time WebSocket messaging |
| ActivityService | 7992 | Watch history, upload analytics, voting |
| AdsAndRevenueService | 8991 | Banner ad upload and retrieval |
| TrendingService | 9090 | Trending content computation trigger |
| TaskExecuterService | 7110 | Background tasks (e.g., batch embedding updates) |
| ConnectionService | 8001 | User-to-user connections |
| Inference (Python) | 9000 | CLIP-based text/video vectorization |
| AudioToText (Python) | 9018 | Whisper-based audio transcription |

## Getting Started

### Prerequisites

- Go 1.24+
- Python 3.12+
- PostgreSQL
- FFmpeg (for video processing)

### Setup

1. Clone the repo:
   ```bash
   git clone https://github.com/ambiguousphoton/TurbowzBackend.git
   cd TurbowzBackend
   ```

2. Set up the Python inference environment:
   ```bash
   cd Services/InferenceService
   python -m venv inv
   source inv/bin/activate
   pip install -r requirements.txt  # install dependencies (torch, whisper, fastapi, etc.)
   ```

3. Configure your PostgreSQL connection in `repository/db.go`.

### Running the Services

**Python Inference Services:**
```bash
cd Services/InferenceService
source inv/bin/activate
uvicorn Inference:app --host 0.0.0.0 --port 9000
uvicorn AudioToText:app --host 0.0.0.0 --port 9018
```

**Go Services:**
```bash
go run Services/UserData/UserDataService.go                 # :8100
go run Services/ServerDataReceive/ServerDataReceive.go      # :8080
go run Services/ServerDataStream/ServerDataStream.go        # :8091
go run Services/ServerDataSearch/ServerDataSearch.go        # :8082
go run Services/VideoMetaDataService/GetUpdateVideoMD.go    # :7999
go run Services/EcoDataGetUpdate/EcoDataService.go          # :7011
go run Services/EventService/EventService.go                # :7002
go run Services/CommentService/CommentService.go            # :7200
go run Services/FollowUserService/FollowUserService.go      # :8010
go run Services/RecommendationService/Recommendation.go     # :8007
go run Services/ImageReturnService/ImageReturner.go         # :8088
go run Services/CommunicationService/Communicate.go         # :8280
go run Services/ActivityService/Activity.go                 # :7992
go run Services/AdsAndRevenueService/Ads.go                 # :8991
go run Services/TrendingService/TrendingTrigger.go          # :9090
go run Services/TaskExecuterService/TaskExecuter.go         # :7110
go run Services/ConnectionService/ConnectionService.go      # :8001
go run Services/SocketConnectionService/SocketConnection.go # :8181
```

## API Reference

### Authentication

**Create Account**
```bash
curl -X POST http://localhost:8100/create-new-account \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "user_handle=<handle>" \
  -d "user_profile_name=<name>" \
  -d "email=<email>" \
  -d "phoneNumber=<phone>" \
  -d "password=<password>"
```
Optional fields: `userDescription`, `fromLocation`, `gender`, `userDateOfBirth`  
Unique constraints: `phoneNumber`, `email`, `user_handle`

**Login**
```bash
curl -X POST http://localhost:8100/authenticate \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "user_handle=<handle>" \
  -d "password=<password>"
```
Returns: `{"token": "<jwt>", "userID": "<id>"}`

### Videos

| Endpoint | Method | Port | Description |
|---|---|---|---|
| `/upload` | POST | 8080 | Upload video (multipart: `video`, `title`, `info`, `tags`, `user_name`) |
| `/get-video-stream/{id}` | GET | 8091 | Stream video (HLS) |
| `/vmd?video_id=` | GET | 7999 | Get video metadata |
| `/view` | POST | 7999 | Record a view |
| `/luv` | POST | 7999 | Like/unlike a video |
| `/get-videos-score?video_id=` | GET | 7999 | Get video quality/AI scores |
| `/get-trending-videos?limit=&offset=&userID=` | GET | 7999 | Trending videos |
| `/get-saved-videos?limit=&offset=` | GET | 7999 | User's saved videos (auth required) |
| `/search?keyword=&limit=&offset=` | GET | 8082 | Search videos |
| `/search-video-with?userID=` | GET | 8082 | Get videos by user |

### Ecos (Image + Text Posts)

| Endpoint | Method | Port | Description |
|---|---|---|---|
| `/eco-upload` | POST | 8080 | Upload eco (multipart: `eco_text`, `tags`, `images`) |
| `/emd?eco_id=` | GET | 7011 | Get eco metadata |
| `/luv` | POST | 7011 | Like/unlike an eco |
| `/get-echo-score?echo_id=` | GET | 7011 | Get eco quality/AI scores |
| `/get-trending-ecos?limit=&offset=` | GET | 7011 | Trending ecos |
| `/search-eco-by-user?userID=` | GET | 8082 | Get ecos by user |

### Events

| Endpoint | Method | Port | Description |
|---|---|---|---|
| `/event-upload` | POST | 8080 | Create event (multipart: `event_title`, `event_description`, `event_start_time`, `event_end_time`, `tags`, `images`) |
| `/event-md?event_id=` | GET | 7002 | Get event metadata |
| `/increment-event-viewcount?event_id=` | GET | 7002 | Increment event views |

### Social

| Endpoint | Method | Port | Description |
|---|---|---|---|
| `/follow` | POST | 8010 | Follow a user (`followeeID`) |
| `/unfollow` | POST | 8010 | Unfollow a user (`followeeID`) |
| `/get-followers?checkID=` | GET | 8010 | Get followers list |
| `/get-followees?checkID=` | GET | 8010 | Get following list |
| `/get-following-info?userID=&requesterID=` | GET | 8010 | Follower/following counts + follow status |
| `/push-comment` | POST | 7200 | Post a comment (`parentVideoID`, `commentText`) |
| `/get-comment?videoID=&limit=&offset=` | GET | 7200 | Get comments on a video |
| `/add-connection` | POST | 8001 | Add user connection (`contactID`) |

### User

| Endpoint | Method | Port | Description |
|---|---|---|---|
| `/get-user?userID=` | GET | 8100 | Get user info |
| `/search-users?keyword=` | GET | 8100 | Search users |
| `/update-profile` | POST | 8100 | Update profile fields (auth required) |
| `/save-video?videoID=` | POST | 8100 | Save/unsave a video |
| `/save-eco?ecoID=` | POST | 8100 | Save/unsave an eco |
| `/video-saved-status?videoID=` | POST | 8100 | Check video save status |
| `/eco-saved-status?ecoID=` | POST | 8100 | Check eco save status |
| `/get-turbomax-status?userID=` | GET | 8100 | Check Turbomax subscription |
| `/pfp-upload` | POST | 8080 | Upload profile photo |

### Media

| Endpoint | Method | Port | Description |
|---|---|---|---|
| `/i?img=<uuid>` | GET | 8088 | Get thumbnail image |
| `/e?img=<url>&index=` | GET | 8088 | Get eco image |
| `/pfp?user_id=` | GET | 8088 | Get profile photo |

### Recommendations & Activity

| Endpoint | Method | Port | Description |
|---|---|---|---|
| `/recommend?video_id=&page=&limit=` | GET | 8007 | Related videos |
| `/recommend-videos-for-user?user_id=&page=&limit=` | GET | 8007 | Personalized recommendations |
| `/get-user-watch-history?page=&limit=` | GET | 7992 | Watch history (auth required) |
| `/delete-my-history?userID=` | GET | 7992 | Delete watch history |
| `/get-activity-data?userID=` | GET | 7992 | Upload analytics |
| `/post-video-vote` | POST | 7992 | Vote on video quality/AI usage |
| `/post-echo-vote` | POST | 7992 | Vote on eco quality/AI usage |

### AI Inference (Python)

| Endpoint | Method | Port | Description |
|---|---|---|---|
| `/vectorize-video/` | POST | 9000 | Generate CLIP embeddings for a video |
| `/vectorize-user/` | POST | 9000 | Generate CLIP embeddings for a user profile |
| `/audio-to-text/` | POST | 9018 | Transcribe audio using Whisper |

### WebSocket Messaging

Connect to `ws://localhost:8280/connect-with-socket-server` with `Authorization` header.

Send:
```json
{
  "destinationID": "<user_id>",
  "messageText": "Hello!",
  "links": ""
}
```

### Ads

| Endpoint | Method | Port | Description |
|---|---|---|---|
| `/upload-b-ads` | POST | 8991 | Upload banner ad (multipart: `title`, `redirect_url`, `image`) |
| `/get-b-ads?page=&limit=` | GET | 8991 | Get banner ads |

## Database Schema

### Core Tables
- `video_data` — UID, Title, Information, Date, UUID (uploader ID)
- Separate tables for views and likes
- User accounts with unique handle, email, phone

### Storage
- `MediaData/videos/` — Video files
- `MediaData/images/` — Image files

## Performance Notes

- Whisper (base, CPU, int8): ~42s latency for audio-to-text
- AV1 codec considered for future video compression optimization

## Project Structure

```
├── Services/
│   ├── ActivityService/          # Watch history, analytics, voting
│   ├── AdsAndRevenueService/     # Banner ads
│   ├── CommentService/           # Video comments
│   ├── CommunicationService/     # WebSocket messaging
│   ├── ConnectionService/        # User connections
│   ├── EcoDataGetUpdate/         # Eco CRUD
│   ├── EventService/             # Events
│   ├── FollowUserService/        # Follow system
│   ├── ImageReturnService/       # Image serving
│   ├── InferenceService/         # Python AI (CLIP, Whisper)
│   ├── RecommendationService/    # Content recommendations
│   ├── ServerDataReceive/        # Upload handling
│   ├── ServerDataSearch/         # Search
│   ├── ServerDataStream/         # Video streaming
│   ├── SocketConnectionService/  # Socket connections
│   ├── TaskExecuterService/      # Background tasks
│   ├── TrendingService/          # Trending computation
│   ├── UserData/                 # User management
│   └── VideoMetaDataService/     # Video metadata
├── authenticator/                # JWT auth middleware
├── models/                       # Data models
├── repository/                   # Database layer
├── Monitoring/                   # Latency monitoring
└── Docs/                         # Documentation
```
