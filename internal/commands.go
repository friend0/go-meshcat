package internal

type Command struct {
	Type string `json:"type" msgpack:"type"`
	Path string `json:"path" msgpack:"path"`
}

type SetObject struct {
	Command
	Object ThreeObject `json:"object" msgpack:"object"`
}

type SetTransform struct {
	Command
	Matrix []float32 `json:"matrix" msgpack:"matrix"`
}

type Delete struct {
	Command
}

type SetAnimation struct {
	Animations interface{}
	Command
	Options AnimationOptions
}

type SetProperty struct {
	Value interface{}
	Command
	SetProperty string
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
