package ui

import (
	"fmt"
	"os"
	"strconv"

	"github.com/manifoldco/promptui"
)

type TrackResult struct {
	Title    string
	Artist   string
	Album    string
	Year     string
	Bitrate  string
	Duration string
}

type PaginatedDisplay struct {
	Theme Theme
}

func NewPaginatedDisplay(theme Theme) *PaginatedDisplay {
	return &PaginatedDisplay{Theme: theme}
}

// RenderPaginatedResults shows 5 tracks at a time with 'Next 5' / 'Prev 5' controls
func (pd *PaginatedDisplay) RenderPaginatedResults(artist string, results []TrackResult) (*TrackResult, error) {
	pageSize := 5
	currentPage := 0
	totalItems := len(results)

	for {
		start := currentPage * pageSize
		end := start + pageSize
		if end > totalItems {
			end = totalItems
		}

		pageItems := results[start:end]

		// Render Table for current 5 items
		table := NewTable(pd.Theme,
			Column{Title: "#", Width: 3},
			Column{Title: "Title", Width: 25},
			Column{Title: "Album", Width: 20},
			Column{Title: "Year", Width: 6},
			Column{Title: "Bitrate", Width: 8},
		)

		for i, item := range pageItems {
			globalIdx := start + i + 1
			table.AddRow(
				strconv.Itoa(globalIdx),
				item.Title,
				item.Album,
				item.Year,
				item.Bitrate,
			)
		}

		// Clear screen effect & show banner
		fmt.Printf("\033[H\033[2J")
		fmt.Println(TrueColor(fmt.Sprintf("🎤 Search Results for '%s' (Showing %d-%d of %d)", artist, start+1, end, totalItems), pd.Theme.Primary))
		table.Render(os.Stdout)

		// Build Prompt Choices
		var menuOptions []string
		for i, item := range pageItems {
			globalIdx := start + i + 1
			menuOptions = append(menuOptions, fmt.Sprintf("[%d] Download: %s - %s", globalIdx, item.Title, item.Album))
		}

		if end < totalItems {
			menuOptions = append(menuOptions, "➡️  Next 5 Results...")
		}
		if currentPage > 0 {
			menuOptions = append(menuOptions, "⬅️  Previous 5 Results...")
		}
		menuOptions = append(menuOptions, "❌ Cancel Search")

		prompt := promptui.Select{
			Label: "Select Action",
			Items: menuOptions,
			Size:  10,
		}

		idx, selectedOpt, err := prompt.Run()
		if err != nil {
			return nil, err
		}

		if selectedOpt == "➡️  Next 5 Results..." {
			currentPage++
			continue
		} else if selectedOpt == "⬅️  Previous 5 Results..." {
			currentPage--
			continue
		} else if selectedOpt == "❌ Cancel Search" {
			return nil, fmt.Errorf("search canceled by user")
		}

		// Selected a track from the page
		selectedTrack := pageItems[idx]
		return &selectedTrack, nil
	}
}
