package unstructured

import "strings"

// Encoding is a type that represents an encoding.
type Encoding string

// Encoding constants.
const (
	EncodingUTF8          Encoding = "utf_8"
	EncodingISO88591      Encoding = "iso_8859_1"
	EncodingISO88596      Encoding = "iso_8859_6"
	EncodingISO88598      Encoding = "iso_8859_8"
	EncodingASCII         Encoding = "ascii"
	EncodingBig5          Encoding = "big5"
	EncodingUTF16         Encoding = "utf_16"
	EncodingUTF16Be       Encoding = "utf_16_be"
	EncodingUTF16Le       Encoding = "utf_16_le"
	EncodingUTF32         Encoding = "utf_32"
	EncodingUTF32Be       Encoding = "utf_32_be"
	EncodingUTF32Le       Encoding = "utf_32_le"
	EncodingEUCJIS2004    Encoding = "euc_jis_2004"
	EncodingEUCJISX0213   Encoding = "euc_jisx0213"
	EncodingEUCJP         Encoding = "euc_jp"
	EncodingEUCKR         Encoding = "euc_kr"
	EncodingGb18030       Encoding = "gb18030"
	EncodingSHIFTJIS      Encoding = "shift_jis"
	EncodingSHIFTJIS2004  Encoding = "shift_jis_2004"
	EncodingSHIFTJISX0213 Encoding = "shift_jisx0213"
)

// String implements the fmt.Stringer interface, canonicalizing the encoding name.
func (e Encoding) String() string {
	s := strings.TrimSpace(string(e))
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, "_", "-")

	switch s {
	case "iso_8859_6_i", "iso_8859_8_i",
		"iso_8859_6_e", "iso_8859_8_e":
		s = s[:len(s)-2]
	}

	return s
}
