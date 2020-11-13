package likapi

import (
	"fmt"
	"strings"
)

func SetToPart(drive DataDriver, to DataDriver, part string) int {
	id := drive.GetPage().ContinueToPage(to.GetPage())
	url := BuildUrl(id, part)
	drive.SetResponse(url, "_topart")
	return id
}

func BuildUrl(id int, part string) string {
	url := part
	if url == "" {
		url = "/"
	}
	if strings.Index(url, "?_sp=") < 0 && strings.Index(url, "&_sp=") < 0 {
		if strings.Index(url, "?") < 0 {
			url += "?"
		} else {
			url += "&"
		}
		url += "_sp=" + fmt.Sprint(id)
	}
	return url
}
