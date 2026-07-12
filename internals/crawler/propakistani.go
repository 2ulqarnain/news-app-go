package crawler

import (
	"GoNewsScrapper/internals/database"
	"context"
	"log"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func CrawlProPakistaniPk(ctx context.Context, timeout time.Duration) ([]database.CreateNewsParams, error) {
	url := "https://propakistani.pk"

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		//chromedp.ExecPath("/usr/bin/chromium-browser"),

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

	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()
	ctxNew, cancel := chromedp.NewContext(allocCtx)
	defer cancel()
	ctxWithTimeout, cancel := context.WithTimeout(ctxNew, timeout*time.Second)
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
				{URLPattern: "*://*/*.webp"},
				{URLPattern: "*://*/*.svg"},
				{URLPattern: "*://*/*.woff"},
				{URLPattern: "*://*/*.woff2"},
				{URLPattern: "*://*/*.ttf"},
				{URLPattern: "*://*/*.otf"},
				{URLPattern: "*://*/*.css"},
				{URLPattern: "*://*doubleclick/*"},
				{URLPattern: "*://*.youtube.*/*"},
				{URLPattern: "*://*googletagmanager*/*"},
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
		var html string

		if e := chromedp.Run(ctx,
			chromedp.OuterHTML("html", &html),
		); e == nil {
			log.Println(html)
		}

		return nil, err
	}

	log.Printf("crawling finished, crawled %d news", len(news))
	return news, nil
}
