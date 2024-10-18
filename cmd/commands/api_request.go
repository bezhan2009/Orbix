package commands

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func ApiRequest() {
	reader := bufio.NewReader(os.Stdin)

	// Ввод URL
	fmt.Print("Enter URL: ")
	url, _ := reader.ReadString('\n')
	url = strings.TrimSpace(url)

	// Ввод метода (GET, POST, PUT и т.д.)
	fmt.Print("Enter method (GET, POST, PUT, DELETE): ")
	method, _ := reader.ReadString('\n')
	method = strings.TrimSpace(method)

	// Ввод заголовков (headers)
	headers := make(map[string]string)
	for {
		fmt.Print("Enter the title (in Key:Value format) or press Enter to complete: ")
		headerLine, _ := reader.ReadString('\n')
		headerLine = strings.TrimSpace(headerLine)
		if headerLine == "" {
			break
		}
		headerParts := strings.SplitN(headerLine, ":", 2)
		if len(headerParts) == 2 {
			headers[strings.TrimSpace(headerParts[0])] = strings.TrimSpace(headerParts[1])
		} else {
			fmt.Println("The header format is incorrect. Please use the Key:Value format")
		}
	}

	// Ввод тела запроса (если метод POST или PUT)
	var bodyData []byte
	if method == "POST" || method == "PUT" {
		fmt.Print("Enter the request body (in JSON or text format): ")
		body, _ := reader.ReadString('\n')
		bodyData = []byte(body)
	}

	// Создание клиента и запроса
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(bodyData))
	if err != nil {
		fmt.Println("Error creating the request:", err)
		return
	}

	// Установка заголовков
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Выполнение запроса
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error when executing the request:", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// Чтение и вывод ответа
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading the response:", err)
		return
	}

	// Форматирование JSON-ответа
	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, body, "", "    ")
	if err != nil {
		fmt.Println("Ошибка при форматировании JSON:", err)
		return
	}

	fmt.Println("Response status:", resp.Status)
	fmt.Println("Response body:\n", prettyJSON.String())
}
