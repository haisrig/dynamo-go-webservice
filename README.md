# dynamo-go-webservice
Go Web service uploads data into AWS DynamoDB. It has 3 endpoints.
1. /employees/     -- Its a GET request. To verify server is up or not.
2. /employees/list -- Its a GET request. Gets the **employees** data from the DynamoDB.
3. /employeed/add  -- Its a PUT request. Reads the employee data in JSON format and upload to DynamoDB.

# Dependencies
AWS DyanmoDB module
`go get github.com/aws/aws-sdk-go/`

# Test scripts
## HTTP Get request with sequential processing
```
func testWithoutConcurrency() {
	start := time.Now()
	for i := 1; i <= 1000; i++ {
		http.Get("http://localhost:8000/employees/list")
	}
	fmt.Println("Without Concurrency", time.Now().Sub(start))
}
```

## HTTP Get request with concurrent processing
```
func testWithConcurrency() {
	wg := sync.WaitGroup{}
	start := time.Now()
	for i := 1; i <= 1000; i++ {
		wg.Add(1)
		go func(i int) {
			http.Get("http://localhost:8000/employees/list")
			wg.Done()
			fmt.Println("Completed ", i)
		}(i)
	}
	wg.Wait()
	fmt.Println("With Concurrency", time.Now().Sub(start))
}
```

## HTTP PUT request with sequential processing
```
func testPostWithoutConcurrency(c string) {
	start := time.Now()
	for i := 1; i <= 1000; i++ {
		requestBody, _ := json.Marshal(map[string]string{
			"Empid":    "101" + strconv.Itoa(i) + c,
			"Name":     "ABC" + strconv.Itoa(i),
			"Age":      "35",
			"Location": "India",
		})
		req, _ := http.NewRequest(http.MethodPut, "http://localhost:8000/employees/add", bytes.NewBuffer(requestBody))
		_, _ = (&http.Client{}).Do(req)
	}
	fmt.Println("Post Without Concurrency", time.Now().Sub(start))
}
```

## HTTP PUT request with concurrent processing
```
func testPostWithConcurrency(c string) {
	start := time.Now()
	wg := sync.WaitGroup{}
	for i := 1; i <= 1000; i++ {
		wg.Add(1)
		go func(i int) {
			fmt.Println("Sending data to DynamoDB -", i)
			requestBody, _ := json.Marshal(map[string]string{
				"Empid":    "101" + strconv.Itoa(i) + c,
				"Name":     "ABC" + strconv.Itoa(i),
				"Age":      "35",
				"Location": "India",
			})
			req, _ := http.NewRequest(http.MethodPut, "http://localhost:8000/employees/add", bytes.NewBuffer(requestBody))
			_, _ = (&http.Client{}).Do(req)
			wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Println("Post With Concurrency", time.Now().Sub(start))
}
```