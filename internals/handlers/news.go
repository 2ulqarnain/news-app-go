package handlers

import (
	"GoNewsScrapper/internals/service"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type news interface {
	CreateBulkNews(w http.ResponseWriter, r *http.Request)
	GetAllNews(w http.ResponseWriter, r *http.Request)
	SearchNewsByTitle(w http.ResponseWriter, r *http.Request)
}

type News struct {
	svc *service.News
}

func (n *News) Crawl(w http.ResponseWriter, r *http.Request) {
	var err error
	go func() {
		if err = n.svc.CrawlNews(); err != nil {
			log.Println(err)
		}
	}()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("Crawled Successfully!"))
}

func (n *News) GetAllNews(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	news, err := n.svc.GetAllNews(ctx)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(news)
}

func (n *News) SearchNewsByTitle(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func NewNews(svc *service.News) *News {
	return &News{svc: svc}
}
