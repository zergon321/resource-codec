package codec

// AnimationKey adds an animation prefix
// to the key so it's identified as the
// one storing an animation.
func AnimationKey(key string) string {
	return "anim:" + key
}

// TagKey adds a tag prefix to the key
// so it's identified as the one storing
// other key that match the tag.
func TagKey(key string) string {
	return "tag:" + key
}
