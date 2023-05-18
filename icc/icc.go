package icc

type IccEmbedder interface {
	EmbedIccProfile(icc_profile []byte) error
}

// Embed ICC profile to image.
func EmbedIccProfile(profile_name string, target IccEmbedder) error {

	raw_profile_bytes, err := get_icc_profile(profile_name) // Get ICC profile bytes.
	if err != nil {
		return err
	}

	err = target.EmbedIccProfile(raw_profile_bytes) // Embed ICC profile to target image.
	return err
}
