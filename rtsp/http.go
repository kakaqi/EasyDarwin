package rtsp

import (
	"fmt"
	"github.com/penggy/EasyGoLib/utils"
	"io/ioutil"
	"net/http"
	"strings"
)

type Media struct {
	Url        string
	HttpClient *http.Client
	Chan       chan []string
}

func NewsSyncStream() *Media {
	media := &Media{
		Url: utils.Conf().Section("rtsp").Key("sync_stream_url").MustString(""),
		HttpClient: &http.Client{
			Timeout: 0,
		},
		Chan: make(chan []string, 100),
	}

	media.Run()
	return media

}

func (this *Media) Run() {
	go func() {
		for {
			select {
			case data := <-this.Chan:
				this.SyncStream(data[0], data[1], data[2])
			}
		}
	}()
}
func (this *Media) SyncStream(path string, type_ string, time string) {
	api := "/live/app/stream/sync"
	params := fmt.Sprintf("?path=%s&type=%s&time=%s", path, type_, time)
	fmt.Println(params)
	url := this.Url + api + params
	_, err := this.sendGet(url)
	if err != nil {
		fmt.Println(err)
	}
	return
}
func (this *Media) sendGet(url string) (body []byte, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	resp, err := this.HttpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	if i := strings.Index(resp.Status, "200"); i == -1 {
		err = fmt.Errorf("请求失败:%s", body)
		return
	}

	return
}
