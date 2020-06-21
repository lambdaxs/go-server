package sms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/lambdaxs/go-server/lib/code"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type RLCouldClient struct {
	Host      string //https://app.cloopen.com:8883
	Sid       string //8a48b551516c09cd0151801ae8c02537
	AuthToken string //70fab7ca87d843b984adac6bbaa136fe
	AppID     string //8aaf07086df12352016df26a6def0234
	AppToken  string //18670cf40f94eee3cfe4670b3e447eba
	Client    http.Client
}

func NewRLCouldClient(host, sid, authToken, appID, appToken string) *RLCouldClient {
	return &RLCouldClient{
		Host:      host,
		Sid:       sid,
		AuthToken: authToken,
		AppID:     appID,
		AppToken:  appToken,
		Client: http.Client{
			Timeout: time.Second * 3,
		},
	}
}

type SMSContent struct {
	To         string   `json:"to"`
	AppID      string   `json:"appId"`
	TemplateID string   `json:"templateId"`
	Datas      []string `json:"datas"`
}

func (c *SMSContent) JSONBuffer() []byte {
	buf, _ := json.Marshal(c)
	return buf
}

type SMSResponse struct {
	StatusCode string `json:"statusCode"` //"000000"
}

//登陆验证短信
func (s *RLCouldClient) SendLoginCode(phone string, code string) error {
	return s.Send([]string{phone}, []string{code, "2分钟"}, "481362")
}

func (s *RLCouldClient) Send(to []string, datas []string, templateID string) error {
	timestamp := time.Now().Format("20060102150405")
	sig := strings.ToUpper(code.MD5Str(fmt.Sprintf("%s%s%s", s.Sid, s.AuthToken, timestamp)))
	content := SMSContent{
		To:         strings.Join(to, ","),
		AppID:      s.AppID,
		TemplateID: templateID,
		Datas:      datas,
	}
	req, err := http.NewRequest("POST",
		fmt.Sprintf("%s/2013-12-26/Accounts/%s/SMS/TemplateSMS?sig=%s", s.Host, s.Sid, sig),
		bytes.NewBuffer(content.JSONBuffer()))
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	req.Header.Set("Authorization", code.Base64Str(fmt.Sprintf("%s:%s", s.Sid, timestamp)))
	resp, err := s.Client.Do(req)
	if err != nil {
		return err
	}
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	resp.Body.Close()
	fmt.Println(string(buf))

	result := SMSResponse{}
	if err := json.Unmarshal(buf, &result); err != nil {
		return err
	}
	if result.StatusCode == "000000" {
		return nil
	}
	return fmt.Errorf("短信发送失败:%s", result.StatusCode)
}
