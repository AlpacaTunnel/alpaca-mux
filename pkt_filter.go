package alpacamux

import (
	"time"
)

const RATE_LIMIT = 160000 // pps

// TODO: bool is one byte, how about use 1 bit struct?
type PktFilter struct {
	Limit  uint32
	Latest uint32
	Mark0  []bool
	Mark1  []bool
	Mark2  []bool
}

func (filter *PktFilter) Init() {
	filter.Limit = RATE_LIMIT
	filter.Mark0 = make([]bool, filter.Limit)
	filter.Mark1 = make([]bool, filter.Limit)
	filter.Mark2 = make([]bool, filter.Limit)
}

func Abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func (filter *PktFilter) IsValid(timestamp, sequence uint32) bool {
	if sequence >= filter.Limit {
		log.Debug("Pkt sequence number exceeded limit")
		return false
	}

	if Abs(int64(timestamp)-time.Now().Unix()) > 2592000 {
		log.Debug("Peer timestamp shifts beyond 30 days")
		return false
	}

	if int(timestamp-filter.Latest) < -600 {
		log.Debug("Pkt delayed more than 600s, treat as invalid")
		return false
	}

	if filter.isDuplicated(timestamp, sequence) {
		log.Debug("Pkt is duplicated")
		return false
	}

	return true
}

func (filter *PktFilter) isDuplicated(timestamp, sequence uint32) bool {
	diff := int64(timestamp) - int64(filter.Latest)

	if diff > 2 {
		filter.Latest = timestamp
		// make new slice to reset all to false
		filter.Mark0 = make([]bool, filter.Limit)
		filter.Mark1 = make([]bool, filter.Limit)
		filter.Mark2 = make([]bool, filter.Limit)

		filter.Mark0[sequence] = true
		return false
	}

	if diff == 2 {
		filter.Latest = timestamp
		filter.Mark2 = filter.Mark0
		filter.Mark1 = make([]bool, filter.Limit)
		filter.Mark0 = make([]bool, filter.Limit)

		filter.Mark0[sequence] = true
		return false
	}

	if diff == 1 {
		filter.Latest = timestamp
		filter.Mark2 = filter.Mark1
		filter.Mark1 = filter.Mark0
		filter.Mark0 = make([]bool, filter.Limit)

		filter.Mark0[sequence] = true
		return false
	}

	if diff == 0 {
		if filter.Mark0[sequence] {
			return true
		} else {
			filter.Mark0[sequence] = true
			return false
		}
	}

	if diff == -1 {
		if filter.Mark1[sequence] {
			return true
		} else {
			filter.Mark1[sequence] = true
			return false
		}
	}

	if diff == -2 {
		if filter.Mark2[sequence] {
			return true
		} else {
			filter.Mark2[sequence] = true
			return false
		}
	}

	// if diff < -2, do nothing and treat it as not dup
	return false
}
