---
title: "Your First Database"
chapter: 3

course_key: "JvSfm3LDdJuD6jxD1uwoC-D1MbUFY_nTnLLRk9dJarqHkFVsPTwaDIE2We6m0SXtNCZyxb5QWGRsTwdsYElXgQ"
key: "FP3LGh7I_V0WvSTRQU_2CLSqpROJRzMl5BJC5tiZzrQALpI-pA7NiS_Fetxfi2lltrRwVGaJd05kfe9hUD4wkg"
---

Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nunc nec dui a massa tempus sollicitudin. In lobortis auctor ullamcorper. Aliquam erat volutpat. Morbi tincidunt nunc eget risus pretium, non porta felis ultrices. Aenean quam eros, sollicitudin vestibulum posuere nec, tincidunt quis lacus. Morbi laoreet ex lacinia ipsum malesuada feugiat id et enim. Vivamus sapien dolor, dictum vitae facilisis eget, dapibus a quam. Quisque vel purus in urna sagittis efficitur. Etiam sodales, felis vel iaculis vestibulum, mauris velit consectetur enim, vel rhoncus tortor sem in est. Suspendisse rhoncus consequat mi euismod efficitur. Phasellus congue nibh vitae nulla tempus scelerisque. Nunc sit amet lorem mi. Integer ut est tempus, faucibus eros nec, convallis leo. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos.

Integer gravida massa a nunc finibus, sed pretium massa vestibulum. Etiam dictum eros id felis rutrum eleifend. Donec auctor diam eget porttitor luctus. Curabitur non felis purus. Vestibulum maximus ipsum sed orci facilisis convallis. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Maecenas et ultricies diam, vel lobortis nisi. Donec porttitor efficitur lacus a eleifend. Sed sed tristique urna. Phasellus metus ante, elementum eu ipsum id, scelerisque iaculis ipsum. Nulla suscipit augue eget nulla fermentum, sit amet pretium augue lobortis. Fusce non nibh sed ligula imperdiet convallis. Phasellus maximus pellentesque efficitur. Maecenas tincidunt a velit eu aliquet. In hac habitasse platea dictumst. Donec suscipit sed ipsum ac fermentum.

## Here is some code to see it render:

```golang
func (db *SQLiteDatabase) GetAllTutorialsPaginated(page, elements int) ([]*models.TutorialModel, error) {
	query := `SELECT id, title, slug, description, thumbnail_url, banner_url, content, published, author_id, file_checksum, file_key, created_at, updated_at FROM tutorials ORDER BY updated_at DESC, title ASC LIMIT ? OFFSET ?;`

	offset := (page - 1) * elements

	var tutorials []*models.TutorialModel

	rows, err := db.connection.Query(query, elements, offset)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var tutorial models.TutorialModel
		var published int

		if err := rows.Scan(&tutorial.ID, &tutorial.Title, &tutorial.Slug, &tutorial.Description, &tutorial.ThumbnailURL, &tutorial.BannerURL, &tutorial.Content, &published, &tutorial.AuthorID, &tutorial.FileChecksum, &tutorial.FileKey, &tutorial.CreatedAt, &tutorial.UpdatedAt); err != nil {
			return nil, err
		}

		if published == 1 {
			tutorial.Published = true
		}

		tutorials = append(tutorials, &tutorial)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tutorials, nil
}
```

Suspendisse nec hendrerit quam. Nam imperdiet arcu augue, nec fermentum ipsum viverra quis. Phasellus bibendum bibendum vestibulum. Cras molestie aliquet orci sed blandit. Vivamus non lacus dignissim, interdum augue quis, rhoncus nisl. Quisque risus odio, consectetur id nisl eu, lobortis venenatis mauris. Donec dictum vitae odio ac viverra. Suspendisse potenti. Cras ac mauris eget sapien semper eleifend id et dui.

Donec sagittis sapien eu diam sagittis mollis. Donec rhoncus nunc et nibh fringilla laoreet. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Nullam tincidunt odio eu ex semper semper. Pellentesque dapibus urna ac purus luctus facilisis. Proin facilisis sem at nisl blandit condimentum. In purus velit, dictum ut tempor varius, pharetra vitae arcu. Nulla vestibulum faucibus dui, eu ullamcorper sapien luctus eu. Nulla justo mi, rhoncus eget volutpat et, porta eget enim. Mauris sollicitudin tempus tincidunt. Cras pharetra consectetur venenatis. Phasellus fermentum sollicitudin velit, nec suscipit felis auctor tincidunt.

Duis eu nulla a augue venenatis gravida eget in elit. Nulla sit amet congue augue, id tristique lorem. Proin nec posuere dui, quis sollicitudin sapien. Maecenas mauris turpis, posuere ut lacus ut, luctus malesuada lorem. Aliquam molestie odio at mauris mollis ullamcorper. Aliquam vitae sem molestie, viverra quam gravida, rhoncus eros. Suspendisse non augue eget dolor condimentum condimentum. Nunc quis semper risus. Maecenas vel vestibulum mauris. Maecenas in velit et orci blandit pretium a sed lorem. Curabitur sodales hendrerit dui, non maximus lorem vehicula et. In porttitor rhoncus dolor eget congue. Donec elementum euismod mi sed elementum. Etiam nec finibus nunc.

Mauris at augue at justo commodo imperdiet nec quis magna. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nulla in ornare neque, ut sagittis mi. Nunc cursus sem tellus, quis vulputate mauris finibus non. Nulla velit quam, volutpat eu diam at, ullamcorper auctor ante. Nullam maximus a mauris in dignissim. Nullam et elementum felis. Sed nec elit porttitor, sodales nibh a, mattis enim. Vestibulum finibus mi dui, sed faucibus tortor ultricies sit amet. Curabitur eget enim dignissim, ultrices eros nec, rhoncus tellus. Sed porttitor aliquet eros a dapibus. Praesent aliquet sapien eros, vel rhoncus mi suscipit at.

Nulla eleifend erat enim, sed sagittis neque rhoncus ut. Ut maximus tincidunt massa sit amet egestas. Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas. Praesent pharetra purus elementum massa elementum, quis tristique enim eleifend. In id diam sit amet leo tempus gravida. Nunc quis lorem et lacus ultrices euismod. Integer porta porta tempor. Maecenas vitae turpis commodo, lacinia diam nec, interdum risus. In sed ultricies nunc, vel feugiat dui. Fusce posuere id neque quis semper. Vestibulum ac tristique enim, in aliquam turpis. In hac habitasse platea dictumst. Pellentesque scelerisque arcu nulla, consectetur luctus sapien euismod sit amet. Phasellus urna lacus, pellentesque et est a, lacinia aliquam mauris. Sed purus orci, semper nec ante condimentum, euismod pellentesque orci. Praesent in malesuada tellus.

## Here is another title.

In hac habitasse platea dictumst. Duis nulla nulla, gravida at mi at, cursus cursus tellus. Donec finibus finibus blandit. Aliquam lobortis non turpis ut lacinia. Phasellus mattis mollis risus, in finibus ligula blandit vel. Vestibulum ornare dictum justo non tristique. Ut tristique ullamcorper felis et accumsan. Ut in orci lacus. Suspendisse potenti. Aliquam vitae purus justo. Donec est metus, laoreet ut erat sed, pretium varius urna. Curabitur ullamcorper, lorem quis molestie gravida, massa augue varius velit, a dignissim metus nisi a ante. Praesent eget ex elementum, tincidunt sem facilisis, viverra justo. Aenean pretium eros nunc, eget hendrerit lacus tincidunt ut. Morbi tempor finibus faucibus. Aliquam sagittis, erat sed auctor pellentesque, ante elit pulvinar sapien, ac scelerisque nibh nisi dictum metus.

Integer ex quam, fermentum sit amet dictum non, varius id augue. Etiam aliquam lacus volutpat, lacinia justo at, molestie libero. In in nisl lectus. Vestibulum fringilla elit nec sapien malesuada scelerisque. Nulla iaculis leo porta, tincidunt purus id, viverra nisl. Integer a eleifend eros. Nam ullamcorper sollicitudin turpis, nec molestie ex commodo eu. Curabitur pharetra justo ut leo viverra dignissim. Sed eget aliquam nibh. Donec feugiat blandit ipsum non iaculis. Vestibulum tincidunt tristique ligula in sollicitudin. Aenean bibendum magna dui, dignissim vestibulum sapien pulvinar eu. Nulla lobortis at augue in faucibus. Nullam urna nulla, imperdiet et ex ut, varius fringilla lorem.

Nunc non faucibus erat, eu volutpat risus. Vestibulum a lorem ex. Sed efficitur tortor sapien, a maximus dui accumsan at. Vivamus id augue in lorem tempor fringilla. Phasellus eu feugiat quam. Curabitur eu gravida nulla. Proin efficitur blandit ante, quis eleifend neque euismod vel. Suspendisse ut tellus risus. Quisque suscipit nisl non orci suscipit pretium. Donec et nisl suscipit quam lobortis ornare. Nam lacinia, eros eget dignissim auctor, lacus ipsum hendrerit ante, ut vulputate urna est ut odio. Aenean eu quam fringilla, suscipit erat vitae, condimentum sem. Nullam ut tortor ex.
