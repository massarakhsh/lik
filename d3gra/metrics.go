package d3gra

import (
	"fmt"
	"time"

	"github.com/massarakhsh/lik"
	"github.com/massarakhsh/lik/likdom"
	"github.com/massarakhsh/lik/metric"
)

func BuildPanel(id string, sx, sy int, distance time.Duration, series lik.Lister) likdom.Domer {
	options := lik.BuildSet("width", sx, "height", sy)
	to := time.Now()
	from := to.Add(-distance)
	options.SetValues("to", to.UnixMilli(), "from", from.UnixMilli())
	options.SetValues("series", series)

	div := likdom.BuildDiv("id", id, "width", sx, "hight", sy)
	script := div.BuildItem("script")
	script.BuildString(fmt.Sprintf("let options = %s;\n", options.SerializeJavascript()))
	script.BuildString(fmt.Sprintf("draw_charts('%s', options);\n", id))
	return div
}

func BuildSeria(metro *metric.MetricValue, to time.Time, step time.Duration, count int) lik.Seter {
	seria := lik.BuildSet()
	data := seria.AddList("data")
	history := metro.GetListPath(to, step, count)
	for n := 0; n < count; n++ {
		ni := count - 1 - n
		at := to.Add(-time.Duration(ni) * step)
		if ni >= 0 && ni < len(history) {
			elm := data.AddItemSet()
			elm.SetValue("date", at.UnixMilli())
			elm.SetValue("value", history[ni])
		}
	}
	return seria
}
