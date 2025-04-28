package neo4j

// Node represents a node in the Neo4j graph database
type Node struct {
	label         string
	properties    map[string]interface{}
	relationships map[string]interface{}
}

// NewNode creates a new Node instance
func NewNode() *Node {
	return &Node{
		properties:    make(map[string]interface{}),
		relationships: make(map[string]interface{}),
	}
}

// SetLabel sets the label for the node
func (n *Node) SetLabel(label string) {
	n.label = label
}

// GetLabel returns the node's label
func (n *Node) GetLabel() string {
	return n.label
}

// AddProperty adds a property to the node
func (n *Node) AddProperty(key string, value interface{}) {
	n.properties[key] = value
}

// RemoveProperty removes a property from the node
func (n *Node) RemoveProperty(key string) {
	delete(n.properties, key)
}

// AddRelationship adds a relationship to the node
func (n *Node) AddRelationship(relationship string, value interface{}) {
	n.relationships[relationship] = value
}

// RemoveRelationship removes a relationship from the node
func (n *Node) RemoveRelationship(relationship string) {
	delete(n.relationships, relationship)
}

// GetProperties returns a copy of the node's properties
func (n *Node) GetProperties() map[string]interface{} {
	propertiesCopy := make(map[string]interface{})
	for k, v := range n.properties {
		propertiesCopy[k] = v
	}
	return propertiesCopy
}

// GetRelationships returns a copy of the node's relationships
func (n *Node) GetRelationships() map[string]interface{} {
	relationshipsCopy := make(map[string]interface{})
	for k, v := range n.relationships {
		relationshipsCopy[k] = v
	}
	return relationshipsCopy
}
