package constants

type SportType string

const (
	Running  SportType = "running"
	Cycling  SportType = "cycling"
	Swimming SportType = "swimming"
	Strength SportType = "strength"
	Yoga     SportType = "yoga"
	JumpRope SportType = "jump_rope"
)

var METValues = map[SportType]float64{
	Running:  9.8,
	Cycling:  7.5,
	Swimming: 8.0,
	Strength: 6.0,
	Yoga:     3.0,
	JumpRope: 12.0,
}
