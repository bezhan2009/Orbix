package commands

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	ReadJsonUtils "goCmd/cmd/commands/Read/utils"
	"goCmd/system"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func handleRequestFileExpansion(expansion string) bool {
	return expansion == "json" ||
		expansion == "xml" ||
		expansion == "yaml" ||
		expansion == "yml" ||
		expansion == "toml" ||
		expansion == "ini" ||
		expansion == "jsonl" ||
		expansion == "json5" ||
		expansion == "avro" ||
		expansion == "bin" ||
		expansion == "parquet" ||
		expansion == "msgpack" ||
		expansion == "txt"
}

func ApiRequest() {
	reader := bufio.NewReader(os.Stdin)
	colors := system.GetColorsMap()

	// Ввод URL
	fmt.Print(colors["yellow"]("Enter URL: "))
	url, _ := reader.ReadString('\n')
	url = strings.TrimSpace(url)

	// Ввод метода (GET, POST, PUT и т.д.)
	fmt.Print(colors["yellow"]("Enter method (GET, POST, PUT, DELETE): "))
	method, _ := reader.ReadString('\n')
	method = strings.TrimSpace(strings.ToUpper(method))

	// Ввод заголовков (headers)
	headers := make(map[string]string)
	for {
		fmt.Print(colors["yellow"]("Enter the title (in Key:Value format) or press Enter to complete: "))
		headerLine, _ := reader.ReadString('\n')
		headerLine = strings.TrimSpace(headerLine)
		if headerLine == "" {
			break
		}
		headerParts := strings.SplitN(headerLine, ":", 2)
		if len(headerParts) == 2 {
			headers[strings.TrimSpace(headerParts[0])] = strings.TrimSpace(headerParts[1])
		} else {
			fmt.Println(colors["red"]("The header format is incorrect. Please use the Key:Value format"))
		}
	}

	// Ввод тела запроса (если метод POST или PUT)
	var bodyData []byte
	if method == "POST" || method == "PUT" {
		fmt.Print(colors["yellow"]("Enter the request body (in JSON or text format or file name): "))
		body, _ := reader.ReadString('\n')
		bodySplit := strings.Split(body, ".")
		if handleRequestFileExpansion(strings.TrimSpace(strings.ToLower(bodySplit[1]))) {
			Data, err := ReadJsonUtils.File(strings.TrimSpace(body))
			if err != nil {
				fmt.Println(colors["red"]("Error reading json file:", err))
			}

			body = string(Data)
		}
		bodyData = []byte(body)
	}

	// Создание клиента и запроса
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(bodyData))
	if err != nil {
		fmt.Println(colors["red"]("Error creating the request:", err))
		return
	}

	// Установка заголовков
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Выполнение запроса
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(colors["red"]("Error when executing the request:", err))
		return
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			fmt.Println(colors["red"]("Error when closing resp. body:", err))
		}
	}(resp.Body)

	// Чтение и вывод ответа
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(colors["red"]("Error reading the response:", err))
		return
	}

	// Форматирование JSON-ответа
	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, body, "", "    ")
	if err != nil {
		fmt.Println(colors["red"]("Error when formatting JSON:", err))
		fmt.Println(colors["yellow"](string(body)))
		return
	}

	fmt.Println(colors["magenta"]("Response status:", resp.Status))
	fmt.Println(colors["magenta"]("\nResponse body:\n", prettyJSON.String()))
}
