package handlers

import (
	"GoNewsScrapper/internals/service"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type News struct {
	svc *service.News
}

func (n *News) Crawl(w http.ResponseWriter, _ *http.Request) {
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
	_, _ = w.Write([]byte("Crawled Successfully!"))
}

func (n *News) GetAllNews(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	news, err := n.svc.GetAllNews(ctx)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	_ = json.NewEncoder(w).Encode(news)
}

func (n *News) SearchNewsByTitle(_ http.ResponseWriter, _ *http.Request) {
	//TODO implement me
	panic("implement me")
}

func NewNews(svc *service.News) *News {
	return &News{svc: svc}
}
