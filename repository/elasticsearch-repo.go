package repository

import (
	"context"
	"fmt"
	"crypto/rand"
	"log"
	"time"
	"reflect"
	"math/big"
	"strconv"

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
			elastic.SetURL("http://elastic:elastic@ec2-3-238-134-215.compute-1.amazonaws.com:9200/"),
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

	bulk := elasticClient.Bulk()
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	for i := 0; i < 100; i++ {
		
		
		b := make([]rune, 10)
		for i := range b {
			number, _ := rand.Int(rand.Reader, big.NewInt(52))
			number2, _ := strconv.Atoi(number.String())
			b[i] = letters[number2]
		}

		index, _ := rand.Prime(rand.Reader, 64)
		index_string := index.String()

		var postNew = entity.Post{Message: post.Message, Created: time.Now() }
		postNew.Keys = map[string]string{
			string(b):"asdf",
		}
	

		req := elastic.NewBulkIndexRequest()
		req.OpType("index") // set type to "index" document
		req.Type(indexType)
		req.Index(indexName)
		req.Id(index_string)
		req.Doc(postNew)

		bulk = bulk.Add(req)
		
	}
	log.Println("NewBulkIndexRequest().NumberOfActions():", bulk.NumberOfActions())
	bulkResp, err := bulk.Do(ctx)
	
	if err != nil {
		panic(err)
	}

	indexed := bulkResp.Indexed()

	log.Println("nbulkResp.Indexed():", indexed)
	log.Println("bulkResp.Indexed() TYPE:", reflect.TypeOf(indexed))

	t := reflect.TypeOf(indexed)

	log.Println("nt:", t)
	log.Println("NewBulkIndexRequest().NumberOfActions():", bulk.NumberOfActions())

	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		log.Println("nbulkResp.Indexed() METHOD NAME:", i, method.Name)
		log.Println("bulkResp.Indexed() method:", method)
	}

	return post, nil
}

func (*repo) FindAll() ([]entity.Post, error) {
	ctx := context.Background()

	for {
		elasticClient, err = elastic.NewClient(
			elastic.SetURL("http://elastic:elastic@ec2-3-238-134-215.compute-1.amazonaws.com:9200/"),
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
