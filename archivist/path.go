package archivist

type path struct {
	Display string
	Edit    string
	Save    string
}

func newPath(baseurl string) path {
	baseEditUrl := baseurl + "/edit"
	return path{
		Display: baseurl,
		Edit:    baseEditUrl,
	}
}
