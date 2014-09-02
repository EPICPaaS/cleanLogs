// cleanLogs project main.go
package main

import (
	"flag"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

type FileRecord struct {
	FullPath    string
	ModuleName  string
	Port        string
	FetchedTime string
}

type FRList []FileRecord

func (lst FRList) Len() int {
	return len(lst)
}

func (lst FRList) Less(i, j int) bool {
	return lst[i].FetchedTime > lst[j].FetchedTime // 按值排序在
}

func (lst FRList) Swap(i, j int) {
	lst[i], lst[j] = lst[j], lst[i]
}

func GetFetchRecord(filePath string) *FileRecord {
	filename := path.Base(filePath)
	fns := strings.Split(filename, ".")
	flevels := strings.Split(filePath, "/")
	size := len(flevels)
	if size < 3 || len(fns) < 2 {
		return nil
	}
	moduleName := flevels[size-2]
	//appID := flevels[size-3]

	ns := fns[0]
	i := strings.LastIndex(ns, "_")
	port := ns[i+1:]
	fetchedTime := fns[len(fns)-1]
	return &FileRecord{filePath, moduleName, port, fetchedTime}

}

func GetFilelist(path1 string) []string {
	var all []string
	filepath.Walk(path1, func(path2 string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		baseName := path.Base(path2)
		if strings.Contains(baseName, ".fetched.") {
			all = append(all, path2)
		}
		return nil
	})

	return all
}

func main() {
	//ddir := "/home/paas/paas/logs"
	ddir := "/home/yourchanges/paas_home_dev/logs"
	ddir = "/home/yourchanges/paas_home_dev/logs/services/todone-0.0.2/todone_srv"
	dnum := 3

	dir := flag.String("dir", "", "the paas logs dir")
	num := flag.Int("num", -1, "the number to keep")
	flag.Parse()

	if *dir != "" {
		ddir = *dir
	}

	if *num != -1 {
		dnum = *num
	}

	list := GetFilelist(ddir)
	allMap := make(map[string]FRList, 200)

	for _, f := range list {
		r := GetFetchRecord(f)
		key := r.ModuleName + ":::" + r.Port
		if v, ok := allMap[key]; ok {
			//有，就追加
			v = append(v, *r)
			allMap[key] = v
		} else {
			//没有，就新建
			allMap[key] = FRList{*r}
		}

	}
	for k, v := range allMap {
		log.Println("handle " + k)
		sort.Sort(v)
		i := 0
		for _, r := range v {
			log.Println(r)
			i = i + 1
			if i > dnum {
				log.Println("delete " + r.FullPath)
				os.Remove(r.FullPath)
			}
		}
	}

	log.Println(dnum)
}
