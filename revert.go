package revert

import "github.com/eden-framework/courier"

type RevertFunc func(id uint64, meta ...courier.Metadata) error

var revertFuncMapper = map[string]RevertFunc{}

func RegisterRevertFunc(funcID string, revertFunc RevertFunc) {
	revertFuncMapper[funcID] = revertFunc
}

type ResponseRevertID interface {
	GetRevertID() uint64
}

type Revert struct {
	processSequence []string
	processResult   map[string]ResponseRevertID
}

func NewRevert() *Revert {
	return &Revert{
		processResult: map[string]ResponseRevertID{},
	}
}

func (r *Revert) Do(funcID string, handler func() (ResponseRevertID, error)) (err error) {
	r.processSequence = append(r.processSequence, funcID)

	defer func() {
		if err != nil {
			// rollback
			for _, processor := range r.processSequence {
				if revertFunc, ok := revertFuncMapper[processor]; ok {
					if prevResp, ok := r.processResult[processor]; ok {
						_ = revertFunc(prevResp.GetRevertID())
						// TODO retry
					}
				}
			}
		}
	}()

	var resp ResponseRevertID
	resp, err = handler()
	if err == nil && resp != nil {
		r.processResult[funcID] = resp
	}

	return
}
