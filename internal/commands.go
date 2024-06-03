package internal

type Command struct {
	Type string `json:"type" msgpack:"type"`
	Path string `json:"path" msgpack:"path"`
}

type SetObject struct {
	Command
	Object Scene `json:"object" msgpack:"object"`
}

type SetTransform struct {
	Command
	Matrix []float32 `json:"matrix" msgpack:"matrix"`
}

type Delete struct {
	Command
}

type SetAnimation struct {
	Command
	Animations interface{}
	Options    AnimationOptions
}

type SetProperty struct {
	Command
	SetProperty string
	Value       interface{}
}

type CaptureImage struct {
	Command
	Xres int `json:"xres" msgpack:"xres"`
	Yres int `json:"yres" msgpack:"yres"`
}

type AnimationOptions struct {
	Play        bool
	Repetitions int
}
