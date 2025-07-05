package myipld

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestNewMyNode(t *testing.T) {
	t.Run("valid data", func(t *testing.T) {
		data := map[string]interface{}{"key": "value"}
		node, err := NewMyNode(data)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		// Verify data marshalling
		var unmarshaled map[string]interface{}
		if err := json.Unmarshal(node.Data, &unmarshaled); err != nil {
			t.Errorf("Data not valid JSON: %v", err)
		}
		if unmarshaled["key"] != "value" {
			t.Errorf("Data corruption. Expected 'value', got '%v'", unmarshaled["key"])
		}

		// Verify CID was computed
		if node.Cid == (MyCID{}) {
			t.Error("CID not computed for new node")
		}

		// Verify rawData cache
		if node.rawData == nil {
			t.Error("rawData not cached after creation")
		}
	})

	t.Run("invalid data", func(t *testing.T) {
		// Channel can't be JSON marshaled
		invalid := make(chan int)
		_, err := NewMyNode(invalid)
		if err == nil {
			t.Fatal("Expected error for invalid data, got nil")
		}
		if !strings.Contains(err.Error(), "failed to marshal") {
			t.Errorf("Unexpected error: %v", err)
		}
	})
}

func TestAddLink(t *testing.T) {
	// Create parent node
	parent, _ := NewMyNode(map[string]string{"type": "parent"})
	initialCID := parent.Cid

	// Create child node
	child, _ := NewMyNode(map[string]string{"type": "child"})

	// Add link
	err := parent.AddLink("child-link", child.Cid)
	if err != nil {
		t.Fatalf("AddLink failed: %v", err)
	}

	// Verify link was added
	if len(parent.Links) != 1 {
		t.Fatalf("Expected 1 link, got %d", len(parent.Links))
	}
	link := parent.Links[0]
	if link.Name != "child-link" || link.Cid != child.Cid {
		t.Errorf("Link mismatch. Expected {child-link %v}, got %v", child.Cid, link)
	}

	// Verify CID changed
	if parent.Cid == initialCID {
		t.Error("CID not updated after adding link")
	}

	// Verify rawData updated
	if parent.rawData == nil {
		t.Fatal("rawData not updated after AddLink")
	}
}

func TestToBytes(t *testing.T) {
	t.Run("with cached data", func(t *testing.T) {
		node, _ := NewMyNode("test")
		originalRaw := node.rawData

		data, err := node.ToBytes()
		if err != nil {
			t.Fatalf("ToBytes failed: %v", err)
		}

		if !compareBytes(data, originalRaw) {
			t.Error("Returned bytes don't match cached rawData")
		}
	})

	t.Run("without cached data", func(t *testing.T) {
		node, _ := NewMyNode("test")
		node.rawData = nil // Force recomputation

		data, err := node.ToBytes()
		if err != nil {
			t.Fatalf("ToBytes failed: %v", err)
		}

		if data == nil {
			t.Fatal("No data returned from ToBytes")
		}

		// Should match recomputed structure
		var reconstructed struct {
			Data  json.RawMessage
			Links []MyLink
		}
		if err := json.Unmarshal(data, &reconstructed); err != nil {
			t.Errorf("Serialized data not valid JSON: %v", err)
		}
	})
}

func TestFromBytes(t *testing.T) {
	original, _ := NewMyNode(map[string]interface{}{"id": 123})
	original.AddLink("self", original.Cid)

	// Serialize node
	data, _ := original.ToBytes()

	// Deserialize
	reconstructed, err := FromBytes(data)
	if err != nil {
		t.Fatalf("FromBytes failed: %v", err)
	}

	// Verify data integrity
	if !compareBytes(reconstructed.Data, original.Data) {
		t.Error("Deserialized data doesn't match original")
	}

	// Verify links integrity
	if len(reconstructed.Links) != 1 {
		t.Fatalf("Expected 1 link, got %d", len(reconstructed.Links))
	}
	if reconstructed.Links[0] != original.Links[0] {
		t.Error("Deserialized link doesn't match original")
	}

	// Verify CID
	if reconstructed.Cid != original.Cid {
		t.Error("Deserialized CID doesn't match original")
	}

	t.Run("invalid data", func(t *testing.T) {
		_, err := FromBytes([]byte("{invalid json}"))
		if err == nil {
			t.Fatal("Expected error for invalid JSON, got nil")
		}
		if !strings.Contains(err.Error(), "failed to unmarshal") {
			t.Errorf("Unexpected error: %v", err)
		}
	})
}

func TestCidConsistency(t *testing.T) {
	node1, _ := NewMyNode("identical")
	node2, _ := NewMyNode("identical")
	if node1.Cid != node2.Cid {
		t.Errorf("CID mismatch for identical nodes:\n%v\nvs\n%v", node1.Cid, node2.Cid)
	}
	node2.AddLink("new-link", node1.Cid)
	if node1.Cid == node2.Cid {
		t.Error("Different nodes have same CID after modification")
	}
}

func compareBytes(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}