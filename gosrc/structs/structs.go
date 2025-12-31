package structs

type ContextStruct struct{}

///////////////////////////////

// Post data 구조체
type LoginData struct { //로그인 데이터
	Email string `json:"email"`
	PWD   string `json:"pwd"`
}

type LogoutDataForm struct { //로그아웃 데이터
	Email string `json:"email"`
	Refresh_Token   string `json:"refresh_token"`
}

type LogoutResult struct { //로그아웃 결과
	Success bool `json:"success"`
}

//////////////////////////////////////////////////////
// 기기

type LoginDeviceData struct { //기기 로그인 데이터
	DeviceId string `json:"device_id"`
}

type LoginDevice struct { // 로그인
	Items LoginDeviceDataForm `json:"item"`
}

type LoginDeviceDataForm struct {
	Success bool   `json:"success"`
	DeviceId string `json:"device_id"`
	TokenPair TokenPair `json:"token_pair"`
}

///////////////////////////////

// 페이지 구조체
type Page struct {
	Title string
}

type LoginPage struct { // 사용자 로그인 페이지
	Items LoginDataForm `json:"item"`
}

// /////////////////////////////		// 페이지 정보 구조체

type LoginDataForm struct { // 로그인 시도 대답 폼
	Success bool   `json:"success"`
	Email string `json:"email"`
	TokenPair TokenPair `json:"token_pair"`
	//Name  string `json:"name"`
}

//////////////////////////////// jwt 토큰 구조체

type RenewedAccessToken struct {	// 리프레시 토큰으로 엑세스 토큰 재발급 대답 폼
	Success bool   `json:"success"`
	Email string `json:"email"`
	TokenPair TokenPair `json:"token_pair"`
}

type TokenPair struct {	// 토큰 페어
	Access_Token  string `json:"access_token"`
	Refresh_Token string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type RefreshTokenRequest struct {	//웹 리프레시 토큰 요청
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenAppRequest struct {	//앱 리프레시 토큰 요청
	RefreshToken string `json:"refresh_token"`
}

////////////////////////////////////////////////////////////////////////////////////////////
//Web

type DashboardPage struct { // 대쉬보드 페이지
	List    []DeviceDetail         `json:"list"`
	Appoint DeviceAppointItem      `json:"appoint"`
	Daily   []DeviceDailyListItem  `json:"daily"`
	Graph   []DeviceDailyGraphItem `json:"graph"`
}

type DeviceListPage struct { // 하우스 기본정보 페이지
	Count int            `json:"count"`
	Items []DeviceDetail `json:"item"`
}

type DeviceAppointListPage struct { // 하우스 예약 정보 페이지
	Items DeviceAppointItem `json:"items"`
}

type DeviceAppointItem struct {
	Devices []DeviceSimple  `json:"devices"`
	Id      string          `json:"id"`
	Auto    int             `json:"auto"`
	Items   []DeviceAppoint `json:"items"`
}

type DeviceDailyListPage struct { // 일자별 센서값 조회 페이지
	Count int                   `json:"count"`
	Items []DeviceDailyListItem `json:"items"`
}

type DeviceDailyListItem struct { // 일자별 센서값 조회 구조체
	Id           string  `json:"id"`
	Temp         float32 `json:"temp"`
	Humi         float32 `json:"humi"`
	Co2_Sensor   float32 `json:"co2_sensor"`
	Ph_Ec_Sensor float32 `json:"ph_ec_sensor"`
	Time         string  `json:"time"`
}

type DeviceGraphTimeData struct { //센서값 그래프 조회 데이터
	StartTime string `json:"starttime"`
	EndTime   string `json:"endtime"`
	Id        string `json:"id"`
}

type DeviceDailyGraphPage struct { // 센서값 그래프 조회 페이지
	ListItem []DeviceDailyGraphListItem `json:"listitem"`
	Items    []DeviceDailyGraphItem     `json:"item"`
}

type DeviceDailyGraphListItem struct { // 센서값 그래프 기기 리스트 조회 구조체
	Id string `json:"id"`
}

type DeviceDailyGraphItem struct { // 센서값 그래프 조회 구조체
	Temp         float32 `json:"temp"`
	Humi         float32 `json:"humi"`
	Co2_Sensor   float32 `json:"co2_sensor"`
	Ph_Ec_Sensor float32 `json:"ph_ec_sensor"`
	Time         string  `json:"time"`
}

type CMDAllOnoff struct { // onoff 명령
	Id    string `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

type CMDAllOnoffConf struct { // onoff 명령
	Id     string `json:"id"`
	Result bool   `json:"result"`
}

type CMDRebootSoftware struct { // 소프트웨어 재시작 명령
	Id string `json:"id"`
}

type CMDRebootSoftwareConf struct { // 소프트웨어 재시작 명령
	Id     string `json:"id"`
	Result bool   `json:"result"`
}

type CMDRebootHardware struct { // 하드웨어 재시작 명령
	Id string `json:"id"`
}

type CMDRebootHardwareConf struct { // 하드웨어 재시작 명령
	Id     string `json:"id"`
	Result bool   `json:"result"`
}

////////////////////////////////////////////////////////////////////////////////////////////
//App

type LoginApp struct { // 로그인
	Items LoginAppDataForm `json:"item"`
}

type LoginAppDataForm struct {
	Success bool   `json:"success"`
	Email string `json:"email"`
	TokenPair TokenPair `json:"token_pair"`
}

type SendPushTitleAndContents struct {
	Ko string `json:"ko"`
	En string `json:"en"`
}

type SendPush struct {
	App_Id             string                   `json:"app_id"`
	IsAndroid          bool                     `json:"isAndroid"`
	Include_Player_IDs []string                 `json:"include_player_ids"`
	Contents           SendPushTitleAndContents `json:"contents"`
	Headings           SendPushTitleAndContents `json:"headings"`
}

type UserAlarmList struct {
	Id      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Time    string `json:"time"`
	Read    bool   `json:"read"`
}

type UserAlarmListReq struct {
	User string `json:"user"`
}

type UserAlarmListConf struct {
	User  string          `json:"user"`
	Items []UserAlarmList `json:"items"`
}

type ReportOneSignalId struct {
	Email string `json:"email"`
	Os_Id string `json:"os_id"`
}

type ReportOneSignalIdConf struct {
	Os_Id  string `json:"os_id"`
	Result bool   `json:"result"`
}

type DeviceDetail struct {
	Id           string  `json:"id"`
	Comp         bool    `json:"comp"`
	Pump         bool    `json:"pump"`
	Temp         float32 `json:"temp"`
	Humi         float32 `json:"humi"`
	F_Fan        bool    `json:"f_fan"`
	E_Fan        bool    `json:"e_fan"`
	Aircon       bool    `json:"aircon"`
	Co2_Sensor   float32 `json:"co2_sensor"`
	C02_Creator  bool    `json:"co2_creator"`
	Ph_Ec_Sensor float32 `json:"ph_ec_sensor"`
	Led          bool    `json:"led"`
	Cctv         string  `json:"cctv"`
	Auto         int     `json:"auto"`
	B_Time       string  `json:"b_time"`
}

type Boot struct {
	Id       string          `json:"id"`
	State    BootState       `json:"state"`
	Appoints []DeviceAppoint `json:"appoints"`
}

type BootState struct {
	Comp         bool    `json:"comp"`
	Pump         bool    `json:"pump"`
	Temp         float32 `json:"temp"`
	Humi         float32 `json:"humi"`
	F_Fan        bool    `json:"f_fan"`
	E_Fan        bool    `json:"e_fan"`
	Aircon       bool    `json:"aircon"`
	Co2_Sensor   float32 `json:"co2_sensor"`
	C02_Creator  bool    `json:"co2_creator"`
	Ph_Ec_Sensor float32 `json:"ph_ec_sensor"`
	Led          bool    `json:"led"`
	Cctv         string  `json:"cctv"`
	Auto         bool    `json:"auto"`
}

type BootConf struct {
	Id     string `json:"id"`
	Result bool   `json:"result"`
}

type DeviceSimpleList struct {
	Items []DeviceSimple `json:"items"`
}

type DeviceSimple struct {
	Id   string  `json:"id"`
	Temp float32 `json:"temp"`
	Humi float32 `json:"humi"`
	Cctv string  `json:"cctv"`
	Auto int     `json:"auto"`
}

type CMDComp struct {
	Id    string `json:"id"`
	Value bool   `json:"value"`
}

type CMDCompConf struct {
	Id     string `json:"id"`
	Result bool   `json:"result"`
}

type CMDPump struct {
	Id    string `json:"id"`
	Value bool   `json:"value"`
}

type CMDPumpConf struct {
	Id     string `json:"id"`
	Result bool   `json:"result"`
}

type CMDFFan struct {
	Id    string `json:"id"`
	Value bool   `json:"value"`
}

type CMDFFanConf struct {
	Id     string `json:"id"`
	Result bool   `json:"result"`
}

type CMDEFan struct {
	Id    string `json:"id"`
	Value bool   `json:"value"`
}

type CMDEFanConf struct {
	Id     string `json:"id"`
	Result bool   `json:"result"`
}

type CMDAircon struct {
	Id    string `json:"id"`
	Value bool   `json:"value"`
}

type CMDAirconConf struct {
	Id     string `json:"id"`
	Result bool   `json:"result"`
}

type CMDCo2Creator struct {
	Id    string `json:"id"`
	Value bool   `json:"value"`
}

type CMDCo2CreatorConf struct {
	Id     string `json:"id"`
	Result bool   `json:"result"`
}

type CMDLED struct {
	Id    string `json:"id"`
	Value bool   `json:"value"`
}

type CMDLEDConf struct {
	Id     string `json:"id"`
	Result bool   `json:"result"`
}

type REPORTComp struct {
	Id    string `json:"id"`
	Value bool   `json:"value"`
}

type REPORTCompConf struct {
	Id     string `json:"id"`
	Result bool   `json:"result"`
}

type REPORTPump struct {
	Id    string `json:"id"`
	Value bool   `json:"value"`
}

type REPORTPumpConf struct {
	Id     string `json:"id"`
	Result bool   `json:"result"`
}

type REPORTFFan struct {
	Id    string `json:"id"`
	Value bool   `json:"value"`
}

type REPORTFFanConf struct {
	Id     string `json:"id"`
	Result bool   `json:"result"`
}

type REPORTEFan struct {
	Id    string `json:"id"`
	Value bool   `json:"value"`
}

type REPORTEFanConf struct {
	Id     string `json:"id"`
	Result bool   `json:"result"`
}

type REPORTAircon struct {
	Id    string `json:"id"`
	Value bool   `json:"value"`
}

type REPORTAirconConf struct {
	Id     string `json:"id"`
	Result bool   `json:"result"`
}

type REPORTCo2Creator struct {
	Id    string `json:"id"`
	Value bool   `json:"value"`
}

type REPORTCo2CreatorConf struct {
	Id     string `json:"id"`
	Result bool   `json:"result"`
}

type REPORTLED struct {
	Id    string `json:"id"`
	Value bool   `json:"value"`
}

type REPORTLEDConf struct {
	Id     string `json:"id"`
	Result bool   `json:"result"`
}

type REPORTTemp struct {
	Id    string  `json:"id"`
	Value float32 `json:"value"`
}

type REPORTTempConf struct {
	Id     string `json:"id"`
	Result bool   `json:"result"`
}

type REPORTHumi struct {
	Id    string  `json:"id"`
	Value float32 `json:"value"`
}

type REPORTHumiConf struct {
	Id     string `json:"id"`
	Result bool   `json:"result"`
}

type REPORTCo2Sensor struct {
	Id    string  `json:"id"`
	Value float32 `json:"value"`
}

type REPORTCo2SensorConf struct {
	Id     string `json:"id"`
	Result bool   `json:"result"`
}

type REPORTPHECSensor struct {
	Id    string  `json:"id"`
	Value float32 `json:"value"`
}

type REPORTPHECSensorConf struct {
	Id     string `json:"id"`
	Result bool   `json:"result"`
}

type REPORTAuto struct {
	Id    string `json:"id"`
	Value int    `json:"value"`
}

type REPORTAutoConf struct {
	Id     string `json:"id"`
	Result bool   `json:"result"`
}

type REPORTAllSensor struct {
	Id    string               `json:"id"`
	Value REPORTAllSensorItems `json:"value"`
}

type REPORTAllSensorConf struct {
	Id     string `json:"id"`
	Result bool   `json:"result"`
}

type REPORTAllSensorItems struct {
	Id           string  `json:"id"`
	Temp         float32 `json:"temp"`
	Humi         float32 `json:"humi"`
	Co2_Sensor   float32 `json:"co2_sensor"`
	Ph_Ec_Sensor float32 `json:"ph_ec_sensor"`
	Time         string  `json:"time"`
}

type Alarm struct {
	Id    string            `json:"id"`
	Title string            `json:"title"`
	Body  string            `json:"body"`
	Data  map[string]string `json:"data"`
}

type AlarmConf struct {
	Id     string `json:"id"`
	Result bool   `json:"result"`
}

type Heartbeat struct{}

type HeartbeatConf struct {
	CurrentTime string `json:"currentTime" validate:"required,ISO8601date"`
}

// //////////////////////////////////////////// 예약 //////////////////////////////////////////////////////////////////////////////////////////

type DeviceAppointAllRequest struct {
	Id string `json:"id"`
}

type DeviceAppointList struct {
	Id    string          `json:"id"`
	Auto  int             `json:"auto"`
	Items []DeviceAppoint `json:"items"`
}

type DeviceAppoint struct {
	Name  string `json:"name"`
	Start string `json:"startt"`
	Stop  string `json:"stopt"`
	Onoff bool   `json:"onoff"`
}

type ACMDComPump struct {
	Id    string `json:"id"`
	Start string `json:"startt"`
	Stop  string `json:"stopt"`
	Onoff bool   `json:"onoff"`
}

type ACMDFFan struct {
	Id    string `json:"id"`
	Start string `json:"startt"`
	Stop  string `json:"stopt"`
	Onoff bool   `json:"onoff"`
}

type ACMDEFan struct {
	Id    string `json:"id"`
	Start string `json:"startt"`
	Stop  string `json:"stopt"`
	Onoff bool   `json:"onoff"`
}

type ACMDAircon struct {
	Id    string `json:"id"`
	Start string `json:"startt"`
	Stop  string `json:"stopt"`
	Onoff bool   `json:"onoff"`
}

type ACMDCo2Creator struct {
	Id    string `json:"id"`
	Start string `json:"startt"`
	Stop  string `json:"stopt"`
	Onoff bool   `json:"onoff"`
}

type ACMDLED struct {
	Id    string `json:"id"`
	Start string `json:"startt"`
	Stop  string `json:"stopt"`
	Onoff bool   `json:"onoff"`
}
