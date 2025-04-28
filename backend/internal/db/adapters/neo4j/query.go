package neo4j

import (
	stringutils "backend/internal/utils/string"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// NeoQueryBuilder builds Neo4j queries
type NeoQueryBuilder struct {
	queryType       string
	nodes           []*NodeQueryBuilder
	withReturn      bool
	updateQuery     *UpdateQueryBuilder
	relationBuilder strings.Builder
}

// NewNeoQueryBuilder creates a new NeoQueryBuilder
func NewNeoQueryBuilder(queryType string) *NeoQueryBuilder {
	return &NeoQueryBuilder{
		queryType: queryType,
		nodes:     make([]*NodeQueryBuilder, 0),
	}
}

// WithNode adds a node to the query
func (nqb *NeoQueryBuilder) WithNode(node *Node) *NodeQueryBuilder {
	nodeQueryBuilder := NewNodeQueryBuilder(node, nqb)
	nqb.nodes = append(nqb.nodes, nodeQueryBuilder)
	return nodeQueryBuilder
}

// WithReturn sets whether the query should return values
func (nqb *NeoQueryBuilder) WithReturn(withReturn bool) *NeoQueryBuilder {
	nqb.withReturn = withReturn
	return nqb
}

// WithUpdateQuery creates an update query builder
func (nqb *NeoQueryBuilder) WithUpdateQuery() *UpdateQueryBuilder {
	nqb.updateQuery = NewUpdateQueryBuilder(nqb)
	return nqb.updateQuery
}

// Relation creates a relation query builder
func (nqb *NeoQueryBuilder) Relation(a rune, b rune, relation string, isBiDirectional bool) *RelationQueryBuilder {
	return NewRelationQueryBuilder(a, b, relation, isBiDirectional, nqb, &nqb.relationBuilder)
}

// Build constructs the full Cypher query
func (nqb *NeoQueryBuilder) Build() string {
	var builder strings.Builder
	returnChars := make([]string, 0)

	for _, node := range nqb.nodes {
		builder.WriteString(strings.ToUpper(nqb.queryType))
		builder.WriteString(" (")

		returnChars = append(returnChars, string(node.GetNodeTag()))

		builder.WriteString(fmt.Sprintf("%c:%s", node.GetNodeTag(), strings.ToLower(node.GetNode().GetLabel())))

		if node.IsWithProperties() {
			builder.WriteString("{")
			properties := node.GetNode().GetProperties()
			propCount := 0
			for typeName, typeValue := range properties {
				builder.WriteString(fmt.Sprintf("%s: '%v'", typeName, typeValue))
				propCount++
				if propCount < len(properties) {
					builder.WriteString(", ")
				}
			}
			builder.WriteString("}")
		}

		builder.WriteString(")")

		if node.GetRelationshipMatch() != "" {
			returnChar := strings.ToLower(stringutils.GenerateRandomString(1))
			returnChars = append(returnChars, returnChar)
			builder.WriteString(fmt.Sprintf("-[:%s]->(%s)", node.GetRelationshipMatch(), returnChar))
		}
		builder.WriteString("  ")
	}

	// Remove trailing spaces
	queryString := builder.String()
	queryString = strings.TrimSuffix(queryString, "  ")

	relationQuery := nqb.relationBuilder.String()
	if relationQuery != "" {
		queryString += relationQuery
	}

	if nqb.updateQuery != nil {
		queryString += nqb.updateQuery.GetQuery()
	}

	if nqb.withReturn {
		queryString += " RETURN "
		for i, label := range returnChars {
			queryString += fmt.Sprintf("%s", label)
			if i < len(returnChars)-1 {
				queryString += ", "
			}
		}
	}

	return queryString
}

// NodeQueryBuilder builds node queries
type NodeQueryBuilder struct {
	node              *Node
	nodeTag           rune
	withProperties    bool
	relationshipMatch string
	queryBuilder      *NeoQueryBuilder
}

// NewNodeQueryBuilder creates a new NodeQueryBuilder
func NewNodeQueryBuilder(node *Node, queryBuilder *NeoQueryBuilder) *NodeQueryBuilder {
	rand.Seed(time.Now().UnixNano())
	return &NodeQueryBuilder{
		node:           node,
		nodeTag:        rune('a' + rand.Intn(26)),
		withProperties: false,
		queryBuilder:   queryBuilder,
	}
}

// WithProperties sets whether to include properties
func (nqb *NodeQueryBuilder) WithProperties(withProperties bool) *NodeQueryBuilder {
	nqb.withProperties = withProperties
	return nqb
}

// WithTag sets the node tag
func (nqb *NodeQueryBuilder) WithTag(tag rune) *NodeQueryBuilder {
	nqb.nodeTag = tag
	return nqb
}

// RelatesTo sets the relationship match
func (nqb *NodeQueryBuilder) RelatesTo(relation string) *NeoQueryBuilder {
	nqb.relationshipMatch = relation
	return nqb.queryBuilder
}

// Build returns the query builder
func (nqb *NodeQueryBuilder) Build() *NeoQueryBuilder {
	return nqb.queryBuilder
}

// GetNode returns the node
func (nqb *NodeQueryBuilder) GetNode() *Node {
	return nqb.node
}

// GetNodeTag returns the node tag
func (nqb *NodeQueryBuilder) GetNodeTag() rune {
	return nqb.nodeTag
}

// IsWithProperties returns whether to include properties
func (nqb *NodeQueryBuilder) IsWithProperties() bool {
	return nqb.withProperties
}

// GetRelationshipMatch returns the relationship match
func (nqb *NodeQueryBuilder) GetRelationshipMatch() string {
	return nqb.relationshipMatch
}

// RelationQueryBuilder builds relationship queries
type RelationQueryBuilder struct {
	a               rune
	b               rune
	relation        string
	isBiDirectional bool
	builder         *NeoQueryBuilder
	relationBuilder *strings.Builder
}

// NewRelationQueryBuilder creates a new RelationQueryBuilder
func NewRelationQueryBuilder(a, b rune, relation string, isBiDirectional bool, builder *NeoQueryBuilder, relationBuilder *strings.Builder) *RelationQueryBuilder {
	return &RelationQueryBuilder{
		a:               a,
		b:               b,
		relation:        relation,
		isBiDirectional: isBiDirectional,
		builder:         builder,
		relationBuilder: relationBuilder,
	}
}

// Create adds a create relationship query
func (rqb *RelationQueryBuilder) Create() *NeoQueryBuilder {
	rqb.relationBuilder.WriteString(fmt.Sprintf(" MERGE (%c)-[:%s]->(%c)", rqb.a, rqb.relation, rqb.b))
	if rqb.isBiDirectional {
		rqb.relationBuilder.WriteString(fmt.Sprintf(" MERGE (%c)-[:%s]->(%c)", rqb.b, rqb.relation, rqb.a))
	}
	return rqb.builder
}

// Remove adds a remove relationship query
func (rqb *RelationQueryBuilder) Remove() *NeoQueryBuilder {
	rqb.relationBuilder.WriteString(fmt.Sprintf(" MATCH (%c)-[rel:%s]->(%c) DELETE rel", rqb.a, rqb.relation, rqb.b))
	if rqb.isBiDirectional {
		rqb.relationBuilder.WriteString(fmt.Sprintf(" WITH %c, %c", rqb.a, rqb.b))
		rqb.relationBuilder.WriteString(fmt.Sprintf(" MATCH (%c)-[rel:%s]->(%c) DELETE rel", rqb.b, rqb.relation, rqb.a))
	}
	return rqb.builder
}

// UpdateQueryBuilder builds update queries
type UpdateQueryBuilder struct {
	queryBuilder *NeoQueryBuilder
	builder      strings.Builder
}

// NewUpdateQueryBuilder creates a new UpdateQueryBuilder
func NewUpdateQueryBuilder(queryBuilder *NeoQueryBuilder) *UpdateQueryBuilder {
	return &UpdateQueryBuilder{
		queryBuilder: queryBuilder,
	}
}

// Update adds an update query
func (uqb *UpdateQueryBuilder) Update(oldChar rune, newNode *Node) *UpdateQueryBuilder {
	for key, value := range newNode.GetProperties() {
		uqb.builder.WriteString(" SET ")
		uqb.builder.WriteString(fmt.Sprintf("%c.%s = '%v'", oldChar, key, value))
	}
	return uqb
}

// GetQuery returns the update query string
func (uqb *UpdateQueryBuilder) GetQuery() string {
	return uqb.builder.String()
}

// Build returns the query builder
func (uqb *UpdateQueryBuilder) Build() *NeoQueryBuilder {
	return uqb.queryBuilder
}
