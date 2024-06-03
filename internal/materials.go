package internal

type Material interface {
	NewObject(o *Scene)
}

type LambertMaterial struct {
	Uuid               string  `json:"uuid" msgpack:"uuid"`
	Type               string  `json:"type" msgpack:"type"`
	Color              int     `json:"color" msgpack:"color"`
	Reflectivity       float32 `json:"reflectivity" msgpack:"reflectivity"`
	Side               int     `json:"side" msgpack:"side"`
	Transparent        bool    `json:"transparent" msgpack:"transparent"`
	Opacity            float32 `json:"opacity" msgpack:"opacity"`
	Linewidth          float32 `json:"linewidth" msgpack:"linewidth"`
	Wireframe          bool    `json:"wireframe" msgpack:"wireframe"`
	WireframeLinewidth float32 `json:"wireframe_linewidth" msgpack:"wireframe_linewidth"`
	VertexColors       int     `json:"vertex_colors" msgpack:"vertex_colors"`
}

const DEFAULT_MATERIAL = "396744fe-087d-11ec-9957-7fbbaaa96777"

func NewLambertMaterial() LambertMaterial {
	return LambertMaterial{
		Uuid:               DEFAULT_MATERIAL,
		Type:               "MeshLambertMaterial",
		Color:              16711935,
		Reflectivity:       0.5,
		Side:               2,
		Transparent:        false,
		Opacity:            0.5,
		Linewidth:          1.0,
		Wireframe:          false,
		WireframeLinewidth: 1.0,
		VertexColors:       0,
	}
}
func (l LambertMaterial) NewObject(o *Scene) {
	o.Materials = append(o.Materials, l)
}
