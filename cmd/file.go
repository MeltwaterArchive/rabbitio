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

package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/meltwater/rabbitio/rmq"
)

// FileInput is nice
type FileInput struct {
	queue []string
}

// Path is directory path for consumed RabbitMQ messages
type Path struct {
	name      string
	batchSize int
}

// NewFileInput creates a FileInput from the specified directory
func NewFileInput(path string) *FileInput {
	fi, err := os.Stat(path)
	if err != nil {
		log.Fatalln(err)
	}

	var f *FileInput
	q := []string{}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		files, err := ioutil.ReadDir(path)
		if err != nil {
			log.Fatalf("Couldn't get directory or file: %s", err)
		}
		log.Printf("Found %d file(s) in %s", len(files), path)
		for _, f := range files {
			q = append(q, filepath.Join(path, f.Name()))
		}
	case mode.IsRegular():
		q = append(q, path)
	}

	f = &FileInput{
		queue: q,
	}

	return f
}

func writeFile(b []byte, dir, file string) {
	filePath := filepath.Join(dir, file)
	err := ioutil.WriteFile(filePath, b, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Wrote %d bytes to %s", len(b), filePath)
}

// Send delivers messages to the channel
func (f *FileInput) Send(messages chan rmq.Message) {
	var num int

	// loop over the queued up files
	for _, file := range f.queue {
		// open file from the queue
		fh, err := os.Open(file)
		if err != nil {
			log.Fatalf("failed to open file: %s", err)
		}
		// and clean up afterwards
		defer fh.Close()

		tarNum, err := UnPack(fh, messages)
		if err != nil {
			log.Fatalf("Failed to unpack: %s ", err)
		}
		log.Printf("Extracted %d Messages from tarball: %s", tarNum, file)
		num = num + tarNum
	}
	// when all files are read, close
	close(messages)
	log.Printf("Total %d Messages from tarballs", num)

}

// NewFileOutput creates a Path to output files in from RabbitMQ
func NewFileOutput(path string, batchSize int) *Path {
	return &Path{
		name:      path,
		batchSize: batchSize,
	}
}

// Receive will handle messages and save to path
func (p *Path) Receive(messages chan Message) {

	// create new TarballBuilder
	builder := NewTarballBuilder(p.batchSize)

	builder.Pack(messages, p.name)

}
