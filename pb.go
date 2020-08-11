package main

import (
	"fmt"
	"log"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"

	"github.com/kyleconroy/sourcecode/core"
)

func pb() {
	m := &jsonpb.Marshaler{Indent: "  "}

	{
		// Option 1: Use one_of
		tx := &core.TxOneOfSource{
			Amount: 100,
			Source: &core.TxOneOfSource_Card{
				Card: &core.Card{Id: "card_123"},
			},
		}

		// Type switch for source information
		switch source := tx.Source.(type) {
		case *core.TxOneOfSource_Ach:
			fmt.Println("ach", source.Ach)
		case *core.TxOneOfSource_Card:
			fmt.Println("card", source.Card)
		}

		// Or use the getter methods on the Source message
		if card := tx.GetCard(); card != nil {
			fmt.Println("card", card)
		} else if ach := tx.GetAch(); ach != nil {
			fmt.Println("ach", ach)
		}

		out, _ := m.MarshalToString(tx)
		fmt.Printf("TxOneOfSource json:\n%s\n", out)
	}

	{
		// Option 2: Use any
		// https://medium.com/@pokstad/sending-any-any-thing-in-golang-with-protobuf-3-95f84838028d

		// I gave up with this example. Using the Any type with protocol
		// buffers is complicated, especially considering that you need to
		// serialize the source to set it.

		serialized, err := proto.Marshal(&core.Card{Id: "card_123"})
		if err != nil {
			log.Fatalf("could not serialize card: %s", err)
		}

		tx := &core.TxAnySource{
			Amount: 100,
			Source: &any.Any{
				TypeUrl: "conroy.org/types/Card",
				Value:   serialized,
			},
		}

		// try to unmarshal the Card
		card := &core.Card{}
		if err := ptypes.UnmarshalAny(tx.Source, card); err == nil {
			fmt.Println(card)
		}

		// try to unmarshal the Ach
		ach := &core.ACH{}
		if err := ptypes.UnmarshalAny(tx.Source, ach); err == nil {
			fmt.Println(ach)
		}

		out, err := m.MarshalToString(tx)
		if err != nil {
			log.Fatalf("could not serialize tx: %s", err)
		}
		fmt.Printf("TxAnySource json:\n%s\n", out)
	}
}
