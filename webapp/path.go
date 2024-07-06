package webapp

type path struct {
	Display string
	Edit    string
}

func newPath(baseurl string) path {
	baseEditUrl := baseurl + "/edit"
	return path{
		Display: baseurl,
		Edit:    baseEditUrl,
	}
}
