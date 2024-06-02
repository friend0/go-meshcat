package internal

import (
	"github.com/google/uuid"
)

type Geometry interface{}

type Material interface{}

type Box struct {
	Uuid   string  `json:"uuid" msgpack:"uuid"`
	Type   string  `json:"type" msgpack:"type"`
	Width  float32 `json:"width" msgpack:"width"`
	Height float32 `json:"height" msgpack:"height"`
	Depth  float32 `json:"depth" msgpack:"depth"`
}

func (o *Object) metadata() {
	o.Version = "4.5"
	o.Type = "Object"
}

func (b *Box) hydrateObject() Object {
	b.Uuid = uuid.New().String()
	b.Type = "BoxGeometry"

	obj := Object{}
	obj.metadata()
	obj.Object.Type = b.Type
	obj.Object.Uuid = b.Uuid
	obj.Object.Matrix = []float32{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}
	obj.Geometries = append(obj.Geometries, b)

	return obj
}

type ObjectMetadata struct {
	Version string `json:"version" msgpack:"version"`
	Type    string `json:"type" msgpack:"type"`
}

type ObjectParameters struct {
	Uuid         string    `json:"uuid" msgpack:"uuid"`
	Type         string    `json:"type" msgpack:"type"`
	GeometryUUID string    `json:"geometry" msgpack:"geometry"`
	MaterialUUID string    `json:"material" msgpack:"material"`
	Matrix       []float32 `json:"matrix" msgpack:"matrix"`
}

type Object struct {
	ObjectMetadata
	Geometries []Geometry       `json:"geometry" msgpack:"geometry"`
	Materials  []interface{}    `json:"material" msgpack:"material"`
	Object     ObjectParameters `json:"object" msgpack:"object"`
}
