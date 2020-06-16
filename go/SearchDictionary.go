package odict

import (
	"encoding/json"
	"os"

	"github.com/blevesearch/bleve"
)

func SearchDictionary(dictionary Dictionary, queryStr string) []DictionaryEntry {
	indexPath := getIndexPath(dictionary)
	_, err := os.Stat(indexPath)

	if os.IsNotExist(err) {
		panic("Index path does not exist. Did you call LoadDictionary() first?")
	}

	index, openErr := bleve.Open(indexPath)

	defer index.Close()

	Check(openErr)

	query := bleve.NewMatchQuery(queryStr)
	search := bleve.NewSearchRequest(query)
	searchResults, searchErr := index.Search(search)

	Check(searchErr)

	hits := searchResults.Hits
	entries := make([]DictionaryEntry, len(hits))

	for i := range hits {
		entry := &DictionaryEntry{}
		hitID := hits[i].ID
		b, err := index.GetInternal([]byte(hitID))

		if err != nil {
			panic(err)
		}

		json.Unmarshal(b, &entry)

		entries[i] = *entry
	}

	return entries
}
