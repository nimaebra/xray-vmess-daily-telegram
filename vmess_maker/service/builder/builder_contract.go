package builder

import "xray-telegram/entity"

type Builder struct {
	ServerIP         string
	Setting          entity.Setting
	newVmess         entity.VmessJson
	StringConfigZero string
}

func NewBuilder() *Builder {

	return &Builder{}
}
