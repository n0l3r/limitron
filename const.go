package limitron

const (
	// Rate limiting algorithms
	AlgorithmLeakyBucket   = "leaky_bucket"
	AlgorithmTokenBucket   = "token_bucket"
	AlgorithmFixedWindow   = "fixed_window"
	AlgorithmSlidingWindow = "sliding_window"

	// Store types
	StoreTypeMemory = "memory"
	StoreTypeRedis  = "redis"
)
