package commons

import "golang.org/x/net/html"

func GoQueryRemoveComments(n *html.Node) {
	for c := n.FirstChild; c != nil; {
		next := c.NextSibling
		if c.Type == html.CommentNode {
			n.RemoveChild(c)
		} else {
			GoQueryRemoveComments(c)
		}
		c = next
	}
}
