// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package spdagvm

import (
	"bytes"
	"testing"

	"github.com/ava-labs/gecko/ids"
	"github.com/ava-labs/gecko/snow"
	"github.com/ava-labs/gecko/utils/crypto"
	"github.com/ava-labs/gecko/utils/units"
	"github.com/ava-labs/gecko/utils/wrappers"
)

// Ensure transaction verification fails when a transaction has
// the wrong chain ID
func TestTxVerifyBadChainID(t *testing.T) {
	genesisTx := GenesisTx(defaultInitBalances)

	ctx := snow.DefaultContextTest()
	ctx.NetworkID = 15
	ctx.ChainID = avaxChainID

	builder := Builder{
		NetworkID: ctx.NetworkID,
		ChainID:   ctx.ChainID,
	}
	tx, err := builder.NewTx( //valid transaction
		/*ins=*/ []Input{
			builder.NewInputPayment(
				/*txID=*/ genesisTx.ID(),
				/*txIndex=*/ 0,
				/*amount=*/ 5*units.Avax,
				/*sigs=*/ []*Sig{builder.NewSig(0 /*=index*/)},
			),
		},
		/*outs=*/ []Output{
			builder.NewOutputPayment(
				/*amount=*/ 3*units.Avax,
				/*locktime=*/ 0,
				/*threshold=*/ 0,
				/*addresses=*/ nil,
			),
		},
		/*signers=*/ []*InputSigner{
			{Keys: []*crypto.PrivateKeySECP256K1R{
				keys[1], // reference to vm_test.go
			}},
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	// Should pass verification when chain ID is correct
	if err := tx.Verify(ctx, txFeeTest); err != nil {
		t.Fatal("Should have passed verification")
	}

	ctx.ChainID = ctx.ChainID.Prefix()

	// Should pass verification when chain ID is wrong
	if err := tx.Verify(ctx, txFeeTest); err != errWrongChainID {
		t.Fatal("Should have failed with errWrongChainID")
	}
}

func TestUnsignedTx(t *testing.T) {
	skBytes := []byte{
		0x98, 0xcb, 0x07, 0x7f, 0x97, 0x2f, 0xeb, 0x04,
		0x81, 0xf1, 0xd8, 0x94, 0xf2, 0x72, 0xc6, 0xa1,
		0xe3, 0xc1, 0x5e, 0x27, 0x2a, 0x16, 0x58, 0xff,
		0x71, 0x64, 0x44, 0xf4, 0x65, 0x20, 0x00, 0x70,
	}
	outputPaymentBytes := []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x30, 0x39, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0xd4, 0x31, 0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x02, 0x51, 0x02, 0x5c, 0x61,
		0xfb, 0xcf, 0xc0, 0x78, 0xf6, 0x93, 0x34, 0xf8,
		0x34, 0xbe, 0x6d, 0xd2, 0x6d, 0x55, 0xa9, 0x55,
		0xc3, 0x34, 0x41, 0x28, 0xe0, 0x60, 0x12, 0x8e,
		0xde, 0x35, 0x23, 0xa2, 0x4a, 0x46, 0x1c, 0x89,
		0x43, 0xab, 0x08, 0x59,
	}
	outputTakeOrLeaveBytes := []byte{
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x30, 0x39, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0xd4, 0x31, 0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x01, 0x51, 0x02, 0x5c, 0x61,
		0xfb, 0xcf, 0xc0, 0x78, 0xf6, 0x93, 0x34, 0xf8,
		0x34, 0xbe, 0x6d, 0xd2, 0x6d, 0x55, 0xa9, 0x55,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xdd, 0xd5,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0xc3, 0x34, 0x41, 0x28, 0xe0, 0x60, 0x12, 0x8e,
		0xde, 0x35, 0x23, 0xa2, 0x4a, 0x46, 0x1c, 0x89,
		0x43, 0xab, 0x08, 0x59,
	}
	inputPaymentBytes := []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x02, 0x03,
		0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b,
		0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13,
		0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b,
		0x1c, 0x1d, 0x1e, 0x1f, 0x00, 0x00, 0x00, 0x09,
		0x00, 0x00, 0x00, 0x00, 0x07, 0x5b, 0xcd, 0x15,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x07,
	}
	chainID := ids.NewID([32]byte{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
		0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
	})

	f := crypto.FactorySECP256K1R{}
	sk, err := f.ToPrivateKey(skBytes)
	if err != nil {
		t.Fatal(err)
	}

	c := Codec{}
	p := wrappers.Packer{Bytes: outputPaymentBytes}
	outputPayment := c.unmarshalOutput(&p)
	if p.Errored() {
		t.Fatal(p.Err)
	}

	p = wrappers.Packer{Bytes: outputTakeOrLeaveBytes}
	outputTakeOrLeave := c.unmarshalOutput(&p)
	if p.Errored() {
		t.Fatal(p.Err)
	}

	p = wrappers.Packer{Bytes: inputPaymentBytes}
	inputPayment := c.unmarshalInput(&p)
	if p.Errored() {
		t.Fatal(p.Err)
	}

	inputPaymentSigner := &InputSigner{
		Keys: []*crypto.PrivateKeySECP256K1R{
			sk.(*crypto.PrivateKeySECP256K1R),
		},
	}

	b := Builder{
		NetworkID: 0,
		ChainID:   chainID,
	}
	tx, err := b.NewTx(
		/*inputs=*/ []Input{inputPayment},
		/*outputs=*/ []Output{outputPayment, outputTakeOrLeave},
		/*signers=*/ []*InputSigner{inputPaymentSigner},
	)
	if err != nil {
		t.Fatal(err)
	}

	unsignedTxBytes, err := c.MarshalUnsignedTx(tx)
	if err != nil {
		t.Fatal(err)
	}

	expected := []byte{
		// codec:
		0x00, 0x00, 0x00, 0x02,
		// networkID:
		0x00, 0x00, 0x00, 0x00,
		// chainID:
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
		0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
		// number of outputs:
		0x00, 0x00, 0x00, 0x02,
		// output payment:
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x30, 0x39, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0xd4, 0x31, 0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x02, 0x51, 0x02, 0x5c, 0x61,
		0xfb, 0xcf, 0xc0, 0x78, 0xf6, 0x93, 0x34, 0xf8,
		0x34, 0xbe, 0x6d, 0xd2, 0x6d, 0x55, 0xa9, 0x55,
		0xc3, 0x34, 0x41, 0x28, 0xe0, 0x60, 0x12, 0x8e,
		0xde, 0x35, 0x23, 0xa2, 0x4a, 0x46, 0x1c, 0x89,
		0x43, 0xab, 0x08, 0x59,
		// output take-or-leave:
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x30, 0x39, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0xd4, 0x31, 0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x01, 0x51, 0x02, 0x5c, 0x61,
		0xfb, 0xcf, 0xc0, 0x78, 0xf6, 0x93, 0x34, 0xf8,
		0x34, 0xbe, 0x6d, 0xd2, 0x6d, 0x55, 0xa9, 0x55,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xdd, 0xd5,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0xc3, 0x34, 0x41, 0x28, 0xe0, 0x60, 0x12, 0x8e,
		0xde, 0x35, 0x23, 0xa2, 0x4a, 0x46, 0x1c, 0x89,
		0x43, 0xab, 0x08, 0x59,
		// number of inputs:
		0x00, 0x00, 0x00, 0x01,
		// input payment:
		0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x02, 0x03,
		0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b,
		0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13,
		0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b,
		0x1c, 0x1d, 0x1e, 0x1f, 0x00, 0x00, 0x00, 0x09,
		0x00, 0x00, 0x00, 0x00, 0x07, 0x5b, 0xcd, 0x15,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x07,
	}
	if !bytes.Equal(unsignedTxBytes, expected) {
		t.Fatalf("Codec.MarshalUnsignedTx returned:\n0x%x\nExpected:\n0x%x", unsignedTxBytes, expected)
	}
}

func TestSignedTx(t *testing.T) {
	skBytes := []byte{
		0x98, 0xcb, 0x07, 0x7f, 0x97, 0x2f, 0xeb, 0x04,
		0x81, 0xf1, 0xd8, 0x94, 0xf2, 0x72, 0xc6, 0xa1,
		0xe3, 0xc1, 0x5e, 0x27, 0x2a, 0x16, 0x58, 0xff,
		0x71, 0x64, 0x44, 0xf4, 0x65, 0x20, 0x00, 0x70,
	}
	outputPaymentBytes := []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x30, 0x39, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0xd4, 0x31, 0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x02, 0x51, 0x02, 0x5c, 0x61,
		0xfb, 0xcf, 0xc0, 0x78, 0xf6, 0x93, 0x34, 0xf8,
		0x34, 0xbe, 0x6d, 0xd2, 0x6d, 0x55, 0xa9, 0x55,
		0xc3, 0x34, 0x41, 0x28, 0xe0, 0x60, 0x12, 0x8e,
		0xde, 0x35, 0x23, 0xa2, 0x4a, 0x46, 0x1c, 0x89,
		0x43, 0xab, 0x08, 0x59,
	}
	outputTakeOrLeaveBytes := []byte{
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x30, 0x39, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0xd4, 0x31, 0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x01, 0x51, 0x02, 0x5c, 0x61,
		0xfb, 0xcf, 0xc0, 0x78, 0xf6, 0x93, 0x34, 0xf8,
		0x34, 0xbe, 0x6d, 0xd2, 0x6d, 0x55, 0xa9, 0x55,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xdd, 0xd5,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0xc3, 0x34, 0x41, 0x28, 0xe0, 0x60, 0x12, 0x8e,
		0xde, 0x35, 0x23, 0xa2, 0x4a, 0x46, 0x1c, 0x89,
		0x43, 0xab, 0x08, 0x59,
	}
	inputPaymentBytes := []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x02, 0x03,
		0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b,
		0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13,
		0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b,
		0x1c, 0x1d, 0x1e, 0x1f, 0x00, 0x00, 0x00, 0x09,
		0x00, 0x00, 0x00, 0x00, 0x07, 0x5b, 0xcd, 0x15,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x07,
	}
	chainID := ids.NewID([32]byte{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
		0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
	})

	f := crypto.FactorySECP256K1R{}
	sk, err := f.ToPrivateKey(skBytes)
	if err != nil {
		t.Fatal(err)
	}

	c := Codec{}
	p := wrappers.Packer{Bytes: outputPaymentBytes}
	outputPayment := c.unmarshalOutput(&p)
	if p.Errored() {
		t.Fatal(p.Err)
	}

	p = wrappers.Packer{Bytes: outputTakeOrLeaveBytes}
	outputTakeOrLeave := c.unmarshalOutput(&p)
	if p.Errored() {
		t.Fatal(p.Err)
	}

	p = wrappers.Packer{Bytes: inputPaymentBytes}
	inputPayment := c.unmarshalInput(&p)
	if p.Errored() {
		t.Fatal(p.Err)
	}

	inputPaymentSigner := &InputSigner{
		Keys: []*crypto.PrivateKeySECP256K1R{
			sk.(*crypto.PrivateKeySECP256K1R),
		},
	}

	b := Builder{
		NetworkID: 0,
		ChainID:   chainID,
	}
	tx, err := b.NewTx(
		/*inputs=*/ []Input{inputPayment},
		/*outputs=*/ []Output{outputPayment, outputTakeOrLeave},
		/*signers=*/ []*InputSigner{inputPaymentSigner},
	)
	if err != nil {
		t.Fatal(err)
	}
	signedTxBytes := tx.Bytes()

	expected := []byte{
		// unsigned transaction:
		0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
		0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
		0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x30, 0x39,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xd4, 0x31,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x02,
		0x51, 0x02, 0x5c, 0x61, 0xfb, 0xcf, 0xc0, 0x78,
		0xf6, 0x93, 0x34, 0xf8, 0x34, 0xbe, 0x6d, 0xd2,
		0x6d, 0x55, 0xa9, 0x55, 0xc3, 0x34, 0x41, 0x28,
		0xe0, 0x60, 0x12, 0x8e, 0xde, 0x35, 0x23, 0xa2,
		0x4a, 0x46, 0x1c, 0x89, 0x43, 0xab, 0x08, 0x59,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x30, 0x39, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0xd4, 0x31, 0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x01, 0x51, 0x02, 0x5c, 0x61,
		0xfb, 0xcf, 0xc0, 0x78, 0xf6, 0x93, 0x34, 0xf8,
		0x34, 0xbe, 0x6d, 0xd2, 0x6d, 0x55, 0xa9, 0x55,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xdd, 0xd5,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0xc3, 0x34, 0x41, 0x28, 0xe0, 0x60, 0x12, 0x8e,
		0xde, 0x35, 0x23, 0xa2, 0x4a, 0x46, 0x1c, 0x89,
		0x43, 0xab, 0x08, 0x59, 0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x02, 0x03,
		0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b,
		0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13,
		0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b,
		0x1c, 0x1d, 0x1e, 0x1f, 0x00, 0x00, 0x00, 0x09,
		0x00, 0x00, 0x00, 0x00, 0x07, 0x5b, 0xcd, 0x15,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x07,
		// signature:
		0x41, 0x3a, 0x8e, 0x30, 0x72, 0x1b, 0xd3, 0xdd,
		0x0f, 0x49, 0x3d, 0x0d, 0xed, 0x82, 0x8b, 0x90,
		0x8c, 0xfb, 0x5d, 0xd5, 0x4b, 0x63, 0x76, 0x42,
		0x99, 0x66, 0xda, 0x10, 0x14, 0x81, 0x89, 0x2d,
		0x22, 0x4b, 0x6c, 0x95, 0x6b, 0x93, 0x05, 0x13,
		0x83, 0x5d, 0xea, 0xa4, 0x44, 0x8f, 0x46, 0xb1,
		0x23, 0x45, 0x47, 0x05, 0xe9, 0xa5, 0x3b, 0xfc,
		0x27, 0x09, 0x21, 0x1a, 0x5c, 0x5a, 0x58, 0xec,
		0x01,
	}

	if !bytes.Equal(signedTxBytes, expected) {
		t.Fatalf("Codec.MarshalTx returned:\n0x%x\nExpected:\n0x%x", signedTxBytes, expected)
	}
}