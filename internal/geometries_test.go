package internal

import (
	"bytes"
	"testing"

	"github.com/vmihailenco/msgpack/v5"
)

func TestGeometrySerialization(t *testing.T) {
	original := &Box{
		SceneElement: SceneElement{
			Uuid: "3d984587-36fc-404e-bc68-3a25574cdb8b",
			Type: "BoxGeometry",
		},
		Width:  8,
		Height: 9,
		Depth:  10,
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
}
