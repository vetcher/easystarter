package commands

import (
	"flag"
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/kpango/glg"
	"github.com/vetcher/easystarter/util"
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
			return "less"
		case "darwin":
			return "less"
		default:
			return "open -a"
		}
	} else {
		return *logsViewerFlag
	}
}

type LogsCommand struct {
}

func (c *LogsCommand) Validate(args ...string) error {
	if len(args) > 0 {
		return nil
	}
	return AtLeastOneArgumentErr
}

func (c *LogsCommand) Exec(args ...string) error {
	cmd := exec.Command(logsViewer, filepath.Join(util.StartupDir(), "logs", fmt.Sprintf("%s.log", args[0])))
	err := cmd.Run()
	if err != nil {
		glg.Warnf("Open logs for %s error: %v", args[0], err)
	}
	return nil
}
