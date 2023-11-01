package token

type Audience string

var (
	USER_APP     Audience = "user_app"
	DRIVER_APP   Audience = "driver_app"
	ADMIN_PORTAL Audience = "admin_portal"
)

func ParseAudience(v string) *Audience {
	aud := Audience(v)
	switch aud {
	case USER_APP, DRIVER_APP, ADMIN_PORTAL:
		return &aud
	default:
		return nil
	}
}

type Platform string

const (
	WEBSITE Platform = "website"
	ANDROID Platform = "android"
	IOS     Platform = "ios"
)

func ParsePlatform(v string) *Platform {
	p := Platform(v)
	switch p {
	case WEBSITE, ANDROID, IOS:
		return &p
	default:
		return nil
	}
}
