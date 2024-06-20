package internal

import (
	"bytes"
	"log"
	"testing"

	"github.com/vmihailenco/msgpack/v5"
)

func TestSetObjectSerialization(t *testing.T) {

	original := SetObject{
		Command: Command{
			Type: "set_object",
			Path: "/vehicles/vehicle_0",
		},
		Object: Scene{
			Metadata: SceneMetadata{
				Version: 4.5,
				Type:    "BufferGeometry",
			},
			Geometries: []Geometry{},
			Materials:  []Material{},
			Object:     SceneObject{},
		},
	}

	// Serialize the object to MessagePack
	var buf bytes.Buffer
	encoder := msgpack.NewEncoder(&buf)
	err := encoder.Encode(original)
	if err != nil {
		t.Fatalf("failed to encode: %v", err)
	}

	// Deserialize the object from MessagePack
	var deserialized SetObject
	decoder := msgpack.NewDecoder(&buf)
	err = decoder.Decode(&deserialized)
	if err != nil {
		t.Fatalf("failed to decode: %v", err)
	}

	// Compare the original and deserialized objects
	if original.Type != deserialized.Type {
		t.Errorf("expected Type %s, got %s", original.Type, deserialized.Type)
	}
	log.Printf("%#v", buf)

}
