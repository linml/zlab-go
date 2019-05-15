package crypto

import (
	"encoding/base64"
	"fmt"

	"zlab/library/nano"
	"zlab/library/nano/session"
	"github.com/xxtea/xxtea-go/xxtea"
)

var xxteaKey = []byte("7AEC4MA152BQE9HWQ7KB")

type Crypto struct {
	Key []byte
}

func NewCrypto() *Crypto {
	return &Crypto{xxteaKey}
}

func (c *Crypto) Inbound(s *session.Session, msg nano.Message) error {
	out, err := base64.StdEncoding.DecodeString(string(msg.Data))
	if err != nil {
		return err
	}

	out = xxtea.Decrypt(out, c.Key)
	if out == nil {
		return fmt.Errorf("decrypt error=%s", err.Error())
	}
	msg.Data = out
	return nil
}

func (c *Crypto) Outbound(s *session.Session, msg nano.Message) error {
	out := xxtea.Encrypt(msg.Data, c.Key)
	msg.Data = []byte(base64.StdEncoding.EncodeToString(out))
	return nil
}
