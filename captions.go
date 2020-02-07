package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"io"
	"log"
	"os"
)

var (
	videoId = flag.String("vid", "1Swj8xHC_Rs", "Youtube video Id")
	tfmt    = flag.String("tfmt", "srt", "caption format , default srt.")
)

func main() {
	flag.Parse()
	client := getClient(youtube.YoutubeForceSslScope)
	service, err := youtube.NewService(context.TODO(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("create service failed:%v\n", err)
	}

	captions, err := List(service, *videoId)
	if err != nil {
		log.Fatal("get caption list failed:", err)
	}

	for _, caption := range captions {
		fmt.Printf("%+v\n\n", caption.Snippet)
		fileName := fmt.Sprintf("%s.%s.%s", caption.Snippet.VideoId, caption.Snippet.Language, *tfmt)
		err = Download(service, caption.Id, fileName)
		if err != nil {
			log.Printf("download caption failed:%v", err)
		}
	}

}

func List(service *youtube.Service, videoId string) (captions []*youtube.Caption, err error) {
	call := service.Captions.List("snippet", videoId)
	resp, err := call.Do()
	if err != nil {
		return nil, err
	}
	captions = resp.Items
	return captions, nil
}

func Download(service *youtube.Service, captionId string, filename string) (err error) {
	call := service.Captions.Download(captionId)
	resp, err := call.Tfmt(*tfmt).Download()
	if err != nil {
		return
	}
	defer resp.Body.Close()
	return download(resp.Body, filename)

}

func download(reader io.Reader, filename string) (err error) {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	_, err = io.Copy(file, reader)
	return
}
