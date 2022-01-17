package teststruct

//go:generate vstruct -struct=Simples,UnionSim,UnionSim2,UnionSimAcc

func lenCov(len uint8) uint8 {
	return len & 0x3F
}

type Simples struct {
	Id      uint32 `vstruct:"ta:tb"`
	NameLen uint8
	Name    []uint8 `vstruct:"repeat:NameLen,access:lenCov"`
}

type UnionSim struct {
	Len uint8
	Arr []Simples `vstruct:"repeat:Len,refer:true"`
	Crc uint16
}
