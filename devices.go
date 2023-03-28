package webcam

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/unix"
)

type CameraInfo struct {
	Path   string
	Driver string
	Card   string
	Bus    string
}

func ListDevices() ([]CameraInfo, error) {
	files, err := os.ReadDir("/dev")
	if err != nil {
		return nil, fmt.Errorf("read /dev: %w", err)
	}

	cameras := make([]CameraInfo, 0)

	for _, f := range files {
		if !strings.HasPrefix(f.Name(), "video") {
			continue
		}

		path := filepath.Join("/dev", f.Name())

		handle, err := unix.Open(path, unix.O_RDWR|unix.O_NONBLOCK, 0666)
		fd := uintptr(handle)

		if fd < 0 || err != nil {
			continue
		}

		caps, err := checkCapabilities(fd)
		unix.Close(handle)

		if err != nil {
			continue
		}

		cameras = append(cameras, caps.getCameraInfo(path))
	}

	return cameras, nil
}
