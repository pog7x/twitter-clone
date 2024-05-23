package fixtures

import (
	"context"

	"github.com/samber/do"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"twitter-clone/internal/repository/dbrepository"
)

type Service struct {
	logger     *logrus.Logger
	trendsRepo *dbrepository.TrendRepository
	usersRepo  *dbrepository.UserRepository
}

func NewFixturesService(i *do.Injector) (*Service, error) {
	return &Service{
		logger:     do.MustInvoke[*logrus.Logger](i),
		trendsRepo: do.MustInvoke[*dbrepository.TrendRepository](i),
		usersRepo:  do.MustInvoke[*dbrepository.UserRepository](i),
	}, nil
}

func (s *Service) CreateTrends(ctx context.Context) {
	_, _ = s.trendsRepo.Create(ctx, dbrepository.CreateTrendPayload{ID: 1, Name: "#Python", TweetCount: 413765})
	_, _ = s.trendsRepo.Create(ctx, dbrepository.CreateTrendPayload{ID: 2, Name: "#Golang", TweetCount: 1901231})
	_, _ = s.trendsRepo.Create(ctx, dbrepository.CreateTrendPayload{ID: 3, Name: "#Rust", TweetCount: 67234432})
}

func (s *Service) CreateUsers(ctx context.Context) {
	pass1, _ := bcrypt.GenerateFromPassword([]byte("pass123"), bcrypt.DefaultCost)
	pass2, _ := bcrypt.GenerateFromPassword([]byte("pass456"), bcrypt.DefaultCost)

	_, _ = s.usersRepo.Create(ctx, dbrepository.CreateUserPayload{
		ID:          1,
		Password:    pass1,
		Username:    "pog7x",
		Name:        "developer",
		Website:     "https://github.com/pog7x",
		PicCover:    "https://ideogram.ai/api/images/direct/T91kUQhETeyPyiyqCOfwcQ.png",
		Pic:         "https://avataaars.io/?avatarStyle=Circle&topType=LongHairFrida&accessoriesType=Round&facialHairType=Blank&clotheType=ShirtVNeck&clotheColor=Gray01&eyeType=Happy&eyebrowType=RaisedExcitedNatural&mouthType=Smile&skinColor=Pale",
		Description: "Just a developer that interested in Golang.",
	})

	_, _ = s.usersRepo.Create(ctx, dbrepository.CreateUserPayload{
		ID:          2,
		Password:    pass2,
		Username:    "pog8x",
		Name:        "web-developer",
		Website:     "https://github.com/pog7x",
		PicCover:    "https://ideogram.ai/api/images/direct/T91kUQhETeyPyiyqCOfwcQ.png",
		Pic:         "https://getavataaars.com/?accessoriesType=Prescription02&avatarStyle=Circle&clotheColor=Blue01&clotheType=GraphicShirt&eyeType=Cry&eyebrowType=Default&facialHairColor=Red&facialHairType=BeardLight&graphicType=SkullOutline&hairColor=Auburn&hatColor=PastelYellow&mouthType=Grimace&skinColor=Tanned&topType=LongHairStraight",
		Description: "Just a developer that interested in JavaScript.",
	})
}
