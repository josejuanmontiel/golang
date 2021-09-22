package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("headless", false))...)
	defer cancel()
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	// // create chrome instance
	// ctx, cancel := chromedp.NewContext(
	// 	context.Background(),
	// 	chromedp.WithLogf(log.Printf),
	// )
	// defer cancel()

	// // create a timeout
	// ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	// defer cancel()
	start := time.Now()
	// navigate to a page, wait for an element, click
	var res1, res2 string
	err := chromedp.Run(ctx,
		// emulation.SetUserAgentOverride("WebScraper 1.0"),
		chromedp.Navigate(`https://www.ocasionplus.com/coches-ocasion`),
		// wait for footer element is visible (ie, page is loaded)
		// chromedp.ScrollIntoView(`footer`),
		// chromedp.WaitVisible(`footer < div`),
		// chromedp.Text(`h1`, &res, chromedp.NodeVisible, chromedp.ByQuery),
		chromedp.Text(`.twelve > div:nth-child(1) > div:nth-child(2) > div:nth-child(1) > div:nth-child(1) > div:nth-child(1) > div:nth-child(2) > div:nth-child(1) > div:nth-child(1) > a:nth-child(1)`, &res1, chromedp.NodeVisible, chromedp.ByQuery),
		chromedp.Text(`.twelve > div:nth-child(1) > div:nth-child(2) > div:nth-child(2) > div:nth-child(1) > div:nth-child(1) > div:nth-child(2) > div:nth-child(1) > div:nth-child(1) > a:nth-child(1) > h2:nth-child(1)`, &res2, chromedp.NodeVisible, chromedp.ByQuery),
		//      	   .twelve > div:nth-child(1) > div:nth-child(2) > div:nth-child(4) > div:nth-child(2) > div:nth-child(1) > div:nth-child(2) > div:nth-child(1) > div:nth-child(1) > a:nth-child(1) > h2:nth-child(1)

	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("1º: '%s'\n", res1)
	fmt.Printf("2º: '%s'\n", res2)

	fmt.Printf("\nTook: %f secs\n", time.Since(start).Seconds())
}
