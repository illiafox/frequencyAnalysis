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

func generateBarItems(values []int) []opts.BarData {
	// Create slice
	items := make([]opts.BarData, len(values))
	for i := range values {
		// Set BarData with slice values
		items[i] = opts.BarData{Value: values[i]}
	}
	return items
}

func makeBar(values []int, xaxis []string) *charts.Bar {
	// Initialize Bar
	bar := charts.NewBar()
	// Set its options
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Frequency Analysis",
			Subtitle: "",
			Link:     "",
			Right:    "45%",
		}),

		charts.WithInitializationOpts(opts.Initialization{
			Width:  "110em",
			Height: "50em",
		}),

		charts.WithXAxisOpts(opts.XAxis{
			Name: "letters",
		}),

		charts.WithYAxisOpts(opts.YAxis{
			Name: "times occurred",
		}),
	)
	bar.SetXAxis(xaxis).AddSeries("Category A", generateBarItems(values))
	return bar
}

func drawBar(str []string, cnt []int, path string) error {
	// Create file
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	// Render page
	return components.NewPage().
		// Add bars
		AddCharts(makeBar(cnt, str)).
		// Write to file
		Render(io.MultiWriter(f))
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Not enough arguments! Usage:\n./app input.txt output.html")
		return
	}
	read, write := os.Args[1], os.Args[2]
	// Read input file
	data, err := ioutil.ReadFile(read)
	if err != nil {
		fmt.Println("Opening file: ", err)
		return
	}

	input := strings.ToLower(string(data))
	counter := make(map[rune]int)

	fmt.Print(log.Prefix(), "Counting... ")

	// Initialize time
	t := time.Now()

	for _, c := range input {
		if unicode.IsLetter(c) {
			// Increment in map
			counter[c]++
		}
	}

	// Create two slices, which are needed for drawBar in the end
	str := make([]string, len(counter))
	cnt := make([]int, len(counter))

	i := 0
	for k, v := range counter {
		// Set values in slices above
		str[i] = string(k)
		cnt[i] = v
		i++
	}

	//	Sort slices in descending order
	sort.Slice(cnt, func(i, j int) bool {
		if cnt[i] > cnt[j] {
			// Swap second slice
			str[i], str[j] = str[j], str[i]
			return true
		}
		return false
	})

	fmt.Println(time.Since(t).Seconds(), "s")
	fmt.Print(log.Prefix(), "Drawing... ")

	// Refresh time and draw bar
	t = time.Now()
	err = drawBar(str, cnt, write)

	fmt.Println(time.Since(t).Seconds(), "s")
	if err != nil {
		log.Println("[ERROR]", err)
	}
}
