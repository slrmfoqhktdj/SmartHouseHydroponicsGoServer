package clientpage

import (
	"encoding/json"
	"log"
	"net/http"
	"smart_house/gosrc/clientpagedbconn"
	"smart_house/gosrc/structs"
)

func Login_Handler(w http.ResponseWriter, r *http.Request) {

	log.Print("Login_Handler")

	msg := structs.LoginData{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("Login_Handler err : %s", err)
		return
	}
	defer r.Body.Close()

	json, _ := json.Marshal(&structs.LoginPage{Items: clientpagedbconn.Login_select(msg.Email, msg.PWD)})

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func Login_App_Handler(w http.ResponseWriter, r *http.Request) {

	log.Print("Login_App_Handler")

	msg := structs.LoginData{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("Login_App_Handler err : %s", err)
		return
	}
	defer r.Body.Close()

	json, _ := json.Marshal(&structs.LoginApp{Items: clientpagedbconn.LoginApp_select(msg.Email, msg.PWD)})

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func Login_Device_Handler(w http.ResponseWriter, r *http.Request) {

	log.Print("Login_Device_Handler")

	msg := structs.LoginDeviceData{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("Login_Device_Handler err : %s", err)
		return
	}
	defer r.Body.Close()

	json, _ := json.Marshal(&structs.LoginDevice{Items: clientpagedbconn.LoginDevice_select(msg.DeviceId)})

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func App_Device_Report_OneSignalID_Handler(w http.ResponseWriter, r *http.Request) {
	log.Print("App_Device_Report_OneSignalID_Handler")

	msg := structs.ReportOneSignalId{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("Login_Handler err : %s", err)
		return
	}
	defer r.Body.Close()

	rb := clientpagedbconn.AppUser_OneSignalID_update(msg.Email, msg.Os_Id)

	json, _ := json.Marshal(&structs.ReportOneSignalIdConf{Os_Id: msg.Os_Id, Result: rb})

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func UserAlarmList_Handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("UserAlarmList_Handler %v", r.URL.Query().Get("user"))

	rb := clientpagedbconn.UserAlarmList_select(r.URL.Query().Get("user"))

	log.Printf("UserAlarmList_Handler %v", rb)
	json, _ := json.Marshal(&structs.UserAlarmListConf{User: r.URL.Query().Get("user"), Items: rb})
	log.Printf("UserAlarmList_Handler")
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
