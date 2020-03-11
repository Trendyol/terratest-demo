package util

import (
	. "github.com/onsi/ginkgo"
	"sync"
)

var BeforeAll = func(beforeAllFunc func()) {
	var once sync.Once

	BeforeEach(func() {
		once.Do(func() {
			beforeAllFunc()
		})
	})
}
