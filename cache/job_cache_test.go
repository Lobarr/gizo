package cache_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/gizo-network/gizo/core/merkletree"
	"github.com/gizo-network/gizo/job"

	"github.com/gizo-network/gizo/cache"
	"github.com/gizo-network/gizo/core"
	"github.com/gizo-network/gizo/mocks"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("JobCache", func() {
	var mockBlockChain *mocks.IBlockChainMock
	var mockBigCache *mocks.IBigCacheMock
	var mockLogger *mocks.ILoggerMock

	BeforeEach(func() {
		mockBlockChain = &mocks.IBlockChainMock{}
		mockBigCache = &mocks.IBigCacheMock{}
		mockLogger = &mocks.ILoggerMock{}
	})

	AfterEach(func() {
		mockBlockChain = nil
		mockBigCache = nil
		mockLogger = nil
	})

	Describe("NewJobCache()", func() {
		It("should create job cache instance", func() {
			mockBlockChain.GetLatest15Func = func() ([]core.Block, error) {
				return []core.Block{}, nil
			}
			mockLogger.DebugFunc = func(...interface{}) error {
				return nil
			}
			jobCache, err := cache.NewJobCache(mockBlockChain, mockBigCache, mockLogger, false)
			Expect(err).To(BeNil())
			Expect(jobCache).ToNot(BeNil())
		})
	})

	Describe("JobCache.Get()", func() {
		Context("given get fails", func() {
			It("should return error", func() {
				expectedErr := "some-error'"
				mockBlockChain.GetLatest15Func = func() ([]core.Block, error) {
					return []core.Block{}, nil
				}
				mockBigCache.GetFunc = func(string) ([]byte, error) {
					return nil, errors.New(expectedErr)
				}
				mockLogger.DebugFunc = func(...interface{}) error {
					return nil
				}
				mockLogger.DebugfFunc = func(string, ...interface{}) error {
					return nil
				}
				jobCache, err := cache.NewJobCache(mockBlockChain, mockBigCache, mockLogger, false)
				_, err = jobCache.Get("")
				Expect(err.Error()).To(Equal(expectedErr))
			})
		})

		It("should get job from cache", func() {
			expectJobId := "some-id"
			expectedJob := job.Job{
				ID: expectJobId,
			}
			expectedJobBytes, err := json.Marshal(expectedJob)
			Expect(err).To(BeNil())
			mockBlockChain.GetLatest15Func = func() ([]core.Block, error) {
				return []core.Block{}, nil
			}
			mockBigCache.GetFunc = func(string) ([]byte, error) {
				return expectedJobBytes, nil
			}
			mockLogger.DebugFunc = func(...interface{}) error {
				return nil
			}
			mockLogger.DebugfFunc = func(string, ...interface{}) error {
				return nil
			}
			jobCache, err := cache.NewJobCache(mockBlockChain, mockBigCache, mockLogger, false)
			Expect(err).To(BeNil())
			resJob, err := jobCache.Get(expectJobId)
			Expect(err).To(BeNil())
			Expect(resJob.GetID()).To(Equal(expectJobId))
		})

		It("should log debug level", func() {
			expectedCalls := 1
			expectJobId := "some-id"
			expectedJob := job.Job{
				ID: expectJobId,
			}
			expectedJobBytes, err := json.Marshal(expectedJob)
			Expect(err).To(BeNil())
			mockBlockChain.GetLatest15Func = func() ([]core.Block, error) {
				return []core.Block{}, nil
			}
			mockBigCache.GetFunc = func(string) ([]byte, error) {
				return expectedJobBytes, nil
			}
			mockLogger.DebugFunc = func(...interface{}) error {
				return nil
			}
			mockLogger.DebugfFunc = func(string, ...interface{}) error {
				return nil
			}
			jobCache, err := cache.NewJobCache(mockBlockChain, mockBigCache, mockLogger, false)
			Expect(err).To(BeNil())
			_, err = jobCache.Get(expectJobId)
			Expect(err).To(BeNil())
			Expect(len(mockLogger.DebugfCalls())).To(Equal(expectedCalls))
		})
	})

	Describe("JobCache.IsFull()", func() {
		It("should be full", func() {
			expectedIsFull := true
			mockBlockChain.GetLatest15Func = func() ([]core.Block, error) {
				return []core.Block{}, nil
			}
			mockBigCache.LenFunc = func() int {
				return cache.MaxCacheLen
			}
			mockLogger.DebugFunc = func(...interface{}) error {
				return nil
			}
			jobCache, err := cache.NewJobCache(mockBlockChain, mockBigCache, mockLogger, false)
			Expect(err).To(BeNil())
			Expect(jobCache.IsFull()).To(Equal(expectedIsFull))
		})

		It("should not be full", func() {
			expectedIsFull := false
			mockBlockChain.GetLatest15Func = func() ([]core.Block, error) {
				return []core.Block{}, nil
			}
			mockBigCache.LenFunc = func() int {
				return 0
			}
			mockLogger.DebugFunc = func(...interface{}) error {
				return nil
			}
			jobCache, err := cache.NewJobCache(mockBlockChain, mockBigCache, mockLogger, false)
			Expect(err).To(BeNil())
			Expect(jobCache.IsFull()).To(Equal(expectedIsFull))
		})
	})

	Describe("JobCache.fill()", func() {
		It("should fill up cache", func() {
			expectedJobId := "test-id"
			expectedCalls := 1
			mockBlockChain.GetLatest15Func = func() ([]core.Block, error) {
				return []core.Block{
					core.Block{
						Jobs: []*merkletree.MerkleNode{
							&merkletree.MerkleNode{
								Job: job.Job{
									ID: expectedJobId,
								},
							},
						},
					},
				}, nil
			}
			mockBigCache.SetFunc = func(string, []byte) error {
				return nil
			}
			mockBigCache.LenFunc = func() int {
				return 0
			}
			mockLogger.DebugFunc = func(...interface{}) error {
				return nil
			}
			mockLogger.DebugfFunc = func(string, ...interface{}) error {
				return nil
			}
			_, err := cache.NewJobCache(mockBlockChain, mockBigCache, mockLogger, false)
			Expect(err).To(BeNil())
			Expect(len(mockBlockChain.GetLatest15Calls())).To(Equal(expectedCalls))
			Expect(mockBigCache.SetCalls()[0].In1).To(Equal(expectedJobId))
		})

		Context("given blockchain fails", func() {
			It("should return error", func() {
				expectedErr := "some-error"
				expectedCalls := 1
				mockBlockChain.GetLatest15Func = func() ([]core.Block, error) {
					return nil, errors.New(expectedErr)
				}
				mockLogger.DebugFunc = func(...interface{}) error {
					return nil
				}
				mockLogger.DebugfFunc = func(string, ...interface{}) error {
					return nil
				}
				_, err := cache.NewJobCache(mockBlockChain, mockBigCache, mockLogger, false)
				Expect(len(mockBlockChain.GetLatest15Calls())).To(Equal(expectedCalls))
				Expect(err.Error()).To(Equal(expectedErr))
			})
		})

		It("should log debug error", func() {
			expectedCalls := 1
			mockBlockChain.GetLatest15Func = func() ([]core.Block, error) {
				return []core.Block{}, nil
			}

			mockLogger.DebugFunc = func(...interface{}) error {
				return nil
			}

			_, err := cache.NewJobCache(mockBlockChain, mockBigCache, mockLogger, false)
			Expect(err).To(BeNil())
			Expect(len(mockLogger.DebugCalls())).To(Equal(expectedCalls))
		})
	})

	Describe("JobCache.StopWatchMode()", func() {
		It("should stop watch mode", func() {
			expectedCalls := 1
			mockBlockChain.GetLatest15Func = func() ([]core.Block, error) {
				return []core.Block{}, nil
			}
			mockLogger.DebugFunc = func(...interface{}) error {
				return nil
			}
			jobCache, err := cache.NewJobCache(mockBlockChain, mockBigCache, mockLogger, true)
			Expect(err).To(BeNil())
			Expect(len(mockLogger.DebugCalls())).To(Equal(expectedCalls))
			jobCache.StopWatchMode()
			Expect(len(mockLogger.DebugCalls())).To(Equal(expectedCalls))
		})
	})

	Describe("JobCache.Set()", func() {
		Context("given cache is full", func() {
			It("should return error", func() {
				expectedKey := "some-key"
				expectedVal := []byte{}
				mockBlockChain.GetLatest15Func = func() ([]core.Block, error) {
					return []core.Block{}, nil
				}
				mockBigCache.LenFunc = func() int {
					return cache.MaxCacheLen
				}

				mockLogger.DebugFunc = func(...interface{}) error {
					return nil
				}
				mockLogger.DebugfFunc = func(string, ...interface{}) error {
					return nil
				}
				jobCache, err := cache.NewJobCache(mockBlockChain, mockBigCache, mockLogger, false)
				err = jobCache.Set(expectedKey, expectedVal)
				Expect(err.Error()).To(Equal(cache.ErrCacheFull.Error()))
			})
		})

		Context("given cache is not full", func() {
			It("should get job from cache", func() {
				expectJobId := "some-id"
				expectedJob := job.Job{
					ID: expectJobId,
				}
				expectedJobBytes, err := json.Marshal(expectedJob)
				Expect(err).To(BeNil())
				mockBlockChain.GetLatest15Func = func() ([]core.Block, error) {
					return []core.Block{}, nil
				}
				mockBigCache.LenFunc = func() int {
					return 0
				}
				mockBigCache.SetFunc = func(string, []byte) error {
					return nil
				}
				mockLogger.DebugFunc = func(...interface{}) error {
					return nil
				}
				mockLogger.DebugfFunc = func(string, ...interface{}) error {
					return nil
				}
				jobCache, err := cache.NewJobCache(mockBlockChain, mockBigCache, mockLogger, false)
				Expect(err).To(BeNil())
				err = jobCache.Set(expectJobId, expectedJobBytes)
				Expect(err).To(BeNil())
				Expect(mockBigCache.SetCalls()[0].In1).To(Equal(expectJobId))
				Expect(mockBigCache.SetCalls()[0].In2).To(Equal(expectedJobBytes))
			})
		})

		It("should log debug level", func() {
			expectJobId := "some-id"
			expectedJob := []byte{}
			expectedCalls := 1
			mockBlockChain.GetLatest15Func = func() ([]core.Block, error) {
				return []core.Block{}, nil
			}
			mockBigCache.LenFunc = func() int {
				return 0
			}
			mockBigCache.SetFunc = func(string, []byte) error {
				return nil
			}
			mockLogger.DebugFunc = func(...interface{}) error {
				return nil
			}
			mockLogger.DebugfFunc = func(string, ...interface{}) error {
				return nil
			}
			jobCache, err := cache.NewJobCache(mockBlockChain, mockBigCache, mockLogger, false)
			Expect(err).To(BeNil())
			err = jobCache.Set(expectJobId, expectedJob)
			Expect(err).To(BeNil())
			Expect(len(mockLogger.DebugfCalls())).To(Equal(expectedCalls))
		})
	})
})

func TestJobCache(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "JobCache")
}
