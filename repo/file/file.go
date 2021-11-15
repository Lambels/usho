package file

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"sync"

	"github.com/Lambels/usho/encoding/base64"
	"github.com/Lambels/usho/repo"
)

var sizeOfLength = 8

func New(path string) (repo.Repo, error) {
	_, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	return &storage{
		path: path,
	}, nil
}

type storage struct {
	path string
	sync.Mutex
}

func (s *storage) New(ctx context.Context, in repo.URLRequest) (repo.URLResponse, error) {
	s.Lock()
	defer s.Unlock()

	var out repo.URLResponse

	f, err := os.OpenFile(s.path, os.O_RDWR, os.ModeAppend)
	if err != nil {
		return out, err
	}

	d, err := io.ReadAll(f)
	if err != nil {
		return out, err
	}

	id := s.GenerateID(ctx, d)

	URL := []byte(in.Intial)

	if err := binary.Write(f, binary.LittleEndian, int64(len(URL))); err != nil {
		return out, err
	}
	if err := binary.Write(f, binary.LittleEndian, id); err != nil {
		return out, err
	}

	if _, err := f.Write(URL); err != nil {
		return out, err
	}

	out.ID = id
	out.Initial = in.Intial
	out.Short = base64.EncodeID(id)

	return out, f.Close()
}

func (s *storage) Get(ctx context.Context, in string) (repo.URLResponse, error) {
	s.Lock()
	defer s.Unlock()

	var out repo.URLResponse

	d, err := ioutil.ReadFile(s.path)
	if err != nil {
		return out, err
	}

	inID, err := base64.DecodeID(in)
	if err != nil {
		return out, err
	}

	for {
		if len(d) == 0 {
			return out, repo.ErrNotFound
		}
		if len(d) < sizeOfLength {
			return out, fmt.Errorf("abnormal number of bytes in %s: %d", s.path, len(d))
		}

		var len uint64
		if err := binary.Read(bytes.NewReader(d[:sizeOfLength]), binary.LittleEndian, &len); err != nil {
			return out, err
		}
		d = d[sizeOfLength:]

		var id uint64
		if err := binary.Read(bytes.NewReader(d), binary.LittleEndian, &id); err != nil {
			return out, err
		}
		d = d[sizeOfLength:]

		URL := string(d[:len])
		if URL == in {
			out.ID = id
			out.Initial = URL
			out.Short = base64.EncodeID(id)

			return out, nil
		} else if inID == id {
			out.ID = id
			out.Initial = URL
			out.Short = base64.EncodeID(id)

			return out, nil
		}
		d = d[len:]
	}
}

func (s *storage) Close() error { return nil }

func (s *storage) GenerateID(ctx context.Context, d []byte) uint64 {
	done := make(chan uint64, 1)
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			default:
				// Generate random int64
				buf := make([]byte, 8)
				rand.Read(buf)
				randId := binary.LittleEndian.Uint64(buf)

				// gorutine checks for id and if id works it pushes id into done
				go func(d []byte, id uint64, done chan uint64) {
					for {
						if len(d) == 0 {
							done <- id
							return
						}
						if len(d) < sizeOfLength {
							return
						}

						var len uint64
						if err := binary.Read(bytes.NewReader(d[:sizeOfLength]), binary.LittleEndian, &len); err != nil {
							return
						}
						d = d[sizeOfLength:]

						var selId uint64
						if err := binary.Read(bytes.NewReader(d[:sizeOfLength]), binary.LittleEndian, &selId); err != nil {
							return
						}
						d = d[sizeOfLength:]
						d = d[len:]

						if id == selId {
							return
						}
					}

				}(d, randId, done)
			}
		}
	}()

	// waits for id
	id := <-done

	// kills spawning gorutine
	cancel()
	return id
}
