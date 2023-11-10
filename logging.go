package main

import (
	"bytes"
	"fmt"
	"net/http"
)

var (
	groupId = "e237mbgin6tbfvj6edbq"
	sinkId  = "e236b97jbqeuf514maq9"
)

func runLogExport() {
	var jsonData = []byte(fmt.Sprintf(`{
        "groupId": "%s",
        "sinkId": "%s",
        "params": {},
        "resultFilename": "log_file",
        "since": "2023-10-27T23:30:00Z",
        "until": "2023-11-5T00:00:00Z"
        }`, groupId, sinkId))
	req, err := http.NewRequest("POST", RUN_LOG_EXPORT_TEMPLATE, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("ОШИБКА (из runLogExport, во время создания запроса):", err.Error())
	}

	req.Header.Set("Authorization", "Bearer "+IAM_TOKEN)

	req.Header.Set("Content-Type", "application/json")

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("ОШИБКА (из runLogExport, после выполнения запроса):", err.Error())
	}

}
