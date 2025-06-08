package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHearX(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "HearX Suite")
}
