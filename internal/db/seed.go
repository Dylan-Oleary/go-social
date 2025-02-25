package db

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"

	"github.com/Dylan-Oleary/go-social/internal/store"
)

var mockUsernames = []string{
	"AlphaWolf", "BetaBear", "CyberNinja", "DarkPhoenix", "EchoFalcon",
	"FrostByte", "GhostRider", "HyperTiger", "IronClaw", "JetStream",
	"KnightHawk", "LunarFox", "MatrixX", "NeonViper", "OmegaRogue",
	"PixelPioneer", "QuantumQuake", "RogueSpecter", "ShadowFang", "ThunderBolt",
	"UltraVortex", "ViperStrike", "WarpRunner", "XenonBlast", "YellowJacket",
	"ZetaStorm", "AquaBlade", "BlazeFury", "CrimsonShade", "DaggerFrost",
	"EmeraldLancer", "FlameHunter", "GlitchMage", "HorizonDrake", "InfernoKnight",
	"JungleWraith", "KineticWave", "LightningRogue", "MidnightEcho", "NovaDrift",
	"OblivionSeeker", "PhantomStrike", "QuasarRanger", "RebelShadow", "StormBreaker",
	"TurbineFalcon", "VectorPhantom", "WarpShifter", "X-Caliber", "ZenithWarrior",
}

var mockPosts = []store.Post{
	{Title: "Morning Motivation", Content: "Start your day with a positive mindset!", Tags: []string{"motivation", "morning", "self-improvement"}},
	{Title: "Tech Trends", Content: "AI is evolving faster than ever.", Tags: []string{"technology", "AI", "trends"}},
	{Title: "Workout Tips", Content: "Stay hydrated and stretch before exercise.", Tags: []string{"fitness", "health", "workout"}},
	{Title: "Book Recommendations", Content: "Check out this thriller novel!", Tags: []string{"books", "reading", "thriller"}},
	{Title: "Coding Challenge", Content: "Solve this algorithm problem today.", Tags: []string{"programming", "coding", "challenge"}},
	{Title: "Healthy Eating", Content: "Try adding more greens to your diet.", Tags: []string{"health", "nutrition", "food"}},
	{Title: "Weekend Plans", Content: "What’s everyone doing this weekend?", Tags: []string{"weekend", "plans", "fun"}},
	{Title: "Movie Night", Content: "Rewatching an old classic tonight!", Tags: []string{"movies", "entertainment", "classic"}},
	{Title: "Travel Goals", Content: "Dreaming of my next destination.", Tags: []string{"travel", "adventure", "dream"}},
	{Title: "Music Vibes", Content: "This playlist is getting me through the day.", Tags: []string{"music", "playlist", "vibes"}},
	{Title: "Daily Reflection", Content: "Taking a moment to appreciate life.", Tags: []string{"self-care", "mindfulness", "life"}},
	{Title: "Game On", Content: "Who’s up for some online matches?", Tags: []string{"gaming", "multiplayer", "fun"}},
	{Title: "Coffee Break", Content: "Nothing beats a fresh cup in the morning.", Tags: []string{"coffee", "morning", "energy"}},
	{Title: "Mindfulness", Content: "Take deep breaths and stay present.", Tags: []string{"mindfulness", "meditation", "calm"}},
	{Title: "Coding Fun", Content: "Experimenting with a new framework!", Tags: []string{"coding", "programming", "development"}},
	{Title: "Productivity Hacks", Content: "Try the Pomodoro technique for focus.", Tags: []string{"productivity", "time-management", "focus"}},
	{Title: "DIY Project", Content: "Building something cool this weekend!", Tags: []string{"DIY", "crafts", "weekend"}},
	{Title: "Gadget Review", Content: "Testing out the latest smartphone.", Tags: []string{"technology", "gadgets", "review"}},
	{Title: "Motivational Quote", Content: "You miss 100% of the shots you don’t take.", Tags: []string{"motivation", "quotes", "success"}},
	{Title: "Fitness Goals", Content: "Hit a new personal best today!", Tags: []string{"fitness", "workout", "progress"}},
	{Title: "Startup Ideas", Content: "Brainstorming the next big thing.", Tags: []string{"business", "startups", "entrepreneurship"}},
	{Title: "Hiking Adventures", Content: "Exploring a new trail this weekend!", Tags: []string{"hiking", "outdoors", "nature"}},
	{Title: "Gaming Update", Content: "Just unlocked a rare achievement!", Tags: []string{"gaming", "achievement", "fun"}},
	{Title: "Cooking Experiment", Content: "Tried making homemade pasta today!", Tags: []string{"cooking", "food", "homemade"}},
	{Title: "Photography Love", Content: "Captured an amazing sunset shot.", Tags: []string{"photography", "nature", "sunset"}},
	{Title: "Life Advice", Content: "Focus on progress, not perfection.", Tags: []string{"self-improvement", "life", "growth"}},
	{Title: "Career Growth", Content: "Working on improving my skills.", Tags: []string{"career", "growth", "learning"}},
	{Title: "Coding Debugging", Content: "Finally fixed that annoying bug!", Tags: []string{"coding", "debugging", "programming"}},
	{Title: "Random Thoughts", Content: "Is time travel theoretically possible?", Tags: []string{"science", "space", "theory"}},
	{Title: "Gadget Wishlist", Content: "Hoping to get a new laptop soon.", Tags: []string{"technology", "wishlist", "gadgets"}},
	{Title: "Music Discovery", Content: "Found an underrated band today!", Tags: []string{"music", "discovery", "bands"}},
	{Title: "Science Breakthrough", Content: "A new exoplanet was discovered!", Tags: []string{"science", "space", "discovery"}},
	{Title: "Life Updates", Content: "Big changes coming soon!", Tags: []string{"life", "updates", "changes"}},
	{Title: "Pet Stories", Content: "My dog did the funniest thing today.", Tags: []string{"pets", "dogs", "funny"}},
	{Title: "Morning Routine", Content: "Trying to build better habits.", Tags: []string{"habits", "morning", "routine"}},
	{Title: "Space Exploration", Content: "Mars mission update is exciting!", Tags: []string{"space", "NASA", "exploration"}},
	{Title: "Life Challenges", Content: "Pushing through tough times.", Tags: []string{"life", "motivation", "challenges"}},
	{Title: "Healthy Mindset", Content: "Your thoughts shape your reality.", Tags: []string{"mindset", "motivation", "positivity"}},
	{Title: "Stock Market News", Content: "Big moves happening today!", Tags: []string{"finance", "stocks", "market"}},
	{Title: "Dream Big", Content: "Anything is possible with effort.", Tags: []string{"motivation", "success", "dreams"}},
	{Title: "Nature Love", Content: "Fresh air makes everything better.", Tags: []string{"nature", "outdoors", "fresh-air"}},
	{Title: "Home Decor", Content: "Redecorating my living space.", Tags: []string{"home", "decor", "design"}},
	{Title: "Favorite Food", Content: "Pizza is always a good idea.", Tags: []string{"food", "favorites", "pizza"}},
	{Title: "Mental Health", Content: "Self-care is important!", Tags: []string{"mental-health", "wellness", "self-care"}},
	{Title: "Science Wonders", Content: "Black holes are mind-boggling.", Tags: []string{"science", "space", "black-holes"}},
	{Title: "Workout Challenge", Content: "Who’s up for a 30-day challenge?", Tags: []string{"fitness", "challenge", "workout"}},
	{Title: "Music Playlist", Content: "This song is on repeat today.", Tags: []string{"music", "playlist", "favorite"}},
	{Title: "Photography Adventures", Content: "Exploring new angles today.", Tags: []string{"photography", "adventure", "creativity"}},
	{Title: "Book Club", Content: "Starting a new book this week.", Tags: []string{"books", "reading", "literature"}},
	{Title: "Hidden Skills", Content: "Discovered a new talent today!", Tags: []string{"skills", "learning", "growth"}},
}

var mockComments = []string{
	"Great post!", "This is so helpful, thanks!", "I totally agree!", "Well said!",
	"Interesting perspective!", "Can you elaborate on that?", "Love this!", "Thanks for sharing!",
	"I learned something new today.", "This made my day!", "Brilliant insight!", "So true!",
	"I’ve been thinking about this too.", "I never thought of it that way!", "Amazing read!",
	"Couldn’t have said it better myself.", "This changed my perspective.", "This is spot on!",
	"Such a good reminder.", "I needed to hear this today.", "Very insightful!", "I have a question about this.",
	"You nailed it!", "Mind blown!", "Incredibly well explained.", "This is so relatable!",
	"I’ll have to try this.", "I appreciate the detail.", "Where can I learn more about this?",
	"Thanks for the advice!", "This is exactly what I was looking for.", "Totally bookmarking this.",
	"Wow, this is amazing!", "I disagree, but I respect your opinion.", "Really well-written!",
	"You make a great point!", "I wish I knew this earlier.", "Very informative!", "Fantastic breakdown!",
	"I appreciate your perspective.", "This is super useful!", "What a great take on this.",
	"I never knew that!", "I have mixed feelings about this.", "This deserves more attention!",
	"I can’t wait to try this.", "You always share great insights!", "This was eye-opening!",
	"This is gold!", "I love how you explained this.", "This should be common knowledge!",
	"I see this differently, but I appreciate your view.", "I had the same thought!",
	"I’m going to share this with my friends.", "Thanks for putting this out there!",
	"This really resonated with me.", "Can you recommend more on this topic?",
	"This is the kind of content I love.", "I completely understand where you’re coming from.",
	"I needed this today!", "This should be taught in schools!", "So much value in this post.",
	"I’ve struggled with this too.", "This made me think.", "Such a refreshing take!",
	"This will help me a lot.", "You explained it so well!", "I was just thinking about this!",
	"This helped me see things differently.", "I’d love to hear more about this.",
	"You’re onto something here.", "This perspective is fascinating!", "This is a game-changer!",
	"Thank you for this!", "This deserves more attention.", "I have a slightly different take on this.",
	"Your content is always top-notch!", "So helpful, as always!", "I’m saving this for later.",
	"I feel the same way!", "This made a lot of sense!", "I never realized this before.",
	"This is exactly what I needed.", "Absolutely love this!", "I think more people should see this.",
	"You really put things into perspective.", "This was very thought-provoking!",
	"This is such an underrated topic!", "I appreciate the clarity!", "Very well articulated!",
	"I have a different experience with this.", "This is so practical and useful!",
	"I’ll definitely be implementing this.", "I had an aha moment reading this!", "This got me thinking.",
	"I completely resonate with this.", "I appreciate the honesty here.", "This is something I’ll keep in mind.",
	"I can relate to this so much!", "I needed this reminder.", "This post stands out!",
	"I’d love to hear more examples!", "Such a simple yet powerful idea.", "Absolutely spot on!",
	"This gave me a new way to look at things.", "Very well put!", "So many great takeaways here!",
}

func Seed(store store.Storage) error {
	log.Println("Beginning database seed...")

	ctx := context.Background()
	users := generateUsers()

	for _, user := range users {
		if err := store.Users.Create(ctx, user); err != nil {
			log.Fatal("Error creating mock user: ", err)
			return err
		}
	}

	posts := generatePosts(100, users)

	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Fatal("Error creating mock post: ", err)
			return err
		}

		comments := generateComments(rand.IntN(5), post, users)

		for _, comment := range comments {
			if err := store.Comments.Create(ctx, comment); err != nil {
				log.Fatal("Error creating mock comment: ", err)
				return err
			}
		}
	}

	log.Println("Finished database seed!")

	return nil
}

func generateComments(num int, post *store.Post, users []*store.User) []*store.Comment {
	comments := make([]*store.Comment, num)

	for i := range num {
		user := users[rand.IntN(len(users))]
		comments[i] = &store.Comment{
			Content: mockComments[rand.IntN(len(mockComments))],
			PostID:  post.ID,
			UserID:  user.ID,
		}
	}

	return comments
}

func generatePosts(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)

	for i := range num {
		user := users[rand.IntN(len(users))]
		mockPost := mockPosts[rand.IntN(len(mockPosts))]
		posts[i] = &store.Post{
			Title:   mockPost.Title,
			Content: mockPost.Content,
			Tags:    mockPost.Tags,
			UserID:  user.ID,
		}
	}

	return posts
}

func generateUsers() []*store.User {
	users := make([]*store.User, len(mockUsernames))

	for i, mockUsername := range mockUsernames {
		users[i] = &store.User{
			Username: mockUsername,
			Email:    fmt.Sprintf("%s@gmail.com", mockUsername),
			Password: "password",
		}
	}

	return users
}
