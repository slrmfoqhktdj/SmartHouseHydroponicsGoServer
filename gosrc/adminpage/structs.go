package adminpage

import (
	"smart_house/gosrc/structs"
)

type IFDashboardPage struct { // 대쉬보드 페이지
	List       []structs.DeviceDetail         `json:"list"`
	Appoint    structs.DeviceAppointItem      `json:"appoint"`
	Daily      []structs.DeviceDailyListItem  `json:"daily"`
	Graph      []structs.DeviceDailyGraphItem `json:"graph"`
	DeviceUser []DeviceUserListItem           `json:"deviceuser"`
}

type UpdmngPage struct { // 로그 관리 페이지
	Count int                 `json:"count"`
	Items []DeviceUpdListItem `json:"item"`
}

type LogmngPage struct { // 로그 관리 페이지
	Count int                 `json:"count"`
	Items []DeviceLogListItem `json:"item"`
}

type DeviceUserPage struct { // 등록 유저 기기 관리 페이지
	Count int                  `json:"count"`
	Items []DeviceUserListItem `json:"item"`
}

type UserPage struct { // 유저 관리 페이지
	Count int            `json:"count"`
	Items []UserListItem `json:"item"`
}

type DeviceUserListItem struct {
	Id     string `json:"id"`
	User   string `json:"user"`
	R_Time string `json:"r_time"`
}

type UserListItem struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type UserInsertItem struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type UserAlertMessage struct {
	User    string `json:"user"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type DeviceUserConf struct { // 등록 유저 기기 관리 페이지 명령 여부
	Item bool `json:"item"`
}

type UserConf struct { // 유저 관리 페이지 명령 여부
	Item bool `json:"item"`
}

type DeviceUpdListItem struct {
	Id     string `json:"id"`
	Auto   int    `json:"auto"`
	B_Time string `json:"b_time"`
}

type DeviceLogListItem struct {
	Id   string `json:"id"`
	Auto int    `json:"auto"`
}

type CMDLogUpload struct {
	Id        string `json:"id"`
	Location  string `json:"location" validate:"required"`
	StartTime string `json:"startTime"`
	StopTime  string `json:"stopTime"`
}

type CMDLogUploadConf struct {
	FileName string `json:"fileName"`
}

type CMDPrgUpdate struct {
	Id       string `json:"id"`
	Location string `json:"location" validate:"required"`
}

type CMDPrgUpdateConf struct {
	Id     string `json:"id"`
	Result bool   `json:"result"`
}
