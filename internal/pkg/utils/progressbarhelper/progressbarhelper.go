package progressbarhelper

import "github.com/schollz/progressbar/v3"

func BarDescribe(bar *progressbar.ProgressBar, description string) {
	if bar != nil {
		bar.Describe(description)
	}
}
func BarClear(bar *progressbar.ProgressBar) {
	if bar != nil {
		_ = bar.Clear()
	}
}
