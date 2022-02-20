// ILLIAFOX 2022lf

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// Generate Bar Items
func generateBarItems(values []int) []opts.BarData {
	// Create Slie
	items := make([]opts.BarData, len(values))
	// Range values
	for i := range values {
		// Set BarData with vaule
		items[i] = opts.BarData{Value: values[i]}
	}
	return items
}

// Create Bar
func makeBar(values []int, xaxis []string) *charts.Bar {
	// Initialize Bar
	bar := charts.NewBar()
	// Set Options
	bar.SetGlobalOptions(
		// Title
		charts.WithTitleOpts(opts.Title{
			Title:    "Frequency Analysis",
			Subtitle: "",
			Link:     "",
			Right:    "45%",
		}),
		// Size
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "110em",
			Height: "50em",
		}),
		// X Axis
		charts.WithXAxisOpts(opts.XAxis{
			Name: "letters",
		}),
		// Y Axis
		charts.WithYAxisOpts(opts.YAxis{
			Name: "times occured",
		}),
	)
	// Set X
	bar.SetXAxis(xaxis).AddSeries("Category A", generateBarItems(values))
	return bar
}

// Draw bar (write to file)
func drawBar(str []string, cnt []int, path string) error {
	// Create file
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	// Render page and return error
	return components.NewPage().
		AddCharts(makeBar(cnt, str)).
		Render(io.MultiWriter(f))
}

func main() {
	// Check arguments
	if len(os.Args) != 3 {
		fmt.Println("Not enough arguments! Usage:\n./app input.txt output.html")
		return
	}
	// Read input file
	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("Opening file: ", err)
		return
	}
	// Lower all chars
	input := strings.ToLower(string(data))
	// Create map
	counter := make(map[rune]int, len(input))
	// Print timer
	fmt.Print(log.Prefix(), "Counting... ")
	// Initialize time
	t := time.Now()
	// Range input
	for _, c := range input {
		// If Letter
		if unicode.IsLetter(c) {
			// Increment in map
			counter[c]++
		}
	}
	// Create two slices
	str := make([]string, 0, len(counter))
	cnt := make([]int, 0, len(counter))
	// Range map
	for k, v := range counter {
		// Append to slices above
		str = append(str, string(k))
		cnt = append(cnt, v)
	}
	// Sort slices in descending order
	sort.Slice(cnt, func(i, j int) bool {
		if cnt[i] > cnt[j] {
			// Swap other slice
			str[i], str[j] = str[j], str[i]
			return true
		}
		return false
	})
	// Print timers
	fmt.Println(time.Since(t).Seconds(), "s")
	fmt.Print(log.Prefix(), "Drawing... ")
	// Reinitialize time
	t = time.Now()
	// draw bar
	err = drawBar(str, cnt, os.Args[2])
	fmt.Println(time.Since(t).Seconds(), "s")
	// Check error
	if err != nil {
		log.Println("[ERROR]", err)
	}

}
