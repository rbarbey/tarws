package util

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

// S3KeyName generates a key name for a tar uploaded to S3. It's based on the
// provided path and looks like this
// `tarws-20160713-0047-home-jane-photos`
func S3KeyName(path string, date time.Time) string {
	formattedDate := date.Format("20060102-1504")
	absolutePath, _ := filepath.Abs(path)

	return fmt.Sprintf("tarws-%s%s", formattedDate, strings.Replace(absolutePath, "/", "-", -1))
}
