package registrar

import (
	"fmt"
	"log"
)

type registry struct {
	subscriberContact map[string]string
}

var subDB registry

func init() {
	subDB = registry{}
	subDB.subscriberContact = make(map[string]string)
}

// RegisterSubscriber stores/updates the contact of the subscriber in the registry
func RegisterSubscriber(sub string, newContact string) {
	if currentContact, ok := subDB.subscriberContact[sub]; ok {
		log.Printf("[REGISTRY] Updating existing subscriber: %v", sub)
		log.Printf("[REGISTRY] %v -> %v", currentContact, newContact)
	} else {
		log.Printf("[REGISTRY] New subscriber: %v", sub)
	}

	subDB.subscriberContact[sub] = newContact
}

// DeregisterSubscriber removes a subscriber from the registry
func DeregisterSubscriber(sub string) {
	log.Printf("[REGISTRY] Deregistering subscriber: %v", sub)
	delete(subDB.subscriberContact, sub)
}

// Contact retrieves the contact information of the subscriber from the registry
// returns empty string if there is no such subscriber
func Contact(sub string) (string, error) {
	if contact, ok := subDB.subscriberContact[sub]; ok {
		return contact, nil
	}

	return "", fmt.Errorf("subscriber is not found in registry: %v", sub)
}
