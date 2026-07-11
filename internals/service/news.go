package service

import (
	"GoNewsScrapper/internals/database"
	"GoNewsScrapper/internals/repository"
	"GoNewsScrapper/internals/utils"
	"context"
	"log"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

type News struct {
	repo *repository.News
}

func NewNews(repo *repository.News) *News {
	return &News{repo: repo}
}

func (n *News) CreateNews(ctx context.Context, news []database.CreateNewsParams) error {
	return n.repo.InsertNewsBulk(ctx, news)
}

func (n *News) GetAllNews(ctx context.Context) ([]database.News, error) {
	return n.repo.GetAllNews(ctx)
}

func (n *News) CrawlNews() error {
	url := "https://propakistani.pk"

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath("/usr/bin/chromium-browser"),

		chromedp.Headless,
		chromedp.DisableGPU,
		chromedp.NoSandbox,
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,

		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-background-networking", true),
		chromedp.Flag("disable-background-timer-throttling", true),
		chromedp.Flag("disable-renderer-backgrounding", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	var news []database.CreateNewsParams

	log.Println("crawling started...")
	err := chromedp.Run(ctxWithTimeout,
		network.Enable(),
		chromedp.ActionFunc(func(ctx context.Context) error {
			return network.SetBlockedURLs().WithURLPatterns([]*network.BlockPattern{
				{URLPattern: "*://*/*.jpg"},
				{URLPattern: "*://*/*.png"},
				{URLPattern: "*://*/*.gif"},
				{URLPattern: "*://*/*.avif"},
			}).Do(ctx)
		}),
		chromedp.Navigate(url),
		chromedp.WaitVisible("div.teams-news"),
		chromedp.Click("input#pp-load-posts", chromedp.ByQuery),
		chromedp.Sleep(2*time.Second),
		chromedp.Click("input#pp-load-posts", chromedp.ByQuery),
		chromedp.Sleep(2*time.Second),
		chromedp.Click("input#pp-load-posts", chromedp.ByQuery),
		chromedp.Sleep(2*time.Second),
		chromedp.Evaluate(`
			Array.from(document.querySelectorAll(".g1-mosaic-item:has(.g1-frame img), div.tnews-inner.relv")).map(el=>{
				let news = {}
				news.newsUrl = el.querySelector("a").href
				news.imageUrl = el.querySelector("img").src
				news.title = el.querySelector("h3, h5")?.textContent
				news.slug = news.newsUrl?.split('/')?.at(-2)
				return news
			})
		`, &news),
	)
	if err != nil {

		return err
	}

	news = utils.RemoveNewsDuplicatesBySlug(news)

	err = n.CreateNews(ctx, news)
	if err != nil {
		return err
	}

	log.Printf("crawling finished, crawled %d news", len(news))
	return nil
}
