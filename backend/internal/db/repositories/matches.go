package repositories

import (
	"backend/internal/db/adapters/neo4j"
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type MatchesRepository struct {
	Repository *neo4j.Repository
}

func NewMatchesRepository(db *neo4j.Repository) (*MatchesRepository, error) {
	ur := &MatchesRepository{}
	ur.Repository = db
	return ur, nil
}

func (matches *MatchesRepository) CreateMatch(uuid1 string, uuid2 string) error {

	node1 := neo4j.NewNode()
	node1.AddProperty("uuid", uuid1)
	node1.SetLabel("User")
	node2 := neo4j.NewNode()
	node2.AddProperty("uuid", uuid2)
	node2.SetLabel("User")

	n1, err := matches.Repository.CreateNode(node1)
	if err != nil {
		log.Error(err)
		return err
	}

	n2, err := matches.Repository.CreateNode(node2)
	if err != nil {
		log.Error(err)
		return err
	}

	err = matches.Repository.CreateRelation(n1, n2, "MATCH", true)
	if err != nil {
		return err
	}

	return nil
}

func (matches *MatchesRepository) GetMatches(userID string) (*[]uuid.UUID, error) {
	node := neo4j.NewNode()
	node.SetLabel("User")
	node.AddProperty("uuid", userID)

	node, err := matches.Repository.GetNode(node)
	if err != nil {
		return nil, err
	}

	if node == nil {
		fmt.Println("NIGGER AINT THERE")
		return nil, nil //User doesn't exist on here
	}

	relations, err := matches.Repository.GetNodeRelations(node, "MATCH")
	if err != nil {
		return nil, err
	}

	var allIds []uuid.UUID

	for _, node := range relations {
		u := node.GetProperties()["uuid"]
		str, ok := u.(string)
		if ok {

			parsed, err := uuid.Parse(str)

			if err != nil {
				//Should never happen
				continue
			}
			allIds = append(allIds, parsed)

		} else {
			fmt.Println("Not a string: ", u)
		}
	}

	return &allIds, nil

}
