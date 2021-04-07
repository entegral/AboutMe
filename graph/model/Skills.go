package model

import (
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	clientmanager "github.com/entegral/aboutme/aws"
	e "github.com/entegral/aboutme/errors"
	"github.com/sirupsen/logrus"
)


type GoSkills struct {
	Frameworks  []*string `json:"frameworks"`
	Paradigms   []*string `json:"paradigms"`
	AwsServices []*string `json:"aws_services"`
	Misc        []*string `json:"misc"`
}

type GoSkillsInput struct {
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Frameworks  []*string `json:"frameworks"`
	Paradigms   []*string `json:"paradigms"`
	AwsServices []*string `json:"aws_services"`
	Misc        []*string `json:"misc"`
}

type JSSkills struct {
	Frameworks  []*string `json:"frameworks"`
	Paradigms   []*string `json:"paradigms"`
	AwsServices []*string `json:"aws_services"`
	Misc        []*string `json:"misc"`
}

type JSSkillsInput struct {
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Frameworks  []*string `json:"frameworks"`
	Paradigms   []*string `json:"paradigms"`
	AwsServices []*string `json:"aws_services"`
	Misc        []*string `json:"misc"`
}

type PythonSkills struct {
	Frameworks []*string `json:"frameworks"`
	Misc       []*string `json:"misc"`
}

type PythonSkillsInput struct {
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Frameworks []*string `json:"frameworks"`
	Misc       []*string `json:"misc"`
}

type Skills struct {
	Js     *JSSkills          `json:"JS"`
	Go     *GoSkills          `json:"Go"`
	Python *PythonSkills      `json:"Python"`
}

type SkillsInput struct {
	FirstName       string        `json:"fName"`
	LastName        string        `json:"lName"`
}

func SkillCompKey(lastname string, firstname string, skillType *string) CompositeKey {
	sk := "skills > skillType:"
	if skillType != nil {
		sk += *skillType
	}
	return CompositeKey{
		Pk: fmt.Sprintf("lName:%s > fName:%s", lastname, firstname),
		Sk: sk,
	}
}

func SkillKey(lastname string, firstname string, skillType *string) map[string]*dynamodb.AttributeValue {
	compKey := SkillCompKey(lastname, firstname, skillType)
	key, err := dynamodbattribute.MarshalMap(compKey)
	if e.Warn("error occurred during marshalling of Skill key",  err) {
		logrus.Warn("skillType", *skillType)
		return nil
	}
	return key
}

func AllSkillsCompKey(lastname string, firstname string) CompositeKey {
	return SkillCompKey(lastname, firstname, nil)
}

var JS string = "JS"
var Go string = "Go"
var Python string = "Python"

func JSSkillsCompKey(lastname string, firstname string) CompositeKey {
	return SkillCompKey(lastname, firstname, &JS)
}
func GoSkillsCompKey(lastname string, firstname string) CompositeKey {
	return SkillCompKey(lastname, firstname, &Go)
}
func PythonSkillsCompKey(lastname string, firstname string) CompositeKey {
	return SkillCompKey(lastname, firstname, &Python)
}

type skillLookup struct {
	Skills	*Skills
	error		error
}

func (input SkillsInput) Get() (*Skills, error) {
	
	var skillTypes = [3]string{"JS", "Go", "Python"}
	var c = make(chan skillLookup, 3)
	var wg sync.WaitGroup
	for _, skillType := range skillTypes {
		wg.Add(1)
		go fetchSkills(input, skillType, c, &wg)
	}
	wg.Wait()
	close(c)
	combinedSkills := Skills{}
	var err error = nil

	for skillResponse := range c {
		if skillResponse.error != nil {
			err = skillResponse.error
		}
		if skillResponse.Skills != nil {
			if skillResponse.Skills.Js != nil {
				combinedSkills.Js = skillResponse.Skills.Js
			}
			if skillResponse.Skills.Go != nil {
				combinedSkills.Go = skillResponse.Skills.Go
			}
			if skillResponse.Skills.Python != nil {
				combinedSkills.Python = skillResponse.Skills.Python
			}
		}
	}
	return &combinedSkills, err
}

func fetchSkills(input SkillsInput, skillType string, c chan<-skillLookup, wg *sync.WaitGroup) {
	defer wg.Done()
	var skillResponse skillLookup
	params := dynamodb.GetItemInput{
		TableName: MeTableName(),
		Key: SkillKey(input.LastName, input.FirstName, &skillType),
	}
	logrus.Info("params", params)
	out, err := clientmanager.ClientMap["default"].Dynamo.GetItem(&params)
	if e.Warn("error in Skill.Get.Query", err) {
		skillResponse.Skills = nil
		skillResponse.error = err
		c <- skillResponse
		return
	}
	logrus.Debug("Skill output item for " + skillType, out.Item)
	if len(out.Item) == 0 {
		skillResponse.Skills = nil
		skillResponse.error = nil
		c <- skillResponse
		return
	}

	logrus.Info("Found out.item", out.Item)
	if skillType == "JS" {
		err = dynamodbattribute.UnmarshalMap(out.Item, skillResponse.Skills)
		if e.Warn("error in Skill.get.UnmarshalMap loop", err) {
			logrus.Warn("error item:", out.Item)
			skillResponse.Skills = nil
			skillResponse.error = err
			c <- skillResponse
			return
		}
	}
	
	logrus.Info("Found out.item", out.Item)
	logrus.Info("Found skill", *skillResponse.Skills)
	c <- skillResponse
	return
}