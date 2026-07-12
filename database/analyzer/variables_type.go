package analyzer

type VariablesType int

const (
	VariablesTypeGlobal VariablesType = iota
	VariablesTypeLocal
	VariablesTypeSession
)

func (v VariablesType) String() string {
	switch v {
	case VariablesTypeGlobal:
		return "GLOBAL"
	case VariablesTypeLocal:
		return "LOCAL"
	case VariablesTypeSession:
		return "SESSION"
	default:
		return ""
	}
}
