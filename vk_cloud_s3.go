package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"log"
	"net/http"
	"strconv"
	// "mime/multipart"
)

const (
	vkCloudHotboxEndpoint = "https://hb.vkcs.cloud"
	defaultRegion         = "ru-msk"
)

func uploadFilesToVkCloudBucket() {
	var fileBytes []byte
	files := fetchUserFiles()
	for _, file := range files {
		fmt.Println("ПУТЬ ФАЙЛА", file.Path)
		fmt.Println("СТАРТ СКАЧИВАНИЯ ОЧЕРЕДНОГО ФАЙЛА")
		fileBytes = downloadFileFromYandexDisk(file.Path)
		uploadToBucket(fileBytes, file.Name)
	}

}

func downloadFileFromYandexDisk(path string) []byte {
	req, _ := http.NewRequest("GET", DOWNLOAD_FILE_URL_TEMPLATE+path, nil)
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

	var dfresp DownloadFileResponse

	json.Unmarshal(respBody, &dfresp)

	fmt.Println("ССЫЛКА ПОЛУЧЕНА", dfresp.HRef)

	//часть 2

	req, _ = http.NewRequest("GET", dfresp.HRef, nil)
	// req.Header.Set("Authorization", "OAuth "+user.Token)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer resp.Body.Close()

	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("ФАЙЛ ПОЛУЧЕН")

	return respBody
}

func uploadToBucketOnlyOneFile(filename string) {
	var fileBytes []byte
	var flag = false
	files := fetchUserFiles()
	for _, file := range files {
		if filename == file.Name {
			flag = true
			fmt.Println("ПУТЬ ФАЙЛА", file.Path)
			fmt.Println("СТАРТ СКАЧИВАНИЯ ФАЙЛА")
			fileBytes = downloadFileFromYandexDisk(file.Path)
			uploadToBucket(fileBytes, file.Name)
		}
	}

	if flag == false {
		sendMessage("Такого файла нет.")
	}
}

func deleteFileFromBucket(filename string) {
	// Создание сессии
	sess, err := session.NewSession()
	if err != nil {
		log.Fatalf("Unable to create session, %v", err)
	}
	// Подключение к сервису Cloud Storage
	svc := s3.New(sess, &aws.Config{Credentials: credentials.NewStaticCredentials("<id>", "<secret>", ""), Region: aws.String(defaultRegion), Endpoint: aws.String(vkCloudHotboxEndpoint)})

	// Удаление объекта из бакета
	bucket := "bot_bucket"
	idStr := strconv.Itoa(user.Id)
	key := idStr + "_" + filename

	if _, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}); err != nil {
		log.Fatalf("Unable to delete object %q from bucket %q, %v\n", key, bucket, err)
		sendMessage("Ошибка удаления файла.")
	} else {
		log.Printf("Object %q deleted from bucket %q\n", key, bucket)
		sendMessage("Файл успешно удалён.")
	}

}

func uploadToBucket(fileBytes []byte, filename string) {

	sess, _ := session.NewSession()

	// myFile := respBody
	file := bytes.NewReader(fileBytes)

	svc := s3.New(sess, &aws.Config{Credentials: credentials.NewStaticCredentials("<id>", "<secret>", ""), Region: aws.String(defaultRegion), Endpoint: aws.String(vkCloudHotboxEndpoint)})

	// Указываем имя бакета и ключ для файла
	bucket := "bot_bucket"
	idStr := strconv.Itoa(user.Id)
	key := idStr + "_" + filename

	fmt.Println("СТАРТ ЗАГРУЗКИ В БАКЕТ")
	// Загружаем файл в бакет
	if _, err := svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   file,
	}); err != nil {
		log.Fatalf("Unable to upload %q to %q, %v\n", key, bucket, err)
		sendMessage("Такого файла нет.")
	} else {
		fmt.Printf("File %q uploaded to bucket %q\n", key, bucket)
	}

}
