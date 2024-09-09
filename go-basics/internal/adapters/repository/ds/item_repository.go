package ds

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	es8 "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/osalomon89/go-basics/internal/core/domain"
)

type elasticSearch struct {
	client *es8.Client
	index  string
}

func NewEsRepository(client *es8.Client) *elasticSearch {
	return &elasticSearch{
		client: client,
	}
}

func (r *elasticSearch) GetAllItems(ctx context.Context, limit int, searchAfter []interface{}) ([]domain.Item, []interface{}, error) {
	query := map[string]interface{}{
		"size": limit,
		"sort": []map[string]string{
			{"created_at": "asc"},
		},
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
	}

	if searchAfter != nil {
		query["search_after"] = searchAfter
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, nil, fmt.Errorf("error encoding query: %s", err)
	}

	req := esapi.SearchRequest{
		Index: []string{r.index},
		Body:  &buf,
	}

	res, err := req.Do(ctx, r.client)
	if err != nil {
		return nil, nil, fmt.Errorf("search request: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, nil, fmt.Errorf("error response: %s", res.String())
	}

	var result struct {
		Hits struct {
			Hits []struct {
				ID     string        `json:"_id"`
				Source domain.Item   `json:"_source"`
				Sort   []interface{} `json:"sort"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, nil, fmt.Errorf("error decoding response: %s", err)
	}

	items := make([]domain.Item, len(result.Hits.Hits))
	for i, hit := range result.Hits.Hits {
		hit.Source.ID = hit.ID
		items[i] = hit.Source
	}

	if len(result.Hits.Hits) == 0 {
		return items, nil, nil
	}

	return items, result.Hits.Hits[len(result.Hits.Hits)-1].Sort, nil
}

func (r *elasticSearch) AddItem(ctx context.Context, item *domain.Item) error {
	item.ID = ""
	now := time.Now().UTC()
	item.CreatedAt = &now
	item.UpdatedAt = &now

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(item); err != nil {
		return fmt.Errorf("error encoding document: %s", err)
	}

	req := esapi.IndexRequest{
		Index: r.index,
		Body:  &buf,
	}

	res, err := req.Do(ctx, r.client)
	if err != nil {
		return fmt.Errorf("insert: request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 409 {
		return errors.New("conflict")
	}

	if res.IsError() {
		return fmt.Errorf("insert: response: %s", res.String())
	}

	var resBody map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		return fmt.Errorf("error parsing the response body: %s", err)
	}

	if id, ok := resBody["_id"].(string); ok {
		item.ID = id
	} else {
		return fmt.Errorf("error: no _id returned in response")
	}

	return nil
}

func (r *elasticSearch) Update(ctx context.Context, itemNew domain.Item) (*domain.Item, error) {
	now := time.Now().UTC()
	itemNew.UpdatedAt = &now
	bdy, err := json.Marshal(itemNew)
	if err != nil {
		return nil, fmt.Errorf("update: marshall: %w", err)
	}

	req := esapi.UpdateRequest{
		Index:      r.index,
		DocumentID: itemNew.ID,
		Body:       bytes.NewReader([]byte(fmt.Sprintf(`{"doc":%s}`, bdy))),
	}

	res, err := req.Do(ctx, r.client)
	if err != nil {
		return nil, fmt.Errorf("update: request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return nil, fmt.Errorf("not found")
	}

	if res.IsError() {
		return nil, fmt.Errorf("update: response: %s", res.String())
	}

	return &itemNew, nil
}

func (r *elasticSearch) Delete(ctx context.Context, id string) error {
	req := esapi.DeleteRequest{
		Index:      r.index,
		DocumentID: id,
	}

	res, err := req.Do(ctx, r.client)
	if err != nil {
		return fmt.Errorf("delete: request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return fmt.Errorf("not found")
	}

	if res.IsError() {
		return fmt.Errorf("delete: response: %s", res.String())
	}

	return nil
}

func (r *elasticSearch) ReadItem(ctx context.Context, id string) (*domain.Item, error) {
	req := esapi.GetRequest{
		Index:      r.index,
		DocumentID: id,
	}

	res, err := req.Do(ctx, r.client)
	if err != nil {
		return nil, fmt.Errorf("error getting document: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return nil, fmt.Errorf("not found")
	}

	if res.IsError() {
		return nil, fmt.Errorf("find one: response: %s", res.String())
	}

	var (
		item domain.Item
		body document
	)
	body.Source = &item

	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("find one: decode: %w", err)
	}

	item.ID = id

	return &item, nil
}

// document represents a single document in Get API response body.
type document struct {
	Source interface{} `json:"_source"`
}
