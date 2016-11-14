package client_test

import (
	"errors"

	"github.com/apihub/apihub"
	"github.com/apihub/apihub/client"
	"github.com/apihub/apihub/client/connection/connectionfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Service", func() {
	var (
		fakeConnection *connectionfakes.FakeConnection
		cli            apihub.Client
		service        apihub.Service
		spec           apihub.ServiceSpec
	)

	BeforeEach(func() {
		var err error
		fakeConnection = new(connectionfakes.FakeConnection)
		cli = client.New(fakeConnection)

		spec = apihub.ServiceSpec{
			Handle:   "my-handle",
			Disabled: true,
			Timeout:  10,
			Backends: []apihub.BackendInfo{
				apihub.BackendInfo{
					Address: "http://server-a",
				},
				apihub.BackendInfo{
					Address: "http://server-b",
				},
			},
		}
		fakeConnection.AddServiceReturns(spec, nil)

		service, err = cli.AddService(spec)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("Handle", func() {
		It("returns service's handle", func() {
			Expect(service.Handle()).To(Equal("my-handle"))
		})
	})

	Describe("Info", func() {
		BeforeEach(func() {
			fakeConnection.FindServiceReturns(spec, nil)
		})

		It("returns service's info", func() {
			info, err := service.Info()
			Expect(err).NotTo(HaveOccurred())
			Expect(info.Handle).To(Equal(spec.Handle))
			Expect(info.Disabled).To(Equal(spec.Disabled))
			Expect(info.Timeout).To(Equal(spec.Timeout))
			Expect(info.Backends).To(ConsistOf(spec.Backends))
		})

		Context("when fails to get the info", func() {
			BeforeEach(func() {
				fakeConnection.FindServiceReturns(apihub.ServiceSpec{}, errors.New("fail to get the info"))
			})

			It("returns an error", func() {
				_, err := service.Info()
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("Backends", func() {
		BeforeEach(func() {
			fakeConnection.FindServiceReturns(spec, nil)
		})

		It("returns service's backends", func() {
			backends, err := service.Backends()
			Expect(err).NotTo(HaveOccurred())
			Expect(len(backends)).To(Equal(2))
		})

		Context("when fails to get the backend list", func() {
			BeforeEach(func() {
				fakeConnection.FindServiceReturns(apihub.ServiceSpec{}, errors.New("fail to get the backend list"))
			})

			It("returns an error", func() {
				_, err := service.Backends()
				Expect(err).To(HaveOccurred())
			})
		})
	})

})