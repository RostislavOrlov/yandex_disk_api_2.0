package main

import (
	"strconv"
)

func printIntro(update *TelegramMessageResponse) {
	msg := "Привет! Это бот для работы с Яндекс Диском! Для дальнейшей работы с ботом необходимо авторизоваться в Яндекс Диске. Для авторизации введите команду /authorization"
	sendMessage(msg)
}

func calculateAvailableDiskSpace(a, b int) string {
	result := a - b
	resultStr := Convert(result)
	return resultStr
}

func Convert(number int) string {
	numberResult := ""
	if number < 1024 {
		numberResult = strconv.Itoa(number) + " байт"
	} else {
		number /= 1024
		if number < 1024 {
			numberResult = strconv.Itoa(number) + " КБайт"
		} else {
			number /= 1024
			if number < 1024 {
				numberResult = strconv.Itoa(number) + " МБайт"
			} else {
				number /= 1024
				numberResult = strconv.Itoa(number) + " ГБайт"
			}
		}
	}
	return numberResult
}
