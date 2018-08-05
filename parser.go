package main

import (
	"os"
	"fmt"
	"golang.org/x/net/html"
	"strings"
)

var game_map = [][]string {
	{"7.", "8.", "9."},
	{"4.", "5.", "6."},
	{"1.", "2.", "3."},
}

func parse_file(file string, out string) {

	game := [][]string{
		{"_", "_", "_"},
		{"_", "_", "_"},
		{"_", "_", "_"},
	}
	var (
		i = 0
		j = 0
	)

	fmt.Printf("Open: %s, out: %s\n", file, out)
	r, err := os.Open(file)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	doc, err := html.Parse(r)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}

	var reason string

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" && a.Val != "../index.html" {
					game[i][j] = "_"

					parse_file(a.Val, game_map[i][j] + out)

					j++
					if j == 3 {
						j = 0;
						i++;
					}
					break
				}
			}
		}
		if n.Type == html.TextNode && (n.Data == "o" || n.Data == "x") {
			game[i][j] = n.Data;
			j++
			if j == 3 {
				j = 0;
				i++;
			}
		}
		if n.Type == html.TextNode && (n.Data == "It's a tie!" || n.Data ==  "Player O wins!" || n.Data == "Player X wins!") {
			reason = n.Data + "\n"
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	w, err := os.Create(out)
	if err != nil {
		fmt.Println(err)
	}
	if reason != "" {
		w.WriteString(reason)
	}
	for i := 0; i < 3; i++ {
		w.WriteString(strings.Join(game[i][0:3], "|") + "\n")
	}
	w.Close()
}

func main() {
	if(len(os.Args) == 2) {
		parse_file(os.Args[1], "game.f1remoon.com.txt")
	}
}
