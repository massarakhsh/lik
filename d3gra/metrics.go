package d3gra

import (
	"fmt"
	"time"

	"github.com/massarakhsh/lik"
	"github.com/massarakhsh/lik/likdom"
	"github.com/massarakhsh/lik/metric"
)

func BuildDiv(id string, sx, sy int, panel lik.Seter, series lik.Lister) likdom.Domer {
	panel.SetValue("sx", sx)
	panel.SetValue("sy", sy)
	panel.SetValue("series", series)

	div := likdom.BuildDiv("id", id, "width", sx, "hight", sy)
	script := div.BuildItem("script")
	script.BuildString(fmt.Sprintf("let options = %s;\n", panel.SerializeJavascript()))
	script.BuildString(fmt.Sprintf("draw_charts('%s', options);\n", id))
	return div
}

func BuildPanel(distance time.Duration) lik.Seter {
	panel := lik.BuildSet()
	to := time.Now()
	from := to.Add(-distance)
	panel.SetValues("to", to.UnixMilli(), "from", from.UnixMilli())
	return panel
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
