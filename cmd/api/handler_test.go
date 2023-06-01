package main

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Test handler", func() {
	When("Acess auth", func() {
		It("should be successful", func() {
			fact := true
			Expect(fact).To(BeTrue())
			Expect(fact).To(BeFalse())
		})
	})
})
