package main

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

const baseURL = "https://hacker-news.firebaseio.com/v0"

type Item struct {
	ID          int    `json:"id"`
	Deleted     bool   `json:"deleted,omitempty"`
	Type        string `json:"type"`
	By          string `json:"by,omitempty"`
	Time        int64  `json:"time,omitempty"`
	Text        string `json:"text,omitempty"`
	Dead        bool   `json:"dead,omitempty"`
	Parent      int    `json:"parent,omitempty"`
	Poll        int    `json:"poll,omitempty"`
	Kids        []int  `json:"kids,omitempty"`
	URL         string `json:"url,omitempty"`
	Score       int    `json:"score,omitempty"`
	Title       string `json:"title,omitempty"`
	Parts       []int  `json:"parts,omitempty"`
	Descendants int    `json:"descendants,omitempty"`
}

func getItem(id int) (*Item, error) {
	resp, err := http.Get(fmt.Sprintf("%s/item/%d.json", baseURL, id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var item Item
	if err := json.Unmarshal(body, &item); err != nil {
		return nil, err
	}
	return &item, nil
}

func getStories(storyType string) ([]int, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%sstories.json", baseURL, storyType))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var ids []int
	if err := json.Unmarshal(body, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

func stripHTML(s string) string {
	s = strings.ReplaceAll(s, "<p>", "\n\n")
	s = strings.ReplaceAll(s, "<i>", "")
	s = strings.ReplaceAll(s, "</i>", "")
	s = strings.ReplaceAll(s, "<b>", "")
	s = strings.ReplaceAll(s, "</b>", "")
	s = html.UnescapeString(s)
	for {
		start := strings.Index(s, "<")
		if start == -1 {
			break
		}
		end := strings.Index(s[start:], ">")
		if end == -1 {
			break
		}
		s = s[:start] + s[start+end+1:]
	}
	return s
}

func main() {
	var jsonOutput bool
	var count int
	var storyType string

	rootCmd := &cobra.Command{
		Use:   "tech-news",
		Short: "Hacker News CLI",
		Long: `tech-news is a high-performance command-line interface for browsing Hacker News.
It leverages Go's concurrency to fetch data in parallel, making it ideal for both
developers and AI agents who need quick access to tech news and discussions.`,
		Example: `  tech-news list --number 20
  tech-news view 8863
  tech-news list --type new --json`,
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List stories from Hacker News",
		Long:  `Fetch a list of stories from the Hacker News front page, new stories, or best stories.`,
		Example: `  tech-news list -n 15
  tech-news list -t best
  tech-news list --json`,
		Run: func(cmd *cobra.Command, args []string) {
			ids, err := getStories(storyType)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			if count > len(ids) {
				count = len(ids)
			}
			ids = ids[:count]

			var wg sync.WaitGroup
			items := make([]*Item, count)
			errs := make([]error, count)

			for i, id := range ids {
				wg.Add(1)
				go func(index, itemID int) {
					defer wg.Done()
					item, err := getItem(itemID)
					items[index] = item
					errs[index] = err
				}(i, id)
			}
			wg.Wait()

			validItems := make([]*Item, 0)
			for i, item := range items {
				if errs[i] == nil && item != nil {
					validItems = append(validItems, item)
				}
			}

			if jsonOutput {
				out, _ := json.MarshalIndent(validItems, "", "  ")
				fmt.Println(string(out))
				return
			}

			fmt.Printf("\n--- Hacker News: %s STORIES ---\n", strings.ToUpper(storyType))
			for i, item := range validItems {
				tm := time.Unix(item.Time, 0).Format("2006-01-02 15:04:05")
				fmt.Printf("%d. %s\n", i+1, item.Title)
				fmt.Printf("   %d points by %s | %s\n", item.Score, item.By, tm)
				if item.URL != "" {
					fmt.Printf("   %s\n", item.URL)
				}
				fmt.Println()
			}
		},
	}

	listCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")
	listCmd.Flags().IntVarP(&count, "number", "n", 10, "Number of stories")
	listCmd.Flags().StringVarP(&storyType, "type", "t", "top", "Type of stories (top, new, best)")

	viewCmd := &cobra.Command{
		Use:   "view [id]",
		Short: "View details of a specific item (story or comment)",
		Long:  `Fetch and display the full details, including text and top-level comments, of a specific Hacker News item ID.`,
		Example: `  tech-news view 8863
  tech-news view 8863 --json`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			id, _ := strconv.Atoi(args[0])
			item, err := getItem(id)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			if jsonOutput {
				out, _ := json.MarshalIndent(item, "", "  ")
				fmt.Println(string(out))
				return
			}

			fmt.Printf("\n--- %s ---\n", strings.ToUpper(item.Type))
			fmt.Printf("%s\n", item.Title)
			if item.Text != "" {
				fmt.Printf("\n%s\n", stripHTML(item.Text))
			}
			if item.URL != "" {
				fmt.Printf("\nURL: %s\n", item.URL)
			}
			fmt.Printf("\nBy: %s | Score: %d\n", item.By, item.Score)

			if len(item.Kids) > 0 {
				fmt.Printf("\nComments (%d):\n", len(item.Kids))
				commentCount := 5
				if len(item.Kids) < commentCount {
					commentCount = len(item.Kids)
				}

				for i := 0; i < commentCount; i++ {
					comment, err := getItem(item.Kids[i])
					if err == nil && comment != nil {
						text := stripHTML(comment.Text)
						if len(text) > 200 {
							text = text[:200] + "..."
						}
						fmt.Printf(">> %s: %s\n", comment.By, text)
					}
				}
			}
		},
	}

	viewCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")

	rootCmd.AddCommand(listCmd, viewCmd)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
