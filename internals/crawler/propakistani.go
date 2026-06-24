package crawler

import (
	"GoNewsScrapper/internals/database"
	"GoNewsScrapper/internals/repository"
	"context"
	"time"

	"github.com/chromedp/chromedp"
)

type Crawler struct {
	repo *repository.News
}

func New(repo *repository.News) *Crawler {
	return &Crawler{repo: repo}
}

func (c *Crawler) CrawlProPakistaniPk() error {

	url := "https://propakistani.pk"

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.Headless,
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var pageTitle string
	var news []database.CreateNewsParams

	err := chromedp.Run(ctxWithTimeout,
		chromedp.Navigate(url),
		chromedp.WaitVisible("div.team-news"),
		chromedp.Title(&pageTitle),
		chromedp.Evaluate(`
			Array.from(document.querySelectorAll('div.tnews-inner.relv')).map(el => {
    		let news_url = el.querySelector('a').href
    		let image_url = el.querySelector('img').src
    		let title = el.querySelector('h5').textContent
    		return {news_url, image_url, title}
		})`, &news),
	)
	if err != nil {
		return err
	}

	err = c.repo.InsertNewsBulk(ctx, news)
	if err != nil {
		return err
	}

	return nil
}
