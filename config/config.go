package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/etsy/mixer/debug"

	"github.com/etsy/mixer/Godeps/_workspace/src/gopkg.in/gcfg.v1"
)

var Config = &configuration{}
var dbg debug.Debug

const debugon = true

func (cfg *configuration) Load() error {
	dbg = debugon
	fdir := filepath.Dir(pathToThisFile())
	cdir := cfg.GetRootDir()

	config := []string{fmt.Sprintf("%s/../config.cfg", fdir), fmt.Sprintf("%s/config.cfg", cdir), "/etc/mixer/config.cfg"}
	paths := make([]string, 0)
	paths = findConfig(config, paths)

	secrets_config := []string{fmt.Sprintf("%s/../config-secrets.cfg", fdir), "/etc/mixer/config-secrets.cfg"}
	paths = findConfig(secrets_config, paths)

	cfg.mtx.Lock()
	defer cfg.mtx.Unlock()

	for _, p := range paths {
		log.Printf("reading config: %s\n", p)
		err := gcfg.ReadFileInto(cfg, p)
		if err != nil {
			return err
		}
	}
	return nil
}

func pathToThisFile() string {
	_, file, _, _ := runtime.Caller(0)
	return file
}

func (cfg *configuration) GetRootDir() string {
	wdir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// hack for when working dir is subdir (such as when running tests)
	root_location := "."
	if filepath.Base(wdir) != "mixer" {
		root_location = ".."
	}

	return fmt.Sprintf("%s/%s", wdir, root_location)
}

func findConfig(locations []string, paths []string) []string {
	for _, c := range locations {
		dbg.Printf("looking for config: %s\n", c)
		if _, err := os.Stat(c); err == nil {
			dbg.Printf(" config exists: %s\n", c)
			paths = append(paths, c)
			return paths
		}
	}
	return paths
}

type configuration struct {
	mtx sync.Mutex

	Server struct {
		Port int
		Url  string
	}

	Userauth struct {
		Header string
	}

	Mail struct {
		Key           string
		EmailAddress  string `gcfg:"email-address"`
		SMTP_Host     string `gcfg:"smtp-host"`
		Port          int
		Domain        string
		AdminUsername string `gcfg:"admin-username"`
	}

	Database map[string]*struct {
		User     string
		Password string
		Hostname string
		Port     int
		Name     string
	}

	Staff struct {
		DatafeedUrl      []string `gcfg:"datafeed-url"`
		DirectoryUrl     string   `gcfg:"directory-url"`
		DefaultAvatarUrl string   `gcfg:"default-avatar-url"`
	}
}
