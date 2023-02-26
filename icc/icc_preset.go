package icc

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"strings"
)

var display_p3 = `H4sIAOOBGmEC/2NgYFJJLCjIYWFgYMjNKykKcndSiIiMUmB/yMAOhLwMYgwKicnFBY4BAT5AJQww
GhV8u8bACKIv64LMOiU1tUm1XsDXYqbw1YuvRJsY8AOulNTiZCD9B4hTkwuKShgYGFOAbOXykgIQ
uwPIFikCOgrIngNip0PYG0DsJAj7CFhNSJAzkH0DyFZIzkgEmsH4A8jWSUIST0diQ+0FAW6XzOKC
nMRKhQBjBqqDktSKEhDtnF9QWZSZnlGi4AgMpVQFz7xkPR0FIwNDcwYGUJhDVH8OBIclo9gZhFjz
fQYG2/3////fjRDz2s/AsBGok2snQkzDgoFBkJuB4cTOgsSiRLAQMxAzpaUxMHxazsDAG8nAIHwB
qCe6OM3YCCzPyOPEwMB67///z2oMDOyTGRj+Tvj///ei////LgZqvsPAcCAPAJ+/MZ4kAgAA
`

var dci_p3 = `H4sIAKSEGmEC/2NgYFJJLCjIYWFgYMjNKykKcndSiIiMUmC/z8DFwMfAy8DBYJGYXFzgGBDgA1TC
AKNRwbdrDIwg+rIuyKxgR76Vv87e/B1+pqSt9/1hewb8gCsltTgZSP8B4vLkgqISBgZGIGZQLi8p
ALFnANkiRUBHAdlrQOx0CPsAiJ0EYV8BqwkJcgayXwDZAskZiSlA9g8gWycJSTwdiQ21FwRkg30D
QlwVggIUTIwNdY10jQwMzBVcnD0VNAKMNRmoDEpSK0D+Y3DOL6gsykzPKFFwBIZaqoJnXrKejoKR
gaEpAwMoDiCqn7JAwvZLKEKsppSBwYrj////pxFiQUA/bb3BwMAtiBBTrwR6E8g+EVCQWJQIdwDT
zFnFacZGYDaj0DMGBs73//9/62Ng4LMCxuWJ////hPz//w9oHmM0A8Od+wBS2VCSJAIAAA==
`

var romm_rgb = `H4sIABSFGmEC/2NgYDJJLCjIYWFgYMjNKykKcndSiIiMUmC/z8DFwMfAy8DBYJ6YXFzgGBDgA1TC
AKNRwbdrDIwg+rIuyKzCJ1HPFj9NFGo8wyVrNSdtLgN+wJWSWpwMpP8AcWlyQVEJAwMjEDMol5cU
gNgzgGyRIqCjgOw1IHY6hH0AxE6CsK+A1YQEOQPZL4BsheSMxBQGBiYOIFsnCUk8HYkNtRcEpIP8
fX0VgJ63UvAM9lcwMjIwstA1sjIyMDRmoDooSa0A+Y/BOb+gsigzPaNEwREYaqkKnnnJejoKQDtN
GRhAcYAetgixMyYMDJ57QSyEmNIcBoZtdqhiHOxAgg2o16EgsSgRLMQMxIxnzkLMBAEBENFQnGZs
BOYyojkWCx8AhPHgOTQCAAA=
`

var adobe_rgb = `H4sIAJfoGmEC/2NgYDJwdHFyZRJgYMjNKykKcndSiIiMUmC/wMDBwM0gzGDMYJ2YXFzgGBDgwwAE
efl5qQwY4Ns1BkYQfVkXZBYDaYAruaCoBEj/AWKjlNTiZAYGRgMgO7u8pAAozjgHyBZJygazN4DY
RSFBzkD2ESCbLx3CvgJiJ0HYT0DsIqAngOwvIPXpYDYTB9gcCFsGxC5JrQDZy+CcX1BZlJmeUaJg
ZGBgoOCYkp+UqhBcWVySmlus4JmXnF9UkF+UWJKaAlQLcR8YCEIUgkJMw9DS0kKTgcoAFA8Q1udA
cPgyip1BiCFAcmlRGZTJyGRMmI8wY44EA4P/UgYGlj8IMZNeBoYFOgwM/FMRYmqGDAwC+gwM++YA
AIzXb+0wAgAA
`

var srgb = `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
EDm0RmdN1pqraxXWJq09EyUZxYw6GY2ODotuiv7A9GPWMWdjvGKqY2ZYLqx9rOdsR3YZe4pjxynl
TMTaxZbGTsbZxe2Jm4p3iC+Pn+a6cCu5LxM8E2oS5hL9Eo8kLiSFJbUm45Kjk0/xZHiJvJ4UlZSs
lIFUg9SC1JE0m7S9aTN8b35DOpS+Jr1TQBX9TPUJdYVbhaMZ9hlVGW8zQzNPZkln8bL6svWzd2RP
5LjnfL0OtY61rjtXLXdz7uh6p/W1G6ANMRu6N2pszN84vslj09HNhM2Jm3/IM8krzXu9JWxLV75y
/qb8sa0eW5sLJAr4BcPbbLfVbEdt527v32G+Y/+OT4XswmtFJkXlRR+KWcXXvjL9quKrhZ2xO/tL
LEsO7sLs4u0a2u2w+2ipdGlO6dge3z3tZfSywrLXe6P2Xi1fVl6zj7BPuG+kwqeic7/m/l37P1TG
V96pcq5qrVaq3lE9d4B9YPCg48GWGuWaopr3h7iH7tZ61LbXadeVH8Yczjj8tD60vvdrxteNDQoN
RQ0fj/COjBwNPNrTaNXY2KTUVNIMNwubp45FHrv5jes3nS2GLbWttNai4+C48Pizb6O/HTrhfaL7
JONky3da31W3UdoK26H27PaZjviOkc7wzoFTK051d9l2tX1v9P2R02qnq87Inik5Szibf3bhXM65
2fOp56cvxF0Y647qfnBx1cXbPQE9/Ze8L1257H75Yq9T77krdldOX7W5euoa41rHdcvr7X0WfW0/
WPzQ1m/Z337D6kbnTeubXQPLB84OOgxeuOV66/Jtr9vX76y8MzAUMnR3OHJ45C777uS9pHsv72fc
n3+w6SH6YeEjqUflj5Ue1/2o92PriOXImVHX0b4nQU8ejLHGnv+U/tOH8fyn5KflE6oTjZNmk6en
3KduPlv9bPx56vP56YKfpX+ufqH74rtfHH/pm1k1M/6S/3Lh1+JX8q+OvF72unvWf/bxm+Q383OF
b+XfHn3HeNf7Puz9xHzmB+yHio96H7s+eX96uJC8sPAbUqUuGEgMAAA=
`

// Decode profile from compressed base64 string.
func decodeProfile(encoded_profile string) ([]byte, error) {

	var err error = nil

	// Decode base64 string.
	b64_dist := make([]byte, base64.StdEncoding.DecodedLen(len(encoded_profile)))
	data_len, err := base64.StdEncoding.Decode(b64_dist, []byte(encoded_profile))
	if err != nil {
		return nil, err
	}
	b64_dist = b64_dist[:data_len]

	var decoded_buf bytes.Buffer
	reader_uncompressed, err := gzip.NewReader(bytes.NewReader(b64_dist))
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(&decoded_buf, reader_uncompressed)
	if err != nil {
		return nil, err
	}

	return decoded_buf.Bytes(), nil
}

// Get profile by name.
func GetEmbeddedProfile(profile_name string) ([]byte, error) {
	var err error = nil
	var raw_profile []byte = nil

	switch strings.ToUpper(profile_name) {

	case "SRGB":
		raw_profile, err = decodeProfile(srgb)

	case "DISPLAY P3":
		raw_profile, err = decodeProfile(display_p3)

	case "DCI P3":
		raw_profile, err = decodeProfile(dci_p3)

	case "ADOBE RGB":
		raw_profile, err = decodeProfile(adobe_rgb)

	case "ROMM RGB":
		raw_profile, err = decodeProfile(romm_rgb)

	default:
		err = fmt.Errorf("profile not found")
	}

	return raw_profile, err
}

// Validate profile using length (0~3 bytes, big endian)
func validateProfile(raw_profile []byte) bool {
	length := binary.BigEndian.Uint32(raw_profile[0:4])
	return uint32(len(raw_profile)) == length
}

func main() {

}
