package main

import (
	"encoding/json"
	// tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var (
	clientId     = os.Getenv("CLIENT_ID")
	clientSecret = os.Getenv("CLIENT_SECRET")
	redirectURI  = "https://oauth.yandex.ru/verification_code"
)

// Функция для получения ссылки для авторизации
func getAuthURL() string {
	authURL := "https://oauth.yandex.ru/authorize?"
	params := url.Values{}
	params.Set("response_type", "code")
	params.Set("client_id", clientId)
	params.Set("redirect_uri", redirectURI)
	return authURL + params.Encode()
}

func getTokenURL() string {
	return "https://oauth.yandex.ru/token"
}

func requestForUserToken() {
	body := strings.NewReader("grant_type=authorization_code&code=" + auth_code + "&client_id=" + clientId + "&client_secret=" + clientSecret)
	req, _ := http.NewRequest("POST", getTokenURL(), body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", string(body.Len()))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		sendMessage("Вы ввели неверный авторизационный код.")
		authUser()
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	panic(err)
	// }

	println("ТОКЕН: " + string(respBody))
	var userToken UserToken
	err = json.Unmarshal(respBody, &userToken)
	if err != nil {
		return
	}

	user.Token = userToken.AccessToken
	saveTokenInKMS_Storage(user)
	fmt.Println("ТОКЕН ПОЛЬЗОВАТЕЛЯ", user.Token)
	// sendMessage("Ваш токен доступа:\n\n" + user.Token + "\n\nОбратитесь к тех. поддержке телеграмм бота и соообщите им этот токен для дальнейшей работы.")
}

func saveTokenInKMS_Storage(user *User) {
	kmsKeyId := createKMSkey(user)
	secretId := createSecret(kmsKeyId, user)
	idStr := strconv.Itoa(user.Id)
	addVersionLockbox(secretId, idStr, user.Token)
}

func preAuth(update *TelegramMessageResponse) {
	_, err := strconv.Atoi(update.Message.Text)
	if len(update.Message.Text) == 7 && err == nil {
		auth_code = update.Message.Text
		fmt.Println("ДО ПОЛУЧЕНИЯ ТОКЕНА")
		requestForUserToken()
		fmt.Println("ТОКЕН ПОЛЬЗОВАТЕЛЯ (ИЗ preAuth)", user.Token)
		postAuth(update)
		// user.Files = fetchUserFiles()
		// fmt.Println(user.Files)
	}

}

func postAuth(update *TelegramMessageResponse) {
	sendMessage("Авторизация прошла успешно. Добро пожаловать!")
	sendMessage(HELP)
}

func fetchUserFiles() []File {
	req, _ := http.NewRequest("GET", FETCH_USER_FILES_URL_TEMPLATE, nil)
	fmt.Println("ТОКЕН ПОЛЬЗОВАТЕЛЯ В ФЕТЧЕ", user.Token)
	req.Header.Set("Authorization", "OAuth "+user.Token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf(err.Error())
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var filesResponse FilesResponse
	files := make([]File, 0, 1024)
	json.Unmarshal(respBody, &filesResponse)

	var f interface{}
	json.Unmarshal(respBody, &f)
	fmt.Println("ПУСТОЙ ИНТЕРФЕЙС", f)

	filesTmp := filesResponse.Items
	for i := 0; i < len(filesTmp); i++ {
		files = append(files, File{Name: filesTmp[i].Name, Path: filesTmp[i].Path})
	}

	fmt.Println("ФАЙЛЫ ПОЛЬЗОВАТЕЛЯ:", files)

	return files
}
