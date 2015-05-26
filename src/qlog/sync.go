package qlog

import (
	"errors"
	"fmt"
	"github.com/qiniu/api/auth/digest"
	"github.com/qiniu/api/rsf"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type QLogSync struct {
	Bucket              string
	SaveBucket          string
	SaveBucketDomain    string
	IsSaveBucketPrivate bool
	Mac                 *digest.Mac
}

//date in format YY-MM-DD
func (this *QLogSync) Sync(date string) (paths []string, err error) {
	prefix := fmt.Sprintf("_log/%s/%s", this.Bucket, date)
	//get log list of the date from qiniu
	client := rsf.New(this.Mac)
	//this app is designed for middle-scale access log mode, items less than 1000
	entries, _, lerr := client.ListPrefix(nil, this.SaveBucket, prefix, "", 1000)
	if lerr != nil && lerr != io.EOF {
		err = errors.New(fmt.Sprintf("error list bucket of the logs, %s", lerr.Error()))
		return
	}
	paths = make([]string, 0)
	for _, entry := range entries {
		//create link and download file
		logDnLink := fmt.Sprintf("%s/%s", this.SaveBucketDomain, entry.Key)
		if this.IsSaveBucketPrivate {
			logDnLink = this.createPrivateDownloadLink(logDnLink)
		}
		//download and save to local file
		localFpath := entry.Key
		dErr := this.downloadFileToLocal(logDnLink, localFpath)
		if dErr != nil {
			err = errors.New(fmt.Sprintf("error downloading log file %s to %s, due to %s",
				logDnLink, localFpath, err.Error()))
			return
		}
		paths = append(paths, localFpath)
	}
	return
}

func (this *QLogSync) createPrivateDownloadLink(logPublicLink string) string {
	linkToSign := fmt.Sprintf("%s?e=%d", logPublicLink, time.Now().Add(time.Hour*24).Unix())
	token := digest.Sign(this.Mac, []byte(linkToSign))
	return fmt.Sprintf("%s&token=%s", linkToSign, token)
}

func (this *QLogSync) downloadFileToLocal(dnLink string, localFpath string) error {
	//check duplicate
	_, sErr := os.Stat(localFpath)
	if sErr == nil {
		return nil
	}
	//download
	fname := filepath.Base(localFpath)
	localFdir := strings.TrimSuffix(localFpath, fname)
	mErr := os.MkdirAll(localFdir, 0775)
	if mErr != nil {
		return errors.New(fmt.Sprintf("mkdir all %s error, due to", localFdir, mErr.Error()))
	}
	//req
	resp, respErr := http.Get(dnLink)
	if respErr != nil {
		return errors.New(fmt.Sprintf("download log error, %s", respErr.Error()))
	}
	defer resp.Body.Close()

	localFp, openErr := os.OpenFile(localFpath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0655)
	if openErr != nil {
		return errors.New(fmt.Sprintf("open local file error, %s", openErr.Error()))
	}
	defer localFp.Close()

	_, cpErr := io.Copy(localFp, resp.Body)
	if cpErr != nil {
		return errors.New(fmt.Sprintf("copy log data to local error, %s", cpErr.Error()))
	}
	return nil
}
