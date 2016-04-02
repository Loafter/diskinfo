package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"net/http"
	"encoding/json"
	"log"
	"os"
	"flag"
	"diskinfo/cloudstat"
)



func PrepareDiskInfo() []cloudstat.DiskInfo {
	if out, err := exec.Command("wmic", "logicaldisk", "get", "size,freespace,caption").Output(); err != nil {
		fmt.Printf("error: %s\n", out)
	} else {
		res := bytes.Split(out, []byte{10})
		di := make([]cloudstat.DiskInfo, 0)
		for _, e := range res[1:] {
			fi := strings.Fields(string(e))
			if !(len(fi) < 3) {
				free, err := strconv.ParseUint(fi[1], 10, 64)
				if err != nil {
					return nil
				}
				total, err := strconv.ParseUint(fi[2], 10, 64)
				if err != nil {
					return nil
				}
				di = append(di, cloudstat.DiskInfo{Name: fi[0], Free: free, Total: total})
			}
		}
		return di

	}
	return nil
}
func SendStatistic(dat cloudstat.HealthData,url string){
	buf,err:=json.Marshal(dat)
	if err!=nil{
		return
	}
	read:=bytes.NewReader(buf)
	req, err := http.NewRequest("POST", url, read)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	resp.Body.Close()

}
func main() {
	statserv := flag.String("statserv", "http://cloudstat.run.aws-usw02-pr.ice.predix.io/sendstat", "a string")
	flag.Parse()
	for{
		var hd cloudstat.HealthData
		hd.ServerName=os.Getenv("ServerCloudId")
		hd.DisksInfo=PrepareDiskInfo()
		SendStatistic(hd,*statserv)
		time.Sleep(time.Second)
	}



}
