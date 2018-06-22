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

	"github.com/meltwater/rabbitio/rmq"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestNewTarballBuilder(t *testing.T) {
	_, err := NewTarballBuilder(1000)
	assert.NoError(t, err, "should not return error when creating a tarball builder")
}

func TestTarballBuilder_GetWriters(t *testing.T) {
	tarball, _ := NewTarballBuilder(1000)
	err := tarball.getWriters()

	assert.NoError(t, err)
}

func TestTarballBuilder_AddFile(t *testing.T) {
	tarball, _ := NewTarballBuilder(1000)
	m := &rmq.Message{Body: []byte("mymessage"), RoutingKey: "rk"}

	err := tarball.addFile(tarball.tar, "file.tgz", m)

	assert.NoError(t, err)
}

func TestTarballBuilder_Pack(t *testing.T) {
	fs = afero.NewMemMapFs()

	tarball, _ := NewTarballBuilder(1)
	ch := make(chan rmq.Message, 1)
	verify := make(chan rmq.Verify, 1)
	fs.MkdirAll("/data", 0755)

	var err error

	ch <- rmq.Message{Body: []byte("mymessage")}
	go func() {
		err = tarball.Pack(ch, "/data", verify)
	}()

	close(ch)

	assert.NoError(t, err, "received no error")
}
