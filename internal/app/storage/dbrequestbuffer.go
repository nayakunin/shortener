package storage

import (
	"time"
)

// MaxRequests is the maximum number of requests that can be sent to the database at a time.
const MaxRequests = 10

// MaxKeysInRequest is the maximum number of keys that can be sent to the database at a time.
const MaxKeysInRequest = 10

// DeleteRequestsTimeout is the maximum time to wait for a request to be sent to the database.
const DeleteRequestsTimeout = 3 * time.Second

// RequestBatch is a batch of requests to the database.
type RequestBatch struct {
	UserID string
	Keys   []string
}

// RequestBuffer is a buffer for requests to the database.
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

// AddRequest adds a request to the buffer
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

// GetRequests returns a slice of requests from the buffer
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
