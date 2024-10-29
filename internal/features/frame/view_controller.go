package frame

import (
	"framer/internal/api"
	"framer/internal/rpc"
	"framer/internal/view"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func FrameViewHandler(r chi.Router) {
	r.Get(view.FramePath, handleGetFrame)
}

func handleGetFrame(w http.ResponseWriter, r *http.Request) {
	user := api.GetUser(r)

	frames, err := rpc.Client.Frame.ListFrames(view.ContextWithToken(r), nil)
	if err != nil {
		view.Render(w, r, view.Message("Something went wrong"), user)
	}

	resp := make([]*GetFrameResponseDto, len(frames.Frames))

	for i, item := range frames.Frames {
		itemResponse := &GetFrameResponseDto{
			ID:          uuid.MustParse(item.Id),
			Title:       item.Title,
			Description: item.Description,
			UserID:      user.ID,
			FrameStatus: int(item.FrameStatus),
		}
		if item.File != nil {
			itemResponse.FileID = uuid.NullUUID{Valid: true, UUID: uuid.MustParse(item.File.Id)}
		}
		resp[i] = itemResponse
	}

	view.Render(w, r, indexPage(resp), user)
}
