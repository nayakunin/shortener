package storage

import (
	"time"
)

const MaxRequests = 10
const MaxKeysInRequest = 10
const DeleteRequestsTimeout = 3 * time.Second

type RequestBatch struct {
	UserID string
	Keys   []string
}

type RequestBuffer struct {
	buffer         chan RequestBatch
	maxRequests    int
	ticker         *time.Ticker
	isBufferFullCh chan struct{}
}

func newRequestBuffer(maxRequests int) *RequestBuffer {
	return &RequestBuffer{
		buffer:         make(chan RequestBatch, maxRequests),
		maxRequests:    maxRequests,
		ticker:         time.NewTicker(DeleteRequestsTimeout),
		isBufferFullCh: make(chan struct{}),
	}
}

func (rb *RequestBuffer) AddRequest(userID string, keys []string) {
	slicedKeys := make([][]string, 0, len(keys)/MaxKeysInRequest+1)

	for i := 0; i < len(keys); i += MaxKeysInRequest {
		end := i + MaxKeysInRequest
		if end > len(keys) {
			end = len(keys)
		}

		slicedKeys = append(slicedKeys, keys[i:end])
	}

	for _, keys := range slicedKeys {
		rb.buffer <- RequestBatch{
			UserID: userID,
			Keys:   keys,
		}

		if len(rb.buffer) == rb.maxRequests-1 {
			rb.isBufferFullCh <- struct{}{}
		}
	}
}

func (rb *RequestBuffer) GetRequests() []RequestBatch {
	requests := make([]RequestBatch, 0, rb.maxRequests)

	for i := 0; i < rb.maxRequests; i++ {
		select {
		case request := <-rb.buffer:
			requests = append(requests, request)
		default:
			return requests
		}
	}

	return requests
}
