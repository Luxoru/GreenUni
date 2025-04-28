package neo4j

import (
	stringutils "backend/internal/utils/string"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type NeoQueryBuilder struct {
	queryType       string
	nodes           []*NodeQueryBuilder
	withReturn      bool
	updateQuery     *UpdateQueryBuilder
	relationBuilder strings.Builder
}

func NewNeoQueryBuilder(queryType string) *NeoQueryBuilder {
	return &NeoQueryBuilder{
		queryType: queryType,
		nodes:     make([]*NodeQueryBuilder, 0),
	}
}

func (nqb *NeoQueryBuilder) WithNode(node *Node) *NodeQueryBuilder {
	nodeQueryBuilder := NewNodeQueryBuilder(node, nqb)
	nqb.nodes = append(nqb.nodes, nodeQueryBuilder)
	return nodeQueryBuilder
}

func (nqb *NeoQueryBuilder) WithReturn(withReturn bool) *NeoQueryBuilder {
	nqb.withReturn = withReturn
	return nqb
}

func (nqb *NeoQueryBuilder) WithUpdateQuery() *UpdateQueryBuilder {
	nqb.updateQuery = NewUpdateQueryBuilder(nqb)
	return nqb.updateQuery
}

func (nqb *NeoQueryBuilder) Relation(a rune, b rune, relation string, isBiDirectional bool) *RelationQueryBuilder {
	return NewRelationQueryBuilder(a, b, relation, isBiDirectional, nqb, &nqb.relationBuilder)
}

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

type NodeQueryBuilder struct {
	node              *Node
	nodeTag           rune
	withProperties    bool
	relationshipMatch string
	queryBuilder      *NeoQueryBuilder
}

func NewNodeQueryBuilder(node *Node, queryBuilder *NeoQueryBuilder) *NodeQueryBuilder {
	rand.Seed(time.Now().UnixNano())
	return &NodeQueryBuilder{
		node:           node,
		nodeTag:        rune('a' + rand.Intn(26)),
		withProperties: false,
		queryBuilder:   queryBuilder,
	}
}

func (nqb *NodeQueryBuilder) WithProperties(withProperties bool) *NodeQueryBuilder {
	nqb.withProperties = withProperties
	return nqb
}

func (nqb *NodeQueryBuilder) WithTag(tag rune) *NodeQueryBuilder {
	nqb.nodeTag = tag
	return nqb
}

func (nqb *NodeQueryBuilder) RelatesTo(relation string) *NeoQueryBuilder {
	nqb.relationshipMatch = relation
	return nqb.queryBuilder
}

func (nqb *NodeQueryBuilder) Build() *NeoQueryBuilder {
	return nqb.queryBuilder
}

func (nqb *NodeQueryBuilder) GetNode() *Node {
	return nqb.node
}

func (nqb *NodeQueryBuilder) GetNodeTag() rune {
	return nqb.nodeTag
}

func (nqb *NodeQueryBuilder) IsWithProperties() bool {
	return nqb.withProperties
}

func (nqb *NodeQueryBuilder) GetRelationshipMatch() string {
	return nqb.relationshipMatch
}

type RelationQueryBuilder struct {
	a               rune
	b               rune
	relation        string
	isBiDirectional bool
	builder         *NeoQueryBuilder
	relationBuilder *strings.Builder
}

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

func (rqb *RelationQueryBuilder) Create() *NeoQueryBuilder {
	rqb.relationBuilder.WriteString(fmt.Sprintf(" MERGE (%c)-[:%s]->(%c)", rqb.a, rqb.relation, rqb.b))
	if rqb.isBiDirectional {
		rqb.relationBuilder.WriteString(fmt.Sprintf(" MERGE (%c)-[:%s]->(%c)", rqb.b, rqb.relation, rqb.a))
	}
	return rqb.builder
}

func (rqb *RelationQueryBuilder) Remove() *NeoQueryBuilder {
	rqb.relationBuilder.WriteString(fmt.Sprintf(" MATCH (%c)-[rel:%s]->(%c) DELETE rel", rqb.a, rqb.relation, rqb.b))
	if rqb.isBiDirectional {
		rqb.relationBuilder.WriteString(fmt.Sprintf(" WITH %c, %c", rqb.a, rqb.b))
		rqb.relationBuilder.WriteString(fmt.Sprintf(" MATCH (%c)-[rel:%s]->(%c) DELETE rel", rqb.b, rqb.relation, rqb.a))
	}
	return rqb.builder
}

type UpdateQueryBuilder struct {
	queryBuilder *NeoQueryBuilder
	builder      strings.Builder
}

func NewUpdateQueryBuilder(queryBuilder *NeoQueryBuilder) *UpdateQueryBuilder {
	return &UpdateQueryBuilder{
		queryBuilder: queryBuilder,
	}
}

func (uqb *UpdateQueryBuilder) Update(oldChar rune, newNode *Node) *UpdateQueryBuilder {
	for key, value := range newNode.GetProperties() {
		uqb.builder.WriteString(" SET ")
		uqb.builder.WriteString(fmt.Sprintf("%c.%s = '%v'", oldChar, key, value))
	}
	return uqb
}

func (uqb *UpdateQueryBuilder) GetQuery() string {
	return uqb.builder.String()
}

func (uqb *UpdateQueryBuilder) Build() *NeoQueryBuilder {
	return uqb.queryBuilder
}
