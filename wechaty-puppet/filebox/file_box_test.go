package filebox

import (
	"errors"
	"strings"
	"testing"
)

func TestFromJSON(t *testing.T) {
	t.Run("FromJSON json.Unmarshal err", func(t *testing.T) {
		if err := FromJSON("abcd").Error(); !strings.Contains(err.Error(), "FromJSON json.Unmarshal") {
			t.Error(err)
		}
	})
	t.Run("FromJSON invalid value boxType", func(t *testing.T) {
		if err := FromJSON("{}").Error(); !strings.Contains(err.Error(), "FromJSON invalid value boxType") {
			t.Error(err)
		}
	})
	t.Run("FromJSON success", func(t *testing.T) {
		jsonText := `{"base64":"RmlsZUJveEJhc2U2NAo=","boxType":1,"md5":"","metadata":null,"name":"test.txt","size":14,"type":1}`
		if err := FromJSON(jsonText).Error(); err != nil {
			t.Error(err)
		}
	})
}

func TestFromBase64(t *testing.T) {
	t.Run("FromBase64 no base64 data", func(t *testing.T) {
		if err := FromBase64("").Error(); !errors.Is(err, ErrNoBase64Data) {
			t.Error(err)
		}
	})
	t.Run("FromBase64 success", func(t *testing.T) {
		fileBox := FromBase64("RmlsZUJveEJhc2U2NAo=")
		if err := fileBox.Error(); err != nil {
			t.Error(err)
		}
	})
}

func TestFromUrl(t *testing.T) {
	t.Run("FromUrl no url", func(t *testing.T) {
		if err := FromUrl("").Error(); !errors.Is(err, ErrNoUrl) {
			t.Error(err)
		}
	})
	t.Run("FromUrl success", func(t *testing.T) {
		fileBox := FromUrl("https://github.com//dchaofei.jpg?t=123")
		if err := fileBox.Error(); err != nil {
			t.Error(err)
		}
		want := "dchaofei.jpg"
		if fileBox.Name != want {
			t.Errorf("got %s want %s", fileBox.Name, want)
		}
	})
}

func TestFromFile(t *testing.T) {
	t.Run("FromFile no path", func(t *testing.T) {
		if err := FromFile("").Error(); !errors.Is(err, ErrNoPath) {
			t.Error(err)
		}
	})
	t.Run("FromFile success", func(t *testing.T) {
		fileBox := FromFile("testdata/dchaofei.txt")
		if err := fileBox.Error(); err != nil {
			t.Error(err)
		}

		wantName := "dchaofei.txt"
		if wantName != fileBox.Name {
			t.Errorf("got %s want %s", fileBox.Name, wantName)
		}
	})
}

func TestFromQRCode(t *testing.T) {
	t.Run("FromQRCode no QR code", func(t *testing.T) {
		if err := FromQRCode("").Error(); !errors.Is(err, ErrNoQRCode) {
			t.Error(err)
		}
	})
	t.Run("FromQRCode success", func(t *testing.T) {
		fileBox := FromQRCode("hello")
		if err := fileBox.Error(); err != nil {
			t.Error(err)
		}
	})
}

func TestFromUuid(t *testing.T) {
	t.Run("FromUuid no uuid", func(t *testing.T) {
		if err := FromUuid("").Error(); !errors.Is(err, ErrNoUuid) {
			t.Error(err)
		}
	})
	t.Run("FromUuid success", func(t *testing.T) {
		fileBox := FromUuid("xxx-xxx-xxx")
		if err := fileBox.Error(); err != nil {
			t.Error(err)
		}
	})
}

func TestFileBox_ToJSON(t *testing.T) {
	t.Run("ToJSON success", func(t *testing.T) {
		const base64Encode = "RmlsZUJveEJhc2U2NAo="
		const base64Filename = "test.txt"
		const want = `{"base64":"RmlsZUJveEJhc2U2NAo=","boxType":1,"md5":"","metadata":null,"name":"test.txt","size":14,"type":1}`
		jsonString, err := FromBase64(base64Encode, WithName(base64Filename)).ToJSON()
		if err != nil {
			t.Error(err)
		}
		if jsonString != want {
			t.Errorf("got【%s】, want [%s]", jsonString, want)
		}

		newBase64, err := FromJSON(want).ToBase64()
		if err != nil {
			t.Error(err)
		}
		if newBase64 != base64Encode {
			t.Errorf("got【%s】, want [%s]", newBase64, base64Encode)
		}
	})

	t.Run("ToJSON for not supported type", func(t *testing.T) {
		if _, err := FromFile("testdata/dchaofei.txt").ToJSON(); err != ErrToJSON {
			t.Error(err)
		}
	})
}
