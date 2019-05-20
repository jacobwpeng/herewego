package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/golang/glog"
)

const (
	noWallProxySitesFilename = "nowallproxy.sites"
	directSitesFilename      = "direct.sites"
)

type SiteSet map[string]bool

func unescapeProxy(proxy string) string {
	return strings.Replace(proxy, "_", " ", -1)
}

func readSitesFromFile(filePath string) SiteSet {
	sites := make(SiteSet)
	file, err := os.Open(filePath)
	if err != nil {
		glog.Warningf("Open %s error: %s", filePath, err)
		return sites
	}

	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		glog.Warningf("ReadAll %s error: %s", filePath, err)
		return sites
	}

	scanner := bufio.NewScanner(bytes.NewReader(fileContent))
	for scanner.Scan() {
		sites[scanner.Text()] = true
	}

	if err := scanner.Err(); err != nil {
		glog.Warningf("Scan %s error: %s", filePath, err)
	}
	return sites
}

func backupSitesToFile(filePath string, sites SiteSet) bool {
	flags := os.O_WRONLY | os.O_TRUNC | os.O_CREATE
	var perm os.FileMode = 0666
	file, err := os.OpenFile(filePath, flags, perm)
	if err != nil {
		glog.Warningf("OpenFile %s error: %s", filePath, err)
		return false
	}
	defer file.Close()
	w := bufio.NewWriter(file)
	for key, _ := range sites {
		fmt.Fprintf(w, "%s\n", key)
	}
	err = w.Flush()
	if err != nil {
		glog.Warningf("Flush error: %s", err)
	}
	return true
}

type State struct {
	templateFile         string
	listeningAddr        string
	noWallProxySitesFile string
	directSitesFile      string

	m                sync.Mutex
	backupChan       chan struct{}
	noWallProxySites SiteSet
	directSites      SiteSet
}

type App struct {
	state *State
}

func NewApp(templateFile string,
	savedSitesDir string,
	listeningAddr string) *App {

	noWallProxySitesFile := path.Join(savedSitesDir, noWallProxySitesFilename)
	directSitesFile := path.Join(savedSitesDir, directSitesFilename)

	state := &State{
		templateFile:         templateFile,
		listeningAddr:        listeningAddr,
		noWallProxySitesFile: noWallProxySitesFile,
		directSitesFile:      directSitesFile,

		backupChan:       make(chan struct{}, 16),
		noWallProxySites: readSitesFromFile(noWallProxySitesFile),
		directSites:      readSitesFromFile(directSitesFile),
	}

	return &App{state: state}
}

func (app *App) Run() {
	http.HandleFunc("/report", app.HandleReport)
	http.HandleFunc("/generate", app.HandleGenerate)
	http.HandleFunc("/backup", app.HandleBackup)
	go app.ScheduleBackupToFile()
	glog.Fatal(http.ListenAndServe(app.state.listeningAddr, nil))
}

func (app *App) ScheduleBackupToFile() {
	tickerChan := time.NewTicker(time.Second * 60).C
	for {
		select {
		case <-app.state.backupChan:
			app.DoBackup()
		case <-tickerChan:
			app.DoBackup()
		}
	}
}

func (app *App) DoBackup() {
	app.state.m.Lock()
	defer app.state.m.Unlock()
	backupSitesToFile(app.state.noWallProxySitesFile, app.state.noWallProxySites)
	backupSitesToFile(app.state.directSitesFile, app.state.directSites)
}

func (app *App) HandleReport(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	siteType := r.Form.Get("type")
	site := r.Form.Get("site")

	app.state.m.Lock()
	defer app.state.m.Unlock()
	var buffer bytes.Buffer
	if siteType == "nowallproxy" {
		delete(app.state.directSites, site)
		app.state.noWallProxySites[site] = true
	} else if siteType == "direct" {
		delete(app.state.noWallProxySites, site)
		app.state.directSites[site] = true
	} else if siteType == "wallproxy" {
		delete(app.state.directSites, site)
		delete(app.state.noWallProxySites, site)
	} else {
		fmt.Fprintf(&buffer, "Invalid site type: %s", siteType)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(buffer.Bytes())
		return
	}

	glog.Infof("Add %s to type %s", site, siteType)
	fmt.Fprintf(&buffer, "Add %s succeed, type %s", site, siteType)
	w.Write(buffer.Bytes())
}

func (app *App) HandleGenerate(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	nowallProxy := unescapeProxy(r.Form.Get("nowallproxy"))
	wallProxy := unescapeProxy(r.Form.Get("wallproxy"))
	direct := unescapeProxy(r.Form.Get("direct"))

	tmpl := template.Must(template.ParseFiles(app.state.templateFile))

	app.state.m.Lock()
	defer app.state.m.Unlock()

	type GenerateData struct {
		WallProxy   string
		NoWallProxy string
		Direct      string
		NoWallSites SiteSet
		DirectSites SiteSet
	}

	data := GenerateData{
		WallProxy:   wallProxy,
		NoWallProxy: nowallProxy,
		Direct:      direct,
		NoWallSites: app.state.noWallProxySites,
		DirectSites: app.state.directSites,
	}

	tmpl.Execute(w, data)
}

func (app *App) HandleBackup(w http.ResponseWriter, r *http.Request) {
	app.state.backupChan <- struct{}{}
	w.Write([]byte("Send backup request succeed"))
}

func main() {
	templateFile := flag.String("template", "NO_DEFAULT", "Template file")
	savedSitesDir := flag.String("dir", "NO_DEFAULT", "Saved sites directory")
	listeningAddr := flag.String("address", ":12345", "Listening address")
	flag.Parse()
	app := NewApp(*templateFile, *savedSitesDir, *listeningAddr)
	app.Run()
}
