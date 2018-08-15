package aliim

import (
	"encoding/json"
	"log"
	"time"

	"github.com/tidwall/gjson"
)

type ImUserInfo struct {
	Userid   string `json:"userid"`
	Password string `json:"password"`
	Name     string `json:"name"`
	IconUrl  string `json:"icon_url"`
}

type UidSucc struct {
	Uids []string `json:"string"`
}
type UidFail struct {
	Uids []string `json:"string"`
}
type FailMsg struct {
	FailMsg []string `json:"string"`
}
type UserAddResponse struct {
	UidSucc UidSucc `json:"uid_succ"`
	UidFail UidFail `json:"uid_fail"`
	FailMsg FailMsg `json:"fail_msg"`
}

type DeleteMsg struct {
	Msg []string `json:"string"`
}
type UserDeleteResponse struct {
	DeleteMsg DeleteMsg `json:"result"`
}

func getCommonParams() map[string]string {

	params := make(map[string]string)
	params["app_key"] = config.AppKey
	params["format"] = "json"
	params["sign_method"] = "md5"
	params["timestamp"] = time.Now().Format("2006-01-02 15:04:05")
	params["v"] = "2.0"
	return params
}

// 导入用户
func SendAddUsers(imUserInfos []ImUserInfo) (success bool, response string) {

	for _, user := range imUserInfos {
		if user.Userid == "" || user.Password == "" {
			return false, "userid or password is required"
		}
	}
	params := getCommonParams()
	params["method"] = OpenImUserAdd

	result, err := json.Marshal(imUserInfos)
	if err != nil {
		return false, err.Error()
	}
	params["userinfos"] = string(result)

	succ, resData := IMPost(params)

	if succ == false {
		return false, response
	}

	type Result struct {
		Result UserAddResponse `json:"openim_users_add_response"`
	}

	var resultResponse Result
	err = json.Unmarshal(resData, &resultResponse)
	if err != nil {
		//log.Println("err   " + err.Error())
		return false, err.Error()
	}

	log.Println("resData   " + string(resData))
	// {"error_response":{"code":29,"msg":"Invalid app Key","sub_code":"isv.appkey-not-exists","request_id":"fwx2977sp1lq"}}
	// {"openim_users_add_response":{"fail_msg":{},"uid_fail":{},"uid_succ":{"string":["5b714d8cc0b6fa682e552300"]}}}
	if gjson.GetBytes(resData, "error_response.code").String() != "" {
		return false, gjson.GetBytes(resData, "error_response.msg").String()
	}
	failMsg := resultResponse.Result.FailMsg
	if len(failMsg.FailMsg) <= 0 {
		return true, "add success"
	}
	return false, failMsg.FailMsg[0]
}

func SendDeleteUsers(userids string) (success bool, response string) {
	if userids == "" {
		return false, "userid is required"
	}
	params := getCommonParams()
	params["method"] = OpenImUserDelete
	params["userids"] = userids

	succ, resData := IMPost(params)
	//log.Println("resData " + string(resData))
	if succ == false {
		return false, response
	}

	type Result struct {
		UserDeleteResponse UserDeleteResponse `json:"openim_users_delete_response"`
	}
	log.Println("resData   " + string(resData))
	var resultResponse Result
	err := json.Unmarshal(resData, &resultResponse)
	if err != nil {
		return false, err.Error()
	}
	return true, "ok"
}

func SendUpdateUsers(imUserInfos []ImUserInfo) (success bool, response string) {
	for _, user := range imUserInfos {
		if user.Userid == "" {
			return false, "userid is required"
		}
	}
	params := getCommonParams()
	params["method"] = OpenImUserUpdate

	result, err := json.Marshal(imUserInfos)
	if err != nil {
		return false, err.Error()
	}
	params["userinfos"] = string(result)

	succ, resData := IMPost(params)

	if succ == false {
		return false, response
	}

	type Result struct {
		Result UserAddResponse `json:"openim_users_update_response"`
	}
	log.Println("resData   " + string(resData))
	var resultResponse Result
	err = json.Unmarshal(resData, &resultResponse)
	if err != nil {
		//log.Println("err   " + err.Error())
		return false, err.Error()
	}
	failMsg := resultResponse.Result.FailMsg
	if len(failMsg.FailMsg) <= 0 {
		return true, "update success"
	}
	return false, failMsg.FailMsg[0]
}

type CustMsg struct {
	FromUser  string   `json:"from_user"`
	ToUsers   []string `json:"to_users"`
	Summary   string   `json:"summary"`
	Data      string   `json:"data"`
	Aps       string   `json:"aps"` // {"alert":"ios apns push"}
	ApnsParam string   `json:"apns_param"`
	Invisible int32    `json:"invisible"`
	FromNick  string   `json:"from_nick"`
}

func SendCustmsgPush(msg *CustMsg) (success bool, response string) {
	params := getCommonParams()
	params["method"] = OpenImCustmsgPush

	result, err := json.Marshal(*msg)
	if err != nil {
		return false, err.Error()
	}
	params["custmsg"] = string(result)

	succ, resData := IMPost(params)
	return succ, string(resData)
}

type ImMsg struct {
	FromUser   string   `json:"from_user"`
	ToUsers    []string `json:"to_users"`
	MsgType    int32    `json:"msg_type"`
	Context    string   `json:"context"`
	MediaAttr  string   `json:"media_attr"`
	FromTaobao int32    `json:"from_taobao"`
}

func SendImmsgPush(msg *ImMsg) (success bool, response string) {
	params := getCommonParams()
	params["method"] = OpenImImmsgPush

	result, err := json.Marshal(*msg)
	if err != nil {
		return false, err.Error()
	}
	params["immsg"] = string(result)

	succ, resData := IMPost(params)
	return succ, string(resData)
}
