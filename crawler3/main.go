package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
)

func main() {
	// var text string
	// GetHttpHtmlContent("https://www.ocasionplus.com/coches-ocasion", "/html/body/div[1]/div[1]/div[1]/div[2]/div[3]/div/div[2]/div[20]/div/div[1]/div[2]/div/div[2]/div[1]/div/div[1]/span/text()", &text)

	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", false), // debug using
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36`),
	}

	//Initialization parameters, first pass an empty data
	options = append(chromedp.DefaultExecAllocatorOptions[:], options...)

	c, _ := chromedp.NewExecAllocator(context.Background(), options...)

	// create context
	chromeCtx, _ := chromedp.NewContext(c, chromedp.WithLogf(log.Printf))
	//Execute an empty task to create a chrome instance in advance
	chromedp.Run(chromeCtx, make([]chromedp.Action, 0, 1)...)

	// force max timeout of 15 seconds for retrieving and processing the data
	var cancel func()
	chromeCtx, cancel = context.WithTimeout(chromeCtx, 1500*time.Second)

	// navigate
	chromedp.Run(chromeCtx, chromedp.Navigate(`https://www.ocasionplus.com/coches-ocasion`))

	// list awesome go projects for the "Selenium and browser control tools."
	listAwesomeGoProjects(chromeCtx)
	chromedp.Run(chromeCtx,
		chromedp.Click(`//*[@id="page2"]`, chromedp.BySearch),
	)

	fmt.Println("-------------------")

	listAwesomeGoProjects(chromeCtx)
	chromedp.Run(chromeCtx,
		chromedp.Click(`//*[@id="page3"]`, chromedp.BySearch),
	)

	defer cancel()
}

// projectDesc contains a url, description for a project.
type projectDesc struct {
	URL, Description string
}

// listAwesomeGoProjects is the highest level logic for browsing to the
// awesome-go page, finding the specified section sect, and retrieving the
// associated projects from the page.
func listAwesomeGoProjects(ctx context.Context) (map[string]projectDesc, error) {

	sib := ".twelve > div > div > div > div > div > div > div > div"

	// get links and description text
	var linksAndDescriptions []*cdp.Node
	if err := chromedp.Run(ctx, chromedp.Nodes(sib+` > a > h2`, &linksAndDescriptions)); err != nil {
		return nil, fmt.Errorf("could not get links and descriptions: %v", err)
	}

	// get project link text
	var price []*cdp.Node
	if err := chromedp.Run(ctx, chromedp.Nodes(sib+" > div > div > div > span", &price)); err != nil {
		return nil, fmt.Errorf("could not get projects: %v", err)
	}

	var text string

	// process data
	for i := 0; i < len(linksAndDescriptions); i++ {
		printNodes(linksAndDescriptions[i].Children, 1)
		fmt.Print(" ")
		travelSubtree(price, chromedp.ByJSPath)

		// var nodes []*cdp.Node
		// var text string
		// id, _ := dom.QuerySelector(price[i].NodeID, "span").Do(ctx)
		// _ = chromedp.Nodes([]cdp.NodeID{id}, &nodes, chromedp.ByNodeID).Do(ctx)
		// chromedp.Text([]cdp.NodeID{id}, &text, chromedp.ByNodeID)
		fmt.Print(text)
	}

	return nil, nil
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
func travelSubtree(nodes []*cdp.Node, opts ...chromedp.QueryOption) chromedp.Tasks {
	return chromedp.Tasks{
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

//Get the data crawled from the website
func GetHttpHtmlContent(url string, selector string, sel interface{}) (string, error) {
	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", true), // debug using
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

	var htmlContent string
	err := chromedp.Run(timeoutCtx,
		chromedp.Navigate(url),
		// chromedp.WaitVisible(selector),
		chromedp.OuterHTML(sel, &htmlContent, chromedp.ByJSPath),
	)
	if err != nil {
		fmt.Println("Run err : %v\n", err)
		return "", err
	}
	fmt.Println(htmlContent)

	return htmlContent, nil
}
