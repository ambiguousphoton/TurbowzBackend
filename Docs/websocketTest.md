user1id - 27
curl -X POST http://localhost:8100/authenticate \
-H "Content-Type: application/x-www-form-urlencoded" \
-d "user_handle=Aditya" \
-d "password=shinebright"  


response : {"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyNywidXNlcl9oYW5kbGUiOiJBZGl0eWEiLCJleHAiOjE3NjU2ODExMzV9.kjrXBMHJu5xqVB5Q5xiGKHJxiSO4rEoWASzCut-voNg","userID":"27"}

user2id - 98
curl -X POST http://localhost:8100/authenticate \
-H "Content-Type: application/x-www-form-urlencoded" \
-d "user_handle=hero" \
-d "password="  

user3id 
curl -X POST http://localhost:8100/create-new-account \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "user_handle=hero" \      
  -d "user_profile_name=honda" \  
  -d "userDescription=strong" \         
  -d "fromLocation=Jhalander" \
  -d "userDateOfBirth=1992-07-21" \
  -d "gender=Male" \
  -d "email=hero@example.com" \
  -d "phoneNumber=900009" \     
  -d "password="    


curl -X POST http://localhost:8100/create-new-account \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "user_handle=a50" \
  -d "user_profile_name=Govinda" \
  -d "userDescription=Lord of Mathura" \
  -d "fromLocation=Vrindavan" \
  -d "userDateOfBirth=1992-07-21" \
  -d "gender=Male" \
  -d "email=krishnajijio@example.com" \
  -d "phoneNumber=83888888880"  -d "password=ooo"
{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxMDcsInVzZXJfaGFuZGxlIjoiYTUwIiwiZXhwIjoxNzY2MTgzNzU3fQ.9Ql83jZ93pa0fVFZ9pmHOinOu7mnrhO4nN_oABTGod0"}

response: {"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo5OSwidXNlcl9oYW5kbGUiOiJoZXJvIiwiZXhwIjoxNzY1NjgzNDQ5fQ.Fkrj1LoFubq__0mqNNvcVTKo8YSfDZFH_KGMqLFN_yg","userID":"27"}




