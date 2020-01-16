package odict

import (
	"bufio"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	schema "github.com/Linguistic/odict/go/.schema"
	"github.com/golang/snappy"
	flatbuffers "github.com/google/flatbuffers/go"
	uuid "github.com/google/uuid"
)

type xmlDefinitionGroup struct {
	XMLName     xml.Name `xml:"group"`
	Definitions []string `xml:"definition"`
	Description string   `xml:"description,attr"`
}

type xmlUsage struct {
	XMLName          xml.Name             `xml:"usage"`
	POS              string               `xml:"pos,attr"`
	DefinitionGroups []xmlDefinitionGroup `xml:"group"`
	Definitions      []string             `xml:"definition"`
}

type xmlEtymology struct {
	XMLName     xml.Name   `xml:"ety"`
	Description string     `xml:"description,attr"`
	Usages      []xmlUsage `xml:"usage"`
}

type xmlEntry struct {
	XMLName     xml.Name       `xml:"entry"`
	Term        string         `xml:"term,attr"`
	Etymologies []xmlEtymology `xml:"ety"`
}

type xmlDictionary struct {
	XMLName xml.Name   `xml:"dictionary"`
	Name    string     `xml:"name,attr"`
	Entries []xmlEntry `xml:"entry"`
}

func xmlToDictionary(file *os.File) xmlDictionary {
	var dictionary xmlDictionary

	byteValue, _ := ioutil.ReadAll(file)
	xml.Unmarshal(byteValue, &dictionary)

	return dictionary
}

func getDefinitionsVectorFromUsage(builder *flatbuffers.Builder, usage xmlUsage) flatbuffers.UOffsetT {
	definitions := usage.Definitions

	var defBuffer []flatbuffers.UOffsetT

	for idx := range definitions {
		defBuffer = append(defBuffer, builder.CreateString(definitions[idx]))
	}

	defCount := len(defBuffer)

	schema.GroupStartDefinitionsVector(builder, defCount)

	for i := defCount - 1; i >= 0; i-- {
		builder.PrependUOffsetT(defBuffer[i])
	}

	return builder.EndVector(defCount)
}

func getDefinitionsVectorFromGroup(builder *flatbuffers.Builder, group xmlDefinitionGroup) flatbuffers.UOffsetT {
	definitions := group.Definitions

	var defBuffer []flatbuffers.UOffsetT

	for idx := range definitions {
		defBuffer = append(defBuffer, builder.CreateString(definitions[idx]))
	}

	defCount := len(defBuffer)

	schema.GroupStartDefinitionsVector(builder, defCount)

	for i := defCount - 1; i >= 0; i-- {
		builder.PrependUOffsetT(defBuffer[i])
	}

	return builder.EndVector(defCount)
}

func getGroupsVector(builder *flatbuffers.Builder, usage xmlUsage) flatbuffers.UOffsetT {
	groups := usage.DefinitionGroups

	var groupBuffer []flatbuffers.UOffsetT

	for idx := range groups {
		group := groups[idx]
		groupID := builder.CreateString(strconv.Itoa(idx))
		groupDescription := builder.CreateString(group.Description)
		groupDefinitions := getDefinitionsVectorFromGroup(builder, group)

		schema.GroupStart(builder)
		schema.GroupAddId(builder, groupID)
		schema.GroupAddDescription(builder, groupDescription)
		schema.EtymologyAddUsages(builder, groupDefinitions)

		groupBuffer = append(groupBuffer, schema.EtymologyEnd(builder))
	}

	groupCount := len(groupBuffer)

	schema.UsageStartGroupsVector(builder, groupCount)

	for i := groupCount - 1; i >= 0; i-- {
		builder.PrependUOffsetT(groupBuffer[i])
	}

	return builder.EndVector(groupCount)
}

func getUsagesVector(builder *flatbuffers.Builder, ety xmlEtymology) flatbuffers.UOffsetT {
	usages := ety.Usages

	var usageBuffer []flatbuffers.UOffsetT

	for idx := range usages {
		usage := usages[idx]
		usageID := builder.CreateString(strconv.Itoa(idx))
		usagePOS := builder.CreateString(usage.POS)
		usageDefinitionGroups := getGroupsVector(builder, usage)
		usageDefinitions := getDefinitionsVectorFromUsage(builder, usage)

		schema.UsageStart(builder)
		schema.UsageAddId(builder, usageID)
		schema.UsageAddPos(builder, usagePOS)
		schema.UsageAddGroups(builder, usageDefinitionGroups)
		schema.UsageAddDefinitions(builder, usageDefinitions)

		usageBuffer = append(usageBuffer, schema.UsageEnd(builder))
	}

	usageCount := len(usageBuffer)

	schema.EtymologyStartUsagesVector(builder, usageCount)

	for i := usageCount - 1; i >= 0; i-- {
		builder.PrependUOffsetT(usageBuffer[i])
	}

	return builder.EndVector(usageCount)
}

func getEtymologiesVector(builder *flatbuffers.Builder, entry xmlEntry) flatbuffers.UOffsetT {
	etymologies := entry.Etymologies

	var etyBuffer []flatbuffers.UOffsetT

	for idx := range etymologies {
		ety := etymologies[idx]
		etyID := builder.CreateString(strconv.Itoa(idx))
		etyDescription := builder.CreateString(ety.Description)
		etyUsages := getUsagesVector(builder, ety)

		schema.EtymologyStart(builder)
		schema.EtymologyAddId(builder, etyID)
		schema.EtymologyAddDescription(builder, etyDescription)
		schema.EtymologyAddUsages(builder, etyUsages)

		etyBuffer = append(etyBuffer, schema.EtymologyEnd(builder))
	}

	etyCount := len(etyBuffer)

	schema.EntryStartEtymologiesVector(builder, etyCount)

	for i := etyCount - 1; i >= 0; i-- {
		builder.PrependUOffsetT(etyBuffer[i])
	}

	return builder.EndVector(etyCount)
}

func getEntriesVector(builder *flatbuffers.Builder, dictionary xmlDictionary) flatbuffers.UOffsetT {
	entries := dictionary.Entries

	var entryBuffer []flatbuffers.UOffsetT

	for idx := range entries {
		entry := entries[idx]
		entryID := builder.CreateString(strconv.Itoa(idx)) // TODO: add prefix
		entryTerm := builder.CreateString(entry.Term)
		entryEtymologies := getEtymologiesVector(builder, entry)

		schema.EntryStart(builder)
		schema.EntryAddId(builder, entryID)
		schema.EntryAddTerm(builder, entryTerm)
		schema.EntryAddEtymologies(builder, entryEtymologies)

		entryBuffer = append(entryBuffer, schema.EntryEnd(builder))
	}

	entryCount := len(entryBuffer)

	schema.DictionaryStartEntriesVector(builder, entryCount)

	for i := entryCount - 1; i >= 0; i-- {
		builder.PrependUOffsetT(entryBuffer[i])
	}

	return builder.EndVector(entryCount)
}

func dictionaryToBytes(dictionary xmlDictionary) []byte {
	builder := flatbuffers.NewBuilder(1024)

	id := builder.CreateString(base64.StdEncoding.EncodeToString([]byte(uuid.New().String())))
	name := builder.CreateString(dictionary.Name)
	entries := getEntriesVector(builder, dictionary)

	schema.DictionaryStart(builder)
	schema.DictionaryAddId(builder, id)
	schema.DictionaryAddName(builder, name)
	schema.DictionaryAddEntries(builder, entries)

	builder.Finish(schema.DictionaryEnd(builder))

	return builder.FinishedBytes()
}

func createODictFile(outputPath string, dictionary xmlDictionary) {
	dictionaryBytes := dictionaryToBytes(dictionary)
	compressed := snappy.Encode(nil, dictionaryBytes)
	file, err := os.Create(outputPath)

	Check(err)

	defer file.Close()

	signature := []byte("ODICT")
	version := Uint16ToBytes(2)
	compressedSize := uint32(len(compressed))
	compressedSizeBytes := Uint32ToBytes(compressedSize)

	writer := bufio.NewWriter(file)

	sigBytes, sigErr := writer.Write(signature)
	versionBytes, versionErr := writer.Write(version)
	contentSizeBytes, contentCountErr := writer.Write(compressedSizeBytes)
	contentBytes, contentErr := writer.Write(compressed)
	total := sigBytes + versionBytes + contentSizeBytes + contentBytes

	Check(sigErr)
	Check(versionErr)
	Check(contentCountErr)
	Check(contentErr)

	Assert(sigBytes == 5, "Signature bytes do not equal 5")
	Assert(versionBytes == 2, "Version bytes do not equal 2")
	Assert(contentSizeBytes == 4, "Content byte count does not equal 4")
	Assert(contentBytes == int(compressedSize), "Content does not equal the computed byte count")

	writer.Flush()

	fmt.Printf("Wrote %d bytes to path: %s\n", total, outputPath)
}

// WriteDictionary generates an ODict binary file given
// a ODXML input file path
func WriteDictionary(inputPath, outputPath string) {
	xmlFile, err := os.Open(inputPath)

	defer xmlFile.Close()

	if err != nil {
		fmt.Println(err)
		return
	}

	createODictFile(outputPath, xmlToDictionary(xmlFile))
}
