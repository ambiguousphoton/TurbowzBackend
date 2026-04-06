
UserData/UserDataService 


Step 1: Verify Email (check availability + send OTP)

curl -X POST http://localhost:8100/verify-email \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "email=krishna@example.com"

Response: {"message": "verification code sent"}
Errors: 400 invalid email format, 409 email already registered


Step 2: Confirm Email OTP

curl -X POST http://localhost:8100/confirm-email \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "email=krishna@example.com" \
  -d "otp=482917"

Response: {"message": "email verified"}
Errors: 401 invalid or expired code


Step 3: Create Account

curl -X POST http://localhost:8100/create-new-account \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "user_handle=krishna_ji" \
  -d "user_profile_name=Govinda" \
  -d "userDescription=Lord of Mathura" \
  -d "fromLocation=Vrindavan" \
  -d "userDateOfBirth=1992-07-21" \
  -d "gender=Male" \
  -d "email=krishna@example.com" \
  -d "phoneNumber=888888888" \
  -d "password=FlutePower2024"

phoneNumber, email, user_handle -> must be unique
optional ->  fromlocation, userdescription, gender, dob

Validation rules:
  user_handle: 3-30 chars, alphanumeric/underscore only
  password: min 8 chars, must have uppercase + lowercase + digit
  phoneNumber: 7-15 digits, optional leading +
  email: valid email format
  user_profile_name: required, max 100 chars





updating user profile 

curl -X POST http://localhost:8100/update-profile \
  -H "Authorization: < Auth token > \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "userDescription=Jai Shree Ram" \
  -d "fromLocation=India" \
... any userfield can be updated


  .. // userid is needed for updating is obtained from auth token
  .. //all things except user_id can be updated individually or grouped






User Authentication

curl -X POST http://localhost:8100/authenticate \
-H "Content-Type: application/x-www-form-urlencoded" \
-d "user_handle=Aditya" \
-d "password=Shinebright1"

  this returns a jwt token
  Errors: 400 if user_handle or password is missing


Uploading video

curl -X POST http://localhost:8080/upload \
  -F "video=@HBchoubay.mp4" \
  -F "title=hosteler bday" \
  -F "info=om namah shivay" \
  -F "user_name=radha" \
  -F 'tags=["fun","college","hostel"]' \
  -H "Authorization: <auth token>"

  Response: {"video_id":123,"status":"processing"}
  Video is processed asynchronously (HLS encoding, thumbnail, vectorization).


Push Comment on a Video

curl -X POST http://localhost:7200/push-comment \
-d "parentVideoID=1" \
-d "commentText=jai ma bhavani" \
-d "parentCommentID=5" \
-H "Authorization: <auth token>"

Required: parentVideoID, commentText, Authorization header
Optional: parentCommentID (for replies to existing comments, omit for top-level comments)

Response: "Commented"


Get Comment on a Video

curl -X GET "http://localhost:7200/get-comment?videoID=46&limit=10&offset=20"



Push Comment on an Eco

curl -X POST http://localhost:7200/push-eco-comment \
-d "parentEcoID=1" \
-d "commentText=hari bol" \
-d "parentCommentID=3" \
-H "Authorization: <auth token>"

Required: parentEcoID, commentText, Authorization header
Optional: parentCommentID (for replies to existing comments, omit for top-level comments)

Response: "Eco Commented"


Get Comments on an Eco

curl -X GET "http://localhost:7200/get-eco-comment?ecoID=1&limit=10&offset=0"
response
{
  "comments": [
    {
      "Comment_id": 1,
      "Commenter_id": 27,
      "Parent_Eco_id": 0,
      "Comment_text": "hari bol",
      "Comment_date": "2025-01-23T10:00:00Z",
      "Commenter_Handle": "krishna_ji",
      "Commenter_Name": "Govinda",
      "Parent_Comment_ID": null
    },
    {
      "Comment_id": 3,
      "Commenter_id": 99,
      "Parent_Eco_id": 0,
      "Comment_text": "radhe radhe",
      "Comment_date": "2025-01-23T10:05:00Z",
      "Commenter_Handle": "hero",
      "Commenter_Name": "Surya",
      "Parent_Comment_ID": 1
    }
  ],
  "pagination": {
    "limit": 10,
    "offset": 0,
    "total": 2
  }
}



follow User

curl -X POST http://localhost:8010/follow \
-d "followeeID=4" \
-H "Authorization: <Auth token>" \




unfollow User

curl -X POST http://localhost:8010/unfollow \
-d "followeeID=4" \
-H "Authorization: <Auth token>" \




Get followers list

curl -X GET "http://localhost:8010/get-followers?checkID=1" 






Get a User's Following list

curl -X GET "http://localhost:8010/get-followees?checkID=27" 



Add connection // auth tokens user <-> contactId user

curl -X POST "http://localhost:8001/add-connection" -d "contactID=4" -H "Authorization: <auth token>"  





Search for Videos with Keyword
 curl -X GET "http://localhost:8082/search?keyword=om&limit=10&offset=0"
  -H "Accept: application/json"

Search for Videos with userId
curl -X GET "http://localhost:8082/search-video-with?userID=27" \
  -H "Accept: application/json"



Get Thumbnail/ Any Image
curl -v "http://localhost:8088/i?img=7dd9f567-27cd-4343-b5bc-8e77641c96bc" --output thumb1.jpg

--output will store it as thumb1.jpg




Get Profile Photo
curl -v "http://localhost:8088/i?img=p{user id}" 

the photo must be saved as "p{user id}.jpg"








Getting Video Meta Data
curl -X GET "http://localhost:7999/vmd?video_id=14"


Updating Views for Video
curl -X Post "http://localhost:7999/view"
-d "video_id=25"
-d "user_id=3"

Get Video 

curl -X GET "http://localhost:8091/get-video-stream/{VideoId}"



Get User Info With userID


curl -s "http://localhost:8100/get-user?userID=27"



Search Users With keyword


curl -s "http://localhost:8100/search-users?keyword=Best Friend of Krishna"          ! wil not work in terminal
/// curl -s "http://localhost:8100/search-users?keyword=Best%20Friend%20of%20Krishna"




Inference / Vectorise text data / Embedding 

curl -X POST "http://0.0.0.0:9000/vectorize-video/" \
-H "Content-Type: application/json" \
-d '{"title": "Sample Video", "description": "This is a test video description", "tags": ["test", "video", "fastapi"], "user_name": "vyoam", "video_id": 32}'



curl -X POST "http://localhost:8000/vectorize-user/" \
  -H "Content-Type: application/json" \
  -d '{
    "user_handle": "Radhe",
    "user_profile_name": "shyama ji",
    "user_description": "Braj ki pyari",
    "from_location": "Braj, UP",
    "user_date_of_birth": "1980-06-12",
    "gender": "female",
    "tags": ["developer", "AI", "tech"],
    "user_id": 1
  }'





AudioToText

curl -X POST "http://127.0.0.1:9018/audio-to-text/" -F "file=@audio.mp3"

Get Related Videos 

curl -X GET "http://localhost:8007/recommend?video_id=44&page=1&limit=5"



Luv a video

curl -X POST "http://localhost:7999/luv" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -H "Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyNywidXNlcl9oYW5kbGUiOiJBZGl0eWEiLCJleHAiOjE3NjUwNzA4ODV9.crWia0pV93P8F_hX3mveWaZW8OuyLQg-r5Ae0T2FXWk" \
  -d "video_id=36"
  
Response: : luv updated, current luvved status: false%  



Eco Post

curl -X POST http://localhost:8080/eco-upload \
  -H "Authorization: Auth" \
  -F "eco_text=Hariom" \                         
  -F "uploader_name=Aditya" \
  -F 'tags=["eco","nature","green"]' \
  -F "images=@dope.jpg"
  -F "images=@dope.jpg"

Search for Ecos with userId
curl -X GET "http://localhost:8082/search-eco-by-user?userID=27" \
  -H "Accept: application/json"

Get Eco Images
curl -X GET "http://localhost:8088/e?img=<eco url>&index=0" --output image.jpg


Get Profile Photo
curl -X GET "http://localhost:8088/pfp?user_id=27" --output image.jpg


Get User Watch History
curl -X GET "http://localhost:7992/get-user-watch-history?page=1&limit=10" \
  -H "Authorization: <auth token>"


Get Recommendations to User based on his Intrest
curl -X GET "http://localhost:8007/recommend-videos-for-user?user_id=27&page=1&limit=5


Luv a eco

curl -X POST "http://localhost:7011/luv" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -H "Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo5OSwidXNlcl9oYW5kbGUiOiJoZXJvIiwiZXhwIjoxNzY1NjgzNDQ5fQ.Fkrj1LoFubq__0mqNNvcVTKo8YSfDZFH_KGMqLFN_yg" \
  -d "eco_id=1"


Check Eco Luv Status

curl -X POST "http://localhost:7011/check-eco-luv-status" -d "eco_id=1" -d "user_ID=27"

Get Eco Meta Data

curl -X GET "http://localhost:7011/emd?eco_id=1"


Update Pfp

curl -X POST http://localhost:8080/pfp-upload \
  -H "Authorization: Auth" \
  -F "images=@dope.jpg"




Getting Following info 

curl -X GET "http://localhost:8010/get-following-info?userID=27&requesterID=27"


  response : {"FollowerCount":0,"FolloweeCount":5,"AlreadyFollowed":false}



User Video Save

curl -X POST "http://localhost:8100/save-video?videoID=12345" \
  -H "Authorization: <your_jwt_token>"


User Eco Save

curl -X POST "http://localhost:8100/save-eco?ecoID=98765" \
  -H "Authorization: <your_jwt_token>"


Saved Status of Eco

curl -X POST "http://localhost:8100/eco-saved-status?ecoID=1" -H "Authorization: Token"
{"saved":false}

Saved Status of Video

curl -X POST "http://localhost:8100/video-saved-status?videoID=1" -H "Authorization: Token"
{"saved":false}

Websocket 

ws://localhost:8280/connect-with-socket-server
Authorization : token
{
  "destinationID": "27",
  "messageText": "Hello there!",
  "links": ""
}

Other person gets
{"destinationID":"27","links":"","messageText":"Hello there!","roomID":"947778e2-bbe8-41f6-b23f-8000da51f5ac","sourceID":"99"}


Get Saved Videos


curl -X GET "http://localhost:7999/get-saved-videos?limit=10&offset=0" \
  -H "Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyNywidXNlcl9oYW5kbGUiOiJBZGl0eWEiLCJleHAiOjE3NjY3MzUxMDJ9.QZizg4NS0lNi7cJRvQtTVk9I4G8ERmi09mNJVxa_De0"



Get Trending Videos

curl -X GET "http://localhost:7999/get-trending-videos?limit=10&offset=0&userID=27" \



Get Trending Ecos

curl -X GET "http://localhost:7011/get-trending-ecos?limit=10&offset=0" 


Get Turbomax Subscription Status
curl -X GET "http://localhost:8100/get-turbomax-status?userID=27"
{"turbomax_active":true}


Delete History

curl -X GET "http://localhost:7992/delete-my-history?userID=27" \

"History Deleted"


Analytics of User Uploads

curl -X GET "http://localhost:7992/get-activity-data?userID=27" 

{"EcoUploads":{"2025-12-13T00:00:00Z":8},"VideoUploads":{"2025-12-05T00:00:00Z":2,"2025-12-07T00:00:00Z":2,"2025-12-10T00:00:00Z":1}}


Voting on Video by User

curl -X POST "http://localhost:7992/post-video-vote" \
-H "Content-Type: application/json" \
-d '{
  "video_id": 48,
  "user_id": 27,
  "quality": 4,
  "ai_usage": 2
}'



Voting on Eco by User

curl -X POST "http://localhost:7992/post-echo-vote" \
-H "Content-Type: application/json" \
-d '{
  "eco_id": 5,
  "user_id": 27,
  "quality": 4,
  "ai_usage": 2
}'



Create Event

curl -X POST "http://localhost:8080/event-upload" \
  -H "Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyNywidXNlcl9oYW5kbGUiOiJBZGl0eWEiLCJleHAiOjE3NzE3MjA5NDF9.7w0mAo7nD31KdbVWAr0mfk61fK4F1RV6RSogz3CLq3M" \
  -F "event_title=Tech Marathon" \
  -F "event_description=Annual Tech Conference 2024" \
  -F "event_start_time=2024-12-01T09:00:00Z" \
  -F "event_end_time=2024-12-01T18:00:00Z" \
  -F 'tags=["tech","anime"]' \
  -F "images=@image.jpg"


Get Event Meta Data

curl -X GET "http://localhost:7002/event-md?event_id=1"
<!-- {"Event_Id":1,"Event_Url":"50a28451-b12d-4b85-b784-e0a4497ccc13","Event_Title":"Untitled Event","Uploader_ID":27,"Uploader_Handle":"Aditya","Uploader_Name":"Surya","Event_Description":"Annual Tech Conference 2024","View_Count":0,"Luv_Count":0,"Comment_Count":0,"Tags":["tech","anime"],"Already_Luved":false,"Images_Count":0,"Saves_Count":0,"Created_At":"2026-01-23T09:34:09.476141+05:30","Event_Start_Time":"2024-12-01T14:30:00+05:30","Event_End_Time":"2024-12-01T23:30:00+05:30"} -->



Get Specific Event Data

curl -X GET "http://localhost:7002/event-md?event_id=1"




Increment Event Views Count

curl -X GET "http://localhost:7002/increment-event-viewcount?event_id=1"



Get Video Score

curl -X GET "http://localhost:7999/get-videos-score?video_id=14"

    response : {"Video_Quality":{"Float64":3,"Valid":true},"Video_AI_Usage":{"Float64":4,"Valid":true}, Total_Qualtiy_Votes, Total_Ai_Votes} 



Get Echo Score

curl -X GET "http://localhost:7011/get-echo-score?echo_id=5"

    Response : {"Echo_Quality":{"Float64":4,"Valid":true},"Echo_AI_Usage":{"Float64":2,"Valid":true},"Total_Qualtiy_Votes":1,"Total_Ai_Votes":1}