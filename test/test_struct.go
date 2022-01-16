package teststruct

//go:generate vstruct -struct=Simples
type Simples struct {
	Id      uint32 `vstruct:"ta:tb"`
	NameLen uint8
	Name    []uint8 `vstruct:"repeat:NameLen"`
}

//go:generate vstruct -struct=UnionSim
type UnionSim struct {
	Len uint8
	Arr []Simples `vstruct:"repeat:Len,refer:true"`
	Crc uint16
}
