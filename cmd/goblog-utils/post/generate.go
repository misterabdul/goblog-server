package post

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/misterabdul/goblog-server/internal/database"
	"github.com/misterabdul/goblog-server/internal/models"
	"github.com/misterabdul/goblog-server/internal/repositories"
	"github.com/misterabdul/goblog-server/pkg/utils"
)

func Generate(ctx context.Context) {
	var (
		dbConn      *mongo.Database
		post        *models.PostModel
		postContent *models.PostContentModel
		postId      primitive.ObjectID
		now         = primitive.NewDateTimeFromTime(time.Now())
		err         error
	)

	if dbConn, err = database.GetDBConnDefault(ctx); err != nil {
		log.Fatal(err)
	}
	defer dbConn.Client().Disconnect(ctx)

	for i := 0; i < 100; i++ {
		postId = primitive.NewObjectID()
		post = &models.PostModel{
			UID:                postId,
			Slug:               fmt.Sprintf("lorem-ipsum-%d", i),
			Title:              fmt.Sprintf("Lorem Ipsum %d", i),
			Description:        "Lorem ipsum dolor sit amet",
			FeaturingImagePath: "./statics/images/image-example.jpg",
			Categories:         []models.CategoryCommonModel{},
			Tags:               []string{"lorem", "ipsum", "dolor", "sit", "amet"},
			Author: models.UserCommonModel{
				FirstName: "Super Admin",
				Username:  "superadmin",
				Email:     "superadmin@example.com",
			},
			PublishedAt: now,
			CreatedAt:   now,
			UpdatedAt:   now,
			DeletedAt:   nil,
		}
		postContent = &models.PostContentModel{
			UID:     postId,
			Content: lipsumMarkdown(),
		}
		if err = repositories.CreatePost(ctx, dbConn, post, postContent); err != nil {
			log.Fatal(err)
		}
	}
	utils.ConsolePrintlnGreen("Generated 100 dummy posts.")
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
` +
		"```c" + `
  #include <stdio.h>

  int main(int argc, char** argv) {
    printf("Hello world\n");
    return 0;
  }` +
		"```" + `

# Image

---

![Example Image Description](./statics/images/image-example.jpg "Example Image Title")

# :fire: Emojies :tada:

---

:smirk: :heart_eyes: :kissing_heart: :kissing_closed_eyes: :flushed: :relieved: :satisfied: :grin: :wink: :stuck_out_tongue_winking_eye: :stuck_out_tongue_closed_eyes: :grinning: :kissing: :kissing_smiling_eyes: :stuck_out_tongue:

	`
}