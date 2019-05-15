package crypto

import (
	"encoding/base64"
	"fmt"

	"zlab/library/nano"
	"zlab/library/nano/session"
	"github.com/xxtea/xxtea-go/xxtea"
)

var xxteaKey = []byte("7AEC4MA152BQE9HWQ7KB")

//Crypto ..
type Crypto struct {
	key []byte
}

//NewCrypto ..
func NewCrypto() *Crypto {
	return &Crypto{xxteaKey}
}

//Inbound ..
func (c *Crypto) Inbound(s *session.Session, msg nano.Message) error {
	out, err := base64.StdEncoding.DecodeString(string(msg.Data))
	if err != nil {
		//logger.Errorf("Inbound Error=%s, In=%s", err.Error(), string(msg.Data))
		return err
	}

	out = xxtea.Decrypt(out, c.key)
	if out == nil {
		return fmt.Errorf("decrypt error=%s", err.Error())
	}
	msg.Data = out
	return nil
}

//Outbound ..
func (c *Crypto) Outbound(s *session.Session, msg nano.Message) error {
	out := xxtea.Encrypt(msg.Data, c.key)
	msg.Data = []byte(base64.StdEncoding.EncodeToString(out))
	return nil
}
