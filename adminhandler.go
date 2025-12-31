package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"smart_house/gosrc/adminpage"
	"smart_house/gosrc/clientpagedbconn"
	"smart_house/gosrc/structs"
	"time"
)

// 대쉬보드 정보
func IF_Dashboard_Page_Handler(w http.ResponseWriter, r *http.Request) {

	log.Print("IF_Dashboard_Page_Handler")

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
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(&adminpage.IFDashboardPage{List: adminpage.Device_list_select(ConnectedDeivceList_Handler(), "0", "3"), Appoint: IF_Device_Web_List_Appoint_Handler(), Daily: adminpage.Device_daily_list_select("0", "3"), Graph: adminpage.Device_daily_graph_select(StartTime, EndTime, Device_Id), DeviceUser: adminpage.Device_user_list_select("0", "3")})
	w.Write(json)
}

// 하우스 기본 정보
func IF_DeviceList_Page_Handler(w http.ResponseWriter, r *http.Request) {

	log.Printf("IF_DeviceList_Page_Handler %v,%v", r.URL.Query().Get("s_order"), r.URL.Query().Get("e_order"))

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(&structs.DeviceListPage{Count: adminpage.Device_list_count_select(), Items: adminpage.Device_list_select(ConnectedDeivceList_Handler(), r.URL.Query().Get("s_order"), r.URL.Query().Get("e_order"))})
	w.Write(json)
}

// 하우스 기본 정보 검색
func IF_DeviceList_Page_Search_Handler(w http.ResponseWriter, r *http.Request) {

	log.Printf("DeviceList_Page_Search_Handler %v,%v", r.URL.Query().Get("column"), r.URL.Query().Get("data"))

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	// json, _ := json.Marshal(&structs.DeviceListPage{Items: clientpagedbconn.Device_list_select_search(r.URL.Query().Get("column"), r.URL.Query().Get("data"))})
	json, _ := json.Marshal(&structs.DeviceListPage{Items: adminpage.Device_list_select_search(ConnectedDeivceList_Handler(), r.URL.Query().Get("column"), r.URL.Query().Get("data"))})
	w.Write(json)
}

// 하우스 예약 정보
func IF_DeviceAppointList_Page_Handler(w http.ResponseWriter, r *http.Request) {

	log.Printf("IF_DeviceAppointList_Page_Handler")

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(&structs.DeviceAppointListPage{Items: IF_Device_Web_List_Appoint_Handler()})
	w.Write(json)
}

// 하우스 등록 유저 기기 정보
func IF_DeviceUserList_Page_Handler(w http.ResponseWriter, r *http.Request) {

	log.Printf("IF_DeviceUserList_Page_Handler %v,%v", r.URL.Query().Get("s_order"), r.URL.Query().Get("e_order"))

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(&adminpage.DeviceUserPage{Count: adminpage.Device_user_list_count_select(), Items: adminpage.Device_user_list_select(r.URL.Query().Get("s_order"), r.URL.Query().Get("e_order"))})
	w.Write(json)
}

// 하우스 등록 유저 기기 정보 검색
func IF_DeviceUserList_Page_Search_Handler(w http.ResponseWriter, r *http.Request) {

	log.Printf("IF_DeviceUserList_Page_Search_Handler %v,%v", r.URL.Query().Get("column"), r.URL.Query().Get("data"))

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(&adminpage.DeviceUserPage{Items: adminpage.Device_user_list_select_search(r.URL.Query().Get("column"), r.URL.Query().Get("data"))})
	w.Write(json)
}

// 유저 정보
func IF_UserList_Page_Handler(w http.ResponseWriter, r *http.Request) {

	log.Printf("IF_UserList_Page_Handler %v,%v", r.URL.Query().Get("s_order"), r.URL.Query().Get("e_order"))

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(&adminpage.UserPage{Count: adminpage.User_list_count_select(), Items: adminpage.User_list_select(r.URL.Query().Get("s_order"), r.URL.Query().Get("e_order"))})
	w.Write(json)
}

// 유저 정보 검색
func IF_UserList_Page_Search_Handler(w http.ResponseWriter, r *http.Request) {

	log.Printf("IF_UserList_Page_Search_Handler %v,%v", r.URL.Query().Get("column"), r.URL.Query().Get("data"))

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(&adminpage.UserPage{Items: adminpage.User_list_select_search(r.URL.Query().Get("column"), r.URL.Query().Get("data"))})
	w.Write(json)
}

// 기기 예약 정보
func IF_Device_Web_List_Appoint_Handler() structs.DeviceAppointItem {

	log.Printf("IF_Device_Web_List_Appoint_Handler")

	var devicelist = adminpage.DeviceSimple_list(ConnectedDeivceList_Handler())

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

// 하우스 업데이트 관리
func IF_Updmng_Page_Handler(w http.ResponseWriter, r *http.Request) {

	log.Printf("IF_Updmng_Page_Handler %v,%v", r.URL.Query().Get("s_order"), r.URL.Query().Get("e_order"))

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(&adminpage.UpdmngPage{Count: adminpage.Upd_list_count_select(), Items: adminpage.Upd_list_select(ConnectedDeivceList_Handler(), r.URL.Query().Get("s_order"), r.URL.Query().Get("e_order"))})
	w.Write(json)
}

// 하우스 업데이트 관리 검색
func IF_Updmng_Page_Search_Handler(w http.ResponseWriter, r *http.Request) {

	log.Printf("IF_Updmng_Page_Search_Handler %v,%v", r.URL.Query().Get("column"), r.URL.Query().Get("data"))

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(&adminpage.UpdmngPage{Items: adminpage.Upd_list_select_search(ConnectedDeivceList_Handler(), r.URL.Query().Get("column"), r.URL.Query().Get("data"))})
	w.Write(json)
}

// 하우스 로그 관리
func IF_Logmng_Page_Handler(w http.ResponseWriter, r *http.Request) {

	log.Printf("IF_Logmng_Page_Handler %v,%v", r.URL.Query().Get("s_order"), r.URL.Query().Get("e_order"))

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(&adminpage.LogmngPage{Count: adminpage.Log_list_count_select(), Items: adminpage.Log_list_select(ConnectedDeivceList_Handler(), r.URL.Query().Get("s_order"), r.URL.Query().Get("e_order"))})
	w.Write(json)
}

// 하우스 로그 관리 검색
func IF_Logmng_Page_Search_Handler(w http.ResponseWriter, r *http.Request) {

	log.Printf("IF_Logmng_Page_Search_Handler %v,%v", r.URL.Query().Get("column"), r.URL.Query().Get("data"))

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(&adminpage.LogmngPage{Items: adminpage.Log_list_select_search(ConnectedDeivceList_Handler(), r.URL.Query().Get("column"), r.URL.Query().Get("data"))})
	w.Write(json)
}

// 하우스 일자별 센서값 정보
func IF_DeviceDailylist_Page_Handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("IF_DeviceDailylist_Page_Handler %v,%v", r.URL.Query().Get("s_order"), r.URL.Query().Get("e_order"))

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(&structs.DeviceDailyListPage{Count: adminpage.Device_daily_list_count_select(), Items: adminpage.Device_daily_list_select(r.URL.Query().Get("s_order"), r.URL.Query().Get("e_order"))})
	w.Write(json)
}

// 하우스 일자별 센서값 검색
func IF_DeviceDailylist_Page_Search_Handler(w http.ResponseWriter, r *http.Request) {

	log.Printf("IF_DeviceDailylist_Page_Search_Handler %v,%v", r.URL.Query().Get("column"), r.URL.Query().Get("data"))

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(&structs.DeviceDailyListPage{Items: adminpage.Device_daily_list_select_search(r.URL.Query().Get("column"), r.URL.Query().Get("data"))})
	w.Write(json)
}

// 하우스 센서값 그래프 조회
func IF_DeviceDailyGraph_Page_Handler(w http.ResponseWriter, r *http.Request) {

	log.Print("IF_DeviceDailyGraph_Page_Handler")

	msg := structs.DeviceGraphTimeData{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("IF_DeviceDailyGraph_Page_Handler err : %s", err)
		return
	}
	defer r.Body.Close()

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(&structs.DeviceDailyGraphPage{ListItem: adminpage.Device_daily_graph_list_select(), Items: adminpage.Device_daily_graph_select(msg.StartTime, msg.EndTime, msg.Id)})
	w.Write(json)
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// 하우스 등록 유저 기기 정보 삽입
func IF_DeviceUser_Register_Handler(w http.ResponseWriter, r *http.Request) {

	log.Print("IF_DeviceUser_Register_Handler")

	msg := adminpage.DeviceUserListItem{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("IF_DeviceUser_Register_Handler err : %s", err)
		return
	}
	defer r.Body.Close()

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(&adminpage.DeviceUserConf{Item: adminpage.Insert_deviceuser(msg.Id, msg.User)})
	w.Write(json)
}

// 하우스 등록 유저 기기 정보 삭제
func IF_DeviceUser_Deleter_Handler(w http.ResponseWriter, r *http.Request) {

	log.Print("IF_DeviceUser_Deleter_Handler")

	msg := adminpage.DeviceUserListItem{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("IF_DeviceUser_Deleter_Handler err : %s", err)
		return
	}
	defer r.Body.Close()

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(&adminpage.DeviceUserConf{Item: adminpage.Delete_deviceuser(msg.Id, msg.User)})
	w.Write(json)
}

// 유저 정보 삽입
func IF_User_Register_Handler(w http.ResponseWriter, r *http.Request) {

	log.Print("IF_User_Register_Handler")

	msg := adminpage.UserInsertItem{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("IF_User_Register_Handler err : %s", err)
		return
	}
	defer r.Body.Close()

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(&adminpage.UserConf{Item: adminpage.Insert_user(msg.Email, msg.Name, msg.Password)})
	w.Write(json)
}

// 유저 정보 삭제
func IF_User_Deleter_Handler(w http.ResponseWriter, r *http.Request) {

	log.Print("IF_User_Deleter_Handler")

	msg := adminpage.UserListItem{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("IF_User_Deleter_Handler err : %s", err)
		return
	}
	defer r.Body.Close()

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(&adminpage.UserConf{Item: adminpage.Delete_user(msg.Email, msg.Name)})
	w.Write(json)
}

// 유저 정보 갱신
func IF_User_Updater_Handler(w http.ResponseWriter, r *http.Request) {

	log.Print("IF_User_Updater_Handler")

	msg := adminpage.UserInsertItem{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("IF_User_Updater_Handler err : %s", err)
		return
	}
	defer r.Body.Close()

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(&adminpage.UserConf{Item: adminpage.Update_user(msg.Email, msg.Name, msg.Password)})
	w.Write(json)
}

// 유저 알람 전송
func IF_User_Alert_Handler(w http.ResponseWriter, r *http.Request) {

	log.Print("IF_User_Alert_Handler")

	msg := adminpage.UserAlertMessage{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("IF_User_Alert_Handler err : %s", err)
		return
	}
	defer r.Body.Close()

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(&adminpage.UserConf{Item: sendOneSignal(msg.Title, msg.Content, nil, []string{adminpage.User_AlertId_select(msg.User)})})
	w.Write(json)
}

// 유저 전체 알람 전송
func IF_User_Whole_Alert_Handler(w http.ResponseWriter, r *http.Request) {

	log.Print("IF_User_Whole_Alert_Handler")

	msg := adminpage.UserAlertMessage{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("IF_User_Whole_Alert_Handler err : %s", err)
		return
	}
	defer r.Body.Close()

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(&adminpage.UserConf{Item: sendWholeOneSignal(msg.Title, msg.Content, nil)})
	w.Write(json)
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// 하우스 로그 업로드 명령
func IF_Device_Web_Commad_Log_Upload_Handler(w http.ResponseWriter, r *http.Request) {

	msg := adminpage.CMDLogUpload{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("IF_Device_Web_Commad_Log_Upload_Handler : %s", err)
		return
	}
	defer r.Body.Close()

	log.Printf("IF_Device_Web_Commad_Log_Upload_Handler msg : %v", msg)
	var returnValue *adminpage.CMDLogUploadConf

	if csms != nil {
		cp, b := csms.Load(msg.Id)
		if b {

			////////////////GetDiagnostics	// 진단 파일 가져오기 (로그 등등)
			// startTime이 없으면 최근/stoptime 기준 3일전까지의 데이터
			// stoptime이 없으면 최근까지의 데이터
			var req Payload = adminpage.CMDLogUpload{
				Id:        msg.Id,
				Location:  msg.Location, // 두개짜리
				StartTime: msg.StartTime,
				StopTime:  msg.StopTime,
			}
			res, err := cp.Call("CMDLOGUPLOAD", req)
			if err != nil {
				log.Printf("aa %s", err)
			}
			log.Print(res) // 리턴값 파일명. 없으면 오류가 났거나 해당 조건에 맞는 진단이 없어 업로드 안함

			returnValue = res.(*adminpage.CMDLogUploadConf)
		}

		if cp == nil {
			log.Print("cp is nil")
			returnValue = &adminpage.CMDLogUploadConf{}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(returnValue)
}

// 하우스 업데이트 명령
func IF_Device_Web_Commad_Program_Update_Handler(w http.ResponseWriter, r *http.Request) {

	msg := adminpage.CMDPrgUpdate{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("IF_Device_Web_Commad_Program_Update_Handler : %s", err)
		return
	}
	defer r.Body.Close()

	log.Printf("IF_Device_Web_Commad_Program_Update_Handler msg : %v", msg)
	var returnValue *adminpage.CMDPrgUpdateConf

	if csms != nil {
		cp, b := csms.Load(msg.Id)
		if b {

			////////////////GetDiagnostics	// 진단 파일 가져오기 (로그 등등)
			// startTime이 없으면 최근/stoptime 기준 3일전까지의 데이터
			// stoptime이 없으면 최근까지의 데이터
			var req Payload = adminpage.CMDPrgUpdate{
				Id:       msg.Id,
				Location: msg.Location, // 두개짜리
			}
			res, err := cp.Call("CMDPRGUPLOAD", req)
			if err != nil {
				log.Printf("aa %s", err)
			}
			log.Print(res) // 리턴값 파일명. 없으면 오류가 났거나 해당 조건에 맞는 진단이 없어 업로드 안함

			returnValue = res.(*adminpage.CMDPrgUpdateConf)
		}

		if cp == nil {
			log.Print("cp is nil")
			returnValue = &adminpage.CMDPrgUpdateConf{}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(returnValue)
}
