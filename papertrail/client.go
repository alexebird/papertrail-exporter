package papertrail

import (
	//"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	client  *http.Client
	baseUrl = "https://papertrailapp.com/api/v1"
)

func papertrailToken() string {
	return os.Getenv("PAPERTRAIL_API_TOKEN")
}

func addPapertrailTokenHeader(req *http.Request) {
	req.Header.Add("X-Papertrail-Token", papertrailToken())
}

func getJson(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", baseUrl+url, nil)
	if err != nil {
		return nil, err
	}
	addPapertrailTokenHeader(req)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil {
		return nil, err
	}

	return body, nil
}

func init() {
	client = &http.Client{}
}
