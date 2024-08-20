package api

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"io"
	"time"
)

type Lib struct{}

func NewLib() *Lib {
	return &Lib{}
}

func (*Lib) SleepMill(t int64) {
	time.Sleep(time.Millisecond * time.Duration(t))
}

type UnData struct {
	FileName string        `json:"fileName"`
	Data     *bytes.Buffer `json:"data"`
}

func (l *Lib) Unzip(data []byte) ([]UnData, error) {
	reader := bytes.NewReader(data)
	r, err := zip.NewReader(reader, int64(len(data)))
	if err != nil {
		return nil, err
	}
	arr := make([]UnData, 0)
	for _, file := range r.File {
		buf1, err := l.unzip(file)
		if err != nil {
			return nil, err
		}
		arr = append(arr, UnData{file.Name, buf1})
	}
	return arr, nil
}

func (*Lib) unzip(file *zip.File) (*bytes.Buffer, error) {
	rc, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	var buf bytes.Buffer
	_, err = io.Copy(&buf, rc)
	if err != nil {
		return nil, err
	}
	return &buf, nil
}

func (*Lib) UnGzip(data []byte) (*bytes.Buffer, error) {
	gz, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer gz.Close()
	// 读取解压后的数据
	body, err := io.ReadAll(gz)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(body), nil
}
