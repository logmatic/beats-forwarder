package output

import (
	cfg "github.com/logmatic/beats-forwarder/config"
	"errors"
	"bytes"
	"net/http"
)

type HTTPClient struct {
	endpoint string
}

func (c *HTTPClient) Init(config *cfg.Config) error {

	if (config.Output.HTTP.Endpoint == nil || *config.Output.HTTP.Endpoint == "" ) {
		return errors.New("No endpoint URL provided.")
	}
	c.endpoint = *config.Output.HTTP.Endpoint

	return nil
}

func (c *HTTPClient) WriteAndRetry(payload []byte) (error) {

	req, err := http.NewRequest("POST", c.endpoint, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	//req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	_, err = client.Do(req)

	if err != nil {
		return err
	}

	return nil

}

func (c *HTTPClient) Connect() (error) {
	return nil
}

func (c *HTTPClient) Close() {

}
