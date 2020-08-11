package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type SourceType string

const (
	SourceTypeCard SourceType = "Card"
)

type Source interface {
	isSource()
}

type Card struct {
	ID     string     `json:"id"`
	Object SourceType `json:"object"`
}

func (c *Card) isSource() {
}

type TxPoly struct {
	Amount int64  `json:"amount"`
	Source Source `json:"source"`
}

type sourceKey struct {
	Object string `json:"object"`
}

func unmarshalSource(b []byte) (Source, error) {
	var key sourceKey
	err := json.Unmarshal(b, &key)
	if err != nil {
		return nil, err
	}
	switch key.Object {
	case "card":
		card := Card{}
		return &card, json.Unmarshal(b, &card)
	}
	return nil, fmt.Errorf("unknown source type: %s", key.Object)
}

func (tx *TxPoly) UnmarshalJSON(b []byte) error {
	type Alias TxPoly
	var a struct {
		Alias
		Source json.RawMessage `json:"source"`
	}
	err := json.Unmarshal(b, &a)
	if err != nil {
		return err
	}
	src, err := unmarshalSource(a.Source)
	if err != nil {
		return err
	}
	*tx = TxPoly(a.Alias)
	tx.Source = src
	return nil
}

const sourcePoly = `{
  "amount": 100,
  "source": {
    "id": "card_123",
	"object": "card"
  }
}`

const sourceFields = `{
  "amount": 100,
  "source": {
	"card": {
	  "id": "card_123"
	}
  }
}`

const sourceAny = `{
  "amount": 100,
  "source": {
	"type": "card",
	"object": {
	  "id": "card_123"
	}
  }
}`

func jsonSource() {
	{
		// Option 1: Single object unmarshalled by looking at the Object key
		var tx TxPoly
		if err := json.Unmarshal([]byte(sourcePoly), &tx); err != nil {
			log.Fatal(err)
		}

		switch source := tx.Source.(type) {
		case *Card:
			fmt.Println("card", source)
		}

		out, err := json.Marshal(tx)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("tx poly:\n%s\n", string(out))
	}
}
