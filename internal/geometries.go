package internal

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/google/uuid"
)

// Define the
type SceneMetadata struct {
	Version float32 `json:"version" msgpack:"version"`
	Type    string  `json:"type" msgpack:"type"`
}

type SceneElement struct {
	Uuid string `json:"uuid" msgpack:"uuid"`
	Type string `json:"type" msgpack:"type"`
}

// SceneObject is a like join, where the tables being joined are the elements
// of Scene.Geometries, and Scene.Materials, respectively.
// The Scene object itself also gets a UUID, and this is where the transformation matrix is spefied.
type SceneObject struct {
	SceneElement
	GeometryUUID string    `json:"geometry" msgpack:"geometry"`
	MaterialUUID string    `json:"material" msgpack:"material"`
	Matrix       []float32 `json:"matrix" msgpack:"matrix,omitempty"`
}

// Scene contains the geometries and materials that have been defined on the
type Scene struct {
	Metadata   SceneMetadata `json:"metadata" msgpack:"metadata"`
	Geometries []Geometry    `json:"geometries" msgpack:"geometries"`
	Materials  []Material    `json:"materials" msgpack:"materials"`
	Object     SceneObject   `json:"object" msgpack:"object"`
}

func NewScene() Scene {
	return Scene{
		Metadata:   default_scene_metadata(),
		Geometries: []Geometry{},
		Materials:  []Material{},
		Object:     SceneObject{},
	}
}
func default_scene_metadata() SceneMetadata {
	return SceneMetadata{
		Version: 4.5,
		Type:    "Object",
	}
}

type Geometry interface {
	get_element() SceneElement
}

type Box struct {
	SceneElement
	Width  float32 `json:"width" msgpack:"width"`
	Height float32 `json:"height" msgpack:"height"`
	Depth  float32 `json:"depth" msgpack:"depth"`
}

func (b Box) get_element() SceneElement {
	return b.SceneElement
}

func NewBox(width, height, depth float32) Box {
	return Box{
		SceneElement: SceneElement{
			Uuid: uuid.NewString(),
			Type: "BoxGeometry",
		},
		Width:  width,
		Height: height,
		Depth:  depth,
	}
}

type MeshGeometry struct {
	SceneElement
	Format string  `json:"format" msgpack:"format"`
	Data   []uint8 `json:"data" msgpack:"data"`
}

func (m MeshGeometry) get_element() SceneElement {
	return m.SceneElement
}

func NewStarling(x, y, z float64) (MeshGeometry, error) {
	wd, _ := os.Getwd()
	data, err := os.ReadFile(path.Join(wd, "/web/meshcat/data/starling1.stl"))
	if err != nil {
		log.Fatal(err)
		return MeshGeometry{}, err
	}
	fmt.Println(data)
	return MeshGeometry{
		SceneElement: SceneElement{
			Uuid: "cef79e52-526d-4263-b595-04fa2705974e",
			Type: "_meshfile_geometry",
		},
		Format: "stl",
		Data:   data,
	}, nil

}

func Objectify[T Geometry](g T) Scene {
	scene_element := g.get_element()
	obj := NewScene()
	obj.Object.GeometryUUID = scene_element.Uuid
	obj.Object.MaterialUUID = DEFAULT_MATERIAL
	obj.Object.Type = "Mesh"
	obj.Object.Uuid = scene_element.Uuid
	// obj.Object.Matrix = []float32{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}
	obj.Geometries = []Geometry{g}
	obj.Materials = []Material{NewLambertMaterial()}
	return obj
}
