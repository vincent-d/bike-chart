package charts

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"

	"github.com/vincent-d/bike-chart/pkg/bikecount"
	"github.com/vincent-d/bike-chart/pkg/ecovisio"
)

func getMapLink(latitude float32, longitude float32) (string, error) {
	base, err := url.Parse("https://www.google.com/maps/search/")
	if err != nil {
		return "", err
	}
	params := url.Values{}
	params.Add("api", "1")
	params.Add("query", fmt.Sprintf("%f,%f", latitude, longitude))
	base.RawQuery = params.Encode()

	return base.String(), nil
}

func CreateCharts(bikeCountData []bikecount.BikeCountData) *components.Page {

	page := components.NewPage()
	page.Initialization.PageTitle = "Bike Count Data"
	for _, countData := range bikeCountData {

		link, _ := getMapLink(countData.Latitude, countData.Longitude)
		serieName := fmt.Sprintf("Bike Count per day\n\tMax: %d on %s, avg: %d\n\tWeek days avg: %d\n\tBest day in avg: %s (%d)\n\tWeekends avg: %d",
			countData.WindowMax, countData.MaxDate.Format("02/01/2006"), countData.OverallAverage,
			countData.WeekdaysAvg,
			countData.BestDay.String(), countData.DaysAverage[countData.BestDay],
			countData.WeekendsAvg)
		//log.Printf("Bike Count per day in %s (max: %d on %s, avg: %d, week days avg: %d, Mon avg: %d, Tue avg: %d, Wed avg: %d, Thu avg: %d, Fri avg: %d, Sat avg: %d, Sun avg: %d)", totem.Name,
		//	max, maxDate.Format("02/01/2006"), totem.DailyAverage, int(weekDays.compute()), int(days[time.Monday].compute()), int(days[time.Tuesday].compute()), int(days[time.Wednesday].compute()),
		//	int(days[time.Thursday].compute()), int(days[time.Friday].compute()), int(days[time.Saturday].compute()), int(days[time.Sunday].compute()))
		// create a new bar instance
		chart := charts.NewBar()
		// set some global options like Title/Legend/ToolTip or anything else
		chart.SetGlobalOptions(
			charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeRoma}),
			charts.WithTitleOpts(opts.Title{
				Title:    countData.Name,
				Link:     link,
				Subtitle: "Bike Counts",
				Top:      "top",
				Left:     "15%",
			}),
			charts.WithTooltipOpts(opts.Tooltip{
				Show: true,
			}),
			charts.WithLegendOpts(opts.Legend{
				Show:  true,
				Right: "10%",
			}),
			charts.WithYAxisOpts(opts.YAxis{
				Name: "Number of bikes",
			}),
			charts.WithXAxisOpts(opts.XAxis{
				Name: "Date",
			}),
			charts.WithDataZoomOpts(opts.DataZoom{
				Start:      float32(countData.Start),
				End:        float32(countData.End),
				XAxisIndex: []int{0},
			}),
		)

		// Put data into instance
		chart.SetXAxis(countData.DisplayDates).
			AddSeries(serieName, countData.Counts)
		page.AddCharts(chart)
	}

	return page
}

func GetChartPage(filter string, start time.Time, end time.Time) (*components.Page, error) {

	totems, err := ecovisio.GetTotemList()
	if err != nil {
		log.Printf("Error when getting totem info: %v", err)
		return nil, err
	}

	if end.IsZero() {
		end = time.Now()
	}
	if start.IsZero() {
		start = end.AddDate(-1, 0, 0)
	}
	log.Printf("Getting counts for totems with name: %s, between %v and %v", filter, start, end)
	totems, err = ecovisio.FindTotems(totems, filter)
	if err != nil {
		return nil, err
	}
	if totems == nil {
		return nil, nil
	}

	//log.Printf("%#v", *totem)
	bikes, err := ecovisio.GetBikesCountFromTotemInfo(totems)
	if err != nil {
		log.Printf("Error when getting count: %v", err)
		return nil, err
	}

	return CreateCharts(bikecount.PrepareBikeCountData(totems, bikes, start, end)), nil
}
