package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var FOLDER_ID = os.Getenv("FOLDER_ID")
var BOT_TOKEN = os.Getenv("BOT_TOKEN")
var IAM_TOKEN = os.Getenv("IAM_TOKEN")
var chatId int
var auth_code string
var myMap map[string]string
var user *User

func sendRequest(user *User, method, urlTemplate string, args ...string) ([]byte, int) {
	req, _ := http.NewRequest(method, urlTemplate+args[0]+args[1], nil)
	req.Header.Set("Authorization", "OAuth "+user.Token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf(err.Error())
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	if resp.StatusCode >= 400 {
		sendMessage("Ошибка: введён неверный путь. Попробуйте снова")
		return make([]byte, 0), resp.StatusCode
	}

	return respBody, resp.StatusCode
}

func authUser() {
	msg := "Пожалуйста, авторизуйтесь в Яндекс Диске, перейдя по ссылке ниже, " +
		"и укажите в сообщении авторизационный код, который выведется в браузере после успешной " +
		"авторизации и предоставлении некоторых прав доступа: "
	sendMessage(msg)
	sendMessage(getAuthURL())
}

func downloadFile(user *User, path string) {
	respBody, code := sendRequest(user, "GET", DOWNLOAD_FILE_URL_TEMPLATE, path, "")
	if code >= 400 {
		return
	}
	var downloadFileResponse DownloadFileResponse
	json.Unmarshal(respBody, &downloadFileResponse)

	sendMessage("Для скачивания файла перейдите по ссылке ниже:\n" + downloadFileResponse.HRef)

}

func fetchFoldersAndViews(path string, user *User) (FetchFoldersAndViewsResponse, int) {
	respBody, code := sendRequest(user, "GET", FETCH_FOLDER_CONTENT_URL_TEMPLATE, path, "")
	if code >= 400 {
		return FetchFoldersAndViewsResponse{}, code
	}

	var folderContentResponse FetchFoldersAndViewsResponse
	err := json.Unmarshal(respBody, &folderContentResponse)
	if err != nil {
		fmt.Println(err.Error())
	}

	return folderContentResponse, code
}

func showContentFromPath(path string, folderContent FetchFoldersAndViewsResponse) {
	items := folderContent.Embedded.Items
	stringOfContentNames := "\n"
	var typeOfContent string
	for i := 0; i < len(items); i += 1 {
		typeOfContent = items[i].Type
		switch typeOfContent {
		case "dir":
			typeOfContent = "папка"
		case "file":
			typeOfContent = "файл"
		}
		stringOfContentNames += "\n" + items[i].Name + " (" + typeOfContent + ")"
	}
	msg := fmt.Sprintf("Содержимое папки '%s':%s", path, stringOfContentNames)

	sendMessage(msg)
}

func viewFoldersAndFiles(user *User, path string) {
	path = "disk:" + path
	folderContent, code := fetchFoldersAndViews(path, user)
	if code >= 400 {
		return
	}

	showContentFromPath(path, folderContent)
}

func createFolder(user *User, path string) {
	path = "disk:" + path
	_, code := sendRequest(user, "PUT", CREATE_FOLDER_URL_TEMPLATE, path, "")
	//json.Unmarshal()
	if code >= 400 {
		return
	}
	sendMessage("Папка успешно создана.")
}

func deleteFolder(user *User, path string) {
	path = "disk:" + path
	_, code := sendRequest(user, "DELETE", DELETE_FOLDER_URL_TEMPLATE, path, "")
	if code >= 400 {
		return
	}
	sendMessage("Папка успешно удалена.")
}

func copyContentFromTo(user *User, from, path string) {
	from = "disk:" + from
	path = "disk:" + path
	_, code := sendRequest(user, "POST", COPY_FOLDER_OR_FILE_URL_TEMPLATE+"from="+from+"&path="+path, "")
	if code >= 400 {
		return
	}
	sendMessage("Файл/папка успешно скопирован(-а).")
}

func moveContentFromTo(user *User, from, path string) {
	from = "disk:" + from
	path = "disk:" + path
	_, code := sendRequest(user, "POST", MOVE_FOLDER_OR_FILE_URL_TEMPLATE+"from="+from+"&path="+path, "")
	if code >= 400 {
		return
	}
	sendMessage("Файл/папка успешно перемещен(-а).")
}

func sendMessage(msg string) {
	// var chatIdInt int = int(chatId)
	chatIdStr := strconv.Itoa(chatId)
	// if err != nil {
	// 	fmt.Println("from sendMessage", err.Error())
	// }
	text := url.QueryEscape(msg)
	fmt.Println(TELEGRAM_API_TEMPLATE + "/sendMessage?chat_id=" + chatIdStr + "&text=" + text)
	req, err := http.NewRequest("POST", TELEGRAM_API_TEMPLATE+BOT_TOKEN+"/sendMessage?chat_id="+chatIdStr+"&text="+text, nil)
	if err != nil {
		fmt.Println("from sendMessage", err.Error())
	}
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("from sendMessage2", err.Error())
	}
}

func ConvertDiskInfo(resp ShowInfoResponse) ShowInfo {
	totalSpace := Convert(resp.TotalSpace)
	usedSpace := Convert(resp.UsedSpace)
	trashSize := Convert(resp.TrashSize)

	return ShowInfo{TotalSpace: totalSpace, UsedSpace: usedSpace, TrashSize: trashSize}
}

func showInfo(user *User) {

	respBody, code := sendRequest(user, "GET", SHOW_INFORMATION_DISK_URL_TEMPLATE, "", "")
	if code >= 400 {
		return
	}
	var showInfoResponse ShowInfoResponse
	var showInfo ShowInfo
	json.Unmarshal(respBody, &showInfoResponse)

	showInfo = ConvertDiskInfo(showInfoResponse)
	sendMessage("Информация о Вашем диске:\n" +
		"Общий объем Диска, выделенный Вам: " + showInfo.TotalSpace + "\n" +
		"Общий объем Диска, доступный Вам: " + calculateAvailableDiskSpace(showInfoResponse.TotalSpace, showInfoResponse.UsedSpace) + "\n" +
		"Объем файлов, уже хранящихся на Диске: " + showInfo.UsedSpace + "\n" +
		"Объем файлов, находящихся в корзине: " + showInfo.TrashSize + "\n")
}

func cleanTrash(user *User) {
	_, code := sendRequest(user, "DELETE", CLEAN_TRASH_URL_TEMPLATE, "", "")
	if code >= 400 {
		return
	}
	sendMessage("Корзина полностью очищена.")

}

func restoreContentFromTrash(user *User, path string) {
	_, code := sendRequest(user, "PUT", RESTORE_CONTENT_FROM_TRASH_URL_TEMPLATE, path, "")
	if code >= 400 {
		return
	}
	sendMessage("Файл/папка успешно восстановлен(-а) из корзины.")
}

func isInUserFiles(user *User, fileName string) bool {
	var flag = false
	for _, file := range user.Files {
		if fileName == file.Name {
			flag = true
			break
		}
	}
	return flag
}

func getFileByFileName(user *User, name string) File {
	var file File
	for _, fileTmp := range user.Files {
		if name == fileTmp.Name {
			file = fileTmp
		}
	}
	return file
}

func processFile(user *User, fileName string) (string, bool) {
	var file File
	if !isInUserFiles(user, fileName) {
		return "Такого файла нет.", false
	} else {
		file = getFileByFileName(user, fileName)
		return file.Path, true
	}
}

func sendMessageHelp() {
	sendMessage(HELP)
}

func sendMessageInvalidPath() {
	sendMessage("Введен неверный путь")
}

func updateProcessing(update *TelegramMessageResponse) {

	listOfSecrets, _ := listLockboxSecrets()

	fmt.Println("СЕКРЕТЫ", listOfSecrets)

	idStr := strconv.Itoa(user.Id)
	for _, secret := range listOfSecrets.Secrets {
		if idStr == secret.Name {
			secretPayload, _ := getLockboxPayloadAsToken(user, secret.CurrentVersion.SecretId)
			user.Token = secretPayload.Entries[0].TextValue
			fmt.Println("SECRET PAYLOAD", secretPayload)
			break
		}
	}

	// if !flag {
	// 	sendMessage("Вы не авторизованы.\n")
	// 	authUser()
	// }

	if len(update.Message.Text) == 7 && user.Token == "" {
		preAuth(update)
	} else if update.Message.Text == "/start" {
		printIntro(update)
	} else if update.Message.Text == "/authorization" {
		authUser()
	} else if strings.Contains(update.Message.Text, "/команда") {
		if user.Token != "" {
			splitMsg := strings.Split(update.Message.Text, " ")
			switch splitMsg[0] {

			case "/команда_выгрузить_файлы_диск_бакет":
				uploadFilesToVkCloudBucket()
				sendMessage("Все файлы успешно выгружены из Диска в бакет.")

			case "/команда_выгрузить_файл_диск_бакет":
				uploadToBucketOnlyOneFile(splitMsg[1])
				sendMessage("Файл успешно загружен в вк бакет.")

			case "/команда_удалить_файл_бакет":
				deleteFileFromBucket(splitMsg[1])

			case "/команда_зайти_в_папку":
				if splitMsg[1] == "" {
					sendMessageInvalidPath()
					return
				}
				viewFoldersAndFiles(user, splitMsg[1])

			case "/команда_создать_папку":
				if splitMsg[1] == "" {
					sendMessageInvalidPath()
					return
				}
				createFolder(user, splitMsg[1])

			case "/команда_удалить_папку":
				if splitMsg[1] == "" {
					sendMessageInvalidPath()
					return
				}
				deleteFolder(user, splitMsg[1])

			case "/команда_копирование_файла_или_папки":
				if splitMsg[1] == "" || splitMsg[2] == "" {
					sendMessageInvalidPath()
					return
				}
				copyContentFromTo(user, splitMsg[1], splitMsg[2])

			case "/команда_перемещение_файла_или_папки":
				if splitMsg[1] == "" || splitMsg[2] == "" {
					sendMessageInvalidPath()
					return
				}
				moveContentFromTo(user, splitMsg[1], splitMsg[2])

			case "/команда_скачать_файл":
				if splitMsg[1] == "" {
					sendMessageInvalidPath()
					return
				}
				downloadFile(user, splitMsg[1])

			case "/команда_посмотреть_информацию":
				showInfo(user)

			case "/команда_очистить_корзину":
				cleanTrash(user)

			case "/команда_восстановить_из_корзины":
				if splitMsg[1] == "" {
					sendMessageInvalidPath()
					return
				}
				restoreContentFromTrash(user, splitMsg[1])

			case "/команда_файл":
				if splitMsg[1] == "" {
					return
				}
				filePath, exists := processFile(user, splitMsg[1])
				if exists {
					sendMessage("Путь в Яндекс Диске для вашего файла: " + filePath)
				} else {
					sendMessage(filePath)
				}
			default:
				sendMessage("Введена неизвестная команда.")
				sendMessageHelp()
			}
		} else {
			sendMessage("Вы не авторизованы.\n")
			authUser()
		}
	} else {
		sendMessage("Введена неизвестная команда.")
		sendMessageHelp()
	}
}

type Response struct {
	StatusCode int         `json:"statusCode"`
	Body       interface{} `json:"body"`
}

func Main(ctx context.Context, req []byte) (*Response, error) {

	runLogExport()

	var resp TelegramMessageResponse

	err := json.Unmarshal(req, &myMap)
	if err != nil {
		fmt.Println("ОШИБКА", err.Error())
	}

	fmt.Println("МАПА", myMap)
	fmt.Println("ТЕЛО МАПЫ", myMap["body"])

	err = json.Unmarshal([]byte(myMap["body"]), &resp)
	if err != nil {
		fmt.Println("Ошибка при распарсивании строки:", err)
	}

	fmt.Println()
	fmt.Println("СТРУКТУРА", resp)

	chatId = resp.Message.Chat.Id

	user = &User{Id: resp.Message.From.Id, Name: resp.Message.From.Username, Token: "", ChatId: chatId}

	if user.Id == 1334199506 {
		updateProcessing(&resp)
	}

	return &Response{
		StatusCode: 200,
		Body:       "Успешно",
	}, nil

}

func Main1(ctx context.Context, req []byte) (*Response, error) {

	var resp TelegramMessageResponse

	err := json.Unmarshal(req, &myMap)
	if err != nil {
		fmt.Println("ОШИБКА", err.Error())
	}

	fmt.Println("МАПА", myMap)
	fmt.Println("ТЕЛО МАПЫ", myMap["body"])

	err = json.Unmarshal([]byte(myMap["body"]), &resp)
	if err != nil {
		fmt.Println("Ошибка при распарсивании строки:", err)
	}

	fmt.Println()
	fmt.Println("СТРУКТУРА", resp)

	chatId = resp.Message.Chat.Id

	user = &User{Id: resp.Message.From.Id, Name: resp.Message.From.Username, Token: "", ChatId: chatId}

	// updateProcessing(&resp)
	if user.Id == 1334199506 && resp.Message.Text == "/начали" {
		listOfSecrets, _ := listLockboxSecrets()

		fmt.Println("СЕКРЕТЫ", listOfSecrets)

		idStr := strconv.Itoa(user.Id)
		for _, secret := range listOfSecrets.Secrets {
			if idStr == secret.Name {
				secretPayload, _ := getLockboxPayloadAsToken(user, secret.CurrentVersion.SecretId)
				user.Token = secretPayload.Entries[0].TextValue
				fmt.Println("SECRET PAYLOAD", secretPayload)
				break
			}
		}

		// files := fetchUserFiles()
		// fmt.Println("ЫЫЫЫ", files)
		// fileBytes := downloadFileFromYandexDisk("disk:/Тестовая папка/Хлебные крошки.mp4")
		// uploadToBucket(fileBytes, "id_test")
		// uploadFilesToVkCloudBucket()
		// fileBytes, filename := downloadFileFromBucket("Зима.jpg")
		// sendFile(fileBytes, filename)
		// uploadToBucketOnlyOneFile("Зима.jpg")
		deleteFileFromBucket("Зима.jpg")
	} else {
		sendMessage("В данный момент бот не работает!")
	}

	return &Response{
		StatusCode: 200,
		Body:       "Успешно",
	}, nil
}
