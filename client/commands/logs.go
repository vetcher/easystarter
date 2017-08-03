package commands

import (
	"flag"
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/kpango/glg"
	"github.com/vetcher/easystarter/backend/util"
)

var (
	logsViewerFlag = flag.String("open-logs", "", "Show log file in this application")
	logsViewer     = chooseDefaultLogsViewer()
)

func chooseDefaultLogsViewer() string {
	if *logsViewerFlag == "" {
		switch runtime.GOOS {
		case "windows":
			return "notepad"
		case "linux":
			return "gedit"
		case "darwin":
			return "TextEditor"
		default:
			return "open -a"
		}
	} else {
		return *logsViewerFlag
	}
}

type LogsCommand struct {
	svcName string
}

func (c *LogsCommand) Validate(args ...string) error {
	if len(args) > 0 {
		c.svcName = args[0]
		return nil
	}
	return AtLeastOneArgumentErr
}

func (c *LogsCommand) Exec() error {
	go func() {
		cmd := exec.Command(logsViewer, filepath.Join(util.StartupDir(), "logs", fmt.Sprintf("%s.log", c.svcName)))
		err := cmd.Run()
		if err != nil {
			glg.Warnf("Open logs for %s error: %v", c.svcName, err)
		}
	}()
	return nil
}
