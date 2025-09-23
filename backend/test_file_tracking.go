package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	fmt.Println("Testing file upload with database tracking...")

	// Test file upload
	testFile := "test_resume.pdf"

	// Create a test file if it doesn't exist
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		file, err := os.Create(testFile)
		if err != nil {
			fmt.Printf("Error creating test file: %v\n", err)
			return
		}
		file.WriteString("This is a test resume file for database tracking")
		file.Close()
		defer os.Remove(testFile)
	}

	// Test upload
	fmt.Println("1. Testing file upload...")
	err := uploadFile(testFile, "http://localhost:8080/api/v1/upload")
	if err != nil {
		fmt.Printf("Upload failed: %v\n", err)
		return
	}

	// Test file listing
	fmt.Println("2. Testing file listing...")
	err = listFiles("http://localhost:8080/api/v1/files/")
	if err != nil {
		fmt.Printf("List files failed: %v\n", err)
		return
	}

	// Test file stats
	fmt.Println("3. Testing file stats...")
	err = getFileStats("http://localhost:8080/api/v1/files/stats")
	if err != nil {
		fmt.Printf("Get file stats failed: %v\n", err)
		return
	}

	fmt.Println("All tests completed successfully!")
}

func uploadFile(filename string, url string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	fw, err := w.CreateFormFile("file", filepath.Base(filename))
	if err != nil {
		return err
	}

	if _, err = io.Copy(fw, file); err != nil {
		return err
	}

	w.Close()

	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("Upload response: %s\n", string(body))
	return nil
}

func listFiles(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("list files failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("Files list: %s\n", string(body))
	return nil
}

func getFileStats(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("get stats failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("File stats: %s\n", string(body))
	return nil
}
