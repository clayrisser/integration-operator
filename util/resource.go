package util

type ResourceUtil struct{}

func NewResourceUtil() *ResourceUtil {
	return &ResourceUtil{}
}

func (u *ResourceUtil) CreateResource(resource string) {}

func (u *ResourceUtil) UpdateResource(resource string) {}
