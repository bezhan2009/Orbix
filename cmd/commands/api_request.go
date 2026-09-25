package commands

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	ReadJsonUtils "goCmd/cmd/commands/Read/utils"
	"goCmd/system"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ============================================================
// TYPES
// ============================================================

type apiRequestResult struct {
	Index      int
	Status     string
	StatusCode int
	Body       []byte
	Duration   time.Duration
	Err        error
}

// ============================================================
// FILE EXTENSIONS
// ============================================================

func handleRequestFileExpansion(expansion string) bool {
	switch expansion {
	case "json",
		"xml",
		"yaml",
		"yml",
		"toml",
		"ini",
		"jsonl",
		"json5",
		"avro",
		"bin",
		"parquet",
		"msgpack",
		"txt":
		return true
	default:
		return false
	}
}

// ============================================================
// SIMPLE FLAG
//
// Example:
//
// --loading-only
// --no-expanded
// ============================================================

func hasFlag(args []string, flag string) bool {
	expected := "--" + flag

	for _, arg := range args {
		if arg == expected {
			return true
		}
	}

	return false
}

// ============================================================
// INTEGER FLAG
//
// Supports:
//
// --times=20
//
// AND
//
// --times 20
// ============================================================

func getIntFlag(args []string, flagName string) (int, bool) {
	flag := "--" + flagName
	prefix := flag + "="

	for i, arg := range args {

		// -----------------------------------------
		// Example:
		// --times=20
		// -----------------------------------------

		if strings.HasPrefix(arg, prefix) {
			value := strings.TrimPrefix(arg, prefix)

			number, err := strconv.Atoi(value)

			if err != nil || number <= 0 {
				return 0, false
			}

			return number, true
		}

		// -----------------------------------------
		// Example:
		// --times 20
		// -----------------------------------------

		if arg == flag && i+1 < len(args) {
			number, err := strconv.Atoi(args[i+1])

			if err != nil || number <= 0 {
				return 0, false
			}

			return number, true
		}
	}

	return 0, false
}

// ============================================================
// EXECUTE ONE REQUEST
// ============================================================

func executeAPIRequest(
	client *http.Client,
	index int,
	method string,
	url string,
	headers map[string]string,
	bodyData []byte,
) apiRequestResult {

	start := time.Now()

	// IMPORTANT:
	//
	// Every goroutine gets its own http.Request and body reader.
	// Do NOT reuse the same request between goroutines.

	req, err := http.NewRequest(
		method,
		url,
		bytes.NewReader(bodyData),
	)

	if err != nil {
		return apiRequestResult{
			Index:    index,
			Duration: time.Since(start),
			Err: fmt.Errorf(
				"error creating request: %w",
				err,
			),
		}
	}

	// -----------------------------------------
	// Headers
	// -----------------------------------------

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// -----------------------------------------
	// Execute
	// -----------------------------------------

	resp, err := client.Do(req)

	if err != nil {
		return apiRequestResult{
			Index:    index,
			Duration: time.Since(start),
			Err: fmt.Errorf(
				"error executing request: %w",
				err,
			),
		}
	}

	defer resp.Body.Close()

	// -----------------------------------------
	// Read response
	// -----------------------------------------

	responseBody, err := io.ReadAll(resp.Body)

	if err != nil {
		return apiRequestResult{
			Index:      index,
			Status:     resp.Status,
			StatusCode: resp.StatusCode,
			Duration:   time.Since(start),
			Err: fmt.Errorf(
				"error reading response: %w",
				err,
			),
		}
	}

	return apiRequestResult{
		Index:      index,
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
		Body:       responseBody,
		Duration:   time.Since(start),
	}
}

// ============================================================
// PROGRESS BAR
// ============================================================

func printProgress(
	completed int,
	total int,
	successful int,
	failed int,
	start time.Time,
) {

	if total <= 0 {
		return
	}

	// -----------------------------------------
	// Percentage
	// -----------------------------------------

	percent := float64(completed) / float64(total) * 100

	// -----------------------------------------
	// Progress bar
	// -----------------------------------------

	barWidth := 30

	filled := int(
		float64(barWidth) *
			float64(completed) /
			float64(total),
	)

	if filled > barWidth {
		filled = barWidth
	}

	bar := strings.Repeat("█", filled) +
		strings.Repeat("░", barWidth-filled)

	// -----------------------------------------
	// Requests per second
	// -----------------------------------------

	elapsed := time.Since(start).Seconds()

	var requestsPerSecond float64

	if elapsed > 0 {
		requestsPerSecond =
			float64(completed) / elapsed
	}

	// -----------------------------------------
	// Print on same line
	// -----------------------------------------

	fmt.Printf(
		"\rProgress: [%s] %d/%d %.0f%% | %.2f req/s | success: %d | failed: %d",
		bar,
		completed,
		total,
		percent,
		requestsPerSecond,
		successful,
		failed,
	)
}

// ============================================================
// API REQUEST COMMAND
// ============================================================

func ApiRequest(commandArgs []string) {
	reader := bufio.NewReader(os.Stdin)
	colors := system.GetColorsMap()

	// ========================================================
	// DEBUG
	//
	// Uncomment if you need to inspect arguments:
	//
	// fmt.Printf("DEBUG commandArgs: %#v\n", commandArgs)
	// ========================================================

	// ========================================================
	// URL
	// ========================================================

	var url string

	if len(commandArgs) > 0 &&
		!strings.HasPrefix(commandArgs[0], "--") {

		url = commandArgs[0]

		fmt.Println(
			colors["yellow"](
				"URL:",
				url,
			),
		)

	} else {

		fmt.Print(
			colors["yellow"](
				"Enter URL: ",
			),
		)

		url, _ = reader.ReadString('\n')
		url = strings.TrimSpace(url)
	}

	if url == "" {
		fmt.Println(
			colors["red"](
				"URL cannot be empty.",
			),
		)

		return
	}

	// ========================================================
	// HTTP METHOD
	// ========================================================

	var method string

	if len(commandArgs) > 1 &&
		!strings.HasPrefix(commandArgs[1], "--") {

		method = strings.ToUpper(
			strings.TrimSpace(commandArgs[1]),
		)

		fmt.Println(
			colors["yellow"](
				"Method:",
				method,
			),
		)

	} else {

		fmt.Print(
			colors["yellow"](
				"Enter method (GET, POST, PUT, PATCH, DELETE): ",
			),
		)

		method, _ = reader.ReadString('\n')

		method = strings.ToUpper(
			strings.TrimSpace(method),
		)
	}

	if method == "" {
		method = http.MethodGet
	}

	// ========================================================
	// FLAGS
	// ========================================================

	isExpanded := true

	// Support multiple aliases because your earlier command used:
	//
	// --no-extanded
	//
	// Recommended:
	//
	// --no-expanded

	for _, arg := range commandArgs {
		switch arg {

		case "--no-expanded",
			"--not-expanded",
			"--no-extanded":

			isExpanded = false
		}
	}

	loadingOnly := hasFlag(
		commandArgs,
		"loading-only",
	)

	// ========================================================
	// HEADERS
	// ========================================================

	headers := make(map[string]string)

	if isExpanded {

		for {

			fmt.Print(
				colors["yellow"](
					"Enter the header (Key:Value) or press Enter to complete: ",
				),
			)

			headerLine, _ :=
				reader.ReadString('\n')

			headerLine =
				strings.TrimSpace(headerLine)

			if headerLine == "" {
				break
			}

			headerParts :=
				strings.SplitN(
					headerLine,
					":",
					2,
				)

			if len(headerParts) != 2 {

				fmt.Println(
					colors["red"](
						"The header format is incorrect. Use Key:Value",
					),
				)

				continue
			}

			key :=
				strings.TrimSpace(
					headerParts[0],
				)

			value :=
				strings.TrimSpace(
					headerParts[1],
				)

			headers[key] = value
		}
	}

	// ========================================================
	// BODY
	// ========================================================

	var bodyData []byte

	if method == http.MethodPost ||
		method == http.MethodPut ||
		method == http.MethodPatch {

		fmt.Print(
			colors["yellow"](
				"Enter request body (JSON/text/file name): ",
			),
		)

		body, _ :=
			reader.ReadString('\n')

		body =
			strings.TrimSpace(body)

		// -----------------------------------------
		// Is it a file?
		// -----------------------------------------

		if body != "" {

			fileInfo, fileErr :=
				os.Stat(body)

			if fileErr == nil &&
				!fileInfo.IsDir() {

				extension :=
					strings.TrimPrefix(
						strings.ToLower(
							filepath.Ext(body),
						),
						".",
					)

				if handleRequestFileExpansion(
					extension,
				) {

					data, err :=
						ReadJsonUtils.File(body)

					if err != nil {

						fmt.Println(
							colors["red"](
								"Error reading request body file:",
								err,
							),
						)

						return
					}

					bodyData = data

				} else {

					// Existing file but unsupported
					// extension: use entered text.
					bodyData =
						[]byte(body)
				}

			} else {

				// Raw JSON / XML / text
				bodyData =
					[]byte(body)
			}
		}
	}

	// ========================================================
	// REQUEST COUNT
	// ========================================================

	requestCount, hasTimes :=
		getIntFlag(
			commandArgs,
			"times",
		)

	if !hasTimes {

		fmt.Print(
			colors["yellow"](
				"How many times should the request be sent? [1]: ",
			),
		)

		value, _ :=
			reader.ReadString('\n')

		value =
			strings.TrimSpace(value)

		if value == "" {

			requestCount = 1

		} else {

			parsed, err :=
				strconv.Atoi(value)

			if err != nil ||
				parsed <= 0 {

				fmt.Println(
					colors["red"](
						"Invalid request count. Using 1.",
					),
				)

				requestCount = 1

			} else {

				requestCount = parsed
			}
		}
	}

	// ========================================================
	// PARALLEL REQUEST COUNT
	// ========================================================

	parallelCount, hasParallel :=
		getIntFlag(
			commandArgs,
			"parallel",
		)

	if !hasParallel {

		fmt.Print(
			colors["yellow"](
				fmt.Sprintf(
					"Maximum parallel requests? [%d]: ",
					requestCount,
				),
			),
		)

		value, _ :=
			reader.ReadString('\n')

		value =
			strings.TrimSpace(value)

		if value == "" {

			parallelCount =
				requestCount

		} else {

			parsed, err :=
				strconv.Atoi(value)

			if err != nil ||
				parsed <= 0 {

				parallelCount =
					requestCount

			} else {

				parallelCount =
					parsed
			}
		}
	}

	// -----------------------------------------
	// Cannot have more workers than requests
	// -----------------------------------------

	if parallelCount > requestCount {
		parallelCount =
			requestCount
	}

	if parallelCount <= 0 {
		parallelCount = 1
	}

	// ========================================================
	// HTTP TRANSPORT
	// ========================================================

	transport := &http.Transport{

		MaxIdleConns: parallelCount * 2,

		MaxIdleConnsPerHost: parallelCount,

		MaxConnsPerHost: parallelCount,

		IdleConnTimeout: 30 * time.Second,
	}

	// ========================================================
	// HTTP CLIENT
	// ========================================================

	client := &http.Client{
		Transport: transport,

		// Change if necessary.
		Timeout: 60 * time.Second,
	}

	// ========================================================
	// START INFORMATION
	// ========================================================

	fmt.Println()

	fmt.Println(
		colors["yellow"](
			fmt.Sprintf(
				"Sending %d request(s) with max %d parallel request(s)...",
				requestCount,
				parallelCount,
			),
		),
	)

	if loadingOnly {

		fmt.Println(
			colors["yellow"](
				"Mode: loading only",
			),
		)

	}

	fmt.Println()

	// ========================================================
	// CHANNELS
	// ========================================================

	jobs :=
		make(chan int)

	results :=
		make(
			chan apiRequestResult,
			requestCount,
		)

	// ========================================================
	// WORKER GROUP
	// ========================================================

	var wg sync.WaitGroup

	// ========================================================
	// CREATE WORKERS
	// ========================================================

	for workerID := 1; workerID <= parallelCount; workerID++ {

		wg.Add(1)

		go func() {

			defer wg.Done()

			for requestIndex := range jobs {

				result :=
					executeAPIRequest(
						client,
						requestIndex,
						method,
						url,
						headers,
						bodyData,
					)

				results <- result
			}

		}()
	}

	// ========================================================
	// SEND JOBS
	// ========================================================

	go func() {

		for i := 1; i <= requestCount; i++ {

			jobs <- i
		}

		close(jobs)
	}()

	// ========================================================
	// CLOSE RESULTS WHEN ALL WORKERS FINISH
	// ========================================================

	go func() {

		wg.Wait()

		close(results)

	}()

	// ========================================================
	// PROCESS RESULTS
	// ========================================================

	successCount := 0
	errorCount := 0

	completedCount := 0

	totalStart :=
		time.Now()

	var totalRequestDuration time.Duration

	var fastestRequest time.Duration
	var slowestRequest time.Duration

	// ========================================================
	// RECEIVE RESULTS
	// ========================================================

	for result := range results {

		completedCount++

		totalRequestDuration +=
			result.Duration

		// -----------------------------------------
		// Fastest / slowest
		// -----------------------------------------

		if fastestRequest == 0 ||
			result.Duration < fastestRequest {

			fastestRequest =
				result.Duration
		}

		if result.Duration >
			slowestRequest {

			slowestRequest =
				result.Duration
		}

		// -----------------------------------------
		// Count success / error
		// -----------------------------------------

		if result.Err != nil {

			errorCount++

		} else if result.StatusCode >= 200 &&
			result.StatusCode < 400 {

			successCount++

		} else {

			// HTTP request itself worked,
			// but API returned 4xx / 5xx.
			errorCount++
		}

		// ====================================================
		// LOADING ONLY
		// ====================================================

		if loadingOnly {

			printProgress(
				completedCount,
				requestCount,
				successCount,
				errorCount,
				totalStart,
			)

			continue
		}

		// ====================================================
		// NORMAL OUTPUT
		// ====================================================

		fmt.Println(
			colors["magenta"](
				fmt.Sprintf(
					"================ REQUEST #%d ================",
					result.Index,
				),
			),
		)

		// -----------------------------------------
		// Request failed
		// -----------------------------------------

		if result.Err != nil {

			fmt.Println(
				colors["red"](
					"Error:",
					result.Err,
				),
			)

			fmt.Println(
				colors["yellow"](
					"Duration:",
					result.Duration,
				),
			)

			fmt.Println()

			continue
		}

		// -----------------------------------------
		// Response status
		// -----------------------------------------

		fmt.Println(
			colors["magenta"](
				"Response status:",
				result.Status,
			),
		)

		fmt.Println(
			colors["magenta"](
				"Duration:",
				result.Duration,
			),
		)

		// -----------------------------------------
		// Empty body
		// -----------------------------------------

		if len(result.Body) == 0 {

			fmt.Println(
				colors["yellow"](
					"\nResponse body: empty",
				),
			)

			fmt.Println()

			continue
		}

		// -----------------------------------------
		// Try pretty JSON
		// -----------------------------------------

		var prettyJSON bytes.Buffer

		err :=
			json.Indent(
				&prettyJSON,
				result.Body,
				"",
				"    ",
			)

		if err == nil {

			fmt.Println(
				colors["magenta"](
					"\nResponse body:\n",
					prettyJSON.String(),
				),
			)

		} else {

			// HTML/XML/text/etc.
			fmt.Println(
				colors["yellow"](
					"\nResponse body:\n",
					string(result.Body),
				),
			)
		}

		fmt.Println()
	}

	// ========================================================
	// NEW LINE AFTER PROGRESS BAR
	// ========================================================

	if loadingOnly {
		fmt.Println()
		fmt.Println()
	}

	// ========================================================
	// TOTAL TIME
	// ========================================================

	totalExecutionTime :=
		time.Since(totalStart)

	// ========================================================
	// AVERAGE REQUEST DURATION
	// ========================================================

	var averageDuration time.Duration

	if completedCount > 0 {

		averageDuration =
			totalRequestDuration /
				time.Duration(completedCount)
	}

	// ========================================================
	// REQUESTS PER SECOND
	// ========================================================

	var requestsPerSecond float64

	if totalExecutionTime.Seconds() > 0 {

		requestsPerSecond =
			float64(completedCount) /
				totalExecutionTime.Seconds()
	}

	// ========================================================
	// SUMMARY
	// ========================================================

	fmt.Println(
		colors["yellow"](
			"================ SUMMARY ================",
		),
	)

	fmt.Println(
		colors["yellow"](
			fmt.Sprintf(
				"Total requests: %d",
				requestCount,
			),
		),
	)

	fmt.Println(
		colors["yellow"](
			fmt.Sprintf(
				"Successful: %d",
				successCount,
			),
		),
	)

	fmt.Println(
		colors["yellow"](
			fmt.Sprintf(
				"Failed: %d",
				errorCount,
			),
		),
	)

	fmt.Println(
		colors["yellow"](
			fmt.Sprintf(
				"Parallel workers: %d",
				parallelCount,
			),
		),
	)

	fmt.Println(
		colors["yellow"](
			fmt.Sprintf(
				"Requests/sec: %.2f",
				requestsPerSecond,
			),
		),
	)

	fmt.Println(
		colors["yellow"](
			fmt.Sprintf(
				"Average response time: %s",
				averageDuration,
			),
		),
	)

	fmt.Println(
		colors["yellow"](
			fmt.Sprintf(
				"Fastest response: %s",
				fastestRequest,
			),
		),
	)

	fmt.Println(
		colors["yellow"](
			fmt.Sprintf(
				"Slowest response: %s",
				slowestRequest,
			),
		),
	)

	fmt.Println(
		colors["yellow"](
			fmt.Sprintf(
				"Total execution time: %s",
				totalExecutionTime,
			),
		),
	)
}
