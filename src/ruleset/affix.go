package ruleset

type EffectPrefix int

const (
	PrefixNone EffectPrefix = iota
	PrefixExpensive
)

type EffectSuffix int

const (
	SuffixNone EffectSuffix = iota
)
