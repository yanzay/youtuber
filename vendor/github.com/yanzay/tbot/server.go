package tbot

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Server is a telegram bot server. Looks and feels like net/http.
type Server struct {
	bot         *tgbotapi.BotAPI
	mux         Mux
	middlewares []Middleware
}

type Middleware func(HandlerFunction) HandlerFunction

// NewServer creates new Server with Telegram API Token
// and default /help handler
func NewServer(token string) (*Server, error) {
	tbot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	server := &Server{
		bot: tbot,
		mux: NewDefaultMux(),
	}

	server.HandleFunc("/help", server.HelpHandler)

	return server, nil
}

func (s *Server) AddMiddleware(mid Middleware) {
	s.middlewares = append(s.middlewares, mid)
}

// ListenAndServe starts Server, returns error on failure
func (s *Server) ListenAndServe() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := s.bot.GetUpdatesChan(u)
	if err != nil {
		return err
	}
	for update := range updates {
		go s.processMessage(update.Message)
	}
	return nil
}

// HandleFunc delegates HandleFunc to the current Mux
func (s *Server) HandleFunc(path string, handler HandlerFunction, description ...string) {
	s.mux.HandleFunc(path, handler, description...)
}

// Handle is a shortcut for HandleFunc to reply just with static text,
// "description" is for "/help" handler.
func (s *Server) Handle(path string, reply string, description ...string) {
	f := func(m Message) {
		m.Reply(reply)
	}
	s.HandleFunc(path, f, description...)
}

func (s *Server) HandleFile(handler HandlerFunction, description ...string) {
	s.mux.HandleFile(handler, description...)
}

// HandleDefault delegates HandleDefault to the current Mux
func (s *Server) HandleDefault(handler HandlerFunction, description ...string) {
	s.mux.HandleDefault(handler, description...)
}

func (s *Server) processMessage(message *tgbotapi.Message) {
	if message == nil {
		return
	}
	log.Printf("[TBot] %s %s: %s", message.From.FirstName, message.From.LastName, message.Text)
	var handler *Handler
	var data MessageVars
	if message.Document != nil {
		url, err := s.bot.GetFileDirectURL(message.Document.FileID)
		if err != nil {
			log.Printf("[TBot] Error: %s", err.Error())
			return
		}
		handler = s.mux.FileHandler()
		data = map[string]string{"url": url}
	} else {
		message.Text = s.trimBotName(message.Text)
		handler, data = s.mux.Mux(message.Text)
	}
	if handler == nil {
		return
	}
	f := handler.f
	for _, mid := range s.middlewares {
		f = mid(f)
	}
	m := Message{*message, data, make(chan *ReplyMessage), make(chan struct{})}
	go func() {
		f(m)
		m.close <- struct{}{}
	}()
	for {
		select {
		case reply := <-m.replies:
			err := s.dispatchMessage(message.Chat, reply)
			if err != nil {
				log.Printf("Error dispatching message: %q", err)
			}
		case <-m.close:
			return
		}
	}
}

func (s *Server) dispatchMessage(chat *tgbotapi.Chat, reply *ReplyMessage) error {
	_, err := s.bot.Send(reply.msg)
	return err
}

func (s *Server) trimBotName(message string) string {
	parts := strings.SplitN(message, " ", 2)
	command := parts[0]
	command = strings.TrimSuffix(command, "@"+s.bot.Self.UserName)
	command = strings.TrimSuffix(command, "@"+s.bot.Self.FirstName)
	parts[0] = command
	return strings.Join(parts, " ")
}
