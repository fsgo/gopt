package internal

import (
	"math/rand/v2"
	"os"
	"path/filepath"
	"strconv"

	"github.com/fsgo/gopt/internal/xos"
)

func goBinTMPDir() string {
	return filepath.Join(os.TempDir(), "fsgo", "gopt", "gobin")
}

func mv(from string, to string) error {
	if _, err := os.Stat(to); err == nil {
		tmp := filepath.Join(goBinTMPDir(), "back_"+filepath.Base(to)+"_"+strconv.Itoa(rand.Int()))
		_ = xos.MoveFile(to, tmp)
	}
	return xos.MoveFile(from, to)
}

func cleanTmpFiles() {
	_ = os.RemoveAll(goBinTMPDir())
}
