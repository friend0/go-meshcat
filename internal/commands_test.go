package internal

import (
	"bytes"
	"log"
	"testing"

	"github.com/vmihailenco/msgpack/v5"
)

func TestSetObjectSerialization(t *testing.T) {
	original := &SetObject{
		Command: Command{
			Type: "set_object",
			Path: "/tmp/geom",
		},
		Object: NewBox(1, 1, 1),
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
	// if original.Name != deserialized.Name {
	// 	t.Errorf("expected Name %s, got %s", original.Name, deserialized.Name)
	// }
	// if !bytes.Equal(original.Data, deserialized.Data) {
	// 	t.Errorf("expected Data %v, got %v", original.Data, deserialized.Data)
	// }
}
