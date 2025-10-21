
UserData/UserDataService 

creating a new user

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
  -d "password=flutePower2024"

phoneNumber, email, user_handle -> must be unique
optional ->  fromlocation, userdescription, gender, dob





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
-d "password=shinebright" \

  this returns a jwt token




Uploading a Video

curl -X POST http://localhost:8080/upload \
-F "video=@HBchoubay.mp4" \
-F "title=hosteler bday" \
-F "info=om namah shivay" \
-H "Authorization: < auth token >" \


Push Comment on a Video

curl -X POST http://localhost:7200/push-comment \
-d "parentVideoID=1" \ 
-d "commentText=jai ma bhavani" \
-H "Authorization: <auth token>" \


Get Comment on a Video

curl -X GET "http://localhost:7200/get-comment?videoID=14"




follow User

curl -X POST http://localhost:8010/follow \
-d "followeeID=4" \
-H "Authorization: <Auth token>" \




unfollow User

curl -X POST http://localhost:8010/unfollow \
-d "followeeID=4" \
-H "Authorization: <Auth token>" \




Get followers list

curl -X GET "http://localhost:8010/get-followers?checkID=1" -H "Authorization: <Auth token>"






Get a User's Following list

curl -X GET "http://localhost:8010/get-followees?checkID=27" -H "Authorization: <Auth token>"



Add connection // auth tokens user <-> contactId user

curl -X POST "http://localhost:8001/add-connection" -d "contactID=4" -H "Authorization: <auth token>"  





Search for Videos with Keyword
curl -X GET "http://localhost:8082/search?keyword=om" \
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
-d "video_id=14"

Get Video 

curl -X GET "http://localhost:8091/get-video-stream/{VideoId}"



Get User Info With userID


curl -s "http://localhost:8100/get-user?userID=27"



Search Users With keyword


curl -s "http://localhost:8100/search-users?keyword=Best Friend of Krishna"          ! wil not work in terminal
/// curl -s "http://localhost:8100/search-users?keyword=Best%20Friend%20of%20Krishna"




Inference / Vectorise text data / Embedding 

curl -X POST "http://0.0.0.0:9000/vectorize/" \
-H "Content-Type: application/json" \
-d '{
  "title": "Sample Video",
  "description": "This is a test video description",
  "tags": ["test", "video", "fastapi"],
  "user_name": "vyoam"
}'

AudioToText

curl -X POST "http://127.0.0.1:9018/audio-to-text/" -F "file=@audio.mp3"


---------------------------------------------------------------------------------------------------

Websocket 


Connection Request
    
      {
        "Type": "Connection-Request",
        "ReceiverID": "95"
      }
      Header: 
        Autherisation: ""



-------

Connection Approval 

    {
        "Type": "Connection-Approval",
        "ReceiverID" : "27",
        "Status": "Accepted"
    }
        Header: 
          Autherisation: ""

--------

Message


    {
      "Type": "Message",
      "RoomID": "9ff1ea97-a619-47cb-bdc1-5b812ca654fc",
      "Text": "om om !"
    }
        Header: 
          Autherisation: ""

