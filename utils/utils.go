package utils

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	// "strconv"
	"time"
	"bytes"
	"path/filepath"
	"fmt"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/url"
	"net/http"
)

const (
	accessKeyID     = "your-access-key-id"
	accessKeySecret = "your-access-key-secret"
	bucketName      = "your-bucket-name"
	endpoint        = "your-endpoint"
	directory       = "screenshot"
	webhook         = "your webhook" // 填入机器人的webhook
	secret          = "your secret"  




)


func UploadToOSS(data []byte) (string, error) {
	client, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		return "", err
	}

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return "", err
	}

	// 使用当前时间戳生成唯一的文件名
	currentTime := time.Now().Format("20060102_150405")
	fileName := filepath.Join(directory, "screenshot_"+currentTime+".png")

	// fileName := directory + "/screenshot_" + strconv.FormatInt(time.Now().Unix(), 10) + ".png"
	err = bucket.PutObject(fileName, bytes.NewReader(data))
	if err != nil {
		return "", err
	}

	// 获取文件URL
	url, err := bucket.SignURL(fileName, oss.HTTPGet, 600)
	if err != nil {
		return "", err
	}
	return url, nil

}

func generateSignature(secret string, timestamp int64) (string, error) {
	stringToSign := fmt.Sprintf("%d\n%s", timestamp, secret)
	hmac := hmac.New(sha256.New, []byte(secret))
	_, err := hmac.Write([]byte(stringToSign))
	if err != nil {
		return "", err
	}
	signature := base64.StdEncoding.EncodeToString(hmac.Sum(nil))
	return url.QueryEscape(signature), nil
}

func SendImageToDingtalk(imgurl string,webhook string) {
	timestamp := time.Now().UnixNano() / 1e6
	sign, err := generateSignature(secret, timestamp)
	if err != nil {
		fmt.Println("生成签名失败:", err)
		return
	}

	webhookURL := fmt.Sprintf("%s&timestamp=%d&sign=%s", webhook, timestamp, sign)
	fmt.Println("webhookurl",webhookURL)

	data := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"title": "test",
			"text":  fmt.Sprintf("![image](%s)", imgurl),
		},
	}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		fmt.Println("JSON序列化失败:", err)
		return
	}

	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(dataBytes))
	if err != nil {
		fmt.Println("创建请求失败:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("发送请求失败:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		if result["errcode"].(float64) == 0 {
			fmt.Println("图片发送成功")
		} else {
			fmt.Println("图片发送失败:", result)
		}
	} else {
		fmt.Println("图片发送失败, HTTP状态码:", resp.StatusCode)
	}
}