package bricksplanner

import (
	"github.com/gluster/glusterd2/pkg/api"
	"github.com/gluster/glusterd2/plugins/device/deviceutils"
)

// defaultLeastArbiterSize is the size (in KB) the arbiter brick will be assigned to if the brick size is less than 100M.
const defaultLeastArbiterSize = 100000

type replicaSubvolPlanner struct {
	subvolSize       uint64
	replicaCount     int
	arbiterCount     int
	brickSize        uint64
	arbiterBrickSize uint64
}

func (s *replicaSubvolPlanner) Init(req *api.VolCreateReq, subvolSize uint64) {
	s.subvolSize = subvolSize
	s.replicaCount = req.ReplicaCount
	s.arbiterCount = req.ArbiterCount
	s.brickSize = s.subvolSize
	// TODO: the size is calculated in KB, should be changed as per the default.
	// default avgFileSize needs to be changed from 1M to 64K as well.
	// Now we are receiving 1M as default from cli so having it as 1M here as well,
	var avgFileSize uint64 = 1024
	if req.AverageFileSize != 0 {
		avgFileSize = deviceutils.MbToKb(req.AverageFileSize)
	}
	arbiterSize := uint64((4.0) * (float64(subvolSize) / float64(avgFileSize)))
	// Assigning arbiter brick size to be bricksize if its lesser than 100M
	if arbiterSize < defaultLeastArbiterSize {
		if defaultLeastArbiterSize > subvolSize {
			arbiterSize = subvolSize
		} else {
			arbiterSize = defaultLeastArbiterSize
		}
	}
	s.arbiterBrickSize = arbiterSize
}

func (s *replicaSubvolPlanner) BricksCount() int {
	return s.replicaCount + s.arbiterCount
}

func (s *replicaSubvolPlanner) BrickSize(idx int) uint64 {
	if idx == (s.replicaCount) && s.arbiterCount > 0 {
		return s.arbiterBrickSize
	}

	return s.brickSize
}

func (s *replicaSubvolPlanner) BrickType(idx int) string {
	if idx == (s.replicaCount) && s.arbiterCount > 0 {
		return "arbiter"
	}

	return "brick"
}

func init() {
	subvolPlanners["replicate"] = &replicaSubvolPlanner{}
}
