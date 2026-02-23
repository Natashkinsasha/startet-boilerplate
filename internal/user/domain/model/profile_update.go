package model

// ProfileUpdate describes a partial update to a profile's JSONB columns.
// It supports both set (overwrite) and increment operations on numeric fields.
type ProfileUpdate struct {
	NumberSets map[string]float64
	NumberIncr map[string]float64
	StringSets map[string]string
}

func NewProfileUpdate() *ProfileUpdate {
	return &ProfileUpdate{}
}

func (u *ProfileUpdate) SetNumber(key string, val float64) *ProfileUpdate {
	if u.NumberSets == nil {
		u.NumberSets = make(map[string]float64)
	}
	u.NumberSets[key] = val
	return u
}

func (u *ProfileUpdate) IncrNumber(key string, delta float64) *ProfileUpdate {
	if u.NumberIncr == nil {
		u.NumberIncr = make(map[string]float64)
	}
	u.NumberIncr[key] = delta
	return u
}

func (u *ProfileUpdate) SetString(key string, val string) *ProfileUpdate {
	if u.StringSets == nil {
		u.StringSets = make(map[string]string)
	}
	u.StringSets[key] = val
	return u
}
