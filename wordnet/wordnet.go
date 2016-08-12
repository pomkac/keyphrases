package wordnet

/*
* Implements some operations on a WordNet dict
 */

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gelembjuk/keyphrases/helper"
)

const nounIndex = "index.noun"
const verbIndex = "index.verb"
const adjIndex = "index.adj"
const advIndex = "index.adv"

const nounData = "data.noun"
const verbData = "data.verb"
const adjData = "data.adj"
const advData = "data.adv"

type WordNet struct {
	DictLocationDirectory string
	senseCache            map[string]string
	fileHandles           map[int]*os.File
	dataFileHandles       map[int]*os.File
}

type indexRecord struct {
	Lemma       string
	Pos         string
	Offsets     []string
	SenseCnt    int
	PCnt        int
	TagSenseCnt int
	Found       bool
}

type wordReference struct {
	Lemma  string
	Pos    string
	Offset string
}

type dataRecord struct {
	Lemma    string
	Pos      string
	Words    []string
	WordsCnt int
	Found    bool
}

func (this *WordNet) init() error {
	// check if directory exists
	indexpath := this.DictLocationDirectory + nounIndex

	if _, err := os.Stat(indexpath); os.IsNotExist(err) {
		return errors.New("WordNet directory is not set or not found")
	}

	this.fileHandles = map[int]*os.File{}
	this.senseCache = map[string]string{}

	return nil
}

func (this *WordNet) SetDictDirectory(path string) error {
	indexpath := path + nounIndex

	if _, err := os.Stat(indexpath); os.IsNotExist(err) {
		return errors.New("WordNet directory is not set or not found")
	}

	this.DictLocationDirectory = path

	return nil
}

func (this *WordNet) GetWord(word string) (string, error) {
	err := this.init()

	if err != nil {
		return "", err
	}

	return "", nil
}

func (this *WordNet) GetWordOptions(word string) ([]string, error) {
	wordmap, err := this.GetWordOptionsMap(word)

	if err != nil {
		return []string{}, err
	}

	if len(wordmap) == 0 {
		return []string{}, nil
	}
	// sort by values and return sorted

	options := helper.KeysSortedByValuesReverse(wordmap)

	return options, nil
}

func (this *WordNet) GetWordOptionsMap(word string) (map[string]int, error) {
	options := map[string]int{}

	err := this.init()

	if err != nil {
		return options, err
	}

	for i := 1; i <= 4; i++ {
		r, e := this.getRecordForWord(word, i)

		if e != nil {
			return options, e
		}

		if r.Found {
			count := len(r.Offsets)
			options[r.Pos] = count
		}
	}

	return options, nil
}

func (this *WordNet) GetWordSynonims(word string) ([]string, error) {
	return this.GetWordSences(word, "syns")
}

func (this *WordNet) GetWordSences(word string, sensetype string) ([]string, error) {
	options := []string{}

	err := this.init()

	if err != nil {
		return options, err
	}

	// local structure to keep list of senses and references

	wordreferences := []wordReference{}

	// try 4 supported Pos to get list of offsets
	for i := 1; i <= 4; i++ {
		r, e := this.getRecordForWord(word, i)

		if e != nil {
			return options, e
		}

		if r.Found {
			for _, o := range r.Offsets {
				ref := wordReference{Lemma: r.Lemma, Pos: r.Pos, Offset: o}

				wordreferences = append(wordreferences, ref)
			}
		}
	}

	for _, ref := range wordreferences {
		fmt.Printf("%s, %d, %s\n", ref.Lemma, ref.Pos, ref.Offset)
	}
	os.Exit(1)
	return options, nil
}

func (this *WordNet) getRecordForWord(word string, source int) (indexRecord, error) {
	result := indexRecord{}

	// get file handle
	fhandle, err := this.getFileHandle(source)

	if err != nil {
		return result, err
	}

	// find a string in a file starting with thos word
	scanner := bufio.NewScanner(fhandle)

	searchstring := word + " "

	for scanner.Scan() {

		line := scanner.Text()

		if strings.Index(line, searchstring) == 0 {
			// word is found. Parse a line and return
			tokens := strings.SplitN(line, " ", 5)

			if len(tokens) < 4 {
				return result, errors.New("Unable to parse WordNet dict line")
			}

			result.Lemma = tokens[0]
			result.Pos = tokens[1]
			result.SenseCnt, _ = strconv.Atoi(tokens[2])
			result.PCnt, _ = strconv.Atoi(tokens[3])

			rline := tokens[4]

			// drop some part of a string
			for i := 0; i < result.PCnt; i++ {
				tmpres := strings.SplitN(rline, " ", 2)
				rline = tmpres[1]
			}

			tmpres := strings.Split(rline, " ")

			result.TagSenseCnt, _ = strconv.Atoi(tmpres[1])

			if len(tmpres) > 2 {
				for i := 2; i < len(tmpres); i++ {
					if tmpres[i] != "" {
						result.Offsets = append(result.Offsets, tmpres[i])
					}
				}
			}

			result.Found = true

			return result, nil
		}

	}

	if err = scanner.Err(); err != nil {
		return result, err
	}

	return result, nil
}
func (this *WordNet) dataLookup(source int, offset string) (dataRecord, error) {
	result := dataRecord{}
	// get file handle
	fhandle, err := this.getDataFileHandle(source)

	if err != nil {
		return result, err
	}

	seekoffset, err2 := strconv.Atoi(offset)

	if err != nil {
		return result, err2
	}

	_, err = fhandle.Seek(int64(seekoffset), 0)

	if err != nil {
		return result, err
	}

	reader := bufio.NewReader(fhandle)

	line, _ := reader.ReadString('\n')

	fmt.Println(line)

	os.Exit(1)

	return result, nil
}
func (this *WordNet) getFileHandle(source int) (*os.File, error) {
	if handle, ok := this.fileHandles[source]; ok {
		return handle, nil
	}
	// open this file
	var filepath string

	switch {
	case source == 1:
		filepath = nounIndex
	case source == 2:
		filepath = verbIndex
	case source == 3:
		filepath = adjIndex
	case source == 4:
		filepath = advIndex
	default:
		return nil, errors.New("Unknown word type in index file open")
	}

	filepath = this.DictLocationDirectory + filepath

	f, err := os.Open(filepath)

	if err != nil {
		return nil, err
	}

	this.fileHandles[source] = f

	return this.fileHandles[source], nil
}

func (this *WordNet) getDataFileHandle(source int) (*os.File, error) {
	if handle, ok := this.dataFileHandles[source]; ok {
		return handle, nil
	}
	// open this file
	var filepath string

	switch {
	case source == 1:
		filepath = nounData
	case source == 2:
		filepath = verbData
	case source == 3:
		filepath = adjData
	case source == 4:
		filepath = advData
	default:
		return nil, errors.New("Unknown word type in data file open")
	}

	filepath = this.DictLocationDirectory + filepath

	f, err := os.Open(filepath)

	if err != nil {
		return nil, err
	}

	this.dataFileHandles[source] = f

	return this.dataFileHandles[source], nil
}

func (this *WordNet) Free() {
	for i, f := range this.fileHandles {
		f.Close()
		delete(this.fileHandles, i)
	}
	for i, f := range this.dataFileHandles {
		f.Close()
		delete(this.dataFileHandles, i)
	}
}
