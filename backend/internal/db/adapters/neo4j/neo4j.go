package neo4j

import (
	"backend/internal/db"
	"backend/internal/utils/concurrency"
	"errors"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

// Container holds the MySQL connection and thread pool. Implements Database
type Container struct {
	database neo4j.Driver
	pool     *concurrency.ThreadPool
}

// Configurations holds MySQL configuration including auth and DB name.
type Configurations struct {
	Authentication *db.URIConfigurations
	DatabaseName   string
}

// GetAuthenticationConfigurations returns a copy of the authentication config.
func (config Configurations) GetAuthenticationConfigurations() db.AuthenticationConfigurations {
	return config.Authentication.GetAuthenticationConfigurations()
}

// Name returns the name of the SQL driver.
func (neoDatabase *Container) Name() string {
	return "neo4j"
}

// Connect initializes the MySQL connection and thread pool.
func (neoDatabase *Container) Connect(config Configurations) error {
	uri := config.Authentication.URI
	username := config.Authentication.AuthConfig.Username
	password := config.Authentication.AuthConfig.Password

	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))

	if err != nil {
		return err
	}

	neoDatabase.database = driver
	//neoDatabase.pool = concurrency.NewThreadPool(5, 5)
	//neoDatabase.pool.Start()
	return nil
}

// GetThreadPool returns the internal thread pool.
func (neoDatabase *Container) GetThreadPool() *concurrency.ThreadPool {
	return neoDatabase.pool
}

// Close closes the database connection
func (neoDatabase *Container) Close() error {
	return neoDatabase.database.Close()
}

// Repository represents a MySQL repository backed by a Container.
type Repository struct {
	Database *Container
}

// NewRepository creates a new Neo4j repository instance
func NewRepository(driver neo4j.Driver) *Repository {
	return &Repository{
		Database: &Container{
			database: driver,
		},
	}
}

// CreateNode creates a new node in the database
func (repo *Repository) CreateNode(node *Node) (*Node, error) {
	session := repo.createSession()
	defer session.Close()

	query := NewNeoQueryBuilder("MERGE").
		WithNode(node).
		WithProperties(true).
		Build().
		WithReturn(true).
		Build()

	result, err := repo.executeQuery(session, query)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, errors.New("no node created")
	}

	return result[0], nil
}

// UpdateNode updates an existing node in the database
func (repo *Repository) UpdateNode(oldNode *Node, newNode *Node) (*Node, error) {
	session := repo.createSession()
	defer session.Close()

	query := NewNeoQueryBuilder("MATCH").
		WithNode(oldNode).
		WithProperties(true).
		WithTag('a').
		Build().
		WithUpdateQuery().
		Update('a', newNode).
		Build().
		WithReturn(true).
		Build()

	result, err := repo.executeQuery(session, query)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, errors.New("no nodes updated")
	}

	return result[0], nil
}

// GetNode retrieves a node from the database
func (repo *Repository) GetNode(node *Node) (*Node, error) {
	session := repo.createSession()
	defer session.Close()

	query := NewNeoQueryBuilder("MATCH").
		WithNode(node).
		WithProperties(true).
		Build().
		WithReturn(true).
		Build()

	result, err := repo.executeQuery(session, query)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	if len(result) > 1 {
		return nil, errors.New("more than one node found")
	}

	return result[0], nil
}

// GetNodeRelations retrieves nodes related to the given node by the specified relation
func (repo *Repository) GetNodeRelations(node *Node, relation string) ([]*Node, error) {
	session := repo.createSession()
	defer session.Close()

	query := NewNeoQueryBuilder("MATCH").
		WithNode(node).
		WithProperties(true).
		MatchRelationship(relation).
		Build().
		WithReturn(true).
		Build()

	fmt.Println(query)

	return repo.executeQuery(session, query)
}

// CreateRelation creates a relationship between two nodes
func (repo *Repository) CreateRelation(nodeA *Node, nodeB *Node, relation string, isBiDirectional bool) error {
	session := repo.createSession()
	defer session.Close()

	query := NewNeoQueryBuilder("MATCH").
		WithNode(nodeA).
		WithProperties(true).
		WithTag('a').
		Build().
		WithNode(nodeB).
		WithProperties(true).
		WithTag('b').
		Build().
		Relation('a', 'b', relation, isBiDirectional).
		Create().
		WithReturn(true).
		Build()

	fmt.Println(query)

	_, err := repo.executeQuery(session, query)
	return err
}

// RemoveRelation removes a relationship between two nodes
func (repo *Repository) RemoveRelation(nodeA *Node, nodeB *Node, relation string, isBiDirectional bool) error {
	session := repo.createSession()
	defer session.Close()

	query := NewNeoQueryBuilder("MATCH").
		WithNode(nodeA).
		WithProperties(true).
		WithTag('a').
		Build().
		WithNode(nodeB).
		WithProperties(true).
		WithTag('b').
		Build().
		Relation('a', 'b', relation, isBiDirectional).
		Remove().
		WithReturn(true).
		Build()

	_, err := repo.executeQuery(session, query)
	return err
}

// executeQuery executes a Cypher query and returns the resulting nodes
func (repo *Repository) executeQuery(session neo4j.Session, query string) ([]*Node, error) {
	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		records, err := tx.Run(query, nil)
		if err != nil {
			return nil, err
		}

		var nodes []*Node
		for records.Next() {
			record := records.Record()
			node := NewNode()

			// Process each key in the record
			for i, _ := range record.Keys {
				value := record.Values[i]
				if neoNode, ok := value.(neo4j.Node); ok {
					// Set the label from neo4j node
					for _, label := range neoNode.Labels {
						node.SetLabel(label)
					}

					// Add properties from neo4j node
					for k, v := range neoNode.Props {
						node.AddProperty(k, v)
					}
				}
			}

			nodes = append(nodes, node)
		}

		return nodes, nil
	})

	if err != nil {
		return nil, err
	}

	if nodes, ok := result.([]*Node); ok {
		return nodes, nil
	}

	return nil, errors.New("failed to parse query result")
}

// createSession creates a new Neo4j session
func (repo *Repository) createSession() neo4j.Session {
	return repo.Database.database.NewSession(neo4j.SessionConfig{
		DatabaseName: "neo4j",
		AccessMode:   neo4j.AccessModeWrite,
	})
}

// Add method to NodeQueryBuilder to support MatchRelationship
func (nqb *NodeQueryBuilder) MatchRelationship(relation string) *NodeQueryBuilder {
	nqb.relationshipMatch = relation
	return nqb
}
