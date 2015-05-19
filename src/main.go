package main

import (
	"fmt"
	"qlog"
)

func main() {
	l := qlog.QLog{}
	l.Parse(`14.159.100.74 - - [02/Apr/2015:23:59:59 +0800] "GET /@/sports/1.0.7YDXJv22_1.0.7_build-20150330113749_b690_i446_s699.bin HTTP/1.1" 403 935 "-" "AndroidDownloadManager/4.4.4 (Linux; U; Android 4.4.4; MI 4LTE Build/KTU84P)" "http://yi-version.qiniudn.com" V1`)
	fmt.Println(l)
}
