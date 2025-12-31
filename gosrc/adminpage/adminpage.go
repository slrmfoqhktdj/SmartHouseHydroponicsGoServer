package adminpage

import (
	"encoding/json"
	"log"
	"net/http"
	"smart_house/gosrc/structs"
)

func Login_Handler(w http.ResponseWriter, r *http.Request) {

	log.Print("Admin_Login_Handler")

	msg := structs.LoginData{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("Admin_Login_Handler err : %s", err)
		return
	}
	defer r.Body.Close()

	json, _ := json.Marshal(&structs.LoginPage{Items: Login_select(msg.Email, msg.PWD)})

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
