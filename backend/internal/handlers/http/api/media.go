package api

import (
	"fmt"

	"twitter-clone/internal/domain/response"
	"twitter-clone/internal/repository/dbrepository"

	"github.com/kataras/iris/v12"
)

func CreateMediaHandler(ctx iris.Context, mediaRepo *dbrepository.MediaRepository) {
	maxSize := ctx.Application().ConfigurationReadOnly().GetPostMaxMemory()

	err := ctx.Request().ParseMultipartForm(maxSize)
	if err != nil {
		ctx.StopWithError(iris.StatusInternalServerError, err)
		return
	}

	form := ctx.Request().MultipartForm
	files := form.File["file"]

	if len(files) < 1 {
		response.SendErrorResponse(ctx, iris.StatusBadRequest, "no files to save")
		return
	}

	file := files[0]
	filePath := "/uploads/" + file.Filename

	_, err = ctx.SaveFormFile(file, fmt.Sprintf(".%s", filePath))
	if err != nil {
		response.SendErrorResponse(ctx, iris.StatusBadRequest, err.Error())
		return
	}

	media, err := mediaRepo.Create(ctx, dbrepository.CreateMediaPayload{Link: filePath})
	if err != nil {
		response.SendErrorResponse(ctx, iris.StatusInternalServerError, err.Error())
		return
	}

	response.SendOkResponse(ctx, iris.Map{"media_id": media.ID})
}
