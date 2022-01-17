package teststruct

type UnionSim2 struct {
	Len uint8
	Arr Simples `vstruct:"refer:true"`
	Crc uint16  `vstruct:"crc16:true"`
}
