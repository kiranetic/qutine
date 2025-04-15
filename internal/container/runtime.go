package container

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

// RunContainer extracts a tarball into tmpfs and runs a binary
func RunContainer(imageTarPath, command string) error {
	tmpMount := filepath.Join("/tmp", fmt.Sprintf("qutine-%d", os.Getpid()))

	if err := os.MkdirAll(tmpMount, 0700); err != nil {
		return fmt.Errorf("mkdir failed: %v", err)
	}
	if err := syscall.Mount("tmpfs", tmpMount, "tmpfs", 0, "size=100M"); err != nil {
		return fmt.Errorf("failed to mount tmpfs: %v", err)
	}

	rootfs := filepath.Join(tmpMount, "rootfs")
	if err := os.Mkdir(rootfs, 0755); err != nil {
		return fmt.Errorf("mkdir rootfs failed: %v", err)
	}
	if err := extractTar(imageTarPath, rootfs); err != nil {
		return fmt.Errorf("failed to extract tar: %v", err)
	}

	// Debug listing
	fmt.Println("=== Extracted files ===")
	filepath.Walk(rootfs, func(path string, info os.FileInfo, err error) error {
		if err == nil {
			fmt.Println(path)
		}
		return nil
	})
	fmt.Println("=======================")

	targetBin := filepath.Join(rootfs, command)
	if _, err := os.Stat(targetBin); os.IsNotExist(err) {
		fmt.Printf("DEBUG: %s does not exist\n", targetBin)
		return err
	}

	cmd := exec.Command(targetBin)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = rootfs

	return cmd.Run()
}

func extractTar(tarPath, dest string) error {
	file, err := os.Open(tarPath)
	if err != nil {
		return err
	}
	defer file.Close()

	tr := tar.NewReader(file)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		target := filepath.Join(dest, hdr.Name)
		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, os.FileMode(hdr.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(f, tr); err != nil {
				f.Close()
				return err
			}
			f.Close()
		}
	}
	return nil
}
