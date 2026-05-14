package db

import (
	"context"
	"database/sql"
	"log"
	"math/rand"
	"social/internal/store"

	"github.com/brianvoe/gofakeit/v6"
)

func Seed(store store.Storage, db *sql.DB) {
	ctx := context.Background()

	users := generateUsers(100)
	tx, _ := db.BeginTx(ctx, nil)

	for _, user := range users {
		if err := store.Users.Create(ctx, tx, user); err != nil {
			_ = tx.Rollback()
			log.Println("Error creating user:", err)
			return
		}
	}

	tx.Commit()

	posts := generatePosts(200, users)

	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("Error creating post:", err)
			return
		}
	}

	comments := generateComments(500, users, posts)

	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Println("Error creating comment:", err)
			return
		}
	}

	log.Println("Seeding complete")
}

func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)

	for i := 0; i < num; i++ {
		users[i] = &store.User{
			Username: gofakeit.Username(), // e.g., "coolcoder99"
			Email:    gofakeit.Email(),    // e.g., "sarah.smith@gmail.com"
			Role: store.Role{
				Name: "user",
			},
		}
	}

	return users
}

func generatePosts(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)

	for i := 0; i < num; i++ {

		user := users[rand.Intn(len(users))]

		randomTags := []string{
			gofakeit.Word(),
			gofakeit.Word(),
		}

		posts[i] = &store.Post{
			UserID:  user.ID,
			Title:   gofakeit.Sentence(5),
			Content: gofakeit.Paragraph(1, 4, 10, "\n"),
			Tags:    randomTags,
		}
	}

	return posts
}

func generateComments(num int, users []*store.User, posts []*store.Post) []*store.Comment {
	cms := make([]*store.Comment, num)

	for i := 0; i < num; i++ {
		cms[i] = &store.Comment{
			PostID:  posts[rand.Intn(len(posts))].ID,
			UserID:  users[rand.Intn(len(users))].ID,
			Content: gofakeit.Sentence(8),
		}
	}

	return cms
}
