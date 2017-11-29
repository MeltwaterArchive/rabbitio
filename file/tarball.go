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
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/meltwater/rabbitio/rmq"
	"github.com/pborman/uuid"
)

// TarballBuilder build tarballs from the stream of incoming docs
// and spits out tarballs into a channel
type TarballBuilder struct {
	lock    sync.Mutex
	tarSize int
	wg      sync.WaitGroup
	buf     *bytes.Buffer
	gzip    *gzip.Writer
	tar     *tar.Writer
}

// NewTarballBuilder created a TarballBuilder
func NewTarballBuilder(tarSize int) *TarballBuilder {
	t := &TarballBuilder{
		tarSize: tarSize,
	}
	err := t.getWriters()
	if err != nil {
		log.Fatal(err)
	}
	return t
}

// get a new set of writers to write to
func (t *TarballBuilder) getWriters() (err error) {
	t.lock.Lock()

	t.buf = new(bytes.Buffer)
	t.gzip, err = gzip.NewWriterLevel(t.buf, gzip.BestCompression)
	t.tar = tar.NewWriter(t.gzip)

	t.lock.Unlock()
	return err
}

// add a new file to the tarball writer
func (t *TarballBuilder) addFile(tw *tar.Writer, name string, m *rmq.Message) error {
	header := new(tar.Header)
	header.Name = name
	header.Size = int64(len(m.Body))
	header.Mode = 0644
	header.ModTime = time.Now()
	header.Xattrs = make(map[string]string)
	header.Xattrs["amqp.routingKey"] = m.RoutingKey

	for k, v := range m.Headers {
		switch v.(type) {
		case string:
			header.Xattrs[k] = v.(string)
		}
	}

	if err := tw.WriteHeader(header); err != nil {
		return err
	}
	if _, err := tw.Write(m.Body); err != nil {
		return err
	}
	return nil
}

// UnPack will decompress and send messages out on channel from file
func UnPack(wg *sync.WaitGroup, file *os.File, messages chan rmq.Message) (n int, err error) {

	// wrap fh in a gzip reader
	gr, err := gzip.NewReader(file)

	// adds tar reader in the gzip
	tr := tar.NewReader(gr)

	// loop over the files in the tarball
	for {
		hdr, terr := tr.Next()
		if terr == io.EOF {
			// end of tar archive
			break
		}
		if terr != nil {
			return n, terr
		}
		wg.Add(1)

		// create a Buffer to work on
		// TODO: reuse if GC pressure is a problem
		buf := bytes.NewBuffer(make([]byte, 0, hdr.Size))

		// copy the doc from the tarball to our buffer
		if _, err = io.Copy(buf, tr); err != nil {
			return n, err
		}

		// generate and push the message to the output channel
		messages <- *rmq.NewMessageFromAttrs(buf.Bytes(), hdr.Xattrs)
		n++
	}
	return n, err
}

// Pack messages from the channel into the directory
func (t *TarballBuilder) Pack(messages chan rmq.Message, dir string, verify chan rmq.Verify) {

	t.wg.Add(1)

	docNum := 0
	fileNum := 0

	var deliveryTag uint64
	for doc := range messages {
		deliveryTag = doc.DeliveryTag

		docNum++
		if docNum >= t.tarSize {

			fileNum++
			t.tar.Flush()
			t.tar.Close()
			t.gzip.Close()

			// writes to tarball here when reached the t.tarSize
			err := writeFile(t.buf.Bytes(), dir, fmt.Sprintf("%d_messages_%d.tgz", fileNum, docNum))
			if err != nil {
				log.Fatal(err)
			}
			verify <- rmq.Verify{MultiAck: true, Tag: doc.DeliveryTag}

			err = t.getWriters()
			if err != nil {
				log.Fatal(err)
			}
			docNum = 0
		}

		if err := t.addFile(t.tar, uuid.New(), &doc); err != nil {
			log.Fatalln(err)
		}
	}
	t.tar.Flush()
	t.tar.Close()
	t.gzip.Close()

	fileNum++

	// writes to tarball here when not reached the t.tarSize
	err := writeFile(t.buf.Bytes(), dir, fmt.Sprintf("%d_messages_%d.tgz", fileNum, docNum))
	if err != nil {
		log.Fatal(err)

	}

	// Does not ack the messages unless it is repeated, not sure why yet..
	// Might want to change using delivery ack interface
	verify <- rmq.Verify{MultiAck: true, Tag: deliveryTag}
	verify <- rmq.Verify{MultiAck: true, Tag: deliveryTag}

	t.wg.Done()
	close(verify)
	log.Print("tarball writer closing")
}

// CloseWaiter waits for the wg and then closes
func (t *TarballBuilder) CloseWaiter(out chan []byte) {
	t.wg.Wait()
	close(out)
}
