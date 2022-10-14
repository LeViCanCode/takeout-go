package takeout

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/html"
	"io"
	"net/http"
	"os"
	"path"
)

type (
	//Client The base client for Takeout
	Client struct {
		//token the token from Takeout
		token string
		//baseUrl the base url used by the whole lib
		baseUrl string
		//httpClient the http Client that should be used by the library
		httpClient *http.Client
		//isLoggedIn checked if the client has been logged in (token checked)
		isLoggedIn bool
		// logger
		logger *log.Logger
	}
	//ClientOptions the options for the client
	ClientOptions struct {
		Token      string
		Debug      bool
		HttpClient *http.Client
	}
)

// New Create a new Client
func (c Client) New(options ClientOptions) Client {
	c.token = options.Token
	c.baseUrl = "https://takeout.bysourfruit.com/api"
	c.logger = log.New()
	c.logger.SetLevel(log.ErrorLevel)
	if options.Debug {
		c.logger.SetLevel(log.DebugLevel)
	}
	c.httpClient = options.HttpClient
	return c
}

func (c Client) Login() (*Client, error) {
	// Make the body according to https://github.com/Takeout-bysourfruit/takeout.js/blob/main/src/index.js#L32
	body, err := json.Marshal(struct {
		Token string `json:"token"`
	}{Token: c.token})
	if err != nil {
		return nil, err
	}
	// Actually make the http request
	resp, err := c.httpClient.Post(c.baseUrl+"/auth/verify", "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, FailedToVerifyToken
	}
	defer resp.Body.Close()
	c.isLoggedIn = true
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	c.logger.Debug("Verified token, " + string(respBody))
	return &c, nil
}

// GetLocalTemplate Returns minified html code from a file
func (c Client) GetLocalTemplate(p string) (string, error) {
	data, err := os.ReadFile(path.Clean(p))
	if err != nil {
		return "", err
	}
	h, err := minifyHTML(data)
	if err != nil {
		return "", err
	}
	return h, nil
}

// GetCloudTemplate Returns minified template from the Takeout cdn
func (c Client) GetCloudTemplate(name string) (string, error) {
	resp, err := c.httpClient.Get(fmt.Sprintf("https://cdn-takeout.bysourfruit.com/cloud/read?name=%s&token=%s", name, c.token))
	if err != nil {
		return "", nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Debugf("Failed got %v instead of 200 status code", resp.StatusCode)
		return "", FailedToGetCloudTemplate
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	// minify out of good measure (realistically I don't think you have to) - Max
	return minifyHTML(respBody)
}

func minifyHTML(data []byte) (string, error) {
	m := minify.New()
	m.AddFunc("text/html", html.Minify)
	return m.String("text/html", string(data))
}

// SendEmail
