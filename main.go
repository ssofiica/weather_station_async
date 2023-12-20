package main

import (
 "github.com/gin-gonic/gin"
 "fmt"
 "math/rand"
 "net/http"
 "time"
 "bytes"
 "encoding/json"
)

func randomStatus() int {
	time.Sleep(10 * time.Second) // Задержка на 10 секунд
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(30) - 10
}

const key string = "a4e0oinhl932as15"

type Result struct {
	Temperature int `json:"temperature"`
	Key string `json:"key"`
}

func performPUTRequest(url string, data Result) (*http.Response, error) {
	// Сериализация структуры в JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// Создание PUT-запроса
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	// Выполнение запроса
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return resp, nil
}

func SendStatus(pk string, url string) {
	fmt.Println(url)
	// Выполнение расчётов с randomStatus
	result := randomStatus()
	data := Result{Temperature: result, Key: key,}
	// Отправка PUT-запроса к основному серверу
	_, err := performPUTRequest(url, data)
	if err != nil {
		fmt.Println("Error sending status:", err)
		return
	}

	fmt.Println("Status sent successfully for pk:", pk)
}

func main() {
	router:= gin.Default()
	// Обработчик POST-запроса для set_temperature
	router.POST("/set_temperature/", func(c *gin.Context) {
		// Получение значения "id" из запроса
		id := c.PostForm("id")
		// Запуск горутины для отправки статуса
		go SendStatus(id, "http://localhost:8000/requests/temperature/"+id+"/", ) // Замените на ваш реальный URL
	
		c.JSON(http.StatusOK, gin.H{"message": "Status update initiated"})
	})
	router.Run(":8080")
}