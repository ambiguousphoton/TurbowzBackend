
Secure Curls

curl -X GET "http://localhost:7992/secured-get-user-watch-history?password=OMAngOMAngOMang&page=1&limit=20&userID=27"

curl "http://localhost:7110/update-user-embeddings?userID=*&OMAngOMAngOMang=OMAngOMAngOMang"

---> eifjcbrhrhkkjrtfnhtibcifcniithbfclhtvibebtgj
{"message":"All user embeddings updated successfully","count":43}






curl -X POST http://localhost:8991/upload-b-ads \
  -H "Authorization: " \
  -F "title=Sample Banner" \
  -F "redirect_url=https://example.com" \
  -F "image=@image.jpg"


fetch ads
curl "http://localhost:8991/get-b-ads?page=1&limit=10"
