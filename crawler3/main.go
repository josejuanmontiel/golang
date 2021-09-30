package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

func main() {
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

	sib := ".twelve > div > div > div > div > div > div > div > div > a"
	// price := ".twelve > div > div > div > div > div > div > div > div > div > div > div > span"

	// get project link text
	var projects []*cdp.Node
	if err := chromedp.Run(ctx, chromedp.Nodes(sib+" > *", &projects)); err != nil {
		return nil, fmt.Errorf("could not get projects: %v", err)
	}

	// get links and description text
	var linksAndDescriptions []*cdp.Node
	if err := chromedp.Run(ctx, chromedp.Nodes(sib+` > h2`, &linksAndDescriptions)); err != nil {
		return nil, fmt.Errorf("could not get links and descriptions: %v", err)
	}
	// process data
	for i := 0; i < len(linksAndDescriptions); i++ {
		printNodes(linksAndDescriptions[i].Children, 1)
	}

	return nil, nil
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
