package sets

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

const (
	// LORO min < x < max
	LORO = iota
	// LORC min < x <= max
	LORC
	// LCRO min <= x < max
	LCRO
	// LCRC min <= x <=max
	LCRC
	//LCRI min <= x
	LCRI
	// LORI min < x
	LORI
	// LIRC x <= max
	LIRC
	// LIRO x < max
	LIRO
)

type Bound int

const (
	UpperBound Bound = iota
	LowerBound
)

// INF 无穷大
const INF = 10000000

// NINF 负无穷大
const NINF = -10000000

// CntRange 数量范围
type CntRange struct {
	min  int
	max  int
	mode int
}

// NewRange ... ps:str [1, 3]
func NewRange(str string) (*CntRange, error) {
	if strings.Contains(strings.ToLower(str), "inf") {
		return newInfRange(str)
	}
	return newNormRange(str)
}

func newNormRange(str string) (*CntRange, error) {
	ss := strings.Split(str, ",")
	if len(ss) != 2 {
		return nil, errors.New("format err, ps [1, 3]")
	}
	ss[0] = strings.TrimSpace(ss[0])
	ss[1] = strings.TrimSpace(ss[1])

	var lopen, lclose, ropen, rclose bool
	var min, max, mode int
	if strings.Contains(ss[1], "]") {
		rstr := strings.Trim(ss[1], "]")
		num, err := strconv.Atoi(rstr)
		if err != nil {
			return nil, err
		}
		max, rclose = num, true
	} else if strings.Contains(ss[1], ")") {
		rstr := strings.Trim(ss[1], ")")
		num, err := strconv.Atoi(rstr)
		if err != nil {
			return nil, err
		}
		max, ropen = num, true
	} else {
		return nil, errors.New("not ) or ]")
	}

	if strings.Contains(ss[0], "[") {
		lstr := strings.Trim(ss[0], "[")
		num, err := strconv.Atoi(lstr)
		if err != nil {
			return nil, err
		}
		min, lclose = num, true
	} else if strings.Contains(ss[0], "(") {
		lstr := strings.Trim(ss[0], "(")
		num, err := strconv.Atoi(lstr)
		if err != nil {
			return nil, err
		}
		min, lopen = num, true
	} else {
		return nil, errors.New("not ( or [")
	}

	if min > max {
		return nil, fmt.Errorf("err min=%d > max=%d", min, max)
	}

	if lclose && rclose {
		mode = LCRC
	} else if lclose && ropen {
		mode = LCRO
	} else if lopen && rclose {
		mode = LORC
	} else {
		mode = LORO
	}

	return &CntRange{min: min, max: max, mode: mode}, nil
}

func newInfRange(str string) (*CntRange, error) {
	ss := strings.Split(str, ",")
	if len(ss) != 2 {
		return nil, errors.New("format err, ps [1, 3]")
	}
	ss[0] = strings.TrimSpace(ss[0])
	ss[1] = strings.TrimSpace(ss[1])

	var min, max, mode int
	if strings.Contains(strings.ToLower(ss[0]), "inf") { // 左侧无穷
		if strings.Contains(ss[1], "]") {
			rstr := strings.Trim(ss[1], "]")
			num, err := strconv.Atoi(rstr)
			if err != nil {
				return nil, err
			}
			min, max, mode = NINF, num, LIRC
		} else if strings.Contains(ss[1], ")") {
			rstr := strings.Trim(ss[1], ")")
			num, err := strconv.Atoi(rstr)
			if err != nil {
				return nil, err
			}
			min, max, mode = NINF, num, LIRO
		} else {
			return nil, errors.New("not ) or ]")
		}
	} else { //  右侧为无穷
		if strings.Contains(ss[0], "[") {
			lstr := strings.Trim(ss[0], "[")
			num, err := strconv.Atoi(lstr)
			if err != nil {
				return nil, err
			}
			min, max, mode = num, INF, LCRI
		} else if strings.Contains(ss[0], "(") {
			lstr := strings.Trim(ss[0], "(")
			num, err := strconv.Atoi(lstr)
			if err != nil {
				return nil, err
			}
			min, max, mode = num, INF, LORI
		} else {
			return nil, errors.New("not ( or [")
		}
	}

	return &CntRange{min: min, max: max, mode: mode}, nil
}

// InRange ...
func (r *CntRange) InRange(x int) bool {
	switch r.mode {
	case LORO:
		if r.min < x && x < r.max {
			return true
		}
	case LORC:
		if r.min < x && x <= r.max {
			return true
		}
	case LCRO:
		if r.min <= x && x < r.max {
			return true
		}
	case LCRC:
		if r.min <= x && x <= r.max {
			return true
		}
	case LCRI:
		if r.min <= x {
			return true
		}
	case LORI:
		if r.min < x {
			return true
		}
	case LIRC:
		if x <= r.max {
			return true
		}
	case LIRO:
		if x < r.max {
			return true
		}
	}

	return false
}

// Rand 生成改范围内的随机数
func (r *CntRange) Rand() int {
	min, max := r.getMinMax()
	return min + rand.Intn(max+1-min)
}

// GetMinMax 获得边界
func (r *CntRange) getMinMax() (min, max int) {
	switch r.mode {
	case LORO:
		min, max = r.min+1, r.max-1
	case LORC:
		min, max = r.min+1, r.max
	case LCRO:
		min, max = r.min, r.max-1
	case LCRC:
		min, max = r.min, r.max
	case LCRI:
		min, max = r.min, r.max
	case LORI:
		min, max = r.min+1, r.max
	case LIRC:
		min, max = r.min, r.max
	case LIRO:
		min, max = r.min, r.max-1
	}
	return min, max
}

// RandWithArgs data, mode 用data 根据mode替换上界 或 下界
func (r *CntRange) RandWithArgs(data int, mode Bound) int {
	min, max := r.getMinMax()
	switch mode {
	case LowerBound:
		min = data
	case UpperBound:
		max = data
	}
	if max < min {
		return max
	}
	return min + rand.Intn(max+1-min)
}

func (r *CntRange) Max() int {
	return r.max
}
