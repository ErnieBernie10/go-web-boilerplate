package api

const (
	LoginApiPath    = "/api/auth/login"
	RefreshApiPath  = "/api/auth/refresh"
	RegisterApiPath = "/api/auth/register"

	GetFramesApiPath    = "/api/frames"
	GetFrameApiPath     = "/api/frames/{id}"
	PostFrameApiPath    = "/api/frames"
	PutFrameApiPath     = "/api/frames/{id}"
	DeleteFrameApiPath  = "/api/frames/{id}"
	UploadFileApiPath   = "/api/file/{filename}"
	DownloadFileApiPath = "/api/file/{id}"
)
