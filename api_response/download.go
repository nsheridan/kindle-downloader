package api_response

type USBDownloadResponse struct {
	DownloadViaUSB struct {
		Success bool   `json:"success"`
		URL     string `json:"URL"`
	} `json:"DownloadViaUSB"`
}
