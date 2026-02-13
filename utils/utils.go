package utils

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"
)

func LoadEnv(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line[0] == '#' {
			continue
		}
		key, val := func() (string, string) {
			x := strings.Split(line, "=")
			return x[0], x[1]
		}()
		os.Setenv(key, val)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func HumanMessageSize(bytes uint, si bool, dp int) string {
	var thresh uint
	if si {
		thresh = 1000
	} else {
		thresh = 1024
	}

	if bytes < thresh {
		return fmt.Sprintf("%d B", bytes)
	}

	var units [8]string
	if si {
		units = [8]string{"KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}
	} else {
		units = [8]string{"KiB", "MiB", "GiB", "TiB", "PiB", "EiB", "ZiB", "YiB"}
	}

	x := float64(bytes)
	u := -1

	for x >= float64(thresh) && u < len(units)-1 {
		x /= float64(thresh)
		u++
	}

	return fmt.Sprintf("%.*f %s", dp, x, units[u])
}

func PrintBanner(path string) {
	b, err := os.ReadFile(path)
	if err != nil {
		return
	}
	fmt.Printf("\x1b[0;93m%s\x1b[0m\n", string(b))
}

func DecodeBase64(b []byte) ([]byte, error) {
	b, err := base64.StdEncoding.DecodeString(string(b))
	if err != nil {
		return nil, nil
	}
	return b, nil
}
