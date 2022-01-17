package teststruct

//go:generate vstruct -struct=Simples,UnionSim,UnionSim2

type Simples struct {
	Id      uint32 `vstruct:"ta:tb"`
	NameLen uint8
	Name    []uint8 `vstruct:"repeat:NameLen"`
}

type UnionSim struct {
	Len uint8
	Arr []Simples `vstruct:"repeat:Len,refer:true"`
	Crc uint16
}
