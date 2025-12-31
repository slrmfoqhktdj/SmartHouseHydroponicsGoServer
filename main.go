package main

import (
	"bytes"
	"encoding/json"
	"encoding/pem"
	"crypto/x509"
	"fmt"
	"os"
	"runtime"
	"strings"

	// "io"
	"log"
	"math/rand"
	"net/http"

	// "os"

	"time"

	// "github.com/dgrijalva/jwt-go"
	// "google.golang.org/api/option"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"smart_house/gosrc/adminpage"

	"smart_house/gosrc/clientpage"
	"smart_house/gosrc/clientpagedbconn"
	"smart_house/gosrc/device_handler"
	jwt "smart_house/gosrc/jwt"
	"smart_house/gosrc/structs"

	// "github.com/rs/cors"
	"go.uber.org/zap"
	// firebase "firebase.google.com/go/v4"
	// "firebase.google.com/go/v4/messaging"
)

// var (
// 	fcmClient *messaging.Client
// 	ctx       context.Context = context.Background()
// )

// func fcmInit() {
// 	opt := option.WithCredentialsFile("./fcm_auth.json")
// 	app, err := firebase.NewApp(ctx, nil, opt)
// 	if err != nil {
// 		log.Fatalf("error initializing app: %v\n", err)
// 	}

// 	fcmClient, err = app.Messaging(ctx)
// 	if err != nil {
// 		log.Fatalf("error getting Messaging client: %v\n", err)
// 	}

// }

// func sendFCM(deviceid, tile, body string, data map[string]string) bool {
// 	message := &messaging.Message{
// 		Data:         data,
// 		Notification: &messaging.Notification{Title: tile, Body: body},
// 		// Topic:        deviceid,
// 		Android: &messaging.AndroidConfig{
// 			Priority: "high",
// 			Notification: &messaging.AndroidNotification{
// 				Title: tile,
// 				Body:  body,
// 				Color: "#ff0000",
// 				Sound: "default",
// 			},
// 		},
// 		Token: "fR_DsYgpT6SorRSWOZv-6G:APA91bFo6jux_KYS2Y9R1q7lrAZqZMDXGfiG5LjwC9caGqF2Madj2Kwrspffo64ICYfsj9qDp2KgASMdCjE7Ukccgnm-FFnV7N87-j2696NUbpMI7rKrOdA",
// 	}

// 	response, err := fcmClient.Send(ctx, message)
// 	if err != nil {
// 		log.Println("메시지 전송 실패 ", err)
// 		return false
// 	}

// 	log.Println("메시지 전송 성공 ", response)
// 	return true
// }

// fcmMsg := &messaging.Message{
// 	Data: map[string]string{
// 		"deviceid": req.Id,
// 		"time":     "2:45",
// 		"title":    "test",
// 		"body":     "testbody",
// 	},
// 	// Notification: &messaging.Notification{
// 	// 	Title: "test1",
// 	// 	Body:  "test1body",
// 	// },
// 	Android: &messaging.AndroidConfig{
// 		Notification: &messaging.AndroidNotification{
// 			Title: "test2",
// 			Body:  "test2body",
// 		},
// 		Data: map[string]string{
// 			"deviceid": req.Id,
// 			"time":     "2:45",
// 		},
// 		Priority: "high",
// 	},

// 	Topic: req.Id,
// 	// Token: "fLQG9KSPSqm7yuqOyWx4Fb:APA91bFzV155r8o6N8_r8wNy0YuaP1wGKEkH3DGcZGOTt0iE9RHZAefwxQOWyLa128pjVkNMl7xxm4FauHu4vmr9kCcQVaA0pHzYVEqx0dd9On1uKWFUKK4",
// }

// fcmRes, err := fcmClient.Send(ctx, fcmMsg)
// if err != nil {
// 	log.Print(err)
// }

// log.Print("Successfully sent message:", fcmRes)

func sendOneSignal(tile, body string, data map[string]string, deviceids []string) bool {
	// oneSignal api key
	// os_v2_app_7haiev6at5gavhyhtpadr4odsatd7w5r6izeutmopezwlz5yqwlh3kfmfuuyzhhlxeihqgxcwqsjzts4cgr2aqzsrjbmo4el6le56fi

	//app id
	//f9c08257-c09f-4c0a-9f07-9bc038f1c390

	d, _ := json.Marshal(&structs.SendPush{App_Id: os.Getenv("ONESIGNALAPPID"), IsAndroid: true, Include_Player_IDs: deviceids, Contents: structs.SendPushTitleAndContents{Ko: tile, En: tile}, Headings: structs.SendPushTitleAndContents{Ko: body, En: body}})

	buff := bytes.NewBuffer(d)

	req, err := http.NewRequest("POST", os.Getenv("ONESIGNALRESTFULURL"), buff)
	if err != nil {
		fmt.Println("푸시 에러:", err)
		return false
	}

	req.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Accept":        {"application/json"},
		"Authorization": {strings.Join([]string{"Bearer ", os.Getenv("ONESIGNALAPIKEY")}, " ")},
	}

	client := &http.Client{}

	// fmt.Println("요청 :", req)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("푸시 에러:", err)
		return false
	}

	// fmt.Println("푸시 :", resp.Status)

	if resp.StatusCode != http.StatusOK {
		fmt.Println("푸시 전송 실패:", resp)
		return false
	}

	defer resp.Body.Close()

	return true

}

func sendWholeOneSignal(tile, body string, data map[string]string) bool {
	// oneSignal api key
	// os_v2_app_7haiev6at5gavhyhtpadr4odsatd7w5r6izeutmopezwlz5yqwlh3kfmfuuyzhhlxeihqgxcwqsjzts4cgr2aqzsrjbmo4el6le56fi

	//app id
	//f9c08257-c09f-4c0a-9f07-9bc038f1c390

	d, _ := json.Marshal(&structs.SendPush{App_Id: os.Getenv("ONESIGNALAPPID"), IsAndroid: true, Contents: structs.SendPushTitleAndContents{Ko: tile, En: tile}, Headings: structs.SendPushTitleAndContents{Ko: body, En: body}})

	buff := bytes.NewBuffer(d)

	req, err := http.NewRequest("POST", os.Getenv("ONESIGNALRESTFULURL"), buff)
	if err != nil {
		fmt.Println("푸시 에러:", err)
		return false
	}

	req.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Accept":        {"application/json"},
		"Authorization": {strings.Join([]string{"Bearer ", os.Getenv("ONESIGNALAPIKEY")}, " ")},
	}

	client := &http.Client{}

	// fmt.Println("요청 :", req)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("푸시 에러:", err)
		return false
	}

	// fmt.Println("푸시 :", resp.Status)

	if resp.StatusCode != http.StatusOK {
		fmt.Println("푸시 전송 실패")
		return false
	}

	defer resp.Body.Close()

	return true

}

var logd *zap.SugaredLogger

var csms *Server

//////////////////////////////////////////////////////////////////////////////////////////////

// 대쉬보드 정보
func Dashboard_Page_Handler(w http.ResponseWriter, r *http.Request) {

	log.Print("Dashboard_Page_Handler,%v", r.URL.Query().Get("user"))

	// 한국 시간대 로드
	loc, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		fmt.Println("시간대 로딩 에러:", err)
		return
	}

	// 현재 시간을 한국 시간대로 설정
	now := time.Now().In(loc)

	Device_Id := "전체"
	StartTime := fmt.Sprintf("%s %s", now.AddDate(0, 0, -3).Format("2006-01-02"), "00:00:00")
	EndTime := fmt.Sprintf("%s %s", now.Format("2006-01-02"), "23:59:59")
	user := r.URL.Query().Get("user")
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(&structs.DashboardPage{List: clientpagedbconn.Device_list_select(ConnectedDeivceList_Handler(), user, "0", "3"), Appoint: Device_Web_List_Appoint_Handler(user), Daily: clientpagedbconn.Device_daily_list_select(user, "0", "3"), Graph: clientpagedbconn.Device_daily_graph_select(StartTime, EndTime, Device_Id)})
	w.Write(json)
}

// 하우스 기본 정보
func DeviceList_Page_Handler(w http.ResponseWriter, r *http.Request) {

	log.Printf("DeviceList_Page_Handler %v,%v,%v", r.URL.Query().Get("user"), r.URL.Query().Get("s_order"), r.URL.Query().Get("e_order"))

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(&structs.DeviceListPage{Count: clientpagedbconn.Device_list_count_select(r.URL.Query().Get("user")), Items: clientpagedbconn.Device_list_select(ConnectedDeivceList_Handler(), r.URL.Query().Get("user"), r.URL.Query().Get("s_order"), r.URL.Query().Get("e_order"))})
	w.Write(json)
}

// 하우스 기본 정보 검색
func DeviceList_Page_Search_Handler(w http.ResponseWriter, r *http.Request) {

	log.Printf("DeviceList_Page_Search_Handler %v,%v,%v", r.URL.Query().Get("user"), r.URL.Query().Get("column"), r.URL.Query().Get("data"))

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	// json, _ := json.Marshal(&structs.DeviceListPage{Items: clientpagedbconn.Device_list_select_search(r.URL.Query().Get("column"), r.URL.Query().Get("data"))})
	json, _ := json.Marshal(&structs.DeviceListPage{Items: clientpagedbconn.Device_list_select_search(ConnectedDeivceList_Handler(), r.URL.Query().Get("user"), r.URL.Query().Get("column"), r.URL.Query().Get("data"))})
	w.Write(json)
}

// 하우스 예약 정보
func DeviceAppointList_Page_Handler(w http.ResponseWriter, r *http.Request) {

	log.Printf("DeviceAppointList_Page_Handler %v", r.URL.Query().Get("user"))

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(&structs.DeviceAppointListPage{Items: Device_Web_List_Appoint_Handler(r.URL.Query().Get("user"))})
	w.Write(json)
}

// 하우스 예약 정보 검색
func DeviceAppointList_Page_Search_Handler(w http.ResponseWriter, r *http.Request) {

	log.Printf("DeviceAppointList_Page_Search_Handler %v", r.URL.Query().Get("id"))

	did := r.URL.Query().Get("id")

	var returnValue *structs.DeviceAppointList

	if csms != nil {
		cp, b := csms.Load(did)
		if b {
			var req Payload = structs.DeviceAppointAllRequest{
				Id: did,
			}
			res, err := cp.Call("AppointList", req)
			if err != nil {
				log.Printf("aa %s", err)
			}
			log.Print(res)
			returnValue = res.(*structs.DeviceAppointList)
		}

		if cp == nil {
			log.Print("cp is nil")
			returnValue = &structs.DeviceAppointList{Id: did, Items: clientpagedbconn.Device_list_appoint_select(did), Auto: 2}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(returnValue)
}

// 기기 예약 정보
func Device_Web_List_Appoint_Handler(user string) structs.DeviceAppointItem {

	log.Printf("Device_Web_List_Appoint_Handler %v", user)

	var devicelist = clientpagedbconn.DeviceSimple_list(ConnectedDeivceList_Handler(), user)

	if devicelist == nil {
		return structs.DeviceAppointItem{}
	}

	did := devicelist[0].Id

	var returnValue *structs.DeviceAppointList

	if csms != nil {
		cp, b := csms.Load(did)
		if b {
			var req Payload = structs.DeviceAppointAllRequest{
				Id: did,
			}
			res, err := cp.Call("AppointList", req)
			if err != nil {
				log.Printf("aa %s", err)
			}
			log.Print(res)
			returnValue = res.(*structs.DeviceAppointList)
		}

		if cp == nil {
			log.Print("cp is nil")
			returnValue = &structs.DeviceAppointList{Id: did, Items: clientpagedbconn.Device_list_appoint_select(did), Auto: 2}
		}
	}

	return structs.DeviceAppointItem{Devices: devicelist, Id: returnValue.Id, Auto: returnValue.Auto, Items: returnValue.Items}
}

// 하우스 일자별 센서값 정보
func DeviceDailylist_Page_Handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("DeviceDailylist_Page_Handler %v,%v,%v", r.URL.Query().Get("user"), r.URL.Query().Get("s_order"), r.URL.Query().Get("e_order"))

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(&structs.DeviceDailyListPage{Count: clientpagedbconn.Device_daily_list_count_select(r.URL.Query().Get("user")), Items: clientpagedbconn.Device_daily_list_select(r.URL.Query().Get("user"), r.URL.Query().Get("s_order"), r.URL.Query().Get("e_order"))})
	w.Write(json)
}

// 하우스 일자별 센서값 검색
func DeviceDailylist_Page_Search_Handler(w http.ResponseWriter, r *http.Request) {

	log.Printf("DeviceDailylist_Page_Search_Handler %v,%v,%v", r.URL.Query().Get("user"), r.URL.Query().Get("column"), r.URL.Query().Get("data"))

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(&structs.DeviceDailyListPage{Items: clientpagedbconn.Device_daily_list_select_search(r.URL.Query().Get("user"), r.URL.Query().Get("column"), r.URL.Query().Get("data"))})
	w.Write(json)
}

// 하우스 센서값 그래프 조회
func DeviceDailyGraph_Page_Handler(w http.ResponseWriter, r *http.Request) {

	log.Print("DeviceDailyGraph_Page_Handler %v", r.URL.Query().Get("user"))

	msg := structs.DeviceGraphTimeData{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("DeviceDailyGraph_Page_Handler err : %s", err)
		return
	}
	defer r.Body.Close()

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(&structs.DeviceDailyGraphPage{ListItem: clientpagedbconn.Device_daily_graph_list_select(r.URL.Query().Get("user")), Items: clientpagedbconn.Device_daily_graph_select(msg.StartTime, msg.EndTime, msg.Id)})
	w.Write(json)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// 하우스 onoff 명령
func Device_Web_Commad_Onoff_Handler(w http.ResponseWriter, r *http.Request) {

	msg := structs.CMDAllOnoff{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("Device_Web_Commad_Onoff_Handler err : %s", err)
		return
	}
	defer r.Body.Close()

	log.Printf("Device_Web_Commad_Onoff_Handler msg : %v", msg)
	var returnValue *structs.CMDAllOnoffConf

	if csms != nil {
		cp, b := csms.Load(msg.Id)
		if b {
			var req Payload = structs.CMDAllOnoff{
				Id:    msg.Id,
				Key:   msg.Key,
				Value: msg.Value,
			}
			res, err := cp.Call("CMDALLONOFF", req)
			if err != nil {
				log.Printf("aa %s", err)
			}
			log.Print(res)

			returnValue = res.(*structs.CMDAllOnoffConf)

		}

		if cp == nil {
			log.Print("cp is nil")
			returnValue = &structs.CMDAllOnoffConf{}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(returnValue)
}

// 하우스 소프트웨어 재부팅 명령
func Device_Web_Commad_Reboot_Software_Handler(w http.ResponseWriter, r *http.Request) {

	msg := structs.CMDRebootSoftware{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("Device_Web_Commad_Reboot_Software_Handler err : %s", err)
		return
	}
	defer r.Body.Close()

	log.Printf("Device_Web_Commad_Reboot_Software_Handler msg : %v", msg)
	var returnValue *structs.CMDRebootSoftwareConf

	if csms != nil {
		cp, b := csms.Load(msg.Id)
		if b {
			var req Payload = structs.CMDRebootSoftware{
				Id: msg.Id,
			}
			res, err := cp.Call("CMDREBOOTSOFT", req)
			if err != nil {
				log.Printf("aa %s", err)
			}
			log.Print(res)

			returnValue = res.(*structs.CMDRebootSoftwareConf)

		}

		if cp == nil {
			log.Print("cp is nil")
			returnValue = &structs.CMDRebootSoftwareConf{}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(returnValue)
}

// 하우스 하드웨어 재부팅 명령
func Device_Web_Commad_Reboot_Hardware_Handler(w http.ResponseWriter, r *http.Request) {

	msg := structs.CMDRebootHardware{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("Device_Web_Commad_Reboot_Hardware_Handler err : %s", err)
		return
	}
	defer r.Body.Close()

	log.Printf("Device_Web_Commad_Reboot_Hardware_Handler msg : %v", msg)
	var returnValue *structs.CMDRebootHardwareConf

	if csms != nil {
		cp, b := csms.Load(msg.Id)
		if b {
			var req Payload = structs.CMDRebootHardware{
				Id: msg.Id,
			}
			res, err := cp.Call("CMDREBOOTHARD", req)
			if err != nil {
				log.Printf("aa %s", err)
			}
			log.Print(res)

			returnValue = res.(*structs.CMDRebootHardwareConf)

		}

		if cp == nil {
			log.Print("cp is nil")
			returnValue = &structs.CMDRebootHardwareConf{}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(returnValue)
}

// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 사용자 알람 리스트 정보
func App_UserAlarmList_Handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("App_UserAlarmList_Handler%v", r.URL.Query().Get("user"))

	rb := clientpagedbconn.UserAlarmList_select(r.URL.Query().Get("user"))

	log.Printf("App_UserAlarmList_Handler %v", rb)
	json, _ := json.Marshal(&structs.UserAlarmListConf{User: r.URL.Query().Get("user"), Items: rb})
	log.Printf("App_UserAlarmList_Handler")
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

// 하우스 장치 간략한 정보
func App_Device_List_Simple_Handler(w http.ResponseWriter, r *http.Request) {

	log.Printf("App_Device_List_Simple_Handler %v", r.URL.Query().Get("user"))

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	json, _ := json.Marshal(&structs.DeviceSimpleList{Items: clientpagedbconn.DeviceSimple_list(ConnectedDeivceList_Handler(), r.URL.Query().Get("user"))})
	w.Write(json)
}

// 하우스 장치 자세한 정보
func App_Device_Detail_Handler(w http.ResponseWriter, r *http.Request) {

	log.Printf("App_Device_Detail_Handler %v", r.URL.Query().Get("id"))

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	var j = clientpagedbconn.DeviceDetail_select(ConnectedDeivceList_Handler(), r.URL.Query().Get("id"))

	json, _ := json.Marshal(&structs.DeviceDetail{Id: j.Id, Comp: j.Comp, Pump: j.Pump, Temp: j.Temp, Humi: j.Humi, F_Fan: j.F_Fan, E_Fan: j.E_Fan, Aircon: j.Aircon, Co2_Sensor: j.Co2_Sensor, C02_Creator: j.C02_Creator, Ph_Ec_Sensor: j.Ph_Ec_Sensor, Led: j.Led, Cctv: j.Cctv, Auto: j.Auto, B_Time: j.B_Time})
	w.Write(json)
}

// 하우스 장치 comp 명령어
func App_Device_CMD_Comp_Handler(w http.ResponseWriter, r *http.Request) {

	log.Print("App_Device_CMD_Comp_Handler")

	msg := structs.CMDComp{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("App_Device_CMD_Comp_Handler err : %s", err)
		return
	}

	defer r.Body.Close()

	var returnValue *structs.CMDCompConf

	if csms != nil {
		cp, b := csms.Load(msg.Id)
		if b {
			var req Payload = structs.CMDComp{
				Id:    msg.Id,
				Value: msg.Value,
			}
			res, err := cp.Call("CMDCOMP", req)
			if err != nil {
				log.Printf("aa %s", err)
			}
			log.Print(res)
			returnValue = res.(*structs.CMDCompConf)

		}

		if cp == nil {
			log.Print("cp is nil")
			returnValue = &structs.CMDCompConf{Id: msg.Id, Result: false}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(returnValue)
}

// 하우스 장치 pump 명령어
func App_Device_CMD_Pump_Handler(w http.ResponseWriter, r *http.Request) {

	log.Print("App_Device_CMD_Pump_Handler")

	msg := structs.CMDPump{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("App_Device_CMD_Pump_Handler err : %s", err)
		return
	}

	defer r.Body.Close()

	var returnValue *structs.CMDPumpConf

	if csms != nil {
		cp, b := csms.Load(msg.Id)
		if b {
			var req Payload = structs.CMDPump{
				Id:    msg.Id,
				Value: msg.Value,
			}
			res, err := cp.Call("CMDPUMP", req)
			if err != nil {
				log.Printf("aa %s", err)
			}
			log.Print(res)
			returnValue = res.(*structs.CMDPumpConf)

		}

		if cp == nil {
			log.Print("cp is nil")
			returnValue = &structs.CMDPumpConf{Id: msg.Id, Result: false}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(returnValue)
}

// 하우스 장치 FFan 명령어
func App_Device_CMD_FFan_Handler(w http.ResponseWriter, r *http.Request) {

	log.Print("App_Device_CMD_FFan_Handler")

	msg := structs.CMDFFan{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("App_Device_CMD_FFan_Handler err : %s", err)
		return
	}

	defer r.Body.Close()

	var returnValue *structs.CMDFFanConf

	if csms != nil {
		cp, b := csms.Load(msg.Id)
		if b {
			var req Payload = structs.CMDFFan{
				Id:    msg.Id,
				Value: msg.Value,
			}
			res, err := cp.Call("CMDFFAN", req)
			if err != nil {
				log.Printf("aa %s", err)
			}
			log.Print(res)
			returnValue = res.(*structs.CMDFFanConf)

		}

		if cp == nil {
			log.Print("cp is nil")
			returnValue = &structs.CMDFFanConf{Id: msg.Id, Result: false}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(returnValue)
}

// 하우스 장치 EFan 명령어
func App_Device_CMD_EFan_Handler(w http.ResponseWriter, r *http.Request) {

	log.Print("App_Device_CMD_EFan_Handler")

	msg := structs.CMDEFan{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("App_Device_CMD_EFan_Handler err : %s", err)
		return
	}

	defer r.Body.Close()

	var returnValue *structs.CMDEFanConf

	if csms != nil {
		cp, b := csms.Load(msg.Id)
		if b {
			var req Payload = structs.CMDEFan{
				Id:    msg.Id,
				Value: msg.Value,
			}
			res, err := cp.Call("CMDEFAN", req)
			if err != nil {
				log.Printf("aa %s", err)
			}
			log.Print(res)
			returnValue = res.(*structs.CMDEFanConf)
		}

		if cp == nil {
			log.Print("cp is nil")
			returnValue = &structs.CMDEFanConf{Id: msg.Id, Result: false}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(returnValue)
}

// 하우스 장치 Aircon 명령어
func App_Device_CMD_Aircon_Handler(w http.ResponseWriter, r *http.Request) {

	log.Print("App_Device_CMD_Aircon_Handler")

	msg := structs.CMDAircon{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("App_Device_CMD_Aircon_Handler err : %s", err)
		return
	}

	defer r.Body.Close()

	var returnValue *structs.CMDAirconConf

	if csms != nil {
		cp, b := csms.Load(msg.Id)
		if b {
			var req Payload = structs.CMDAircon{
				Id:    msg.Id,
				Value: msg.Value,
			}
			res, err := cp.Call("CMDAIRCON", req)
			if err != nil {
				log.Printf("aa %s", err)
			}
			log.Print(res)
			returnValue = res.(*structs.CMDAirconConf)

		}

		if cp == nil {
			log.Print("cp is nil")
			returnValue = &structs.CMDAirconConf{Id: msg.Id, Result: false}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(returnValue)
}

// 하우스 장치 Co2Creator 명령어
func App_Device_CMD_Co2Creator_Handler(w http.ResponseWriter, r *http.Request) {

	log.Print("App_Device_CMD_Co2Creator_Handler")

	msg := structs.CMDCo2Creator{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("App_Device_CMD_Co2Creator_Handler err : %s", err)
		return
	}

	defer r.Body.Close()

	var returnValue *structs.CMDCo2CreatorConf

	if csms != nil {
		cp, b := csms.Load(msg.Id)
		if b {
			var req Payload = structs.CMDCo2Creator{
				Id:    msg.Id,
				Value: msg.Value,
			}
			res, err := cp.Call("CMDCO2CREATOR", req)
			if err != nil {
				log.Printf("aa %s", err)
			}
			log.Print(res)
			returnValue = res.(*structs.CMDCo2CreatorConf)

		}

		if cp == nil {
			log.Print("cp is nil")
			returnValue = &structs.CMDCo2CreatorConf{Id: msg.Id, Result: false}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(returnValue)
}

// 하우스 장치 LED 명령어
func App_Device_CMD_LED_Handler(w http.ResponseWriter, r *http.Request) {

	log.Print("App_Device_CMD_LED_Handler")

	msg := structs.CMDLED{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("App_Device_CMD_LED_Handler err : %s", err)
		return
	}

	defer r.Body.Close()

	var returnValue *structs.CMDLEDConf

	if csms != nil {
		cp, b := csms.Load(msg.Id)
		if b {
			var req Payload = structs.CMDLED{
				Id:    msg.Id,
				Value: msg.Value,
			}
			res, err := cp.Call("CMDLED", req)
			if err != nil {
				log.Printf("aa %s", err)
			}
			log.Print(res)
			returnValue = res.(*structs.CMDLEDConf)
		}

		if cp == nil {
			log.Print("cp is nil")
			returnValue = &structs.CMDLEDConf{Id: msg.Id, Result: false}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(returnValue)
}

// 기기 예약 정보
func App_Device_List_Appoint_Handler(w http.ResponseWriter, r *http.Request) {

	did := r.URL.Query().Get("deviceid")

	log.Printf("App_Device_List_Appoint_Handler %v", did)

	var returnValue *structs.DeviceAppointList

	if csms != nil {
		cp, b := csms.Load(did)
		if b {
			var req Payload = structs.DeviceAppointAllRequest{
				Id: did,
			}
			res, err := cp.Call("AppointList", req)
			if err != nil {
				log.Printf("aa %s", err)
			}
			log.Print(res)
			returnValue = res.(*structs.DeviceAppointList)
		}

		if cp == nil {
			log.Print("cp is nil")
			returnValue = &structs.DeviceAppointList{Id: did, Items: clientpagedbconn.Device_list_appoint_select(did), Auto: 2}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(returnValue)
}

// 하우스 예약 compump 명령어
func App_Device_ACMD_ComPump_Handler(w http.ResponseWriter, r *http.Request) {

	msg := structs.ACMDComPump{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("App_Device_ACMD_ComPump_Handler err : %s", err)
		return
	}

	log.Print("App_Device_ACMD_ComPump_Handler")

	defer r.Body.Close()

	var returnValue *structs.CMDCompConf

	if csms != nil {
		cp, b := csms.Load(msg.Id)
		if b {
			var req Payload = structs.ACMDComPump{
				Id:    msg.Id,
				Start: msg.Start,
				Stop:  msg.Stop,
				Onoff: msg.Onoff,
			}
			res, err := cp.Call("ACMDCOMPUMP", req)
			if err != nil {
				log.Printf("aa %s", err)
			}
			log.Print(res)
			returnValue = res.(*structs.CMDCompConf)

			if returnValue.Result == true {
				clientpagedbconn.Device_list_appoint_update(msg.Id, "compump", msg.Start, msg.Stop, msg.Onoff)
			}

		}

		if cp == nil {
			log.Print("cp is nil")
			returnValue = &structs.CMDCompConf{Id: msg.Id, Result: false}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(returnValue)
}

// 하우스  예약  FFan 명령어
func App_Device_ACMD_FFan_Handler(w http.ResponseWriter, r *http.Request) {

	log.Print("App_Device_ACMD_FFan_Handler")

	msg := structs.ACMDFFan{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("App_Device_ACMD_FFan_Handler err : %s", err)
		return
	}

	defer r.Body.Close()

	var returnValue *structs.CMDFFanConf

	if csms != nil {
		cp, b := csms.Load(msg.Id)
		if b {
			var req Payload = structs.ACMDFFan{
				Id:    msg.Id,
				Start: msg.Start,
				Stop:  msg.Stop,
				Onoff: msg.Onoff,
			}
			res, err := cp.Call("ACMDFFAN", req)
			if err != nil {
				log.Printf("aa %s", err)
			}
			log.Print(res)
			returnValue = res.(*structs.CMDFFanConf)

			if returnValue.Result == true {
				clientpagedbconn.Device_list_appoint_update(msg.Id, "ffan", msg.Start, msg.Stop, msg.Onoff)
			}

		}

		if cp == nil {
			log.Print("cp is nil")
			returnValue = &structs.CMDFFanConf{Id: msg.Id, Result: false}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(returnValue)
}

// 하우스  예약  EFan 명령어
func App_Device_ACMD_EFan_Handler(w http.ResponseWriter, r *http.Request) {

	log.Print("App_Device_ACMD_EFan_Handler")

	msg := structs.ACMDEFan{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("App_Device_ACMD_EFan_Handler err : %s", err)
		return
	}

	defer r.Body.Close()

	var returnValue *structs.CMDEFanConf

	if csms != nil {
		cp, b := csms.Load(msg.Id)
		if b {
			var req Payload = structs.ACMDEFan{
				Id:    msg.Id,
				Start: msg.Start,
				Stop:  msg.Stop,
				Onoff: msg.Onoff,
			}
			res, err := cp.Call("ACMDEFAN", req)
			if err != nil {
				log.Printf("aa %s", err)
			}
			log.Print(res)
			returnValue = res.(*structs.CMDEFanConf)

			if returnValue.Result == true {
				clientpagedbconn.Device_list_appoint_update(msg.Id, "efan", msg.Start, msg.Stop, msg.Onoff)
			}

		}

		if cp == nil {
			log.Print("cp is nil")
			returnValue = &structs.CMDEFanConf{Id: msg.Id, Result: false}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(returnValue)
}

// 하우스 예약  Aircon 명령어
func App_Device_ACMD_Aircon_Handler(w http.ResponseWriter, r *http.Request) {

	log.Print("App_Device_ACMD_Aircon_Handler")

	msg := structs.ACMDAircon{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("App_Device_ACMD_Aircon_Handler err : %s", err)
		return
	}

	defer r.Body.Close()

	var returnValue *structs.CMDAirconConf

	if csms != nil {
		cp, b := csms.Load(msg.Id)
		if b {
			var req Payload = structs.ACMDAircon{
				Id:    msg.Id,
				Start: msg.Start,
				Stop:  msg.Stop,
				Onoff: msg.Onoff,
			}
			res, err := cp.Call("ACMDAIRCON", req)
			if err != nil {
				log.Printf("aa %s", err)
			}
			log.Print(res)
			returnValue = res.(*structs.CMDAirconConf)

			if returnValue.Result == true {
				clientpagedbconn.Device_list_appoint_update(msg.Id, "aircon", msg.Start, msg.Stop, msg.Onoff)
			}

		}

		if cp == nil {
			log.Print("cp is nil")
			returnValue = &structs.CMDAirconConf{Id: msg.Id, Result: false}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(returnValue)
}

// 하우스 예약  Co2Creator 명령어
func App_Device_ACMD_Co2Creator_Handler(w http.ResponseWriter, r *http.Request) {

	log.Print("App_Device_ACMD_Co2Creator_Handler")

	msg := structs.ACMDCo2Creator{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("App_Device_ACMD_Co2Creator_Handler err : %s", err)
		return
	}

	defer r.Body.Close()

	var returnValue *structs.CMDCo2CreatorConf

	if csms != nil {
		cp, b := csms.Load(msg.Id)
		if b {
			var req Payload = structs.ACMDCo2Creator{
				Id:    msg.Id,
				Start: msg.Start,
				Stop:  msg.Stop,
				Onoff: msg.Onoff,
			}
			res, err := cp.Call("ACMDCO2CREATOR", req)
			if err != nil {
				log.Printf("aa %s", err)
			}
			log.Print(res)
			returnValue = res.(*structs.CMDCo2CreatorConf)

			if returnValue.Result == true {
				clientpagedbconn.Device_list_appoint_update(msg.Id, "co2creator", msg.Start, msg.Stop, msg.Onoff)
			}

		}

		if cp == nil {
			log.Print("cp is nil")
			returnValue = &structs.CMDCo2CreatorConf{Id: msg.Id, Result: false}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(returnValue)
}

// 하우스  예약 Led 명령어
func App_Device_ACMD_LED_Handler(w http.ResponseWriter, r *http.Request) {

	log.Print("App_Device_ACMD_LED_Handler")

	msg := structs.ACMDLED{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("App_Device_ACMD_LED_Handler err : %s", err)
		return
	}

	defer r.Body.Close()

	var returnValue *structs.CMDLEDConf

	if csms != nil {
		cp, b := csms.Load(msg.Id)
		if b {
			var req Payload = structs.ACMDLED{
				Id:    msg.Id,
				Start: msg.Start,
				Stop:  msg.Stop,
				Onoff: msg.Onoff,
			}
			res, err := cp.Call("ACMDLED", req)
			if err != nil {
				log.Printf("aa %s", err)
			}
			log.Print(res)
			returnValue = res.(*structs.CMDLEDConf)

			if returnValue.Result == true {
				clientpagedbconn.Device_list_appoint_update(msg.Id, "led", msg.Start, msg.Stop, msg.Onoff)
			}

		}

		if cp == nil {
			log.Print("cp is nil")
			returnValue = &structs.CMDLEDConf{Id: msg.Id, Result: false}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(returnValue)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// 원격 연결되어있는 기기 리스트 반환
func ConnectedDeivceList_Handler() []string {

	var returnValue []string

	for k := range csms.chargepoints {
		// log.Printf("csms.chargepoints ----- key = %v \n", k)

		returnValue = append(returnValue, k)
	}

	return returnValue
}

func main() {

	err := godotenv.Load("./info.env")
	if err != nil {
		log.Println(err)
		return
	}

	runtime.NumCPU() // 실행 환경의 물리적 CPU 최대 개수로 설정된다.

	//랜덤값 기준시드 정해주기
	// rand.Seed(time.Now().UnixNano())
	rand.New(rand.NewSource(time.Now().UnixNano()))

	logger, _ := zap.NewDevelopment()
	logd = logger.Sugar()
	defer logd.Sync()

	// set ocpp library's logger to zap logger
	SetLogger(logd)

	//fcm 등록
	// fcmInit()

	// start csms server with default configurations
	csms = NewServer()

	csms.AddSubProtocol("sh1.0")
	// csms.SetCheckOriginHandler(func(r *http.Request) bool { return true })
	csms.SetCheckOriginHandler(SetCheckOriginHanler_H)
	// csms.SetPreUpgradeHandler(customPreUpgradeHandler)	// 접속시 헤더의 Authorization 값으로 id, pw 를 보내 인증하는 방식
	csms.SetCallQueueSize(32)
	// csms.SetTimeoutConfig(ServerTimeoutConfig{OcppWait: 10000}) // 연결 제한시간
	// csms.SetTimeoutConfig(ServerTimeoutConfig{PingWait: 5000}) // 핑 제한시간

	// register charge-point-initiated action handlers
	// csms.After("BootNotification", SendChangeConfigration)
	csms.On("Heartbeat", HeartBeatHandler) // 하트비트

	csms.On("BOOTDEVICE", BootHandler)             // 기기 처음 켜질때 보고 및 등록
	csms.On("REPORTCOMP", COMPHandler)             // 컴프 onoff
	csms.On("REPORTPUMP", PUMPHandler)             // 펌프 onoff
	csms.On("REPORTFFAN", FFANHandler)             // 유동팬 onoff
	csms.On("REPORTEFAN", EFANHandler)             // 배기팬 onoff
	csms.On("REPORTAIRCON", AIRCONHandler)         // 냉난방기 onoff
	csms.On("REPORTCO2CREATOR", CO2CREATORHandler) // co2발생기 onoff
	csms.On("REPORTLED", LEDHandler)               // led onoff
	csms.On("REPORTTEMP", TEMPHandler)             // 온도값
	csms.On("REPORTHUMI", HUMIHandler)             // 습도값
	csms.On("REPORTCO2SENSOR", CO2SENSORHandler)   // co2 센서값
	csms.On("REPORTPHECSENSOR", PHECSENSORHandler) // ph ec 센서값
	csms.On("REPORTAuto", AutoHandler)             //자동 수동 바뀜

	csms.On("ALARM", AlarmHandler) // 기기 알림

	csms.On("REPORTALLSENSOR", ReportAllSensorHandler) // 정기 센서값 보고

	go csms.Start(os.Getenv("CSMSSERVER"), os.Getenv("CSMSSERVER_PATH"), nil)

	//연결이 되어 있는지 확인
	// log.Printf("aa : %s", strconv.FormatBool(csms.IsConnected("device_01")))
	// log.Printf("aa : %s", strconv.FormatBool(csms.IsConnected("device_02")))

	//////////////////////////////////////////////////////////////////////////////////////////////

	mux := mux.NewRouter()

	//////////////////////////////////////////////////////////////////////////////////////////////
	//로그인
	// mux.Path("/").HandlerFunc(adminpage.LoginPage_Handler)

	//웹 고객 로그인
	mux.Path("/login").HandlerFunc(clientpage.Login_Handler).Methods("POST")
	//어플 로그인
	mux.Path("/app/login").HandlerFunc(clientpage.Login_App_Handler).Methods("POST")
	//웹 관리자 로그인
	mux.Path("/if/login").HandlerFunc(adminpage.Login_Handler).Methods("POST")
	//기기 jwt token 발급 로그인
	mux.Path("/device/login").HandlerFunc(clientpage.Login_Device_Handler).Methods("POST")
	//어플 access token 재발급
	mux.Path("/app/jwt/refresh").HandlerFunc(jwt.RefreshToken_App_Handler).Methods("POST")
	//웹 access token 재발급
	mux.Path("/jwt/refresh").HandlerFunc(jwt.RefreshToken_Handler).Methods("POST")

	privateRouter := mux.PathPrefix("/").Subrouter()
	privateRouter.Use(jwt.JwtMiddleware)

	// mux.Path("/login").HandlerFunc(Command_Handler)
	privateRouter.Path("/logout").HandlerFunc(jwt.Logout_Handler).Methods("POST")

	/////////////////////////////////////////////기기/////////////////////////////////////////////////

	//기기 cfssl 재발급
	privateRouter.Path("/device/cfssl/refresh").HandlerFunc(Device_CFSSL_Refresh_Handler).Methods("POST")

	/////////////////////////////////////////////사용자페이지/////////////////////////////////////////////////

	//대쉬보드
	privateRouter.Path("/dashboardlist").HandlerFunc(Dashboard_Page_Handler).Methods("GET")
	//하우스 기본 정보
	privateRouter.Path("/devicelist").HandlerFunc(DeviceList_Page_Handler).Methods("GET")
	//하우스 예약 정보
	privateRouter.Path("/deviceappointlist").HandlerFunc(DeviceAppointList_Page_Handler).Methods("GET")
	//하우스 일자별 센서값 조회
	privateRouter.Path("/devicedailylist").HandlerFunc(DeviceDailylist_Page_Handler).Methods("GET")
	//하우스 센서값 그래프 조회
	privateRouter.Path("/devicedailygraph").HandlerFunc(DeviceDailyGraph_Page_Handler).Methods("POST")
	//사용자 알림 정보
	privateRouter.Path("/user/alarmlist").HandlerFunc(clientpage.UserAlarmList_Handler).Methods("GET")

	//////////////////////////////////사용자페이지 검색///////////////////////////////////

	//하우스 기본 정보 검색
	privateRouter.Path("/devicelist/search").HandlerFunc(DeviceList_Page_Search_Handler).Methods("GET")
	//하우스 예약 정보 검색
	privateRouter.Path("/deviceappointlist/search").HandlerFunc(DeviceAppointList_Page_Search_Handler).Methods("GET")
	//하우스 일자별 센서값 검색
	privateRouter.Path("/devicedailylist/search").HandlerFunc(DeviceDailylist_Page_Search_Handler).Methods("GET")

	//////////////////////////////////사용자페이지 명령///////////////////////////////////

	//하우스 onoff 명령
	privateRouter.Path("/command/nf").HandlerFunc(Device_Web_Commad_Onoff_Handler).Methods("POST")
	//하우스 소프트웨어 재부팅 명령
	privateRouter.Path("/command/rbsw").HandlerFunc(Device_Web_Commad_Reboot_Software_Handler).Methods("POST")
	//하우스 하드웨어 재부팅 명령
	privateRouter.Path("/command/rbhw").HandlerFunc(Device_Web_Commad_Reboot_Hardware_Handler).Methods("POST")

	/////////////////////////////////////////////관리자페이지/////////////////////////////////////////////////

	//대쉬보드
	privateRouter.Path("/if/dashboardlist").HandlerFunc(IF_Dashboard_Page_Handler).Methods("GET")
	//하우스 기본 정보
	privateRouter.Path("/if/devicelist").HandlerFunc(IF_DeviceList_Page_Handler).Methods("GET")
	//하우스 예약 정보
	privateRouter.Path("/if/deviceappointlist").HandlerFunc(IF_DeviceAppointList_Page_Handler).Methods("GET")
	//하우스 등록 유저 기기 정보
	privateRouter.Path("/if/deviceuserlist").HandlerFunc(IF_DeviceUserList_Page_Handler).Methods("GET")
	//유저 정보
	privateRouter.Path("/if/userlist").HandlerFunc(IF_UserList_Page_Handler).Methods("GET")
	//하우스 업데이트 관리
	privateRouter.Path("/if/deviceupdatelist").HandlerFunc(IF_Updmng_Page_Handler).Methods("GET")
	//진단파일 요청 관리
	privateRouter.Path("/if/devicelogcmdpage").HandlerFunc(IF_Logmng_Page_Handler).Methods("GET")
	//하우스 일자별 센서값 조회
	privateRouter.Path("/if/devicedailylist").HandlerFunc(IF_DeviceDailylist_Page_Handler).Methods("GET")
	//하우스 센서값 그래프 조회
	privateRouter.Path("/if/devicedailygraph").HandlerFunc(IF_DeviceDailyGraph_Page_Handler).Methods("POST")

	//////////////////////////////////관리자페이지 검색///////////////////////////////////

	//하우스 기본 정보 검색
	privateRouter.Path("/if/devicelist/search").HandlerFunc(IF_DeviceList_Page_Search_Handler).Methods("GET")
	//하우스 예약 정보 검색
	privateRouter.Path("/if/deviceappointlist/search").HandlerFunc(DeviceAppointList_Page_Search_Handler).Methods("GET")
	//하우스 등록 유저 기기 정보 검색
	privateRouter.Path("/if/deviceuserlist/search").HandlerFunc(IF_DeviceUserList_Page_Search_Handler).Methods("GET")
	//유저 정보 검색
	privateRouter.Path("/if/userlist/search").HandlerFunc(IF_UserList_Page_Search_Handler).Methods("GET")
	//하우스 업데이트 관리 검색
	privateRouter.Path("/if/deviceupdatelist/search").HandlerFunc(IF_Updmng_Page_Search_Handler).Methods("GET")
	//진단파일 요청 관리 검색
	privateRouter.Path("/if/devicelogcmdpage/search").HandlerFunc(IF_Logmng_Page_Search_Handler).Methods("GET")
	//하우스 일자별 센서값 검색
	privateRouter.Path("/if/devicedailylist/search").HandlerFunc(IF_DeviceDailylist_Page_Search_Handler).Methods("GET")

	//////////////////////////////////관리자페이지 수정 및 삽입///////////////////////////////////

	//하우스 등록 유저 기기 정보 삽입
	privateRouter.Path("/if/deviceuser/register").HandlerFunc(IF_DeviceUser_Register_Handler).Methods("POST")
	//하우스 등록 유저 기기 정보 삭제
	privateRouter.Path("/if/deviceuser/deleter").HandlerFunc(IF_DeviceUser_Deleter_Handler).Methods("POST")
	//유저 정보 삽입
	privateRouter.Path("/if/user/register").HandlerFunc(IF_User_Register_Handler).Methods("POST")
	//유저 정보 삭제
	privateRouter.Path("/if/user/deleter").HandlerFunc(IF_User_Deleter_Handler).Methods("POST")
	//유저 정보 업데이트
	privateRouter.Path("/if/user/updater").HandlerFunc(IF_User_Updater_Handler).Methods("POST")
	//유저 알람 전송
	privateRouter.Path("/if/user/alert").HandlerFunc(IF_User_Alert_Handler).Methods("POST")
	//유저 전체 알람 전송
	privateRouter.Path("/if/user/alert/whole").HandlerFunc(IF_User_Whole_Alert_Handler).Methods("POST")

	//////////////////////////////////관리자페이지 명령///////////////////////////////////

	//하우스 log 업로드 명령
	privateRouter.Path("/command/lgup").HandlerFunc(IF_Device_Web_Commad_Log_Upload_Handler).Methods("POST")
	//하우스 업데이트 명령
	privateRouter.Path("/command/udfw").HandlerFunc(IF_Device_Web_Commad_Program_Update_Handler).Methods("POST")

	/////////////////////////////////////////////어플/////////////////////////////////////////////////

	//사용자 알림 정보
	privateRouter.Path("/app/alarmlist").HandlerFunc(App_UserAlarmList_Handler).Methods("GET")
	//하우스 리스트 및 간략 정보
	privateRouter.Path("/app/simple/list").HandlerFunc(App_Device_List_Simple_Handler).Methods("GET")
	//하우스 상태 상세조회
	privateRouter.Path("/app/detail").HandlerFunc(App_Device_Detail_Handler).Methods("GET")
	//어플 onesignal id 보고
	privateRouter.Path("/app/rpt/osid").HandlerFunc(clientpage.App_Device_Report_OneSignalID_Handler).Methods("POST")
	//하우스 컴프 명령
	privateRouter.Path("/app/cmd/comp").HandlerFunc(App_Device_CMD_Comp_Handler).Methods("POST")
	//하우스 펌프 명령
	privateRouter.Path("/app/cmd/pump").HandlerFunc(App_Device_CMD_Pump_Handler).Methods("POST")
	//하우스 유동팬 명령
	privateRouter.Path("/app/cmd/ffan").HandlerFunc(App_Device_CMD_FFan_Handler).Methods("POST")
	//하우스 배기팬 명령
	privateRouter.Path("/app/cmd/efan").HandlerFunc(App_Device_CMD_EFan_Handler).Methods("POST")
	//하우스 냉난방기 명령
	privateRouter.Path("/app/cmd/aircon").HandlerFunc(App_Device_CMD_Aircon_Handler).Methods("POST")
	//하우스 co2발생기 명령
	privateRouter.Path("/app/cmd/co2creator").HandlerFunc(App_Device_CMD_Co2Creator_Handler).Methods("POST")
	//하우스 led 명령
	privateRouter.Path("/app/cmd/led").HandlerFunc(App_Device_CMD_LED_Handler).Methods("POST")
	//기기 예약 정보
	privateRouter.Path("/app/appoint/list").HandlerFunc(App_Device_List_Appoint_Handler).Methods("GET")
	//기기 예약 컴프 명령
	privateRouter.Path("/app/cmd/appoint/compump").HandlerFunc(App_Device_ACMD_ComPump_Handler).Methods("POST")
	//기기 예약 유동팬 명령
	privateRouter.Path("/app/cmd/appoint/ffan").HandlerFunc(App_Device_ACMD_FFan_Handler).Methods("POST")
	//기기 예약 배기팬 명령
	privateRouter.Path("/app/cmd/appoint/efan").HandlerFunc(App_Device_ACMD_EFan_Handler).Methods("POST")
	//기기 예약 냉난방기 명령
	privateRouter.Path("/app/cmd/appoint/aircon").HandlerFunc(App_Device_ACMD_Aircon_Handler).Methods("POST")
	//기기 예약 co2발생기 명령
	privateRouter.Path("/app/cmd/appoint/co2creator").HandlerFunc(App_Device_ACMD_Co2Creator_Handler).Methods("POST")
	//기기 예약 led 명령
	privateRouter.Path("/app/cmd/appoint/led").HandlerFunc(App_Device_ACMD_LED_Handler).Methods("POST")

	//////////////////////////////////////////////////////////////////////////////////////////////

	muxhandler := accessControlMiddleware(mux)

	// n.UseHandler(mux)
	log.Print(http.ListenAndServe(os.Getenv("HTTPSERVER"), muxhandler))

}


func accessControlMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS,PUT")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type ,application/json")

		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})

}

////////////////////////////////////////////////////////////////////////////////////////////////////////

func SetCheckOriginHanler_H(r *http.Request) bool {
	// log.Printf("SetCheckOriginHanler %v	\n", r)
	log.Println("SetCheckOriginHanler")

	for k, v := range csms.chargepoints {
		log.Printf("csms.chargepoints ----- key = %v , value = %v \n", k, v)
	}

	// log.Printf("access : %v", r.Header)
	log.Printf("access : %v", r.Host)

	return true
}

func HeartBeatHandler(cp *ChargePoint, p Payload) Payload {
	req := p.(*structs.Heartbeat)
	log.Printf("\nid: %s\nHeartbeat: %v", cp.Id, req)

	// cp.ResetPingPong(30000) // ping 리셋

	var res Payload = &structs.HeartbeatConf{
		CurrentTime: time.Now().Format("2006-01-02T15:04:05.000Z"),
	}
	return res
}

//////////////////////////////////////////////////////////////////////////////////////////////
type CertResponse struct{
	Ca string `json:"ca"`
	Client  string `json:"client"`
	ClientKey  string `json:"client_key"`
	ExpiryDate  string `json:"expiry_date"`
}

type CertRequest struct{
	DeviceID string `json:"device_id"`
}

//기기 cfssl 재발급
func Device_CFSSL_Refresh_Handler(w http.ResponseWriter, r *http.Request) {

	log.Print("Device_CFSSL_Refresh_Handler")

	msg := CertRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("Device_CFSSL_Refresh_Handler err : %s", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		return
	}
	defer r.Body.Close()

	ca_path := os.Getenv("CFSSL_CA_PATH")
	cert_path := os.Getenv("CFSSL_CERT_PATH")
	key_path := os.Getenv("CFSSL_KEY_PATH")

	ca_Bytes,_ := os.ReadFile(ca_path)
	cert_Bytes,_:= os.ReadFile(cert_path)
	key_Bytes,_ := os.ReadFile(key_path)

	block, _ := pem.Decode(cert_Bytes)
	cert,_:= x509.ParseCertificate(block.Bytes)

	// // 현재 시간이 인증서 유효기간 내에 있는지 확인
	// now := time.Now()
	// if now.Before(cert.NotBefore) || now.After(cert.NotAfter) {
	// 	return false, cert.NotAfter
	// }

	// // 만료 7일 전이면 갱신하도록 설정 (선택사항)
	// if now.Add(7 * 24 * time.Hour).After(cert.NotAfter) {
	// 	return false, cert.NotAfter
	// }

	// 인증서 갱신은 cron job 으로 서버와 분리하여 관리하는게 좋음
	// cmd:= exec.Command("../../../../etc/cfssl","./auto-cfssl.sh","Go")

	// output, err := cmd.Output()

	// if err != nil {
	// 	log.Printf("Device_CFSSL_Refresh_Handler cfssl err : %s", err)
	// }

	// fmt.Printf("Device_CFSSL_Refresh_Handler cfssl output: %s\n", output)

	expiryDate := cert.NotAfter.Format("2006-01-02 15:04:05")


	json, _ := json.Marshal(&CertResponse{Ca: string(ca_Bytes), Client: string(cert_Bytes), ClientKey: string(key_Bytes), ExpiryDate: expiryDate})

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

//////////////////////////////////////////////////////////////////////////////////////////////

func BootHandler(cp *ChargePoint, p Payload) Payload {
	req := p.(*structs.Boot)
	log.Printf("\nid: %s\nBoot: %v", cp.Id, req)
	// csms.Store(cp) // cp 저장

	var res Payload

	// 한국 시간대 로드
	loc, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		fmt.Println("시간대 로딩 에러:", err)
		return structs.BootConf{
			Id:     req.Id,
			Result: false,
		}
	}

	// 현재 시간을 한국 시간대로 설정
	now := time.Now().In(loc)

	B_Time := now.Format("2006-01-02 15:04:05")

	rb := device_handler.Device_Boot_update(req.Id, B_Time, req.State, req.Appoints)

	if rb {
		res = structs.BootConf{
			Id:     req.Id,
			Result: true,
		}
	} else {
		res = structs.BootConf{
			Id:     req.Id,
			Result: false,
		}
	}

	return res

}

func COMPHandler(cp *ChargePoint, p Payload) Payload {
	req := p.(*structs.REPORTComp)
	log.Printf("\nid: %s\nREPORTCOMP: %v", cp.Id, req)
	log.Printf("device comp : %v", req.Value)
	// csms.Store(cp) // cp 저장

	var res Payload

	rb := device_handler.Device_Comp_update(req.Id, req.Value)

	if rb {
		res = structs.REPORTCompConf{
			Id:     req.Id,
			Result: true,
		}
	} else {
		res = structs.REPORTCompConf{
			Id:     req.Id,
			Result: false,
		}
	}

	return res

}

func PUMPHandler(cp *ChargePoint, p Payload) Payload {
	req := p.(*structs.REPORTPump)
	log.Printf("\nid: %s\nREPORTPUMP: %v", cp.Id, req)
	log.Printf("device PUMP : %v", req.Value)
	// csms.Store(cp) // cp 저장

	var res Payload

	rb := device_handler.Device_Pump_update(req.Id, req.Value)

	if rb {
		res = structs.REPORTPumpConf{
			Id:     req.Id,
			Result: true,
		}
	} else {
		res = structs.REPORTPumpConf{
			Id:     req.Id,
			Result: false,
		}
	}

	return res

}

func TEMPHandler(cp *ChargePoint, p Payload) Payload {
	req := p.(*structs.REPORTTemp)
	log.Printf("\nid: %s\nREPORTTEMP: %v", cp.Id, req)
	log.Printf("device TEMP : %f", req.Value)
	// csms.Store(cp) // cp 저장

	var res Payload

	rb := device_handler.Device_Temp_update(req.Id, req.Value)

	if rb {
		res = structs.REPORTTempConf{
			Id:     req.Id,
			Result: true,
		}
	} else {
		res = structs.REPORTTempConf{
			Id:     req.Id,
			Result: false,
		}
	}

	return res

}

func HUMIHandler(cp *ChargePoint, p Payload) Payload {
	req := p.(*structs.REPORTHumi)
	log.Printf("\nid: %s\nREPORTHUMI: %v", cp.Id, req)
	log.Printf("device HUMI : %f", req.Value)
	// csms.Store(cp) // cp 저장

	var res Payload

	rb := device_handler.Device_Humi_update(req.Id, req.Value)

	if rb {
		res = structs.REPORTHumiConf{
			Id:     req.Id,
			Result: true,
		}
	} else {
		res = structs.REPORTHumiConf{
			Id:     req.Id,
			Result: false,
		}
	}

	return res

}

func FFANHandler(cp *ChargePoint, p Payload) Payload {
	req := p.(*structs.REPORTFFan)
	log.Printf("\nid: %s\nREPORTFFAN: %v", cp.Id, req)
	log.Printf("device FFAN : %v", req.Value)
	// csms.Store(cp) // cp 저장

	var res Payload

	rb := device_handler.Device_FFan_update(req.Id, req.Value)

	if rb {
		res = structs.REPORTFFanConf{
			Id:     req.Id,
			Result: true,
		}
	} else {
		res = structs.REPORTFFanConf{
			Id:     req.Id,
			Result: false,
		}
	}

	return res

}

func EFANHandler(cp *ChargePoint, p Payload) Payload {
	req := p.(*structs.REPORTEFan)
	log.Printf("\nid: %s\nREPORTEFAN: %v", cp.Id, req)
	log.Printf("device EFAN : %v", req.Value)
	// csms.Store(cp) // cp 저장

	var res Payload

	rb := device_handler.Device_EFan_update(req.Id, req.Value)

	if rb {
		res = structs.REPORTEFanConf{
			Id:     req.Id,
			Result: true,
		}
	} else {
		res = structs.REPORTEFanConf{
			Id:     req.Id,
			Result: false,
		}
	}

	return res

}

func AIRCONHandler(cp *ChargePoint, p Payload) Payload {
	req := p.(*structs.REPORTAircon)
	log.Printf("\nid: %s\nREPORTAIRCON: %v", cp.Id, req)
	log.Printf("device AIRCON : %v", req.Value)
	// csms.Store(cp) // cp 저장

	var res Payload

	rb := device_handler.Device_Aircon_update(req.Id, req.Value)

	if rb {
		res = structs.REPORTAirconConf{
			Id:     req.Id,
			Result: true,
		}
	} else {
		res = structs.REPORTAirconConf{
			Id:     req.Id,
			Result: false,
		}
	}

	return res

}

func CO2SENSORHandler(cp *ChargePoint, p Payload) Payload {
	req := p.(*structs.REPORTCo2Sensor)
	log.Printf("\nid: %s\nREPORTCO2SENSOR: %v", cp.Id, req)
	log.Printf("device CO2SENSOR : %f", req.Value)
	// csms.Store(cp) // cp 저장

	var res Payload

	rb := device_handler.Device_Co2Sensor_update(req.Id, req.Value)

	if rb {
		res = structs.REPORTCo2SensorConf{
			Id:     req.Id,
			Result: true,
		}
	} else {
		res = structs.REPORTCo2SensorConf{
			Id:     req.Id,
			Result: false,
		}
	}

	return res

}

func CO2CREATORHandler(cp *ChargePoint, p Payload) Payload {
	req := p.(*structs.REPORTCo2Creator)
	log.Printf("\nid: %s\nREPORTCO2CREATOR: %v", cp.Id, req)
	log.Printf("device CO2CREATOR : %v", req.Value)
	// csms.Store(cp) // cp 저장

	var res Payload

	rb := device_handler.Device_Co2Creator_update(req.Id, req.Value)

	if rb {
		res = structs.REPORTCo2CreatorConf{
			Id:     req.Id,
			Result: true,
		}
	} else {
		res = structs.REPORTCo2CreatorConf{
			Id:     req.Id,
			Result: false,
		}
	}

	return res

}

func PHECSENSORHandler(cp *ChargePoint, p Payload) Payload {
	req := p.(*structs.REPORTPHECSensor)
	log.Printf("\nid: %s\nREPORTPHECSENSOR: %v", cp.Id, req)
	log.Printf("device PHECSENSOR : %f", req.Value)
	// csms.Store(cp) // cp 저장

	var res Payload

	rb := device_handler.Device_PHECSensor_update(req.Id, req.Value)

	if rb {
		res = structs.REPORTPHECSensorConf{
			Id:     req.Id,
			Result: true,
		}
	} else {
		res = structs.REPORTPHECSensorConf{
			Id:     req.Id,
			Result: false,
		}
	}

	return res

}

func LEDHandler(cp *ChargePoint, p Payload) Payload {
	req := p.(*structs.REPORTLED)
	log.Printf("\nid: %s\nREPORTLED: %v", cp.Id, req)
	log.Printf("device LED : %v", req.Value)
	// csms.Store(cp) // cp 저장

	var res Payload

	rb := device_handler.Device_LED_update(req.Id, req.Value)

	if rb {
		res = structs.REPORTLEDConf{
			Id:     req.Id,
			Result: true,
		}
	} else {
		res = structs.REPORTLEDConf{
			Id:     req.Id,
			Result: false,
		}
	}

	return res

}

func AutoHandler(cp *ChargePoint, p Payload) Payload {
	req := p.(*structs.REPORTAuto)
	log.Printf("\nid: %s\nAuto: %v", cp.Id, req)
	log.Printf("device Auto : %v", req.Value)
	// csms.Store(cp) // cp 저장

	var res Payload

	rb := device_handler.Device_state_update(req.Id, req.Value)

	if rb {
		res = structs.REPORTAutoConf{
			Id:     req.Id,
			Result: true,
		}
	} else {
		res = structs.REPORTAutoConf{
			Id:     req.Id,
			Result: false,
		}
	}

	return res

}

func ReportAllSensorHandler(cp *ChargePoint, p Payload) Payload {
	req := p.(*structs.REPORTAllSensor)
	log.Printf("\nid: %s\nReportAllSensor: %v", cp.Id, req)
	log.Printf("device all sensor : %v", req.Value)
	// csms.Store(cp) // cp 저장

	var res Payload

	rb := device_handler.Device_state_report_insert(req.Id, req.Value.Temp, req.Value.Humi, req.Value.Co2_Sensor, req.Value.Ph_Ec_Sensor, req.Value.Time)

	if rb {
		res = structs.REPORTAllSensorConf{
			Id:     req.Id,
			Result: true,
		}
	} else {
		res = structs.REPORTAllSensorConf{
			Id:     req.Id,
			Result: false,
		}
	}

	return res

}

// func AlarmHandler2(cp *ChargePoint, p Payload) Payload {
// 	req := p.(*structs.Alarm)
// 	log.Printf("\nid: %s\nAlarm: %v", cp.Id, req)
// 	log.Printf("device Alarm : %v", req.Body)
// 	// csms.Store(cp) // cp 저장

// 	var res Payload

// 	rb := sendFCM(req.Id, req.Title, req.Body, req.Data)
// 	if rb {
// 		res = structs.AlarmConf{
// 			Id:     req.Id,
// 			Result: true,
// 		}
// 	} else {
// 		res = structs.AlarmConf{
// 			Id:     req.Id,
// 			Result: false,
// 		}
// 	}

// 	return res

// }

func AlarmHandler(cp *ChargePoint, p Payload) Payload {
	req := p.(*structs.Alarm)
	log.Printf("\nid: %s\nAlarm: %v", cp.Id, req)
	// log.Printf("device Alarm : %v", req.Body)

	var res Payload

	os_ids, email_ids := device_handler.Device_alarm_user_select(req.Id)

	bb := AlarmInsert(req.Id, req.Title, req.Body, email_ids)

	rb := sendOneSignal(req.Title, req.Body, req.Data, os_ids)
	if rb && bb {
		res = structs.AlarmConf{
			Id:     req.Id,
			Result: true,
		}
	} else {
		res = structs.AlarmConf{
			Id:     req.Id,
			Result: false,
		}
	}

	return res

}

//////////////////////////////////////////////////////////////////////////////////////////////

func AlarmInsert(id string, title string, content string, user []string) bool {

	// 한국 시간대 로드
	loc, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		fmt.Println("시간대 로딩 에러:", err)
		return false
	}

	now := time.Now().In(loc)
	A_Time := now.Format("2006-01-02 15:04:05.000000")

	bb := device_handler.Device_alarm_insert(id, title, content, A_Time, user)

	return bb
}
