package e2e

import (
	"fmt"
	"github.com/onsi/gomega/gexec"
	v1 "github.com/progimage/pkg/models/v1"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

var client *v1.Client
var httpServerSession *gexec.Session

func TestMain(m *testing.M) {
	var err error
	client, err = v1.NewClient("http://localhost:8080/api/v1/")
	if err != nil {
		panic("unable to initialize test client")
	}

	// build and start up the prog image http server to run e2e tests against
	progImagePath, err := gexec.Build("../")
	if err != nil {
		panic("unable to build progimage")
	}

	httpServerSession, err = gexec.Start(exec.Command(filepath.Join(progImagePath, "progimage"), "-port", "8080", "-basePath", os.TempDir()), os.Stdout, os.Stderr)
	if err != nil {
		panic(fmt.Sprintf("unable to run progimage %v", err))
	}

	code := m.Run()
	httpServerSession.Kill()
	os.Exit(code)
}
