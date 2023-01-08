package ruleset

import (
	"fmt"
	"io"

	"github.com/quasilyte/gmath"
)

type Dice struct {
	r     *gmath.Rand
	debug io.Writer
	rolls int
}

func NewDice(r *gmath.Rand, debug io.Writer) *Dice {
	return &Dice{r: r, debug: debug}
}

func (d *Dice) Roll1d6(where, who, why string) int {
	d.rolls++

	result := d.r.IntRange(0, 5)
	if d.debug != nil {
		fmt.Fprintf(d.debug, "[%s] %s rolls 1d6 (%s): %d\n", where, who, why, result+1)
	}

	return result
}
