package layercake_test

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/cloudfoundry-incubator/garden-linux/shed/layercake"
	"github.com/cloudfoundry-incubator/garden-linux/shed/layercake/fake_cake"
	"github.com/cloudfoundry/gunk/command_runner/fake_command_runner"
	. "github.com/cloudfoundry/gunk/command_runner/fake_command_runner/matchers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-golang/lager/lagertest"
)

var _ = Describe("BtrfsCleaningCake", func() {
	var (
		cleaner              *layercake.BtrfsCleaningCake
		runner               *fake_command_runner.FakeCommandRunner
		fakeCake             *fake_cake.FakeCake
		listSubvolumesOutput string
		layerId              = layercake.DockerImageID("the-layer")
		btrfsMountPoint      = "/absolute/btrfs_mount"

		listSubVolumeErr error
		graphDriverErr   error

		removedDirectories []string
	)

	BeforeEach(func() {
		graphDriverErr = nil
		listSubVolumeErr = nil
		removedDirectories = []string{}

		runner = fake_command_runner.New()
		fakeCake = new(fake_cake.FakeCake)
		cleaner = &layercake.BtrfsCleaningCake{
			Cake:            fakeCake,
			Runner:          runner,
			BtrfsMountPoint: btrfsMountPoint,
			RemoveAll: func(dir string) error {
				removedDirectories = append(removedDirectories, dir)
				return nil
			},
			Logger: lagertest.NewTestLogger("test"),
		}

		runner.WhenRunning(fake_command_runner.CommandSpec{
			Path: "btrfs",
			Args: []string{"subvolume", "list", btrfsMountPoint},
		}, func(cmd *exec.Cmd) error {
			_, err := cmd.Stdout.Write([]byte(listSubvolumesOutput))
			Expect(err).NotTo(HaveOccurred())
			return listSubVolumeErr
		})

		fakeCake.PathStub = func(id layercake.ID) (string, error) {
			return "/absolute/btrfs_mount/relative/path/to/" + id.GraphID(), graphDriverErr
		}
	})

	Context("when there are no subvolumes", func() {
		BeforeEach(func() {
			listSubvolumesOutput = "\n"
		})

		It("does not invoke subvolume delete", func() {
			Expect(cleaner.Remove(layerId)).To(Succeed())
			Expect(runner).NotTo(HaveExecutedSerially(fake_command_runner.CommandSpec{
				Path: "btrfs",
				Args: []string{"subvolume", "delete", "/path/to/" + layerId.GraphID()},
			}))
		})

		It("does not delete any directories", func() {
			Expect(cleaner.Remove(layerId)).To(Succeed())
			Expect(removedDirectories).To(BeEmpty())
		})
	})

	Context("when there is a subvolume for the layer, but it does not contain nested subvolumes", func() {
		BeforeEach(func() {
			listSubvolumesOutput = "ID 257 gen 9 top level 5 path relative/path/to/" + layerId.GraphID() + "\n"
		})

		It("does not invoke subvolume delete", func() {
			Expect(cleaner.Remove(layerId)).To(Succeed())
			Expect(runner).NotTo(HaveExecutedSerially(fake_command_runner.CommandSpec{
				Path: "btrfs",
				Args: []string{"subvolume", "delete", "/absolute/btrfs_mount/relative/path/to/" + layerId.GraphID()},
			}))
		})

		It("does not delete any directories", func() {
			Expect(cleaner.Remove(layerId)).To(Succeed())
			Expect(removedDirectories).To(BeEmpty())
		})
	})

	Context("when there is a subvolume for the layer, and it contains nested subvolumes", func() {
		subvolume1 := fmt.Sprintf("%s/relative/path/to/%s/subvolume1", btrfsMountPoint, layerId.GraphID())
		subvolume2 := fmt.Sprintf("%s/relative/path/to/%s/subvolume2", btrfsMountPoint, layerId.GraphID())

		BeforeEach(func() {
			listSubvolumesOutput = fmt.Sprintf(`ID 257 gen 9 top level 5 path relative/path/to/%s
ID 258 gen 9 top level 257 path relative/path/to/%s/subvolume1
ID 259 gen 9 top level 257 path relative/path/to/%s/subvolume2
`, layerId.GraphID(), layerId.GraphID(), layerId.GraphID())
		})

		It("deletes the subvolume", func() {
			Expect(cleaner.Remove(layerId)).To(Succeed())
			Expect(runner).To(HaveExecutedSerially(fake_command_runner.CommandSpec{
				Path: "btrfs",
				Args: []string{"subvolume", "delete", subvolume1},
			}))
			Expect(runner).To(HaveExecutedSerially(fake_command_runner.CommandSpec{
				Path: "btrfs",
				Args: []string{"subvolume", "delete", subvolume2},
			}))
		})

		It("deletes the subvolume directory contents before deleting the subvolume", func() {
			runner.WhenRunning(fake_command_runner.CommandSpec{
				Path: "btrfs",
				Args: []string{"subvolume", "delete", subvolume1},
			}, func(cmd *exec.Cmd) error {
				Expect(removedDirectories).To(ConsistOf(subvolume1))
				return nil
			})

			Expect(cleaner.Remove(layerId)).To(Succeed())
		})

		Context("and the nested subvolumes have nested subvolumes", func() {
			BeforeEach(func() {
				listSubvolumesOutput = fmt.Sprintf(`ID 257 gen 9 top level 5 path relative/path/to/%s
ID 258 gen 9 top level 257 path relative/path/to/%s/subvolume1
ID 259 gen 9 top level 257 path relative/path/to/%s/subvolume1/subsubvol1
`, layerId.GraphID(), layerId.GraphID(), layerId.GraphID())
			})

			It("deletes the subvolumes deepest-first", func() {
				Expect(cleaner.Remove(layerId)).To(Succeed())
				Expect(runner).To(HaveExecutedSerially(fake_command_runner.CommandSpec{
					Path: "btrfs",
					Args: []string{"subvolume", "delete", subvolume1 + "/subsubvol1"},
				}, fake_command_runner.CommandSpec{
					Path: "btrfs",
					Args: []string{"subvolume", "delete", subvolume1},
				}))
			})
		})
	})

	Context("when there is an associated qgroup", func() {
		BeforeEach(func() {
			runner.WhenRunning(fake_command_runner.CommandSpec{
				Path: "btrfs",
				Args: []string{"qgroup", "show", "-f", "/absolute/btrfs_mount/relative/path/to/" + layerId.GraphID()},
			}, func(cmd *exec.Cmd) error {
				_, err := cmd.Stdout.Write([]byte(`qgroupid rfer  excl
-------- ----  ----
0/5      49152 49152
`))
				Expect(err).NotTo(HaveOccurred())
				return listSubVolumeErr
			})
		})

		It("removes the associated qgroup", func() {
			Expect(cleaner.Remove(layerId)).To(Succeed())
			Expect(runner).To(HaveExecutedSerially(fake_command_runner.CommandSpec{
				Path: "btrfs",
				Args: []string{
					"qgroup", "destroy", "0/5", btrfsMountPoint,
				},
			}))
		})
	})

	Context("when running a command fails", func() {
		BeforeEach(func() {
			listSubVolumeErr = errors.New("listing subvolumes failed!")
		})

		It("returns the same error", func() {
			Expect(cleaner.Remove(layerId)).To(MatchError("listing subvolumes failed!"))
		})

		It("doesnt call the delegate", func() {
			cleaner.Remove(layerId)
			Expect(fakeCake.RemoveCallCount()).To(Equal(0))
		})
	})

	Context("when graphdriver fails to get layer path", func() {
		BeforeEach(func() {
			graphDriverErr = errors.New("graphdriver fail!")
		})

		It("returns the same error", func() {
			Expect(cleaner.Remove(layerId)).To(MatchError("graphdriver fail!"))
		})

		It("doesnt call the delegate", func() {
			Expect(fakeCake.RemoveCallCount()).To(Equal(0))
		})
	})

	It("delegates to the delegate", func() {
		Expect(cleaner.Remove(layerId)).To(Succeed())
		Expect(fakeCake.RemoveCallCount()).To(Equal(1))
	})

	Context("when the delegate fails", func() {
		BeforeEach(func() {
			fakeCake.RemoveReturns(errors.New("o no!"))
		})

		It("returns the same error", func() {
			Expect(cleaner.Remove(layerId)).To(MatchError("o no!"))
		})
	})
})
