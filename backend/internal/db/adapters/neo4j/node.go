package neo4j

// Node represents a node in the Neo4j graph database
type Node struct {
	label         string
	properties    map[string]interface{}
	relationships map[string]interface{}
}

func NewNode() *Node {
	return &Node{
		properties:    make(map[string]interface{}),
		relationships: make(map[string]interface{}),
	}
}

func (n *Node) SetLabel(label string) {
	n.label = label
}

func (n *Node) GetLabel() string {
	return n.label
}

func (n *Node) AddProperty(key string, value interface{}) {
	n.properties[key] = value
}

func (n *Node) RemoveProperty(key string) {
	delete(n.properties, key)
}

func (n *Node) AddRelationship(relationship string, value interface{}) {
	n.relationships[relationship] = value
}

func (n *Node) RemoveRelationship(relationship string) {
	delete(n.relationships, relationship)
}

func (n *Node) GetProperties() map[string]interface{} {
	propertiesCopy := make(map[string]interface{})
	for k, v := range n.properties {
		propertiesCopy[k] = v
	}
	return propertiesCopy
}

func (n *Node) GetRelationships() map[string]interface{} {
	relationshipsCopy := make(map[string]interface{})
	for k, v := range n.relationships {
		relationshipsCopy[k] = v
	}
	return relationshipsCopy
}
