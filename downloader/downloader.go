package downloader

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"sync"
	"time"

	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"

	"nsheridan.dev/kindle-downloader/api_response"
)

const (
	devicePage       = "/hz/mycd/digital-console/contentlist/booksPurchases/dateDsc"
	apiEndpoint      = "/hz/mycd/digital-console/ajax"
	downloadEndpoint = "/hz/mycd/ajax"
)

var (
	httpClient = &http.Client{}
)

type failure struct {
	book api_response.ContentItem
	err  error
}

func (f failure) Error() string {
	return fmt.Sprintf("Failed to download %s: %v", f.book.Title, f.err)
}

type Downloader struct {
	Concurrency int
	CSRFToken   string
	AmazonURL   string
	Destination string
	CookieJar   http.CookieJar
}

func (d *Downloader) Download() error {
	httpClient.Jar = d.CookieJar
	device, err := d.getKindleDevice()
	if err != nil {
		return fmt.Errorf("Failed to get Kindle device: %w", err)
	}
	fmt.Println("Kindle device found:", device.DeviceName, device.DeviceSerialNumber)
	books, err := d.getBooks()
	if err != nil {
		return fmt.Errorf("Failed to get books: %w", err)
	}
	fmt.Printf("%d books found\n", len(books))
	return d.downloadAll(books, device)
}

func (d *Downloader) getKindleDevice() (api_response.Device, error) {
	req := d.newReq(http.MethodPost, d.makeURL(apiEndpoint))
	q := req.URL.Query()
	q.Add("activity", "GetDevicesOverview")
	input := map[string]any{"surfaceType": "Desktop"}
	q.Add("activityInput", jsonify(input))
	req.URL.RawQuery = q.Encode()
	resp, err := httpClient.Do(req)
	if err != nil {
		return api_response.Device{}, fmt.Errorf("Failed to get devices overview: %w", err)
	}
	defer resp.Body.Close()
	devices := api_response.GetDevicesOverviewResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&devices); err != nil {
		return api_response.Device{}, fmt.Errorf("Failed to decode devices overview: %w", err)
	}
	if !devices.Success {
		return api_response.Device{}, errors.New("Failed to get devices overview")
	}
	for _, device := range devices.GetDevicesOverview.DeviceList {
		if device.DeviceFamily == "KINDLE" {
			return device, nil
		}
	}
	return api_response.Device{}, errors.New("Kindle device not found in response")
}

func (d *Downloader) getBooks() ([]api_response.ContentItem, error) {
	req := d.newReq(http.MethodPost, d.makeURL(apiEndpoint))
	q := req.URL.Query()
	q.Add("activity", "GetContentOwnershipData")
	input := map[string]any{
		"surfaceType":              "Desktop",
		"contentType":              "Ebook",
		"contentCategoryReference": "booksPurchases",
		"itemStatusList":           []string{"Active"},
		"originTypes":              []string{"Purchase", "Pottermore"},
		"fetchCriteria": map[string]any{
			"sortOrder":         "DESCENDING",
			"sortIndex":         "DATE",
			"startIndex":        0,
			"batchSize":         999,
			"totalContentCount": -1,
		},
	}
	q.Add("activityInput", jsonify(input))
	req.URL.RawQuery = q.Encode()
	resp, err := httpClient.Do(req)
	if err != nil {
		return []api_response.ContentItem{}, fmt.Errorf("Failed to get content overview: %w", err)
	}
	defer resp.Body.Close()
	content := api_response.GetContentOwnershipDataResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&content); err != nil {
		return []api_response.ContentItem{}, fmt.Errorf("Failed to decode content overview: %w", err)
	}
	return content.GetContentOwnershipData.ContentItems, nil
}

func (d *Downloader) downloadAll(books []api_response.ContentItem, device api_response.Device) error {
	var errs []error
	os.MkdirAll(d.Destination, 0755)
	semaphore := make(chan struct{}, d.Concurrency)
	wait := sync.WaitGroup{}
	progress := mpb.New(mpb.WithWaitGroup(&wait), mpb.WithRefreshRate(100*time.Millisecond))
	for _, book := range books {
		bar := progress.AddBar(1, // actual total will be updated when we know the download size
			mpb.PrependDecorators(
				decor.Name(book.Title),
				decor.Percentage(decor.WCSyncSpace),
			),
			mpb.AppendDecorators(
				decor.OnComplete(decor.EwmaETA(decor.ET_STYLE_GO, 30), "done"),
			),
		)
		wait.Add(1)
		go func() {
			semaphore <- struct{}{}
			defer func() {
				<-semaphore
				wait.Done()
			}()
			errs = append(errs, d.downloadSingle(book, device, bar))
		}()
	}
	progress.Wait()
	return errors.Join(errs...)
}

func (d *Downloader) downloadSingle(book api_response.ContentItem, device api_response.Device, bar *mpb.Bar) error {
	// find the download URL for the book
	req := d.newReq(http.MethodPost, d.makeURL(downloadEndpoint))
	q := req.URL.Query()
	input := map[string]any{
		"param": map[string]map[string]any{
			"DownloadViaUSB": {
				"contentName":              book.Asin,
				"encryptedDeviceAccountId": device.DeviceAccountID,
				"originType":               "Purchase",
			},
		},
	}
	q.Add("data", jsonify(input))
	req.URL.RawQuery = q.Encode()
	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("Failed to get download URL: %w", err)
	}
	defer resp.Body.Close()
	data := api_response.USBDownloadResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return fmt.Errorf("Failed to decode download URL: %w", err)
	}
	if !data.DownloadViaUSB.Success {
		return errors.New("Failed to get download URL")
	}

	// download the book
	req = d.newReq(http.MethodGet, data.DownloadViaUSB.URL)
	resp, err = httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("Failed to download file: %w", err)
	}
	defer resp.Body.Close()
	filename := path.Join(d.Destination, book.Title+".azw3")
	out, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("Failed to open file %s for writing: %w", filename, err)
	}
	defer out.Close()

	bar.SetTotal(resp.ContentLength, false)
	proxy := bar.ProxyReader(resp.Body)
	defer proxy.Close()
	_, err = io.Copy(out, proxy)
	return err
}

func (d *Downloader) makeURL(path string) string {
	u, _ := url.JoinPath(d.AmazonURL, path)
	return u
}

// create a new http request with the csrf token
func (d *Downloader) newReq(method string, path string) *http.Request {
	req, _ := http.NewRequest(method, path, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	q := req.URL.Query()
	q.Add("csrfToken", d.CSRFToken)
	req.URL.RawQuery = q.Encode()
	return req
}

// convert a map of params to a json string
func jsonify(params map[string]any) string {
	b, _ := json.Marshal(params)
	return string(b)
}
