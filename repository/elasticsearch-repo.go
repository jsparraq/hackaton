package repository

import (
	"context"
	"fmt"
	"crypto/rand"
	"log"
	"time"
	"reflect"

	"github.com/jsparraq/api-rest/entity"
	"github.com/olivere/elastic/v6"
)

type repo struct{}

var (
	elasticClient *elastic.Client
	err           error
)

const (
	indexName string = "poster"
	indexType string = "post"
)
const mapping = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings":{
		"post":{
			"properties":{
				"message":{
					"type":"text"
				},
				"created":{
					"type":"date"
				}
			}
		}
	}
}`

// NewElasticsearchRepository function
func NewElasticsearchRepository() PostRepository {
	return &repo{}
}

func (*repo) Save(post *entity.Post) (*entity.Post, error) {
	ctx := context.Background()

	for {
		elasticClient, err = elastic.NewClient(
			elastic.SetURL("http://127.0.0.1:9200"),
			elastic.SetSniff(false),
			elastic.SetHealthcheckInterval(10*time.Second),
		)
		if err != nil {
			log.Println(err)
			time.Sleep(3 * time.Second)
		} else {
			break
		}
	}

	exists, err := elasticClient.IndexExists(indexName).Do(ctx)
	if err != nil {
		panic(err)
	}

	if !exists {
		_, err := elasticClient.CreateIndex(indexName).BodyString(mapping).Do(ctx)
		if err != nil {
			panic(err)
		}
	}

	p, _ := rand.Prime(rand.Reader, 64)

	_, err = elasticClient.Index().
		Index(indexName).
		Type(indexType).
		Id(p.String()).
		BodyJson(post).
		Do(ctx)
	if err != nil {
		panic(err)
	}

	return post, nil
}

func (*repo) FindAll() ([]entity.Post, error) {
	ctx := context.Background()

	for {
		elasticClient, err = elastic.NewClient(
			elastic.SetURL("http://127.0.0.1:9200"),
			elastic.SetSniff(false),
			elastic.SetHealthcheckInterval(10*time.Second),
		)
		if err != nil {
			log.Println(err)
			time.Sleep(3 * time.Second)
		} else {
			break
		}
	}

	searchResult, err := elasticClient.Search().
		Index(indexName).   
		Sort("created", true). 
		From(0).Size(1000).  
		Pretty(true).       
		Do(ctx) 
	if err != nil {
		panic(err)
	}

	fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)

	var posts []entity.Post
	var ttyp entity.Post
	for _, item := range searchResult.Each(reflect.TypeOf(ttyp)) {
		if t, ok := item.(entity.Post); ok {
			posts = append(posts, t)
		}
	}

	return posts, nil
}
