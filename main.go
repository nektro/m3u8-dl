package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/nektro/go-util/mbpp"
)

func main() {
	fM3u := flag.String("url", "", "URL to .m3u8 file to download")
	fDir := flag.String("dir", "./data", "Path to directory to save .ts files to")
	fCon := flag.Int("concurrency", 10, "Number of concurrent downloads to run")
	flag.Parse()

	m3u := *fM3u
	if len(m3u) > 0 {
		bys2, _ := fetchBin(m3u, nil)
		tss := parseM3U8(string(bys2))

		dir := *fDir
		dir, _ = filepath.Abs(dir)
		os.MkdirAll(dir, os.ModePerm)

		sdir := filepath.Dir(m3u)
		sdir = strings.Replace(sdir, ":/", "://", 1)

		mbpp.Init(*fCon)

		mbpp.CreateJob("video: "+m3u, func(bar *mbpp.BarProxy) {
			for _, item := range tss {
				bar.AddToTotal(1)
				go mbpp.CreateDownloadJob(sdir+"/"+item, dir+"/"+item, bar)
			}
		})

		mbpp.Wait()
		time.Sleep(time.Second)
		log.Println(mbpp.GetCompletionMessage())
	}
}

func fetchBin(urlS string, headers map[string]string) ([]byte, error) {
	req, _ := http.NewRequest(http.MethodGet, urlS, nil)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	res, _ := http.DefaultClient.Do(req)
	bys, _ := ioutil.ReadAll(res.Body)
	return bys, nil
}

func parseM3U8(in string) []string {
	r := []string{}
	lines := strings.Split(in, "\n")
	for _, l := range lines {
		if len(l) == 0 || l[0] == '#' {
			continue
		}
		r = append(r, l)
	}
	return r
}
