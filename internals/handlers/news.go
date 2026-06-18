package handlers

import (
	"GoNewsScrapper/internals/database"
	"GoNewsScrapper/internals/service"
	"encoding/json"
	"fmt"
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

func (n *News) CreateBulkNews(w http.ResponseWriter, r *http.Request) {
	var news []database.CreateNewsParams
	if err := json.NewDecoder(r.Body).Decode(&news); err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}
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
