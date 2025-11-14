package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os" // (atau "io.Reader" jika file ada di memori)
)

// Struct untuk menampung respon JSON dari Pinata
type PinataResponse struct {
	IpfsHash  string `json:"IpfsHash"`
	PinSize   int    `json:"PinSize"`
	Timestamp string `json:"Timestamp"`
}

// uploadToPinata adalah fungsi internal backend Anda
// (filePath adalah path ke file yang sudah aman & lolos moderasi)
func UploadToPinata(filePath string) (*PinataResponse, error) {

	// 1. Dapatkan Kunci Rahasia Anda dari .env
	// JANGAN PERNAH hardcode kunci Anda
	pinataJWT := os.Getenv("PINATA_JWT_KEY")
	if pinataJWT == "" {
		return nil, fmt.Errorf("PINATA_JWT_KEY tidak ditemukan")
	}

	// 2. Siapkan file untuk di-upload
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("gagal membuka file: %w", err)
	}
	defer file.Close()

	// 3. Buat 'body' request multipart
	// (Ini adalah cara standar untuk 'mengunggah file' via HTTP)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", file.Name())
	log.Println(file.Name())
	if err != nil {
		return nil, err
	}
	io.Copy(part, file)
	writer.Close() // Wajib ditutup agar 'boundary' ditulis

	// 4. Buat request HTTP
	url := "https://api.pinata.cloud/pinning/pinFileToIPFS"
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	// 5. SET HEADER (PENTING)
	// Set 'Content-Type' yang benar
	req.Header.Set("Content-Type", writer.FormDataContentType())
	// Set Kunci JWT Anda
	req.Header.Set("Authorization", "Bearer "+pinataJWT)

	// 6. Jalankan request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 7. Baca respon
	if resp.StatusCode != http.StatusOK {
		// (Baca body error dari Pinata untuk debug)
		return nil, fmt.Errorf("upload gagal, status: %s", resp.Status)
	}

	// 8. Decode respon JSON
	var pinataResp PinataResponse
	if err := json.NewDecoder(resp.Body).Decode(&pinataResp); err != nil {
		return nil, fmt.Errorf("gagal decode respon Pinata: %w", err)
	}

	return &pinataResp, nil
}
