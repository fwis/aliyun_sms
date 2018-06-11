package aliyun_sms

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/fwis/golib/srand"
)

type AliyunSms struct {
	AccessKeyID  string
	AccessSecret string

	SignName     string //签名名称
	TemplateCode string //模板code
}

type SendResponse struct {
	RequestId string
	Code      string
	Message   string
	BizId     string
}

func NewAliyunSms(sign_name string, template_code string, access_key_id string, access_secret string) (*AliyunSms, error) {
	var a AliyunSms
	a.SignName = sign_name
	a.TemplateCode = template_code
	a.AccessKeyID = access_key_id
	a.AccessSecret = access_secret

	return &a, nil
}

func (this *AliyunSms) Send(numbers string, params string) (*SendResponse, error) {
	var request Request
	request.Format = "JSON"
	request.Version = "2017-05-25"
	request.AccessKeyId = this.AccessKeyID
	request.SignatureMethod = "HMAC-SHA1"
	request.Timestamp = time.Now().UTC().Format("2006-01-02T15:04:05Z")
	request.SignatureVersion = "1.0"
	request.SignatureNonce = srand.RandUUID().Hex()
	fmt.Printf("request.SignatureNonce=%#v\n", request.SignatureNonce)

	request.Action = "SendSms"
	request.SignName = this.SignName
	request.TemplateCode = this.TemplateCode
	request.PhoneNumbers = numbers
	request.TemplateParam = params
	request.RegionId = "cn-hangzhou"

	url := request.ComposeUrl("GET", this.AccessSecret)
	var resp *http.Response
	var err error
	resp, err = http.Get(url)
	if err != nil {
		return nil, err
	}
	var b []byte
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	response := SendResponse{}
	err = json.Unmarshal(b, &response)

	return &response, err
}
