package services_test

import (
	"sync"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/williammartin/rwmutextest"
)

var _ = Describe("ErrorsManager", func() {

	em := services.NewErrorsManager()

	It("should concurrently update the counter for that error message", func() {
		expectedTotalCount := 1000
		errorMsg := "foo"
		var wg sync.WaitGroup

		for i := 1; i <= expectedTotalCount; i++ {
			wg.Add(1)
			go func(em services.ErrorsManager, errorMsg string, index int) {
				_ = em.Store(errorMsg)
				wg.Done()
			}(em, errorMsg, i)
		}
		wg.Wait()

		Expect(em.GetCount(errorMsg)).To(Equal(int32(expectedTotalCount)))
	})
})
