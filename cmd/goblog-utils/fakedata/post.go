package fakedata

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/database"
	"github.com/misterabdul/goblog-server/internal/models"
	"github.com/misterabdul/goblog-server/internal/repositories"
	customMongo "github.com/misterabdul/goblog-server/pkg/mongo"
	"github.com/misterabdul/goblog-server/pkg/utils"
)

func GeneratePosts(ctx context.Context) {
	var (
		dbConn            *mongo.Database
		repository        *repositories.PostRepository
		contentRepository *repositories.PostContentRepository
		post              *models.PostModel
		postContent       *models.PostContentModel
		postId            primitive.ObjectID
		now               = primitive.NewDateTimeFromTime(time.Now())
		err               error
	)

	if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
		log.Fatal(err)
	}
	defer dbConn.Client().Disconnect(ctx)
	repository = repositories.NewPostRepository(dbConn)
	contentRepository = repositories.NewPostContentRepository(dbConn)
	for i := 0; i < 200; i++ {
		postId = primitive.NewObjectID()
		post = &models.PostModel{
			UID:                postId,
			Slug:               fmt.Sprintf("lorem-ipsum-%d", i),
			Title:              fmt.Sprintf("Lorem Ipsum %d", i),
			FeaturingImagePath: "./statics/images/image-example.jpg",
			Description:        lipsumParagraph(),
			Categories:         []models.CategoryCommonModel{},
			Tags:               []string{"lorem", "ipsum", "dolor", "sit", "amet"},
			CommentCount:       0,
			PublishedAt:        randNilOrValue(now),
			CreatedAt:          now,
			UpdatedAt:          now,
			DeletedAt:          nil,
			Author: models.UserCommonModel{
				FirstName: "Super Admin",
				Username:  "superadmin",
				Email:     "superadmin@example.com"}}
		postContent = &models.PostContentModel{
			UID:     postId,
			Content: lipsumMarkdown()}
		if err = customMongo.Transaction(ctx, dbConn, false,
			func(sCtx context.Context, dbConn *mongo.Database) (sErr error) {
				if sErr = repository.Save(
					sCtx, post,
				); sErr != nil {
					return sErr
				}
				if sErr = contentRepository.Save(
					sCtx, postContent,
				); sErr != nil {
					return sErr
				}

				return nil
			},
		); err != nil {
			log.Fatal(err)
		}
	}
	utils.ConsolePrintlnGreen("Generated 200 dummy posts.")
}

func randNilOrValue(value interface{}) interface{} {
	if rand.Int()%2 == 0 {
		return value
	}
	return nil
}

func lipsumParagraph() string {
	return "Lorem ipsum dolor sit amet, consectetur adipiscing elit. " +
		"Proin ipsum libero, laoreet at sem vel, ornare dapibus enim. " +
		"Curabitur vestibulum eros ac nisi lobortis, at sollicitudin velit congue. " +
		"In scelerisque feugiat nisi, ut blandit nisi tempus non. " +
		"Vivamus lacinia nisi sit amet aliquam fermentum. " +
		"Aliquam feugiat dui sed dolor accumsan, vel ultrices velit vehicula. " +
		"Mauris elit sapien, interdum in ante ac, iaculis placerat orci. " +
		"Nulla facilisi. Praesent ac auctor arcu. " +
		"Mauris aliquet ultricies enim, viverra aliquet neque lobortis ullamcorper. " +
		"Sed volutpat facilisis lacus nec porttitor. " +
		"Duis feugiat nibh euismod, tincidunt mauris et, laoreet dolor. " +
		"Praesent eu dolor et nisl finibus venenatis. " +
		"Aenean metus eros, malesuada nec diam sit amet, tincidunt facilisis justo."
}

func lipsumMarkdown() string {
	return `
# Headings

---

# Heading 1

## Heading 2

### Heading 3

#### Heading 4

##### Heading 5

###### Heading 6

# Paragraph

---

Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec consequat dictum nulla, ac convallis sapien sodales vel. Mauris quis ullamcorper metus. Sed luctus erat at mauris fringilla vestibulum. Etiam fringilla urna nec scelerisque dignissim. Aenean sit amet risus quis magna lacinia placerat. Praesent condimentum euismod sodales. Ut pretium sagittis velit, elementum faucibus nulla blandit sit amet. Curabitur viverra pulvinar nibh ac porta. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Sed porttitor, ipsum eu varius posuere, erat nisi mollis urna, nec feugiat felis turpis eu leo.
Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec consequat dictum nulla, ac convallis sapien sodales vel. Mauris quis ullamcorper metus. Sed luctus erat at mauris fringilla vestibulum. Etiam fringilla urna nec scelerisque dignissim. Aenean sit amet risus quis magna lacinia placerat. Praesent condimentum euismod sodales. Ut pretium sagittis velit, elementum faucibus nulla blandit sit amet. Curabitur viverra pulvinar nibh ac porta. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Sed porttitor, ipsum eu varius posuere, erat nisi mollis urna, nec feugiat felis turpis eu leo.
Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec consequat dictum nulla, ac convallis sapien sodales vel. Mauris quis ullamcorper metus. Sed luctus erat at mauris fringilla vestibulum. Etiam fringilla urna nec scelerisque dignissim. Aenean sit amet risus quis magna lacinia placerat. Praesent condimentum euismod sodales. Ut pretium sagittis velit, elementum faucibus nulla blandit sit amet. Curabitur viverra pulvinar nibh ac porta. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Sed porttitor, ipsum eu varius posuere, erat nisi mollis urna, nec feugiat felis turpis eu leo.

# Bold

---

**Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec consequat dictum nulla, ac convallis sapien sodales vel. Mauris quis ullamcorper metus. Sed luctus erat at mauris fringilla vestibulum. Etiam fringilla urna nec scelerisque dignissim. Aenean sit amet risus quis magna lacinia placerat. Praesent condimentum euismod sodales. Ut pretium sagittis velit, elementum faucibus nulla blandit sit amet. Curabitur viverra pulvinar nibh ac porta. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Sed porttitor, ipsum eu varius posuere, erat nisi mollis urna, nec feugiat felis turpis eu leo.**

# Italic

---

_Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec consequat dictum nulla, ac convallis sapien sodales vel. Mauris quis ullamcorper metus. Sed luctus erat at mauris fringilla vestibulum. Etiam fringilla urna nec scelerisque dignissim. Aenean sit amet risus quis magna lacinia placerat. Praesent condimentum euismod sodales. Ut pretium sagittis velit, elementum faucibus nulla blandit sit amet. Curabitur viverra pulvinar nibh ac porta. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Sed porttitor, ipsum eu varius posuere, erat nisi mollis urna, nec feugiat felis turpis eu leo._

# Strikethrough

---

~Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec consequat dictum nulla, ac convallis sapien sodales vel. Mauris quis ullamcorper metus. Sed luctus erat at mauris fringilla vestibulum. Etiam fringilla urna nec scelerisque dignissim. Aenean sit amet risus quis magna lacinia placerat. Praesent condimentum euismod sodales. Ut pretium sagittis velit, elementum faucibus nulla blandit sit amet. Curabitur viverra pulvinar nibh ac porta. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Sed porttitor, ipsum eu varius posuere, erat nisi mollis urna, nec feugiat felis turpis eu leo.~

# Blockquote

---

> Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec consequat dictum nulla, ac convallis sapien sodales vel. Mauris quis ullamcorper metus. Sed luctus erat at mauris fringilla vestibulum. Etiam fringilla urna nec scelerisque dignissim. Aenean sit amet risus quis magna lacinia placerat. Praesent condimentum euismod sodales. Ut pretium sagittis velit, elementum faucibus nulla blandit sit amet. Curabitur viverra pulvinar nibh ac porta. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Sed porttitor, ipsum eu varius posuere, erat nisi mollis urna, nec feugiat felis turpis eu leo.

# Ordered List

---

1. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
2. Donec consequat dictum nulla, ac convallis sapien sodales vel.
3. Mauris quis ullamcorper metus.
4. Sed luctus erat at mauris fringilla vestibulum.
5. Etiam fringilla urna nec scelerisque dignissim.
6. Aenean sit amet risus quis magna lacinia placerat.
7. Praesent condimentum euismod sodales.
8. Ut pretium sagittis velit, elementum faucibus nulla blandit sit amet.
9. Curabitur viverra pulvinar nibh ac porta.
10. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae;
11. Sed porttitor, ipsum eu varius posuere, erat nisi mollis urna, nec feugiat felis turpis eu leo.

# Unordered List

---

- Lorem ipsum dolor sit amet, consectetur adipiscing elit.
- Donec consequat dictum nulla, ac convallis sapien sodales vel.
- Mauris quis ullamcorper metus.
- Sed luctus erat at mauris fringilla vestibulum.
- Etiam fringilla urna nec scelerisque dignissim.
- Aenean sit amet risus quis magna lacinia placerat.
- Praesent condimentum euismod sodales.
- Ut pretium sagittis velit, elementum faucibus nulla blandit sit amet.
- Curabitur viverra pulvinar nibh ac porta.
- Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae;
- Sed porttitor, ipsum eu varius posuere, erat nisi mollis urna, nec feugiat felis turpis eu leo.

# Table

---

|                    Lorem ipsum dolor sit amet, consectetur adipiscing elit.                     |                    Lorem ipsum dolor sit amet, consectetur adipiscing elit.                     |
| :---------------------------------------------------------------------------------------------: | :---------------------------------------------------------------------------------------------: |
|                 Donec consequat dictum nulla, ac convallis sapien sodales vel.                  |                 Donec consequat dictum nulla, ac convallis sapien sodales vel.                  |
|                                 Mauris quis ullamcorper metus.                                  |                                 Mauris quis ullamcorper metus.                                  |
|                         Sed luctus erat at mauris fringilla vestibulum.                         |                         Sed luctus erat at mauris fringilla vestibulum.                         |
|                         Etiam fringilla urna nec scelerisque dignissim.                         |                         Etiam fringilla urna nec scelerisque dignissim.                         |
|                       Aenean sit amet risus quis magna lacinia placerat.                        |                       Aenean sit amet risus quis magna lacinia placerat.                        |
|                              Praesent condimentum euismod sodales.                              |                              Praesent condimentum euismod sodales.                              |
|              Ut pretium sagittis velit, elementum faucibus nulla blandit sit amet.              |              Ut pretium sagittis velit, elementum faucibus nulla blandit sit amet.              |
|                            Curabitur viverra pulvinar nibh ac porta.                            |                            Curabitur viverra pulvinar nibh ac porta.                            |
|     Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae;     |     Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae;     |
| Sed porttitor, ipsum eu varius posuere, erat nisi mollis urna, nec feugiat felis turpis eu leo. | Sed porttitor, ipsum eu varius posuere, erat nisi mollis urna, nec feugiat felis turpis eu leo. |

# Checkboxes

---

- [ ] Lorem ipsum dolor sit amet, consectetur adipiscing elit.
- [ ] Lorem ipsum dolor sit amet, consectetur adipiscing elit.
- [ ] Lorem ipsum dolor sit amet, consectetur adipiscing elit.
- [x] Lorem ipsum dolor sit amet, consectetur adipiscing elit.
- [x] Lorem ipsum dolor sit amet, consectetur adipiscing elit.
- [x] Lorem ipsum dolor sit amet, consectetur adipiscing elit.

# Code Block

---

` + "```c" + `
  #include <stdio.h>

  int main(int argc, char** argv) {
    printf("Hello world\n");
    return 0;
  }
` + "```" + `

# Image

---

![Example Image Description](./statics/images/image-example.jpg "Example Image Title")

# :fire: Emojies :tada:

---

:smirk: :heart_eyes: :kissing_heart: :kissing_closed_eyes: :flushed: :relieved: :satisfied: :grin: :wink: :stuck_out_tongue_winking_eye: :stuck_out_tongue_closed_eyes: :grinning: :kissing: :kissing_smiling_eyes: :stuck_out_tongue:

	`
}
