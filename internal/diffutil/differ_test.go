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
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiffer(t *testing.T) {
	is := assert.New(t)

	differ, err := NewDiffer()
	is.NoError(err)
	defer differ.Cleanup()

	oldDir, newDir := differ.subdirs()

	differ.WriteOld("foo", []byte("bar"))
	differ.WriteNew("foo", []byte("bar"))

	is.NoError(differ.Run("diff", nil, nil))

	is.NoError(differ.Cleanup())
	is.False(checkDirectoryExists(oldDir))
	is.False(checkDirectoryExists(newDir))
}

func checkDirectoryExists(dir string) bool {
	_, err := os.Stat(dir)
	return err == nil
}
