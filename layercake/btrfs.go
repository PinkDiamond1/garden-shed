package layercake

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/cloudfoundry-incubator/garden-linux/logging"
	"github.com/cloudfoundry/gunk/command_runner"
	"github.com/pivotal-golang/lager"
)

type BtrfsCleaningCake struct {
	Cake

	Runner          command_runner.CommandRunner
	BtrfsMountPoint string
	RemoveAll       func(dir string) error

	Logger lager.Logger
}

func (c *BtrfsCleaningCake) Remove(id ID) error {
	log := c.Logger.Session("clean-rootfs", lager.Data{"id": id})
	log.Info("start")

	layerPath, err := c.Cake.Path(id)
	if err != nil {
		log.Error("get", err)
		return err
	}

	if err := c.removeSubvols(log, layerPath); err != nil {
		return err
	}

	if err := c.removeQgroup(log, layerPath); err != nil {
		log.Error("remove-qgroup", err)
	}

	defer log.Info("complete")
	return c.Cake.Remove(id)
}

func (c *BtrfsCleaningCake) removeQgroup(log lager.Logger, layerPath string) error {
	log = log.Session("remove-qgroup")
	log.Info("start")

	runner := &logging.Runner{c.Runner, log}

	qgroupInfo, err := c.run(runner, exec.Command("btrfs", "qgroup", "show", "-f", layerPath))
	if err != nil {
		return err
	}

	qgroupInfoLines := strings.Split(qgroupInfo, "\n")
	if len(qgroupInfoLines) != 4 {
		return fmt.Errorf("unexpected qgroup show output: %s", qgroupInfo)
	}

	qgroupid := strings.SplitN(qgroupInfoLines[2], " ", 2)[0]
	_, err = c.run(runner, exec.Command("btrfs", "qgroup", "destroy", qgroupid, c.BtrfsMountPoint))

	if err != nil {
		log.Error("failed", err)
	}

	log.Info("destroyed", lager.Data{"qgroupid": qgroupid})
	return nil
}

func (c *BtrfsCleaningCake) removeSubvols(log lager.Logger, layerPath string) error {
	runner := &logging.Runner{c.Runner, log}

	listSubvolumesOutput, err := c.run(runner, exec.Command("btrfs", "subvolume", "list", c.BtrfsMountPoint))
	if err != nil {
		return err
	}

	subvols := finalColumns(strings.Split(listSubvolumesOutput, "\n"))
	sort.Sort(deepestFirst(subvols))

	for _, subvolume := range subvols {
		subvolumeAbsPath := filepath.Join(c.BtrfsMountPoint, subvolume)

		if strings.Index(subvolumeAbsPath, layerPath) == 0 && subvolumeAbsPath != layerPath {
			log.Info("removing-subvol", lager.Data{"layerPath": layerPath, "subvolumeAbsPath": subvolumeAbsPath})

			c.RemoveAll(subvolumeAbsPath)

			if _, err := c.run(runner, exec.Command("btrfs", "subvolume", "delete", subvolumeAbsPath)); err != nil {
				return err
			}
		}
	}

	return nil
}

func finalColumns(lines []string) []string {
	result := make([]string, 0)
	for _, line := range lines {
		cols := strings.Fields(line)
		if len(cols) == 0 {
			continue
		}

		result = append(result, cols[len(cols)-1])
	}

	return result
}

func (c *BtrfsCleaningCake) run(runner command_runner.CommandRunner, cmd *exec.Cmd) (string, error) {
	var buffer bytes.Buffer
	cmd.Stdout = &buffer
	if err := runner.Run(cmd); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

type deepestFirst []string

func (a deepestFirst) Len() int           { return len(a) }
func (a deepestFirst) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a deepestFirst) Less(i, j int) bool { return len(a[i]) > len(a[j]) }
