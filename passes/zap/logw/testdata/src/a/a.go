package a

func Errorw(msg string, v1 ...interface{}) {}

func do() {
	Errorw("msg", "kkk", 123)
	Errorw("msg", "kkk", 123, "kk", 123)
	var t = []int{1, 2, 3}
	Errorw("msg", "kk", t)
	Errorw("msg", "kk")                 // want "invalid pair provided call to a.Errorw"
	Errorw("msg", 333, 33)              // want "none string key in pair provided call to a.Errorw"
	Errorw("msg", "kk", 321, "kk", 321) // want "duplicate key kk in pair provided call to a.Errorw"
}
