package main

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/componhead/gator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	var username string
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	username = cmd.Args[0]
	ctx := context.Background()
	u, err := s.db.GetUserByName(ctx, username)
	if (err != nil || u == database.User{}) {
		return fmt.Errorf("user %s doesn't exist", username)
	}
	err = s.cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}
	fmt.Printf("Username %s setted\n", username)
	return nil
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.registeredCommands[cmd.Name]
	if !ok {
		return errors.New("command not found")
	}
	return f(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.registeredCommands[name] = f
}

func handlerRegister(s *state, cmd command) error {
	var name string
	if len(cmd.Args) < 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	name = cmd.Args[0]
	ctx := context.Background()
	u, err := s.db.GetUserByName(ctx, name)
	if u.ID != uuid.Nil {
		return fmt.Errorf("user %s already exists", name)
	}
	if err != nil && err.Error() != "sql: no rows in result set" {
		return fmt.Errorf("couldn't get user %s: %w", name, err)
	}
	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	createdUser, err := s.db.CreateUser(ctx, userParams)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}
	s.cfg.SetUser(createdUser.Name)
	fmt.Printf("User created: %+v\n", createdUser)
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't truncate user table: %w", err)
	}
	return nil
}

func handlerUsers(s *state, cmd command) error {
	usrs, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get users: %w", err)
	}
	var userRows []string
	for _, u := range usrs {
		var userRow []string
		userRow = append(userRow, fmt.Sprintf("* %s", u.Name))
		if u.Name == s.cfg.CurrentUserName {
			userRow = append(userRow, " (current)")
		}
		userRows = append(userRows, strings.Join(userRow, ""))
	}
	for _, userRow := range userRows {
		fmt.Println(userRow)
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	rssFeed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", rssFeed)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}
	for _, f := range feeds {
		u, err := s.db.GetUser(context.Background(), f.UserID)
		if err != nil {
			return err
		}
		fmt.Printf("feed name: %v\n", f.Name)
		fmt.Printf("feed url: %v\n", f.Url)
		fmt.Printf("feed user name: %v\n", u.Name)
	}
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <feed_name> <feed_url>", cmd.Name)
	}
	ctx := context.Background()
	u, err := s.db.GetUserByName(ctx, s.cfg.CurrentUserName)
	if err != nil {
		return err
	}
	args := database.AddFeedParams{
		ID:        uuid.New(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    u.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	f, err := s.db.AddFeed(ctx, args)
	if err != nil {
		return err
	}
	fmt.Printf("feed: %+v\n", f)
	return nil
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", res.StatusCode)
	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var rss RSSFeed
	if err := xml.Unmarshal(b, &rss); err != nil {
		return nil, err
	}
	rss.Channel.Description = html.UnescapeString(rss.Channel.Description)
	rss.Channel.Title = html.UnescapeString(rss.Channel.Title)
	for idx, item := range rss.Channel.Item {
		rss.Channel.Item[idx].Description = html.UnescapeString(item.Description)
		rss.Channel.Item[idx].Title = html.UnescapeString(item.Title)
	}
	return &rss, nil
}
