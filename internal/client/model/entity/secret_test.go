package entity

import (
	"testing"
)

func TestSecretData(t *testing.T) {
	// Create an instance of Secret
	secret := &Secret{}

	// Create an instance of SecretData to set
	originalData := SecretData{
		Name:  "TestSecret",
		Value: Password{Pass: "TestPassword"},
		Metadata: map[string]string{
			"env": "test",
		},
	}

	// Set the data
	err := secret.SetData(originalData)
	if err != nil {
		t.Fatalf("Failed to set data: %v", err)
	}

	// Get the data
	retrievedData, err := secret.GetData()
	if err != nil {
		t.Fatalf("Failed to get data: %v", err)
	}

	// Check if the retrieved data matches the original data
	if retrievedData.Name != originalData.Name {
		t.Errorf("Expected Name %s, got %s", originalData.Name, retrievedData.Name)
	}
	if retrievedData.Value != originalData.Value {
		t.Errorf("Expected Value %s, got %s", originalData.Value, retrievedData.Value)
	}
	if len(retrievedData.Metadata) != len(originalData.Metadata) {
		t.Errorf("Expected Metadata length %d, got %d", len(originalData.Metadata), len(retrievedData.Metadata))
	}
	for key, value := range originalData.Metadata {
		if retrievedData.Metadata[key] != value {
			t.Errorf("Expected Metadata[%s] %s, got %s", key, value, retrievedData.Metadata[key])
		}
	}
}
