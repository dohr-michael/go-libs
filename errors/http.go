package errors

import (
	"net/http"
	"github.com/go-chi/render"
	"gopkg.in/asaskevich/govalidator.v4"
)

// HttpResponse renderer type for handling all sorts of errors.
//
// In the best case scenario, the excellent github.com/pkg/errors package
// helps reveal information on the error, setting it on Err, and in the Render()
// method, using it to set the application-specific error code in AppCode.
type HttpResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *HttpResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ToRenderer(err error) render.Renderer {
	switch err.(type) {
	case *notFoundError:
		return NotFoundRenderer
	case *illegalEntityError:
		return InvalidRequestRenderer(err)
	case *govalidator.Error:
		return InvalidRequestRenderer(err)
	case *govalidator.Errors:
		return InvalidRequestRenderer(err)
	}
	return InternalServerErrorRenderer(err)
}

func InvalidRequestRenderer(err error) render.Renderer {
	return &HttpResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func UnprocessableEntityRenderer(err error) render.Renderer {
	return &HttpResponse{
		Err:            err,
		HTTPStatusCode: http.StatusUnprocessableEntity,
		StatusText:     "Unprocessable Entity",
		ErrorText:      err.Error(),
	}
}

var NotFoundRenderer = &HttpResponse{HTTPStatusCode: 404, StatusText: "Resource not found."}

func InternalServerErrorRenderer(err error) render.Renderer {
	return &HttpResponse{
		Err:            err,
		HTTPStatusCode: http.StatusInternalServerError,
		StatusText:     "Internal Server Error",
		ErrorText:      err.Error(),
	}
}
