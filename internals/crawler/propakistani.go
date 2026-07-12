package crawler

import (
	"GoNewsScrapper/internals/database"
	"context"
	"log"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/fetch"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func CrawlProPakistaniPk(ctx context.Context, timeout time.Duration) ([]database.CreateNewsParams, error) {
	url := "https://propakistani.pk"

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath("/usr/bin/google-chrome"),

		chromedp.Headless,
		chromedp.DisableGPU,
		chromedp.NoSandbox,
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,

		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36"),
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
		network.SetCacheDisabled(true),
		fetch.Enable().WithPatterns([]*fetch.RequestPattern{
			{
				URLPattern:   "*",
				RequestStage: fetch.RequestStageRequest,
			},
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			chromedp.ListenTarget(ctx, func(ev interface{}) {
				switch e := ev.(type) {
				case *fetch.EventRequestPaused:
					go func() {
						c := chromedp.FromContext(ctx)
						if c == nil || c.Target == nil {
							return // target not ready, drop safely instead of panicking
						}
						execCtx := cdp.WithExecutor(ctx, c.Target)

						switch e.ResourceType {
						case network.ResourceTypeImage,
							network.ResourceTypeStylesheet,
							network.ResourceTypeFont,
							network.ResourceTypeMedia:
							_ = fetch.FailRequest(e.RequestID, network.ErrorReasonBlockedByClient).Do(execCtx)
						default:
							_ = fetch.ContinueRequest(e.RequestID).Do(execCtx)
						}
					}()
				}
			})
			return nil
		}),
		chromedp.Navigate(url),
		chromedp.WaitVisible("div.teams-news"),
		chromedp.Sleep(100*time.Second),
		chromedp.Click("input#pp-load-posts", chromedp.ByQuery),
		chromedp.Click("input#pp-load-posts", chromedp.ByQuery),
		chromedp.Click("input#pp-load-posts", chromedp.ByQuery),
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
		return nil, err
	}

	log.Printf("crawling finished, crawled %d news", len(news))
	return news, nil
}
