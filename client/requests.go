package client

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/url"
	"strings"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
)

type HttpClient interface {
	GetCookies(u *url.URL) []*http.Cookie
	SetCookies(u *url.URL, cookies []*http.Cookie)
	SetCookieJar(jar http.CookieJar)
	SetProxy(proxyUrl string) error
	GetProxy() string
	SetFollowRedirect(followRedirect bool)
	GetFollowRedirect() bool
	Do(req *http.Request) (*http.Response, error)
	Get(url string) (resp *http.Response, err error)
	Head(url string) (resp *http.Response, err error)
	Post(url, contentType string, body io.Reader) (resp *http.Response, err error)
}

var TLS_DEFAULT = tls_client.DefaultClientProfile

// proxy url format: http://username:password@ip:port
func GetTLS(proxy string, clientProfile tls_client.ClientProfile) HttpClient {
	jar := tls_client.NewCookieJar()
	timeout := 30

	if len(proxy) > 0 {
		timeout = 60
	}

	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(timeout),
		tls_client.WithClientProfile(clientProfile),
		tls_client.WithNotFollowRedirects(),
		tls_client.WithCookieJar(jar),
	}

	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)

	if err != nil {
		log.Println(err)
		return nil
	}
	if len(proxy) > 0 {
		err = client.SetProxy(proxy)

		// println("proxy", proxy)

		if err != nil {
			return nil
		}
	}

	return client
}

func TlsRequest(params TLSParams) ([]byte, error) {

	req, err := http.NewRequest(params.Method, params.Url, params.Body)

	if err != nil {
		return nil, err
	}

	req.Header = params.Headers

	ua := strings.ToLower(req.Header.Get("User-Agent"))

	if strings.Contains(ua, "android") || strings.Contains(ua, "okhttp") {
		params.Client = GetTLS(params.Client.GetProxy(), tls_client.Okhttp4Android11)
	} else if strings.Contains(ua, "ios") || strings.Contains(ua, "iphone") || strings.Contains(ua, "ipad") {
		params.Client = GetTLS(params.Client.GetProxy(), tls_client.Safari_IOS_15_5)
	}

	resp, err := params.Client.Do(req)

	if err != nil {
		return nil, err
	}

	defer func(body io.ReadCloser) {
		body.Close()
	}(resp.Body)

	readBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if params.ExpectedResponse == resp.StatusCode {
		return readBytes, nil
	} else {
		type response struct {
			StatusCode int    `json:"statusCode"`
			Body       []byte `json:"body"`
		}
		errorMsg := response{StatusCode: resp.StatusCode, Body: readBytes}
		jsonBytes, err := json.Marshal(errorMsg)
		if err != nil {
			return nil, err
		}
		errorMsgString := string(jsonBytes)
		return nil, errors.New(errorMsgString)
	}

}

type FullRequestRes struct {
	StatusCode int    `json:"statusCode"`
	Body       []byte `json:"body"`
}

func TlsFullRequest(params TLSParams) ([]byte, http.Response, error) {

	req, err := http.NewRequest(params.Method, params.Url, params.Body)

	if err != nil {
		return nil, http.Response{}, err
	}

	req.Header = params.Headers

	resp, err := params.Client.Do(req)

	if err != nil {
		return nil, http.Response{}, err
	}

	defer resp.Body.Close()

	readBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, http.Response{}, err
	}

	if params.ExpectedResponse != resp.StatusCode {
		type response struct {
			StatusCode int    `json:"statusCode"`
			Body       []byte `json:"body"`
		}
		errorMsg := response{StatusCode: resp.StatusCode, Body: readBytes}
		jsonBytes, err := json.Marshal(errorMsg)
		if err != nil {
			return nil, http.Response{}, err
		}
		errorMsgString := string(jsonBytes)
		return nil, http.Response{}, errors.New(errorMsgString)
	}

	return readBytes, *resp, nil

}
