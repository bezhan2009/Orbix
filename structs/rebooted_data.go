package structs

type RebootedData struct {
	Username          string
	Prefix            string
	Recover           any
	LoopData          OrbixLoopData
	LoadUserConfigsFn func(echo bool) error
}
