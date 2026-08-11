package api

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/kataras/iris/v12"

	"twitter-clone/config"
	"twitter-clone/internal/domain/response"
	"twitter-clone/internal/repository/dbrepository"
)

// MaxUploadSize caps a single upload. The request body is capped as well, so an
// oversized file is rejected before it is written to disk.
const MaxUploadSize = 8 << 20 // 8 MiB

// allowedImageTypes maps a sniffed content type to the extension the file is
// stored under. Types outside this list are rejected; in particular SVG is not
// accepted because it can carry scripts.
var allowedImageTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/gif":  ".gif",
	"image/webp": ".webp",
}

var errUnsupportedMedia = errors.New("unsupported file type: expected jpeg, png, gif or webp")

type Media struct {
	MediaID uint64 `json:"media_id"`
	Link    string `json:"link"`
}

func CreateMediaHandler(ctx iris.Context, mediaRepo *dbrepository.MediaRepository, cfg *config.Config) {
	userID, ok := authenticatedUserID(ctx)
	if !ok {
		return
	}

	ctx.SetMaxRequestBodySize(MaxUploadSize)

	file, header, err := ctx.FormFile("file")
	if err != nil {
		response.SendErrorResponse(ctx, iris.StatusBadRequest, "no file to save")
		return
	}
	defer file.Close()

	if header.Size > MaxUploadSize {
		response.SendErrorResponse(
			ctx,
			iris.StatusRequestEntityTooLarge,
			fmt.Sprintf("file must be at most %d MiB", MaxUploadSize>>20),
		)
		return
	}

	extension, err := detectImageExtension(file)
	if err != nil {
		if errors.Is(err, errUnsupportedMedia) {
			response.SendErrorResponse(ctx, iris.StatusBadRequest, err.Error())
			return
		}
		response.SendInternalError(ctx, err)
		return
	}

	// The stored name is generated server side. The client-supplied filename is
	// never used as a path component, so it cannot escape the uploads directory
	// or overwrite somebody else's file.
	name := uuid.NewString() + extension

	link, err := storeUpload(file, cfg.UploadsDirPath, name)
	if err != nil {
		response.SendInternalError(ctx, err)
		return
	}

	media, err := mediaRepo.Create(ctx, dbrepository.CreateMediaPayload{Link: link, OwnerID: userID})
	if err != nil {
		sendRepositoryError(ctx, err)
		return
	}

	// The link is returned as well so the client can reference the file
	// directly, e.g. when setting it as an avatar or a cover image.
	response.SendOkResponse(ctx, Media{MediaID: media.ID, Link: media.Link})
}

// detectImageExtension sniffs the leading bytes of the upload and rewinds it,
// so the stored type is decided by the content rather than by the name or the
// client-supplied Content-Type header.
func detectImageExtension(file multipart.File) (string, error) {
	head := make([]byte, 512)

	n, err := file.Read(head)
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}

	if _, err = file.Seek(0, io.SeekStart); err != nil {
		return "", err
	}

	contentType, _, _ := strings.Cut(http.DetectContentType(head[:n]), ";")

	extension, ok := allowedImageTypes[strings.TrimSpace(contentType)]
	if !ok {
		return "", errUnsupportedMedia
	}

	return extension, nil
}

// storeUpload writes the file into dir under name and returns the public link.
func storeUpload(file multipart.File, dir, name string) (string, error) {
	destination, err := os.OpenFile(filepath.Join(dir, name), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return "", err
	}
	defer destination.Close()

	if _, err = io.Copy(destination, io.LimitReader(file, MaxUploadSize)); err != nil {
		return "", err
	}

	return path.Join("/", filepath.ToSlash(dir), name), nil
}
