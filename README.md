# Amazon Kindle book downloader

Bulk download kindle books from your Amazon account.

## Running

```
% kindle-downloader -help
Usage of kindle-downloader:
  -amazon-url string
    	Amazon Country URL (default "https://www.amazon.co.uk")
  -concurrency int
    	Number of concurrent downloads (default 20)
  -manual-login
    	Manually login to Amazon
  -output string
    	Output directory for downloads (default "books")
```

By default the program will prompt for credentials, if this isn't working using the `-manual-login` flag will allow you to login using the browser.

## Docker Container

A container with everything necessary to run is available:

```
% docker run -ti -v `pwd`/books:/books ghcr.io/nsheridan/kindle-downloader -output /books
```
