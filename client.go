package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func NewHTTPClient(url string) *HTTPClient {
	return &HTTPClient{
		client:  &http.Client{},
		baseUrl: url,
	}
}

func (c *HTTPClient) getLoaderToken() error {
	var token LoaderToken
	url := fmt.Sprintf("%s/%s", os.Getenv("LOADORDER_BASE_URL"), os.Getenv("LOADORDER_TOKEN_ROUTE"))
	req := GetTokenRequest{
		UserName: os.Getenv("LOADER_USER"),
		Pass:     os.Getenv("LOADER_PASS"),
		Project:  os.Getenv("LOADER_PROJECT"),
	}
	jsonBody, _ := json.Marshal(req)
	resp, err := c.doRequest("POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatal("can`t get jwt token", err)
	}

	body, err := io.ReadAll(resp.Body)

	if err := json.Unmarshal(body, &token); err != nil {
		fmt.Println("error to unmarshal body:", err)
	}

	tokenWithClaims, _, err := new(jwt.Parser).ParseUnverified(token.Token, &TokenClaims{})
	if err != nil {
		log.Fatal(err)
	}

	if claims, ok := tokenWithClaims.Claims.(*TokenClaims); ok {
		token.Exp = claims.Exp
	}

	c.token = token
	return nil
}

func (c *HTTPClient) loadOrder(fileName string) LoadOrderResponse {
	fmt.Printf("Client is: %+v\n", c)

	if c.token.Exp < time.Now().Unix() {
		fmt.Println("Token is expired, get new")
		err := c.getLoaderToken()
		if err != nil {
			log.Fatal("error to get jwt token", err)
		}
	}

	url := fmt.Sprintf("%s/%s", os.Getenv("LOADORDER_BASE_URL"), os.Getenv("LOADORDER_UPLOAD_ROUTE"))

	resp, err := c.doRequestWithJWT("POST", url, c.token, fileName)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatal("can`t load file", err)
	}

	responseBody, _ := io.ReadAll(resp.Body)

	var result LoadOrderResponse
	err = json.Unmarshal(responseBody, &result)
	if err != nil {
		fmt.Println("error to unmarshal body from response", err)
	}

	fmt.Printf("Loader result: %+v\n", result)
	return result
}

func (c *HTTPClient) setStatus(taskUrl string, statusID string) {
	url := fmt.Sprintf("%s/transitions", taskUrl)
	body := map[string]any{
		"transition": map[string]string{
			"id": statusID},
	}
	jsonBody, _ := json.Marshal(body)

	resp, err := c.doRequest("POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		fmt.Println("can`t change status", err)
	}
	defer resp.Body.Close()

}

func (c *HTTPClient) addComment(taskUrl string, comment string) {
	url := fmt.Sprintf("%s/comment", taskUrl)
	body := map[string]string{"body": comment}
	jsonBody, _ := json.Marshal(body)

	resp, err := c.doRequest("POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		fmt.Println("can`t add new comment", err)
	}
	defer resp.Body.Close()

}

func (c *HTTPClient) addCommentWithResult(taskUrl string, result LoadOrderResponse) error {
	responseJson, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}
	comment := JiraComment{Body: fmt.Sprintf("Load Order Response:\n%s\n", string(responseJson))}
	url := fmt.Sprintf("%s/comment", taskUrl)

	jsonBody, err := json.Marshal(comment)
	if err != nil {
		return err
	}

	resp, err := c.doRequest("POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		fmt.Println("can`t add new comment", err)
	}
	defer resp.Body.Close()
	return nil
}

func (c *HTTPClient) GetAttachment(fileName string, content string, taskID string) string {
	resp, err := c.doRequest("GET", content, nil)
	if err != nil {
		log.Fatal("error to get attachment", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatal(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("error to read from body", err)
	}

	filePath, err := seveOrderFile(body, fileName, taskID)
	if err != nil {
		log.Fatal("error to seve file", err)
	}

	fmt.Printf("file successfuly saved to: %s\n", filePath)
	return filePath
}

func (c *HTTPClient) GetTaskDescription(taskUrl string) string {
	resp, err := c.doRequest("GET", taskUrl, nil)
	if err != nil {
		fmt.Print("error to get task description", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("response error", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error to read body", err)
	}

	var task Task
	if err := json.Unmarshal(body, &task); err != nil {
		fmt.Println("error to unmarshal body:", err)
	}

	fmt.Printf("%+v\n", task)

	return "done"
}

func (c *HTTPClient) GetTask(taskUrl string) Task {
	resp, err := c.doRequest("GET", taskUrl, nil)
	if err != nil {
		log.Fatal("error to get task", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatal("can`t get task", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("error to read body", err)
	}

	var task Task
	if err := json.Unmarshal(body, &task); err != nil {
		log.Fatal("error to unmarshal body:", err)
	}
	return task
}

func (c *HTTPClient) doRequest(method string, url string, body io.Reader) (*http.Response, error) {
	var (
		headerKey = "X-ExperimentalApi"
		headerVal = "opt-in"
	)

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(headerKey, headerVal)
	req.SetBasicAuth(os.Getenv("JIRA_USER"), os.Getenv("JIRA_PASS"))

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *HTTPClient) doRequestWithJWT(method string, url string, jwt LoaderToken, fileName string) (*http.Response, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fileBase := filepath.Base(fileName)
	part, err := writer.CreateFormFile("scvFile", fileBase)
	if err != nil {
		fmt.Println("error creating form file:", err)
		return nil, fmt.Errorf("error creating form file: %w", err)
	}

	if _, err := io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("error copying file: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("error closing writer: %w", err)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	bearerToken := strings.Join([]string{jwt.Type, jwt.Token}, " ")
	req.Header.Set("Authorization", bearerToken)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *HTTPClient) GetQueues(addUrl string) Queue {
	url := fmt.Sprintf("%s/%s", c.baseUrl, addUrl)

	resp, err := c.doRequest("GET", url, nil)
	if err != nil {
		log.Fatal("error to get Queues", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatal("can`t get queues", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("error to read body", err)
	}

	var queues Queue
	if err := json.Unmarshal(body, &queues); err != nil {
		log.Fatal("error to unmarshal body", err)
	}
	return queues
}
