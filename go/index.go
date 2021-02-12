package odict

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/document"
	"github.com/blevesearch/bleve/v2/index/scorch"
	idx "github.com/blevesearch/bleve_index_api"
	"github.com/schollz/progressbar/v3"
)

func getIndexPath(dictionaryID string) string {
	return filepath.Join(os.TempDir(), "odict", "idx", dictionaryID)
}

// func getIndexMapping() {
// 	// a generic reusable mapping for english text
// 	textFieldMapping := bleve.NewTextFieldMapping()

// 	// a generic reusable mapping for keyword text
// 	keywordFieldMapping := bleve.NewTextFieldMapping()
// 	keywordFieldMapping.Analyzer = keyword.Name

// 	entryMapping := bleve.NewDocumentMapping()

// 	// name
// 	entryMapping.AddFieldMappingsAt("name", englishTextFieldMapping)

// 	// description
// 	beerMapping.AddFieldMappingsAt("description",
// 		englishTextFieldMapping)

// 	beerMapping.AddFieldMappingsAt("type", keywordFieldMapping)
// 	beerMapping.AddFieldMappingsAt("style", keywordFieldMapping)
// 	beerMapping.AddFieldMappingsAt("category", keywordFieldMapping)

// 	breweryMapping := bleve.NewDocumentMapping()
// 	breweryMapping.AddFieldMappingsAt("name", englishTextFieldMapping)
// 	breweryMapping.AddFieldMappingsAt("description", englishTextFieldMapping)

// 	indexMapping := bleve.NewIndexMapping()
// 	indexMapping.AddDocumentMapping("beer", beerMapping)
// 	indexMapping.AddDocumentMapping("brewery", breweryMapping)

// 	indexMapping.TypeField = "type"
// 	indexMapping.DefaultAnalyzer = "en"

// 	return indexMapping, nil
// }

func IndexDictionary(dictionary Dictionary, overwrite bool) string {
	indexPath := getIndexPath(dictionary.ID)
	_, statErr := os.Stat(indexPath)

	if os.IsNotExist(statErr) {
		fmt.Println("Indexing dictionary (this might take some time)...")
		mapping := bleve.NewIndexMapping()
		index, indexErr := bleve.NewUsing(indexPath, mapping, scorch.Name, scorch.Name, nil)

		defer index.Close()

		Check(indexErr)

		totalEntries := len(dictionary.Entries.Iterable)

		bar := progressbar.Default(int64(totalEntries))
		batch := index.NewBatch()
		batchCount := 0
		batchSize := 100

		for key := range dictionary.Entries.Iterable {
			entry := dictionary.Entries.Get(key)
			doc := document.NewDocument(entry.Term)

			mapping.MapDocument(doc, entry)

			enc := EncodeEntry(entry)

			field := document.NewTextFieldWithIndexingOptions("_source", nil, enc, idx.StoreField)

			nd := doc.AddField(field)

			batch.IndexAdvanced(nd)

			batchCount++

			bar.Add(1)

			time.Sleep(time.Millisecond)

			if batchCount >= batchSize {
				idxErr := index.Batch(batch)

				Check(idxErr)

				batch = index.NewBatch()
				batchCount = 0
			}
		}

		idxErr := index.Batch(batch)

		Check(idxErr)

		fmt.Println()
	} else {
		if overwrite {
			println("Purging existing index...")
			os.RemoveAll(indexPath)
			return IndexDictionary(dictionary, false)
		}
	}

	return indexPath
}
