// Code generated by "stringer -type=ActionKind -trimprefix=Action"; DO NOT EDIT.

package ruleset

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ActionNone-0]
	_ = x[ActionAttack-1]
	_ = x[ActionGuard-2]
	_ = x[ActionSkill-3]
	_ = x[ActionMove-4]
}

const _ActionKind_name = "NoneAttackGuardSkillMove"

var _ActionKind_index = [...]uint8{0, 4, 10, 15, 20, 24}

func (i ActionKind) String() string {
	if i < 0 || i >= ActionKind(len(_ActionKind_index)-1) {
		return "ActionKind(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ActionKind_name[_ActionKind_index[i]:_ActionKind_index[i+1]]
}
