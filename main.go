package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"

	"github.com/go-rod/rod"
	browser_defaults "github.com/go-rod/rod/lib/defaults"
	browser_input "github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/proto"

	"golang.org/x/term"

	"nsheridan.dev/kindle-downloader/downloader"
)

func main() {
	var (
		concurrency = 20
		outputDir   = "books"
		baseURL     = "https://www.amazon.co.uk"
		manualLogin = false
		devicePage  = "/hz/mycd/digital-console/contentlist/booksPurchases/dateDsc"
	)

	flag.StringVar(&baseURL, "amazon-url", baseURL, "Amazon Country URL")
	flag.StringVar(&outputDir, "output", outputDir, "Output directory for downloads")
	flag.IntVar(&concurrency, "concurrency", concurrency, "Number of concurrent downloads")
	flag.BoolVar(&manualLogin, "manual-login", false, "Manually login to Amazon")
	flag.Parse()

	if manualLogin {
		// show the browser instead of running headless
		browser_defaults.ResetWith("show")
	}

	browser := rod.New().MustConnect()
	defer browser.Close()

	// start on a page that requires login
	startPage, err := url.JoinPath(baseURL, devicePage)
	if err != nil {
		fmt.Println("Invalid base URL", baseURL)
		os.Exit(1)
	}
	page := browser.MustPage(startPage).MustWaitStable()
	if manualLogin {
		fmt.Print("Please login to Amazon in the browser, press enter when done")
		fmt.Scanln()
	} else {
		handleLogin(page)
	}

	csrfToken := page.MustEval("() => window.csrfToken").String()

	downloader := downloader.Downloader{
		Concurrency: concurrency,
		CSRFToken:   csrfToken,
		AmazonURL:   baseURL,
		CookieJar:   cookieJar(baseURL, page.MustCookies()),
		Destination: outputDir,
	}
	err = downloader.Download()
	if err != nil {
		fmt.Println("Failed to download books:", err)
		os.Exit(1)
	}
}

// create a net/http cookie jar from the browser cookies
func cookieJar(baseURL string, browserCookies []*proto.NetworkCookie) *cookiejar.Jar {
	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse(baseURL)
	cookies := []*http.Cookie{}
	for _, cookie := range browserCookies {
		cookies = append(cookies, &http.Cookie{
			Name:  cookie.Name,
			Value: strings.Trim(cookie.Value, `"`),
		})
	}
	jar.SetCookies(u, cookies)
	return jar
}

// log into the amazon account using the browser
func handleLogin(page *rod.Page) {
	var el *rod.Element
	reader := bufio.NewReader(os.Stdin)

	el = page.MustElement("input[type='email']")
	if el != nil {
		fmt.Print("Enter email: ")
		email, _ := reader.ReadString('\n')
		el.MustInput(email).MustType(browser_input.Enter)
		page.MustWaitStable()
	}

	el = page.MustElement("input[type='password']")
	if el != nil {
		fmt.Print("Enter password: ")
		password, _ := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		el.MustInput(string(password)).MustType(browser_input.Enter)
		page.MustWaitStable()
	}

	el = page.MustElement("#auth-mfa-otpcode")
	if el != nil {
		fmt.Print("Enter OTP: ")
		otp, _ := reader.ReadString('\n')
		el.MustInput(otp).MustType(browser_input.Enter)
		page.MustWaitStable()
	}
}
