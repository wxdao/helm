/*
Copyright The Helm Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package diffutil

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

type Differ struct {
	tmpdir string
}

// NewDiffer create a differ.
func NewDiffer() (*Differ, error) {
	// create base directory
	tmpdir, err := os.MkdirTemp("", "diff-*")
	if err != nil {
		return nil, err
	}

	// create old & new directories
	if err := os.MkdirAll(filepath.Join(tmpdir, "old"), 0700); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Join(tmpdir, "new"), 0700); err != nil {
		return nil, err
	}

	return &Differ{
		tmpdir: tmpdir,
	}, nil
}

// Cleanup deletes temporary files.
func (differ *Differ) Cleanup() error {
	return os.RemoveAll(differ.tmpdir)
}

// WriteOld puts a piece of old content with a name.
func (differ *Differ) WriteOld(name string, data []byte) error {
	oldDir, _ := differ.subdirs()
	return os.WriteFile(filepath.Join(oldDir, name), data, 0600)
}

// WriteNew puts a piece of new content with a name.
func (differ *Differ) WriteNew(name string, data []byte) error {
	_, newDir := differ.subdirs()
	return os.WriteFile(filepath.Join(newDir, name), data, 0600)
}

// Run executes external diff command.
func (differ *Differ) Run(commandPrefix string, stdout io.Writer, stderr io.Writer) error {
	shell := ""
	var args []string
	if runtime.GOOS == "windows" {
		shell = "cmd"
		args = append(args, "/C")
	} else {
		shell = "/bin/sh"
		args = append(args, "-c")
	}
	if envShell := os.Getenv("SHELL"); envShell != "" {
		shell = envShell
	}

	oldDir, newDir := differ.subdirs()
	args = append(args, fmt.Sprintf("%s %s %s", commandPrefix, oldDir, newDir))

	cmd := exec.Command(shell, args...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	return cmd.Run()
}

func (differ *Differ) subdirs() (string, string) {
	return filepath.Join(differ.tmpdir, "old"), filepath.Join(differ.tmpdir, "new")
}
