package version_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/redhatinsights/insights-ingress-go/internal/config"
	"github.com/redhatinsights/insights-ingress-go/internal/logging"
)

func TestInventory(t *testing.T) {
	cfg := config.Get()
	RegisterFailHandler(Fail)
	logging.InitLogger(cfg)
	RunSpecs(t, "Version Suite")
}
