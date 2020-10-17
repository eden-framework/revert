package revert

import "github.com/eden-framework/courier"

type RevertFunc func(id uint64, meta ...courier.Metadata) error

var revertFuncMapper = map[string]RevertFunc{}

func RegisterRevertFunc(funcID string, revertFunc RevertFunc) {
	revertFuncMapper[funcID] = revertFunc
}
