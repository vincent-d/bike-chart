package main

import (
	"os"
	"time"

	"github.com/vincent-d/bike-chart/pkg/charts"
)

func main() {
	page, _ := charts.GetChartPage("", time.Time{}, time.Time{})
	f, _ := os.Create("bar.html")
	page.Render(f)
}
