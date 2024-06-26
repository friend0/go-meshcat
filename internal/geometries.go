package internal

import (
	"fmt"
	"log"
	"os"
	"path"
	"reflect"

	"github.com/google/uuid"
)

// Define the
type SceneMetadata struct {
	Type    string  `json:"type" msgpack:"type"`
	Version float32 `json:"version" msgpack:"version"`
}

type SceneElement struct {
	Uuid string `json:"uuid" msgpack:"uuid"`
	Type string `json:"type" msgpack:"type"`
}

type GenericGeom map[string]interface{}

func (g GenericGeom) get_element() SceneElement {
	if g["uuid"] == nil {
		g["uuid"] = uuid.NewString()
	}
	if g["type"] == nil {
		if g["shape"] != nil {
			g["type"], _ = g["shape"].(string)
		} else {
			if _, ok := g["radius"]; ok {
				g["type"] = "SphereGeometry"
			} else if _, ok := g["width"]; ok {
				g["type"] = "BoxGeometry"
			} else {
				g["type"] = "_meshfile_geometry"
			}
		}
	}
	return SceneElement{
		Uuid: g["uuid"].(string),
		Type: g["type"].(string),
	}
}

func (g GenericGeom) get_matrix() []float32 {
	// assume position comes in as [x, y, z]
	fmt.Printf("in matrix func %v\n", g)
	fmt.Println(g["position"])
	position, ok := g["position"]
	if !ok {
		fmt.Println("Position not found")
		x, ok := g["x"].(float32)
		if !ok {
			x = 0
		}
		y, ok := g["y"].(float32)
		if !ok {
			y = 0
		}
		z, ok := g["z"].(float32)
		if !ok {
			z = 0
		}
		return []float32{1, 0, 0, x, 0, 1, 0, y, 0, 0, 1, z, 0, 0, 0, 1}
	} else {
		if position, err := interfaceToFloatSlice(position); err == nil {
			// Assign the converted value back to the map
			return []float32{1, 0, 0, float32(position[0]), 0, 1, 0, float32(position[1]), 0, 0, 1, float32(position[2]), 0, 0, 0, 1}
		} else {
			fmt.Println("Error converting:", err)
			return []float32{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}
		}
	}
}

// interfaceToFloatSlice converts an interface to a slice of float64
func interfaceToFloatSlice(val interface{}) ([]float64, error) {
	// Use reflection to check if the value is a slice
	v := reflect.ValueOf(val)
	if v.Kind() != reflect.Slice {
		return nil, fmt.Errorf("value is not a slice")
	}

	// Iterate through the slice and convert elements to float64
	floatSlice := make([]float64, v.Len())
	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i).Interface()
		floatVal, ok := elem.(float64)
		if !ok {
			return nil, fmt.Errorf("element %v is not a float32", elem)
		}
		floatSlice[i] = floatVal
	}

	return floatSlice, nil
}

func (geom GenericGeom) init_element() error {
	_type, ok := geom["type"].(string)
	if !ok {
		// todo: determine the type base on attributes
		return fmt.Errorf("Geometry type not found")
	}
	scene_element := SceneElement{
		Uuid: uuid.NewString(),
		Type: _type,
	}
	geom["uuid"] = scene_element.Uuid
	geom["type"] = scene_element.Type
	return nil
}

// Object is a like join, where the tables being joined are the elements
// of Scene.Geometries, and Scene.Materials, respectively.
// The Scene object itself also gets a UUID, and this is where the transformation matrix is spefied.
type Object struct {
	SceneElement
	GeometryUUID string    `json:"geometry" msgpack:"geometry"`
	MaterialUUID string    `json:"material" msgpack:"material"`
	Matrix       []float32 `json:"matrix" msgpack:"matrix,omitempty"`
}

// ThreeObject contains the geometries and materials that have been defined on the
type ThreeObject struct {
	Metadata   SceneMetadata `json:"metadata" msgpack:"metadata"`
	Geometries []Geometry    `json:"geometries" msgpack:"geometries"`
	Materials  []Material    `json:"materials" msgpack:"materials"`
	Object     Object        `json:"object" msgpack:"object"`
}

func NewScene() ThreeObject {
	return ThreeObject{
		Metadata:   default_scene_metadata(),
		Geometries: []Geometry{},
		Materials:  []Material{},
		Object:     Object{},
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
	get_matrix() []float32
	init_element() error
}

type Box struct {
	SceneElement
	Width  float32 `json:"width" msgpack:"width"`
	Height float32 `json:"height" msgpack:"height"`
	Depth  float32 `json:"depth" msgpack:"depth"`
}

func (b *Box) get_element() SceneElement {
	return b.SceneElement
}

func (b *Box) init_element() error {
	b.SceneElement = SceneElement{
		Uuid: uuid.NewString(),
		Type: "BoxGeometry",
	}
	return nil
}

func (b *Box) get_matrix() []float32 {
	return []float32{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}
}

func NewBox(width, height, depth float32) *Box {
	return &Box{
		SceneElement: SceneElement{
			Uuid: uuid.NewString(),
			Type: "BoxGeometry",
		},
		Width:  width,
		Height: height,
		Depth:  depth,
	}
}

type Sphere struct {
	SceneElement
	Radius float32 `json:"radius" msgpack:"radius"`
}

func NewSphere(radius float32, id string) Sphere {
	if id == "" {
		id = uuid.NewString()
	}
	return Sphere{
		SceneElement: SceneElement{
			Uuid: id,
			Type: "SphereGeometry",
		},
		Radius: radius,
	}
}

func (s *Sphere) get_element() SceneElement {
	return s.SceneElement
}

func (s *Sphere) init_element() error {
	s.SceneElement = SceneElement{
		Uuid: uuid.NewString(),
		Type: "SphereGeometry",
	}
	return nil
}

func (s *Sphere) get_matrix() []float32 {
	return []float32{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}
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
	return MeshGeometry{
		SceneElement: SceneElement{
			Uuid: "cef79e52-526d-4263-b595-04fa2705974e",
			Type: "_meshfile_geometry",
		},
		Format: "stl",
		Data:   data,
	}, nil
}

func Objectify[T Geometry](g T) ThreeObject {
	scene_element := g.get_element()
	obj := NewScene()
	obj.Object.GeometryUUID = scene_element.Uuid
	obj.Object.MaterialUUID = DEFAULT_MATERIAL
	obj.Object.Type = "Mesh"
	obj.Object.Uuid = scene_element.Uuid
	obj.Object.Matrix = []float32{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}
	obj.Geometries = []Geometry{g}
	obj.Materials = []Material{NewLambertMaterial()}
	return obj
}
