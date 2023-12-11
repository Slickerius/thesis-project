package omemo

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	b64 "encoding/base64"
	"encoding/xml"
	"log"
	"math/rand"
	"strconv"
	"time"

	"google.golang.org/protobuf/proto"
	"mellium.im/communique/internal/client"
	"mellium.im/communique/internal/client/doubleratchet"
	"mellium.im/communique/internal/client/omemo/protobuf"
	"mellium.im/communique/internal/client/x3dh"
	"mellium.im/xmpp/jid"
	"mellium.im/xmpp/stanza"
)

func SetupClient(c *client.Client, logger *log.Logger) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	keyBundleAnnouncementStanza := WrapKeyBundle(c)

	err := c.UnmarshalIQ(ctx, keyBundleAnnouncementStanza.TokenReader(), nil)

	if err != nil {
		logger.Printf("Error sending key bundle: %q", err)
	}

	// ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
	// defer cancel()

	// deviceAnnouncementStanza := WrapDeviceIds([]Device{
	// 	{ID: "1", Label: "Acer Aspire 3"},
	// }, c)

	// _, err = c.SendIQ(ctx, deviceAnnouncementStanza.TokenReader())

	// if err != nil {
	// 	logger.Printf("Error sending device list: %q", err)
	// }

}

func InitiateKeyAgreement(initialMessage string, c *client.Client, logger *log.Logger, targetJID jid.JID) (*EncryptedMessage, stanza.Message) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	logger.Printf("Fetching key bundle for " + targetJID.String() + " ...")
	fetchBundleStanza := WrapNodeFetch("urn:xmpp:omemo:2:bundles", c.DeviceId, targetJID, c)

	payload, err := c.SendIQ(ctx, fetchBundleStanza.TokenReader())

	if err != nil {
		logger.Printf("Error fetching key bundle: %v", err)
	}

	logger.Printf("Decoding key bundle for " + targetJID.String() + " ...")

	defer func() {
		payload := payload
		if err != nil {
			payload.Close()
		}
	}()

	var currentEl string
	var targetSpkPub, targetSpkSig, targetIdKeyPub []byte
	var opkList []PreKey
	var opkId string
	var spkId string

	for t, _ := payload.Token(); t != nil; t, _ = payload.Token() {
		switch se := t.(type) {
		case xml.StartElement:
			currentEl = se.Name.Local
			if se.Name.Local == "pk" {
				opkId = se.Attr[0].Value
			} else if se.Name.Local == "spk" {
				spkId = se.Attr[0].Value
			}
		case xml.CharData:
			content := string(se[:])
			switch currentEl {
			case "spk":
				targetSpkPub, _ = b64.StdEncoding.DecodeString(content)
			case "spks":
				targetSpkSig, _ = b64.StdEncoding.DecodeString(content)
			case "ik":
				targetIdKeyPub, _ = b64.StdEncoding.DecodeString(content)
			case "pk":
				opkList = append(opkList, PreKey{ID: opkId, Text: content})
			}
		}
	}

	randomIndex := rand.Intn(len(opkList))
	opk := opkList[randomIndex]
	chosenOpkId, err := strconv.Atoi(opk.ID)
	chosenOpkIdUint := uint32(chosenOpkId)

	opkPub, _ := b64.StdEncoding.DecodeString(opk.Text)

	sharedKey, associatedData, ekPub, err := x3dh.CreateInitialMessage(c.IdPrivKey, targetIdKeyPub, opkPub, targetSpkPub, targetSpkSig)
	logger.Print("SHARED KEY")
	logger.Print(sharedKey)
	logger.Print("ASSOCIATED DATA")
	logger.Print(associatedData)
	logger.Print("EK PUB")
	logger.Print(ekPub)

	sess, err := doubleratchet.CreateActive(sharedKey, associatedData, targetIdKeyPub)

	if err != nil {
		logger.Printf("Failed creating double ratchet session: %s", err)
	}

	if err != nil {
		logger.Printf("Failed marshaling OMEMOKeyExchange: %s", err)
	}

	chosenSpkId, _ := strconv.Atoi(string(spkId))
	chosenSpkIdUint := uint32(chosenSpkId)

	jdid := targetJID.Bare().String() + ":" + c.DeviceId

	c.MessageSession[jdid] = sess

	return EncryptMessage(initialMessage, true, &chosenOpkIdUint, &chosenSpkIdUint, ekPub, c, logger, targetJID)
}

func EncryptMessage(initialMessage string, keyExchange bool, opkId *uint32, spkId *uint32, ek []byte, c *client.Client, logger *log.Logger, targetJID jid.JID) (*EncryptedMessage, stanza.Message) {
	jdid := targetJID.Bare().String() + ":" + c.DeviceId
	sess := c.MessageSession[jdid]

	envelope := WrapEnvelope(initialMessage, c)
	envelopeMarshaled, _ := xml.Marshal(envelope)
	envelopeMarshaledEncoded := b64.StdEncoding.EncodeToString(envelopeMarshaled)

	ciphertext, authKey, err := sess.Encrypt([]byte(envelopeMarshaledEncoded))
	ciphertextEncoded := b64.StdEncoding.EncodeToString(ciphertext)

	// Sess.Encrypt already handles structuring similar to OMEMOMessage.proto, so we don't have to use OMEMOMessage again

	if err != nil {
		logger.Printf("Failed encrypting message with double ratchet session: %s", err)
	}

	mac := hmac.New(sha256.New, authKey)
	mac.Write(ciphertext)
	macResult := mac.Sum(nil)

	authenticatedMessage := &protobuf.OMEMOAuthenticatedMessage{
		Mac:     macResult,
		Message: ciphertext,
	}

	var keyElement string

	if keyExchange {
		keyExchangeMessage := &protobuf.OMEMOKeyExchange{
			PkId:    opkId,
			SpkId:   spkId,
			Ik:      c.IdPubKey,
			Ek:      ek,
			Message: authenticatedMessage,
		}

		omemoKeyExchangeMessage, err := proto.Marshal(keyExchangeMessage)

		if err != nil {
			logger.Printf("Failed marshaling OMEMOKeyExchange: %s", err)
		}

		logger.Print("OMEMOKEYEXCHANGE")
		logger.Print(omemoKeyExchangeMessage)

		keyElement = b64.StdEncoding.EncodeToString(omemoKeyExchangeMessage)
	} else {
		authenticatedMessage, err := proto.Marshal(authenticatedMessage)

		if err != nil {
			logger.Printf("Failed marshaling OMEMOAuthenticatedMessage: %s", err)
		}

		logger.Print("OMEMOAUTHMSG")
		logger.Print(authenticatedMessage)

		keyElement = b64.StdEncoding.EncodeToString(authenticatedMessage)
	}

	encrypted, stanzaMessage := WrapEncrypted(targetJID, c.DeviceId, keyElement, ciphertextEncoded, keyExchange, c)

	return encrypted, stanzaMessage
}
