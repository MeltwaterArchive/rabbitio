// Copyright © 2017 Meltwater
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package file

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

// TestNewOutput will make sure we create directory when missing.
// also checks that it is able to
func TestNewOutput(t *testing.T) {
	fs = afero.NewMemMapFs()

	fs.MkdirAll("data", 0755)
	_, err := NewOutput("data/creates_directory", 100)

	assert.Nil(t, err, "should return no error")
}

func TestPath_Receive(t *testing.T) {
	fs = afero.NewMemMapFs()

	p := &Path{
		name:      "mypath",
		batchSize: 100,
	}

	err := p.create()

	assert.Nil(t, err, "should return no error")

}

func TestWriteFile(t *testing.T) {
	fs = afero.NewMemMapFs()

	fs.MkdirAll("datadir", 0755)
	myBytes := []byte("mydatawritten")
	err := writeFile(myBytes, "datadir", "datafile")

	assert.Nil(t, err, "should return no error")
}
