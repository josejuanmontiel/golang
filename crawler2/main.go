package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
)

func main() {

	htmlOcasion, _ := GetHttpHtmlContent(
		"https://www.ocasionplus.com/coches-ocasion",
		"body",
		".twelve > div:nth-child(1) > div:nth-child(2) > div:nth-child(1) > div:nth-child(1) > div:nth-child(1) > div:nth-child(2) > div:nth-child(1) > div:nth-child(1) > a:nth-child(1) > h2:nth-child(1)",
	)

	fmt.Println(htmlOcasion)

	// create a test server to serve the page
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Title</title>
</head>
<body>
<h1 id="title" class="link">
    <a href="https://test.com/helloworld">
        content of h1 1
    </a>
    <span>hello</span> world
</h1>
</body>
</html>
`,
		)
	}))
	defer ts.Close()

	// create chrome instance
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	start := time.Now()
	// navigate to a page, wait for an element, click
	var res string
	err := chromedp.Run(ctx,
		emulation.SetUserAgentOverride("WebScraper 1.0"),
		chromedp.Navigate(`https://github.com`),
		// wait for footer element is visible (ie, page is loaded)
		chromedp.ScrollIntoView(`footer`),
		// chromedp.WaitVisible(`footer < div`),
		chromedp.Text(`h1`, &res, chromedp.NodeVisible, chromedp.ByQuery),
	)
	if err != nil {
		log.Fatal(err)
	}

	// run task list
	err = chromedp.Run(ctx, travelSubtree(ts.URL, `title`, chromedp.ByID))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("h1 contains: '%s'\n", res)
	fmt.Printf("\nTook: %f secs\n", time.Since(start).Seconds())
}

//Get the data crawled from the website
func GetHttpHtmlContent(url string, selector string, sel interface{}) (string, error) {
	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", false), // debug using
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36`),
	}
	//Initialization parameters, first pass an empty data
	options = append(chromedp.DefaultExecAllocatorOptions[:], options...)

	c, _ := chromedp.NewExecAllocator(context.Background(), options...)

	// create context
	chromeCtx, cancel := chromedp.NewContext(c, chromedp.WithLogf(log.Printf))
	//Execute an empty task to create a chrome instance in advance
	chromedp.Run(chromeCtx, make([]chromedp.Action, 0, 1)...)

	//Create a context with a timeout of 40s
	timeoutCtx, cancel := context.WithTimeout(chromeCtx, 40*time.Second)
	defer cancel()

	chromedp.Query("home")

	x := &chromedp.Selector{"", nil, 1, nil, nil}

	var htmlContent string
	err := chromedp.Run(timeoutCtx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(selector),
		// chromedp.OuterHTML(sel, &htmlContent, chromedp.ByQuery),

		chromedp.Nodes(""),
	)
	if err != nil {
		return "", err
	}
	//log.Println(htmlContent)

	return htmlContent, nil
}

// travelSubtree illustrates how to ask chromedp to populate a subtree of a node.
//
// https://github.com/chromedp/chromedp/issues/632#issuecomment-654213589
// @mvdan explains why node.Children is almost always empty:
// Nodes are only obtained from the browser on an on-demand basis.
// If we always held the entire DOM node tree in memory,
// our CPU and memory usage in Go would be far higher.
// And chromedp.FromNode can be used to retrieve the child nodes.
//
// Users get confused sometimes (why node.Children is empty while node.ChildNodeCount > 0?).
// And some users want to travel a subtree of the DOM more easy.
// So here comes the example.
func travelSubtree(pageUrl, of string, opts ...chromedp.QueryOption) chromedp.Tasks {
	var nodes []*cdp.Node
	return chromedp.Tasks{
		chromedp.Navigate(pageUrl),
		chromedp.Nodes(of, &nodes, opts...),
		// ask chromedp to populate the subtree of a node
		chromedp.ActionFunc(func(c context.Context) error {
			// depth -1 for the entire subtree
			// do your best to limit the size of the subtree
			return dom.RequestChildNodes(nodes[0].NodeID).WithDepth(-1).Do(c)
		}),
		// wait a little while for dom.EventSetChildNodes to be fired and handled
		chromedp.Sleep(time.Second),
		chromedp.ActionFunc(func(c context.Context) error {
			printNodes(nodes, 0)
			return nil
		}),
	}
}

func printNodes(nodes []*cdp.Node, indent int) {
	spaces := strings.Repeat(" ", indent)
	for _, node := range nodes {
		fmt.Print(spaces)
		var extra interface{}
		if node.NodeName == "#text" {
			extra = node.NodeValue
		} else {
			extra = node.Attributes
		}
		fmt.Printf("%s: %q\n", node.NodeName, extra)
		if node.ChildNodeCount > 0 {
			printNodes(node.Children, indent+4)
		}
	}
}
