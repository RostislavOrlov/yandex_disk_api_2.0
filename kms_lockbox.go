package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

var TOKEN = os.Getenv("TOKEN")

func sendRequestKMS(user *User, method, urlTemplate string, body []byte) ([]byte, error) {
	req, err := http.NewRequest(method, urlTemplate, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("ОШИБКА (из sendRequestKMS, во время создания запроса):", err.Error())
	}

	req.Header.Set("Authorization", "Bearer "+IAM_TOKEN)
	if method == "POST" {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("ОШИБКА (из sendRequestKMS, после выполнения запроса):", err.Error())
	}

	if resp.StatusCode >= 400 {
		getIAMtoken()
		respBody, err := sendRequestKMS(user, method, urlTemplate, body)
		return respBody, err
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)

	// if method == "POST" {

	//     fmt.Println("ПУТЬ ЗАПРОСАААОАОАОАОАОА", )
	//     var r interface{}
	//     err := json.Unmarshal(respBody, &r)

	//     if err != nil {
	//         fmt.Println("ОШИБКА (из sendRequestKMS, после преобразования тела запроса):", err.Error())
	//     }

	//     fmt.Println("ТЕЛО ЗАПРОСА НА СОЗДАНИЕ KMS КЛЮЧА", r)
	//     fmt.Println(resp.StatusCode)
	//     fmt.Println(resp.Status)
	//     fmt.Println("ВЕСЬ ЗАПРОС", req)
	// }

	return respBody, err
}

func createKMSkey(user *User) string {
	idStr := strconv.Itoa(user.Id)
	var jsonData = []byte(fmt.Sprintf(`{
        "folderId": "%s",
        "name": "%s",
        "defaultAlgorithm": "AES_256"
        }`, FOLDER_ID, idStr))
	respBody, err := sendRequestKMS(user, "POST", KMS_CREATE_KEY_TEMPLATE, jsonData)

	var resp CreateKMSkeyResponse
	err = json.Unmarshal(respBody, &resp)
	if err != nil {
		fmt.Println("ОШИБКА (из createKMSkey, после преобразования тела запроса):", err.Error())
	}

	fmt.Println(resp.Response.PrimaryVersion.KeyId)
	return resp.Response.PrimaryVersion.KeyId
}

func listKMSkeys() (ListKMSkeysResponse, error) {
	req, err := http.NewRequest("GET", KMS_LIST_KEYS_TEMPLATE+"?folderId="+FOLDER_ID, nil)
	if err != nil {
		fmt.Println("ОШИБКА (из listKMSkeys, во время создания запроса):", err.Error())
	}
	req.Header.Set("Authorization", "Bearer "+IAM_TOKEN)

	resp, err := http.DefaultClient.Do(req)
	if resp.StatusCode >= 400 {
		getIAMtoken()
		listOfKeys, err := listKMSkeys()
		return listOfKeys, err
	}
	if err != nil {
		fmt.Println("ОШИБКА (из listKMSkeys, после выполнения запроса):", err.Error())
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ОШИБКА (из listKMSkeys, после чтения тела запроса):", err.Error())
	}

	var listOfKeys ListKMSkeysResponse

	err = json.Unmarshal(respBody, &listOfKeys)
	if err != nil {
		fmt.Println("ОШИБКА (из listKMSkeys, во время преобразования тела запроса):", err.Error())
	}

	return listOfKeys, err
}

func createSecret(kmsKeyId string, user *User) string {
	idStr := strconv.Itoa(user.Id)
	var jsonData = []byte(fmt.Sprintf(`{
        "folderId": "%s",
        "name": "%s",
        "kmsKeyId": "%s",
        }`, FOLDER_ID, idStr, kmsKeyId))
	respBody, _ := sendRequestLockbox("POST", LOCKBOX_CREATE_SECRET_TEMPLATE, jsonData)
	var resp CreateLockboxSecretResponse
	json.Unmarshal(respBody, &resp)

	return resp.Metadata.SecretId
}

func addVersionLockbox(secretId, key, token string) {
	var jsonData = []byte(fmt.Sprintf(`{
        "payloadEntries": [
                {
                    "key": "%s",
                    "textValue": "%s",
                }
            ]
        }`, key, token))
	sendRequestLockbox("POST", LOCKBOX_ADD_VERSION_TEMPLATE+"/"+secretId+":addVersion", jsonData)
}

func sendRequestLockbox(method, urlTemplate string, body []byte) ([]byte, error) {
	req, err := http.NewRequest(method, urlTemplate, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("ОШИБКА (из sendRequestLockbox, во время создания запроса):", err.Error())
	}

	req.Header.Set("Authorization", "Bearer "+IAM_TOKEN)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("ОШИБКА (из sendRequestLockbox, после выполнения запроса):", err.Error())
	}
	if resp.StatusCode >= 400 {
		getIAMtoken()
		respBody, err := sendRequestLockbox(method, urlTemplate, body)
		return respBody, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)

	return respBody, nil
}

func listLockboxSecrets() (ListLockboxSecretsResponse, error) {
	var jsonData = []byte(fmt.Sprintf(`{
        "folderId": "%s"
        }`, FOLDER_ID))
	req, err := http.NewRequest("GET", LOCKBOX_LIST_SECRETS_TEMPLATE, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("ОШИБКА (из listLockboxSecrets, во время создания запроса):", err.Error())
	}
	req.Header.Set("Authorization", "Bearer "+IAM_TOKEN)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("ОШИБКА (из listLockboxSecrets, после выполнения запроса):", err.Error())
	}

	if resp.StatusCode >= 400 {
		getIAMtoken()
		listOfSecrets, err := listLockboxSecrets()
		return listOfSecrets, err
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ОШИБКА (из listLockboxSecrets, после чтения тела запроса):", err.Error())
	}

	var listOfSecrets ListLockboxSecretsResponse

	err = json.Unmarshal(respBody, &listOfSecrets)
	if err != nil {
		fmt.Println("ОШИБКА (из listLockboxSecrets, во время преобразования тела запроса):", err.Error())
	}

	return listOfSecrets, err
}

func getLockboxPayloadAsToken(user *User, secretId string) (GetLockboxPayloadResponse, error) {

	req, err := http.NewRequest("GET", LOCKBOX_GET_PAYLOAD_TEMPLATE+"/"+secretId+"/payload", nil)
	if err != nil {
		fmt.Println("ОШИБКА (из getLockboxPayloadAsToken, во время создания запроса):", err.Error())
	}
	req.Header.Set("Authorization", "Bearer "+IAM_TOKEN)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("ОШИБКА (из getLockboxPayloadAsToken, после выполнения запроса):", err.Error())
	}

	if resp.StatusCode >= 400 {
		getIAMtoken()
		secretPayload, err := getLockboxPayloadAsToken(user, secretId)
		return secretPayload, err
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ОШИБКА (из getLockboxPayloadAsToken, после чтения тела запроса):", err.Error())
	}

	var secretPayload GetLockboxPayloadResponse

	err = json.Unmarshal(respBody, &secretPayload)
	if err != nil {
		fmt.Println("ОШИБКА (из getLockboxPayloadAsToken, во время преобразования тела запроса):", err.Error())
	}

	return secretPayload, err
}

func getIAMtoken() {
	var jsonData = []byte(fmt.Sprintf(`{
        "yandexPassportOauthToken": "%s"
        }`, TOKEN))
	req, err := http.NewRequest("POST", GET_IAM_TOKEN_TEMPLATE, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("ОШИБКА (из getIAMtoken, во время создания запроса):", err.Error())
	}
	req.Header.Set("Authorization", "OAuth "+TOKEN)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("ОШИБКА (из getIAMtoken, после выполнения запроса):", err.Error())
	}

	fmt.Println("ВЕСЬ ЗАПРОС", req)

	defer resp.Body.Close()

	fmt.Println("КОД ОТВЕТА IAM", resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ОШИБКА (из getIAMtoken, после чтения тела запроса):", err.Error())
	}

	var respToken GetIAMtokenResponse

	err = json.Unmarshal(respBody, &respToken)
	if err != nil {
		fmt.Println("ОШИБКА (из getIAMtoken, во время преобразования тела запроса):", err.Error())
	}
	fmt.Println("IAM-TOKEN", respToken.IamToken)
	IAM_TOKEN = respToken.IamToken
}
