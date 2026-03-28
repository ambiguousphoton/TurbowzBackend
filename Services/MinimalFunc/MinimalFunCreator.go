// narayan narayan narayan narayan
package main

import (
	"GoServers/Services/MinimalFunc/Lib"
	"encoding/json"
	"net/http"
)

type CreateRequest struct {
	ID string `json:"id"`
	Code string `json:"code"`
}

func CreateMiniFunc(w http.ResponseWriter, r *http.Request){
	var req CreateRequest
	json.NewDecoder(r.Body).Decode(&req)
	store.Functions[req.ID] = store.Function{
		ID : req.ID,
		Code : req.Code,
	}
	w.Write([]byte("function stored"))
}

func ExecuteMiniFunc(w http.ResponseWriter, r *http.Request){
	
}


func main(){
	http.HandleFunc("createMiniFunc", CreateMiniFunc)
	http.HandleFunc("executeMiniFunc", ExecuteMiniFunc)
	http.ListenAndServe(":5454", nil)
}